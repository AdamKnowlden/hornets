package main

import (
	//"context"
	//"encoding/json"

	"bufio"
	"context"
	"flag"
	"fmt"
	"log"
	"os/signal"
	"strings"
	"syscall"

	"os"
	keys "rpctest/lib/context"
	"rpctest/lib/merkle"
	merkle_dag "rpctest/lib/merkle/dag"

	"github.com/multiformats/go-multibase"

	"rpctest/lib/encryption/rsa"
)

func main() {
	ctx := context.Background()

	//privateKey := flag.String("private", "", "Private key")
	//publicKey := flag.String("public", "", "Public key")
	//address := flag.String("address", "", "Address")
	flag.Parse()

	//ctx = context.WithValue(ctx, keys.PrivateKey, privateKey)
	//ctx = context.WithValue(ctx, keys.PublicKey, publicKey)
	//ctx = context.WithValue(ctx, keys.Address, address)

	RunCommandWatcher(ctx)
}

func RunCommandWatcher(ctx context.Context) {
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	// Start a goroutine to wait for OS signals
	go func() {
		// Wait for a signal to be received on the channel
		<-sigChan

		Cleanup(ctx)
		os.Exit(0)
	}()

	// Create a scanner to read from os.Stdin
	scanner := bufio.NewScanner(os.Stdin)

	// Run the main loop indefinitely
	for {
		// Print a prompt and wait for user input
		scanner.Scan()

		// Get the user's command as a string
		command := strings.TrimSpace(scanner.Text())
		segments := strings.Split(command, " ")

		// Handle the user's command
		switch segments[0] {
		case "help":
			log.Println("Available Commands:")
			log.Println("generate")
			log.Println("dag")
			log.Println("shutdown")
		case "generate":
			GenerateKeys(ctx)
		case "parse":
			ParseKeys(ctx)
		case "dag":
			SendTestDag(ctx)
		case "shutdown":
			log.Println("Shutting down")
			Cleanup(ctx)
			return
		default:
			log.Printf("Unknown command: %s\n", command)
		}
	}
}

func Cleanup(ctx context.Context) {

}

func GenerateKeys(ctx context.Context) {
	privateKey, err := rsa.CreateKeyPair()
	if err != nil {
		fmt.Println("Failed to create private key")
		return
	}

	rsa.SaveKeyPairToFile(privateKey)
}

func SendTestDag(ctx context.Context) {
	dag, err := merkle_dag.CreateMerkleDag("C:/Users/Adam/Documents/Organizations/Akashic Record/nostr2.0/relayer", multibase.Base64)
	if err != nil {
		fmt.Println("Error: ", err)
		os.Exit(1)
	}

	mega, err := merkle.CreateMegaTree(dag)
	if err != nil {
		fmt.Println("Error: ", err)
		os.Exit(1)
	}

	merkle.SendMegaTree(mega)
}

func ParseKeys(ctx context.Context) {
	privateKey, err := rsa.ParsePrivateKeyFromFile("private.key")
	if err != nil {
		fmt.Println("Failed to parse private key")
		return
	}

	publicKey, err := rsa.ParsePublicKeyFromFile("public.pem")
	if err != nil {
		fmt.Println("Failed to parse public key")
		return
	}

	context.WithValue(ctx, keys.PrivateKey, privateKey)
	context.WithValue(ctx, keys.PublicKey, publicKey)

	fmt.Println("Keys have been parsed")
}
