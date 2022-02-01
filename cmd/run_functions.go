package cmd

import (
	"fmt"

	"github.com/smgladkovskiy/go-mutesting/pkg/models"
	"github.com/smgladkovskiy/go-mutesting/pkg/mutator/branch"
	"github.com/smgladkovskiy/go-mutesting/pkg/mutator/expression"
	"github.com/smgladkovskiy/go-mutesting/pkg/mutator/statement"
	"github.com/smgladkovskiy/go-mutesting/pkg/utils"
	"github.com/spf13/pflag"
)

// nolint: lll
func registerFlags() {
	rootCmd.PersistentFlags().SortFlags = false

	rootCmd.PersistentFlags().Bool("help", false, "Show this help message")

	rootCmd.PersistentFlags().Bool("debug", false, "[General] Debug log output")
	rootCmd.PersistentFlags().Bool("do-not-remove-tmp-folder", false, "[General] Do not remove the tmp folder where all mutations are saved to")
	rootCmd.PersistentFlags().Bool("verbose", false, "[General] Verbose log output")

	rootCmd.PersistentFlags().StringSlice("blacklist", []string{}, "[Files] List of MD5 checksums of mutations which should be ignored. Each checksum must end with a new line character")
	rootCmd.PersistentFlags().Bool("list-files", false, "[Files] List found files")
	rootCmd.PersistentFlags().Bool("print-ast", false, "[Files] Print the ASTs of all given files and exit")

	rootCmd.PersistentFlags().StringSlice("disable", []string{}, "[Mutator] Disable mutator by their name or using * as a suffix pattern")
	rootCmd.PersistentFlags().Bool("list-mutators", false, "[Mutator] list-mutators")

	rootCmd.PersistentFlags().String("match", "", "[Filter] Only functions are mutated that confirm to the arguments regex")

	rootCmd.PersistentFlags().String("exec", "", "[Exec] Execute this command for every mutation (by default the built-in exec command is used)")
	rootCmd.PersistentFlags().Bool("no-exec", false, "[Exec] Skip the built-in exec command and just generate the mutations")
	rootCmd.PersistentFlags().Int64("exec-timeout", 100, "[Exec] Sets a timeout for the command execution (in seconds)")
	rootCmd.PersistentFlags().Int("jobs", utils.MaxJobs(), "[Exec] Allow N jobs at once")

	rootCmd.PersistentFlags().Bool("test-recursive", false, "[Test] Defines if the executer should test recursively")
	rootCmd.PersistentFlags().Float64("score", 0, "[Test] Minimal acceptable scores value. If result is less than given, exit code will be non-zero")
}

func registerMutators() (models.MutatorLookup, error) {
	ml := models.MutatorLookup{}
	mutatorItems := []models.MutatorItem{
		{"branch/case", branch.MutatorCase},
		{"branch/else", branch.MutatorElse},
		{"branch/if", branch.MutatorIf},
		{"expression/comparison", expression.MutatorComparison},
		{"expression/remove", expression.MutatorRemoveTerm},
		{"statement/remove", statement.MutatorRemoveStatement},
	}

	for _, mi := range mutatorItems {
		if err := ml.Register(mi.Name, mi.Mutator); err != nil {
			return nil, fmt.Errorf("mutator %s register error: %w", mi.Name, err)
		}
	}

	return ml, nil
}

func fillOpts(flags *pflag.FlagSet, args []string) models.Options {
	var opts models.Options

	opts.General.Debug, _ = flags.GetBool("debug")
	opts.General.Verbose, _ = flags.GetBool("verbose")
	opts.General.DoNotRemoveTmpFolder, _ = flags.GetBool("do-not-remove-tmp-folder")
	opts.General.FailOnly, _ = flags.GetBool("fail-only")

	opts.Files.Blacklist, _ = flags.GetStringSlice("blacklist")
	opts.Files.ListFiles, _ = flags.GetBool("list-files")
	opts.Files.PrintAST, _ = flags.GetBool("print-ast")

	dms, _ := flags.GetStringSlice("disable")

	for _, dm := range dms {
		opts.Mutator.DisableMutators = append(opts.Mutator.DisableMutators, models.MutatorName(dm))
	}

	opts.Mutator.ListMutators, _ = flags.GetBool("list-mutators")

	opts.Filter.Match, _ = flags.GetString("match")

	opts.Exec.Exec, _ = flags.GetString("exec")
	opts.Exec.NoExec, _ = flags.GetBool("no-exec")
	opts.Exec.Timeout, _ = flags.GetInt64("exec-timeout")
	opts.Exec.Jobs, _ = flags.GetInt("jobs")

	opts.Test.Recursive, _ = flags.GetBool("test-recursive")
	opts.Test.Score, _ = flags.GetFloat64("score")

	opts.Remaining.Targets = args

	return opts
}
