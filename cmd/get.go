package cmd

import (
	"fmt"

	"cvaas_cli/internal"

	"github.com/spf13/cobra"
)
var modelFilter string

var workspaceStateFilter string


var getCmd = &cobra.Command{
	Use:   "get",
	Short: "Récupérer des ressources depuis cvaas-cli",
}

var getDevicesCmd = &cobra.Command{
	Use:   "devices",
	Short: "Afficher l'inventaire des devices",
	Run: func(cmd *cobra.Command, args []string) {
		ctx, cancel, conn := internal.Connect(tokenPath, urlPath)
		defer cancel()
		defer conn.Close()

		devices := internal.ReadInventory(ctx, conn, modelFilter)
		for _, d := range devices {
			fmt.Printf("📟 %s (%s) - %s\n", d.Hostname, d.DeviceID, d.Model)
		}
	},
}

var getWorkspacesCmd = &cobra.Command{
	Use:   "workspaces",
	Short: "Afficher les workspaces filtrés par état",
	Run: func(cmd *cobra.Command, args []string) {
		ctx, cancel, conn := internal.Connect(tokenPath, urlPath)
		defer cancel()
		defer conn.Close()

		workspaces := internal.GetWorkspacesByState(ctx, conn, workspaceStateFilter)
		for _, w := range workspaces {
			fmt.Printf("🧪 %s (%s) - State: %s\n", w.DisplayName, w.ID, w.State)
		}
	},
}


func init() {
	getCmd.AddCommand(getDevicesCmd)
	getDevicesCmd.Flags().StringVar(&modelFilter, "model", "", "Filtrer par modèle (ex: cEOSLab)")
	getCmd.AddCommand(getWorkspacesCmd)
	getWorkspacesCmd.Flags().StringVar(&workspaceStateFilter, "state", "NONE", "Filtrer les workspaces par état (UNSPECIFIED, PENDING, SUBMITTED, ABANDONED, CONFLICTS, ROLLED_BACK)")

}