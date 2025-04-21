package internal

import (
	"bufio"
	"context"
	"crypto/tls"
	"fmt"
	"os"
	"strings"
	"time"

	"google.golang.org/grpc"
	// cvgrpc "github.com/aristanetworks/cloudvision-go/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/metadata"
)

// timeout définit la durée maximale pour l'établissement d'une connexion gRPC avec CVaaS.
// Ce délai est utilisé pour éviter les blocages lors d'une tentative de connexion.
const timeout = 30 * time.Second

// Connect établit une connexion gRPC sécurisée avec la plateforme CVaaS,
// en injectant le token d'authentification et l'URL du serveur depuis deux fichiers fournis.
//
// Les métadonnées de type "Authorization: Bearer <token>" sont ajoutées au contexte
// pour permettre l'authentification auprès de CVaaS.
//
// Paramètres :
//   - tokenPath : chemin vers un fichier contenant un token d'accès (une seule ligne).
//   - urlPath : chemin vers un fichier contenant l'URL du serveur CVaaS.
//
// Retourne :
//   - context.Context : contexte enrichi avec métadonnées pour les appels gRPC.
//   - context.CancelFunc : fonction à appeler pour annuler/fermer le contexte.
//   - *grpc.ClientConn : connexion gRPC active vers CVaaS.
//
// Panique :
//   - Si la lecture des fichiers échoue.
//   - Si la connexion gRPC ne peut pas être établie.
func Connect(tokenPath, urlPath string) (context.Context, context.CancelFunc, *grpc.ClientConn) {
	token, err := readLineFromFile(tokenPath)
	if err != nil {
		panic(fmt.Sprintf("Erreur lecture token : %v", err))
	}
	url, err := readLineFromFile(urlPath)
	if err != nil {
		panic(fmt.Sprintf("Erreur lecture URL : %v", err))
	}

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	ctx = metadata.AppendToOutgoingContext(ctx, "Authorization", "Bearer "+token)

	conn, err := grpc.DialContext(ctx, url,
		grpc.WithTransportCredentials(credentials.NewTLS(&tls.Config{})))
	if err != nil {
		panic(fmt.Sprintf("❌ Erreur connexion gRPC : %v", err))
	}
	return ctx, cancel, conn

}

// readLineFromFile lit la première ligne non vide d’un fichier donné et la retourne
// sous forme de chaîne nettoyée (sans espaces ou retours à la ligne).
//
// Utilisé principalement pour lire des tokens ou des URLs à partir de fichiers.
//
// Paramètres :
//   - filename : chemin absolu ou relatif du fichier à lire.
//
// Retourne :
//   - string : contenu de la première ligne du fichier, sans espaces superflus.
//   - error : une erreur si le fichier est vide ou inaccessible.
func readLineFromFile(filename string) (string, error) {
	file, err := os.Open(filename)
	if err != nil {
		return "", err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	if scanner.Scan() {
		return strings.TrimSpace(scanner.Text()), nil
	}
	return "", fmt.Errorf("fichier vide")
}