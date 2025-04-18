package cmd

import (
	// "cvaas_cli/internal"
	// "fmt"
	// "time"

	"github.com/spf13/cobra"
)

var runCmd = &cobra.Command{
	Use:   "run",
	Short: "Exécuter un processus complet dans cvaas-cli",
	Hidden: true, // Pour masquer la cmd dans le help
}

var runProcessCmd = &cobra.Command{
	Use:   "process",
	Short: "Créer workspace, tag, et assigner aux cEOSLab",
	Run: func(cmd *cobra.Command, args []string) {
		// ctx, cancel, conn := internal.Connect(tokenPath, urlPath)
		// defer cancel()
		// defer conn.Close()
	},
}

func init() {
	runCmd.AddCommand(runProcessCmd)
}