package cmd

import (
	"cvaas_cli/internal"
	"fmt"
	"time"

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
		ctx, cancel, conn := internal.Connect(tokenPath, urlPath)
		defer cancel()
		defer conn.Close()

		workspaceID := "ws-008"
		requestID := fmt.Sprintf("req-%d", time.Now().Unix())
		tagLabel := "eos_version"
		tagValue := "EOS-64"
		elementType, elementSubType := 1, 1

		internal.CreateWorkspace(ctx, conn, workspaceID, requestID, "ADDTOTO")
		internal.CreateTag(ctx, conn, workspaceID, tagLabel, tagValue, elementType, elementSubType)

		devices := internal.ReadInventory(ctx, conn, modelFilter, mlagFilter, danzFilter)
		for _, d := range devices {
			if d.Model == "cEOSLab" {
				internal.AssignTagToDevice(ctx, conn, workspaceID, d.DeviceID, tagLabel, tagValue, elementType, elementSubType)
			}
		}
	},
}

func init() {
	runCmd.AddCommand(runProcessCmd)
}