package cmd

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"sync"

	"github.com/smgladkovskiy/go-mutesting/pkg/errs"
	"github.com/smgladkovskiy/go-mutesting/pkg/infection"
	"github.com/smgladkovskiy/go-mutesting/pkg/models"
	"github.com/smgladkovskiy/go-mutesting/pkg/mutation"
	"github.com/smgladkovskiy/go-mutesting/pkg/parser"
	"github.com/smgladkovskiy/go-mutesting/pkg/utils"
	log "github.com/spacetab-io/logs-go/v2"
	"github.com/spf13/cobra"
	"github.com/zimmski/go-tool/importing"
)

const (
	md5ChecksumLength = 32
)

func run(cmd *cobra.Command, args []string) error {
	opts := fillOpts(cmd.PersistentFlags(), args)

	if err := log.Init("test", log.Config{Level: level(opts), Format: "text", NoColor: true}, "go-mutesting", "", os.Stdout); err != nil {
		return fmt.Errorf("log init error: %w", err)
	}

	log.Info().Msgf("run go-mutesting with %d jobs", opts.Exec.Jobs)

	files := importing.FilesOfArgs(opts.Remaining.Targets)
	if len(files) == 0 {
		return errs.ErrNoSuitableSourceFiles
	}

	if !opts.Files.ListFiles && opts.Files.PrintAST {
		return workWithFiles(opts, files)
	}

	mutationBlackList, err := getBlacklist(opts.Files.Blacklist)
	if err != nil {
		return err
	}

	mutators := models.GetMutators(registerMutators(), opts.Mutator.DisableMutators)

	tmpDir, err := ioutil.TempDir("", "go-mutesting-")
	if err != nil {
		panic(err)
	}

	log.Debug().Str("tmpDir", tmpDir).Msg("mutation saved in tmpDir")

	var (
		execs []string
		wg    sync.WaitGroup
		jobs  = utils.GetJobs(opts)
		c     = make(chan string)
		stats = &models.MutationStats{}
	)

	if opts.Exec.Exec != "" {
		execs = strings.Split(opts.Exec.Exec, " ")
	}

	wg.Add(jobs)

	for runningJobs := 0; runningJobs < jobs; runningJobs++ {
		go func(c chan string) {
			for {
				file, more := <-c
				if !more {
					wg.Done()

					return
				}

				if err := mutation.ProcessFile(opts, tmpDir, file, mutators, mutationBlackList, execs, stats); err != nil {
					log.Error().Err(err).Msg("ProcessFile error")

					// Stop execution right away.
					wg.Done()

					return
				}
			}
		}(c)
	}

	for _, file := range files {
		c <- file
	}

	close(c)
	wg.Wait()

	if !opts.General.DoNotRemoveTmpFolder {
		if err := os.RemoveAll(tmpDir); err != nil {
			panic(err)
		}

		log.Debug().Str("tmpDir", tmpDir).Msg("Remove tmpDir")
	}

	if !opts.Exec.NoExec {
		log.Info().Msgf("The mutation score is %f (total is %d mutants: %d killed, %d escaped, %d skipped, %d duplicates)",
			stats.Score(), stats.Total(), stats.MutantsKilled, stats.MutantsEscaped, stats.MutantsSkipped, stats.Duplicated)

		if stats.Score() < opts.Test.Score {
			return errs.ErrTestScoreWasNotReached
		}
	} else {
		return errs.ErrNoExecWasExecuted
	}

	return nil
}

func level(opts models.Options) string {
	lvl := "info"
	if opts.General.Debug {
		lvl = "debug"
	}

	return lvl
}

func getBlacklist(blackList []string) (map[string]struct{}, error) {
	mutationBlackList := make(map[string]struct{})

	if len(blackList) == 0 {
		return mutationBlackList, nil
	}

	for _, f := range blackList {
		c, err := ioutil.ReadFile(f)
		if err != nil {
			return nil, fmt.Errorf("cannot read blacklist file %q: %w", f, err)
		}

		for _, line := range strings.Split(string(c), "\n") {
			if line == "" {
				continue
			}

			if len(line) != md5ChecksumLength {
				return nil, fmt.Errorf("%q is %w", line, errs.ErrNotMD5Checksum)
			}

			mutationBlackList[line] = struct{}{}
		}
	}

	return mutationBlackList, nil
}

func workWithFiles(opts models.Options, files []string) error {
	if opts.Files.ListFiles {
		for _, file := range files {
			log.Info().Str("file", file).Send()
		}

		return nil
	}

	if opts.Files.PrintAST {
		for _, file := range files {
			log.Info().Str("file", file).Send()

			src, _, err := parser.ParseFile(file)
			if err != nil {
				return fmt.Errorf("could not open file %s: %w", file, err)
			}

			infection.Results(src)
		}

		return nil
	}

	return nil
}
