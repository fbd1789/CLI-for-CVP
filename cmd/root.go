package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var (
	tokenPath string
	urlPath   string
)

var rootCmd = &cobra.Command{
	Use:   "cvaas-cli",
	Short: "CLI pour interagir avec Arista CloudVision",
	Long:  "Outil CLI permettant de créer des workspaces, tags et exécuter des opérations via cvaas-cli",
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println("Erreur:", err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().StringVar(&tokenPath, "token", "", "Chemin vers le fichier token")
	rootCmd.PersistentFlags().StringVar(&urlPath, "url", "", "Chemin vers le fichier URL")

	rootCmd.MarkPersistentFlagRequired("token")
	rootCmd.MarkPersistentFlagRequired("url")

	rootCmd.AddCommand(createCmd)
	rootCmd.AddCommand(getCmd)
	rootCmd.AddCommand(runCmd)
}