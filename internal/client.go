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

const timeout = 30 * time.Second

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