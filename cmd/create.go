package cmd

import (
	"bufio"
	"cvaas_cli/internal"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/spf13/cobra"
)

var createCmd = &cobra.Command{
	Use:   "create",
	Short: "Créer des ressources dans cvaas-cli",
	Hidden: true, // Pour masquer la cmd dans le help
}

var createWorkspaceCmd = &cobra.Command{
	Use:   "workspace",
	Short: "Créer un workspace",
	Run: func(cmd *cobra.Command, args []string) {
		workspaceID := prompt("🆔 Workspace ID: ")
		displayName := prompt("📛 Nom du workspace: ")
		requestID := fmt.Sprintf("req-%d", time.Now().Unix())

		ctx, cancel, conn := internal.Connect(tokenPath, urlPath)
		defer cancel()
		defer conn.Close()

		internal.CreateWorkspace(ctx, conn, workspaceID, requestID, displayName)
	},
}

var createTagCmd = &cobra.Command{
	Use:   "tag",
	Short: "Créer un tag",
	Run: func(cmd *cobra.Command, args []string) {
		workspaceID := prompt("🆔 Workspace ID: ")
		label := prompt("🏷️ Label : ")
		value := prompt("💬 Valeur : ")
		elementType := atoi(prompt("🔢 ElementType (ex: 1) : "))
		elementSubType := atoi(prompt("🔢 ElementSubType (ex: 1) : "))

		ctx, cancel, conn := internal.Connect(tokenPath, urlPath)
		defer cancel()
		defer conn.Close()

		internal.CreateTag(ctx, conn, workspaceID, label, value, elementType, elementSubType)
	},
}

func init() {
	createCmd.AddCommand(createWorkspaceCmd)
	createCmd.AddCommand(createTagCmd)
}

func prompt(msg string) string {
	fmt.Print(msg)
	reader := bufio.NewReader(os.Stdin)
	text, _ := reader.ReadString('\n')
	return strings.TrimSpace(text)
}

func atoi(s string) int {
	var i int
	fmt.Sscanf(s, "%d", &i)
	return i
}