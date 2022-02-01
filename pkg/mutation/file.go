package mutation

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"

	"github.com/smgladkovskiy/go-mutesting/pkg/astutil"
	"github.com/smgladkovskiy/go-mutesting/pkg/models"
	"github.com/smgladkovskiy/go-mutesting/pkg/parser"
	"github.com/smgladkovskiy/go-mutesting/pkg/utils"
	log "github.com/spacetab-io/logs-go/v2"
)

func ProcessFile(
	opts models.Options,
	tmpDir, file string,
	mutators []models.MutatorItem,
	mutationBlackList map[string]struct{},
	execs []string,
	stats *models.MutationStats,
) error {
	log.Debug().Msgf("Mutate file %s", file)

	src, fset, pkg, info, err := parser.ParseAndTypeCheckFile(file)
	if err != nil {
		return fmt.Errorf("parse and type check file error: %w", err)
	}

	if err := os.MkdirAll(tmpDir+"/"+filepath.Dir(file), 0o755); err != nil {
		panic(err)
	}

	tmpFile := tmpDir + "/" + file

	originalFile := fmt.Sprintf("%s.original", tmpFile)

	if err := utils.CopyFile(file, originalFile); err != nil {
		panic(err)
	}

	log.Debug().Str("tempFile", originalFile).Msg("save original file to tempFile")

	if opts.Filter.Match != "" {
		m, err := regexp.Compile(opts.Filter.Match)
		if err != nil {
			return fmt.Errorf("match regex is not valid: %w", err)
		}

		for _, f := range astutil.Functions(src) {
			if m.MatchString(f.Name.Name) {
				mutations := Mutate(opts, mutators, mutationBlackList, pkg, info, file, fset, src, f, tmpFile, execs, stats)

				log.Info().Str("file", file).Int("mutations", mutations).Msg("infected")
			}
		}
	} else {
		mutations := Mutate(opts, mutators, mutationBlackList, pkg, info, file, fset, src, src, tmpFile, execs, stats)

		log.Info().Str("file", file).Int("mutations", mutations).Msg("infected")
	}

	return nil
}
