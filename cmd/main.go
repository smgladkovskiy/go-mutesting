package cmd

import (
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:  "go-mutesting [flags] ...targets",
	Args: cobra.MinimumNArgs(1),
	RunE: run,
}

func Execute() error {
	registerFlags()
	registerMutators()

	return rootCmd.Execute()
}
