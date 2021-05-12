package main

import (
	"dap/dap"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use: "dap",
}

var runCmd = &cobra.Command{
	Use: "run",
	Run: func(cmd *cobra.Command, args []string) {
		dap.Run(args[0])
	},
}

func main() {
	rootCmd.AddCommand(runCmd)
	rootCmd.Execute()
}
