package mutation

import (
	"fmt"
	"go/ast"
	"go/token"
	"go/types"
	"os"
	"os/exec"
	"syscall"

	"github.com/smgladkovskiy/go-mutesting/pkg/errs"
	"github.com/smgladkovskiy/go-mutesting/pkg/models"
	"github.com/smgladkovskiy/go-mutesting/pkg/utils"
	"github.com/smgladkovskiy/go-mutesting/pkg/walk"
	log "github.com/spacetab-io/logs-go/v2"
)

func MutateExec(opts models.Options, pkg *types.Package, file string, mutationFile string, execs []string) (models.MutationResult, error) {
	if len(execs) == 0 {
		return MutationsWithoutExecs(opts, pkg, file, mutationFile)
	}

	return MutationsWithExecs(opts, pkg, file, execs, mutationFile)
}

func MutationsWithExecs(
	opts models.Options,
	pkg *types.Package,
	file string,
	execs []string,
	mutationFile string,
) (models.MutationResult, error) {
	log.Debug().Str("cmd", opts.Exec.Exec).Msg("Execute cmd for mutation")

	execCommand := exec.Command(execs[0], execs[1:]...)

	execCommand.Stderr = os.Stderr
	execCommand.Stdout = os.Stdout

	execCommand.Env = append(os.Environ(), []string{
		"MUTATE_CHANGED=" + mutationFile,
		fmt.Sprintf("MUTATE_DEBUG=%t", opts.General.Debug),
		"MUTATE_ORIGINAL=" + file,
		"MUTATE_PACKAGE=" + pkg.Path(),
		fmt.Sprintf("MUTATE_TIMEOUT=%d", opts.Exec.Timeout),
		fmt.Sprintf("MUTATE_VERBOSE=%t", opts.General.Verbose),
	}...)
	if opts.Test.Recursive {
		execCommand.Env = append(execCommand.Env, "TEST_RECURSIVE=true")
	}

	if err := execCommand.Start(); err != nil {
		return models.ResultEmpty, fmt.Errorf("exec command start error: %w", err)
	}

	// TODO timeout here
	if err := execCommand.Wait(); err != nil {
		return models.GetResultStatus(err)
	}

	return models.ResultMutantEscaped, nil
}

func MutationsWithoutExecs(opts models.Options, pkg *types.Package, file string, mutationFile string) (models.MutationResult, error) {
	log.Debug().Msg("Execute built-in exec command for mutation")

	execExitCode := 0

	diff, err := exec.Command("diff", "-u", file, mutationFile).CombinedOutput()
	if err != nil {
		if e, ok := err.(*exec.ExitError); ok {
			execExitCode = e.Sys().(syscall.WaitStatus).ExitStatus()
		} else {
			return models.ResultEmpty, fmt.Errorf("diff command error: %w", err)
		}
	}

	if execExitCode != 0 && execExitCode != 1 {
		log.Warn().Bytes("diff", diff).Msg("Could not execute diff on mutation file")

		return models.ResultEmpty, fmt.Errorf("diff command %w", errs.ErrWrongExitCode)
	}

	defer func() {
		_ = os.Rename(file+".tmp", file)
	}()

	if err := os.Rename(file, file+".tmp"); err != nil {
		panic(err)
	}

	if err := utils.CopyFile(mutationFile, file); err != nil {
		panic(err)
	}

	pkgName := pkg.Path()
	if opts.Test.Recursive {
		pkgName += "/..."
	}

	test, err := exec.Command("go", "test", "-timeout", fmt.Sprintf("%ds", opts.Exec.Timeout), pkgName).CombinedOutput()
	if err != nil {
		return models.GetResultStatus(err)
	}

	log.Debug().Bytes("tests", test).Msg("got test results")

	return models.ResultMutantEscaped, nil
}

func Mutate(
	opts models.Options,
	mutators []models.MutatorItem,
	mutationBlackList map[string]struct{},
	pkg *types.Package,
	info *types.Info,
	file string,
	fset *token.FileSet,
	src ast.Node,
	node ast.Node,
	tmpFile string,
	execs []string,
	stats *models.MutationStats,
) int {
	mutations := 0

	for _, m := range mutators {
		log.Debug().Str("mutator", m.Name).Msg("Running mutator")

		changed := walk.MutateWalk(pkg, info, node, m.Mutator)

		for {
			_, ok := <-changed
			if !ok {
				break
			}

			mutationFile := fmt.Sprintf("%s.%d", tmpFile, mutations)
			checksum, duplicate, err := utils.SaveAST(mutationBlackList, mutationFile, fset, src)
			if err != nil {
				log.Error().Err(err).Msg("INTERNAL ERROR")
				stats.UnknownResults++
			} else if duplicate {
				log.Debug().Str("mutationFile", mutationFile).Msg("Ignore duplicate mutationFile")

				stats.Duplicated++
			} else {
				log.Debug().Str("mutationFile", mutationFile).Str("checksum", checksum).Msg("Save mutation into mutationFile")

				if !opts.Exec.NoExec {
					result, err := MutateExec(opts, pkg, file, mutationFile, execs)

					msg := fmt.Sprintf("%q with checksum %s", mutationFile, checksum)

					switch result {
					case models.ResultMutantEscaped:
						log.Debug().Err(err).Msgf("PASS %s", msg)
						stats.MutantsEscaped++
					case models.ResultMutantKilled:
						log.Debug().Err(err).Msgf("FAIL %s", msg)
						stats.MutantsKilled++
					case models.ResultMutantSkipped:
						log.Debug().Err(err).Msgf("SKIP %s", msg)
						stats.MutantsSkipped++
					case models.ResultEmpty:
						log.Debug().Err(err).Msgf("UNKNOWN result for %s", msg)
						stats.UnknownResults++
					}
				}
			}

			changed <- true

			// Ignore original state
			<-changed
			changed <- true

			mutations++
		}
	}

	return mutations
}
