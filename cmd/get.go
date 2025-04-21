package cmd

import (
	"fmt"
	"os"

	"cvaas_cli/internal"

	"github.com/spf13/cobra"
)

// modelFilter est un flag CLI permettant de filtrer les devices par modèle spécifique
// (ex: "DCS-7280SR", "cEOSLab", etc.).
var modelFilter string

// workspaceStateFilter est un flag CLI permettant de filtrer les workspaces
// selon leur état (ex: "PENDING", "SUBMITTED", "ABANDONED", etc.).
var workspaceStateFilter string

// mlagFilter est un flag CLI indiquant si la commande "devices" doit retourner uniquement
// les équipements avec MLAG activé.
var mlagFilter bool

// danzFilter est un flag CLI indiquant si la commande "devices" doit retourner uniquement
// les équipements avec DANZ activé.
var danzFilter bool

// getCmd est la commande principale `get` du CLI, utilisée pour récupérer
// des ressources depuis la plateforme CVaaS (CloudVision-as-a-Service).
//
// Cette commande regroupe les sous-commandes `get devices` et `get workspaces`.
var getCmd = &cobra.Command{
	Use:   "get",
	Short: "Récupérer des ressources depuis cvaas-cli",
}

// getDevicesCmd est une sous-commande de `get` utilisée pour afficher l'inventaire
// des équipements disponibles sur CVaaS, avec la possibilité de filtrer par modèle,
// MLAG ou DANZ.
//
// Conflit logique : les flags `--mlag` et `--danz` sont mutuellement exclusifs, et
// leur combinaison est bloquée lors de l'exécution.
var getDevicesCmd = &cobra.Command{
	Use:   "devices",
	Short: "Afficher l'inventaire des devices",
	Run: func(cmd *cobra.Command, args []string) {
		// 🔒 Protection : on interdit les deux flags en même temps
		if mlagFilter && danzFilter {
			fmt.Println("❌ Les filtres --mlag et --danz ne peuvent pas être utilisés en même temps.")
			os.Exit(1)
		}
		ctx, cancel, conn := internal.Connect(tokenPath, urlPath)
		defer cancel()
		defer conn.Close()

		devices := internal.ReadInventory(ctx, conn, modelFilter, mlagFilter, danzFilter)

		for _, d := range devices {
			fmt.Printf("📟 %s (%s) - %s\n", d.Hostname, d.DeviceID, d.Model)
		}
	},
}

// getWorkspacesCmd est une sous-commande de `get` utilisée pour afficher les workspaces
// CVaaS filtrés selon un état particulier (ex: "SUBMITTED", "CONFLICTS").
//
// Utilise le flag `--state` pour sélectionner l'état voulu.
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

// init configure les sous-commandes et leurs flags associés pour la commande principale `get`.
func init() {
	getCmd.AddCommand(getDevicesCmd)
	getDevicesCmd.Flags().StringVar(&modelFilter, "model", "", "Filtrer par modèle (ex: cEOSLab)")
	getCmd.AddCommand(getWorkspacesCmd)
	getWorkspacesCmd.Flags().StringVar(&workspaceStateFilter, "state", "NONE", "Filtrer les workspaces par état (UNSPECIFIED, PENDING, SUBMITTED, ABANDONED, CONFLICTS, ROLLED_BACK)")
	getDevicesCmd.Flags().BoolVar(&mlagFilter, "mlag", false, "Afficher uniquement les devices avec MLAG activé")
	getDevicesCmd.Flags().BoolVar(&danzFilter, "danz", false, "Afficher uniquement les devices avec DANZ activé")
	
}