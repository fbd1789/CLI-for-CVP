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

// workspaceName est un flag CLI utilisé pour spécifier le nom du workspace à créer
// via la commande `create workspace`. Ce champ est obligatoire.
var workspaceName string

// createCmd est la commande principale `create` du CLI, utilisée pour créer
// des ressources sur la plateforme CVaaS (comme des workspaces).
var createCmd = &cobra.Command{
	Use:   "create",
	Short: "Créer des ressources dans cvaas-cli",
}

// WorkspaceEntry représente un enregistrement unique d'un workspace CVaaS,
// utilisé pour la tracabilité dans un fichier YAML local.
type WorkspaceEntry struct {
	WorkspaceID   string `yaml:"workspaceID"`
	RequestID     string `yaml:"RequestID"`
	WorkspaceName string `yaml:"workspaceName"`
}

// WorkspaceYAML est une structure regroupant plusieurs workspaces,
// utilisée pour sérialiser et désérialiser les données dans un fichier YAML.
type WorkspaceYAML struct {
	Workspace []WorkspaceEntry `yaml:"workspace"`
}

// createWorkspaceCmd est une sous-commande de `create` permettant de créer
// un nouveau workspace sur la plateforme CVaaS.
//
// La commande génère automatiquement un ID de workspace et un requestID,
// appelle l'API via gRPC, puis enregistre les métadonnées dans un fichier YAML local.
//
// Flag requis :
//   --name : nom du workspace à créer.
//
// Fichier local :
//   Les données du workspace sont stockées dans `data/workspace.yaml`.
//   Si le fichier ou le dossier n'existent pas, ils seront créés automatiquement.
//
// Panique / erreurs gérées :
//   - Si le nom du workspace n'est pas fourni
//   - Si une erreur survient lors de l'appel gRPC, de la sérialisation YAML,
//     ou de l'écriture dans le système de fichiers
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

// init configure la commande `create workspace` avec son flag obligatoire `--name`,
// l'attache à la commande `create`, puis enregistre `create` dans la racine du CLI.
func init() {
	createWorkspaceCmd.Flags().StringVar(&workspaceName, "name", "", "Nom du workspace à créer (obligatoire)")
	createCmd.AddCommand(createWorkspaceCmd)
	rootCmd.AddCommand(createCmd)
}