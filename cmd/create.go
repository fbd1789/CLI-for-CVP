package cmd

import (
	// "bufio"
	"cvaas_cli/internal"
	"fmt"
	"os"
	// "strings"
	"time"
	"path/filepath"
	"gopkg.in/yaml.v2"

	"github.com/spf13/cobra"
)

var workspaceName string

var createCmd = &cobra.Command{
	Use:   "create",
	Short: "Créer des ressources dans cvaas-cli",
}

type WorkspaceEntry struct {
	WorkspaceID   string `yaml:"workspaceID"`
	RequestID     string `yaml:"RequestID"`
	WorkspaceName string `yaml:"workspaceName"`
}

type WorkspaceYAML struct {
	Workspace []WorkspaceEntry `yaml:"workspace"`
}

var createWorkspaceCmd = &cobra.Command{
	Use:   "workspace",
	Short: "Créer un workspace",
	Run: func(cmd *cobra.Command, args []string) {
		if workspaceName == "" {
			fmt.Println("❌ Veuillez spécifier un nom avec --name")
			os.Exit(1)
		}

		workspaceID := fmt.Sprintf("ws-%d", time.Now().Unix())
		requestID := workspaceID

		ctx, cancel, conn := internal.Connect(tokenPath, urlPath)
		defer cancel()
		defer conn.Close()

		fmt.Printf("🆔 Workspace ID généré : %s\n", workspaceID)
		internal.CreateWorkspace(ctx, conn, workspaceID, requestID, workspaceName)


		entry := WorkspaceEntry{
			WorkspaceID:   workspaceID,
			RequestID:     requestID,
			WorkspaceName: workspaceName,
		}

		yamlPath := filepath.Join("data", "workspace.yaml")

		var workspaceFile WorkspaceYAML
		if content, err := os.ReadFile(yamlPath); err == nil {
			_ = yaml.Unmarshal(content, &workspaceFile)
		}

		workspaceFile.Workspace = append(workspaceFile.Workspace, entry)

		savedData, err := yaml.Marshal(&workspaceFile)
		if err != nil {
			fmt.Printf("❌ Erreur encodage YAML : %v\n", err)
			return
		}

		if err := os.MkdirAll("data", os.ModePerm); err != nil {
			fmt.Printf("❌ Erreur création dossier data : %v\n", err)
			return
		}

		if err := os.WriteFile(yamlPath, savedData, 0644); err != nil {
			fmt.Printf("❌ Erreur écriture workspace.yaml : %v\n", err)
			return
		}

		fmt.Println("✅ Workspace sauvegardé dans data/workspace.yaml")
	},
}

func init() {
	createWorkspaceCmd.Flags().StringVar(&workspaceName, "name", "", "Nom du workspace à créer (obligatoire)")
	createCmd.AddCommand(createWorkspaceCmd)
	rootCmd.AddCommand(createCmd)
}