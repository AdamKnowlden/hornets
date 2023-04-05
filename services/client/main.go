package main

import (
	//"context"
	//"encoding/json"

	"context"
	"flag"
	"fmt"

	"os"
	keys "rpctest/lib/context"
	"rpctest/lib/merkle"
	merkle_dag "rpctest/lib/merkle/dag"

	"github.com/multiformats/go-multibase"
)

func main() {
	ctx := context.Background()

	privateKey := flag.String("private", "", "Private key")
	publicKey := flag.String("public", "", "Public key")
	address := flag.String("address", "", "Address")
	flag.Parse()

	ctx = context.WithValue(ctx, keys.PrivateKey, privateKey)
	ctx = context.WithValue(ctx, keys.PublicKey, publicKey)
	ctx = context.WithValue(ctx, keys.Address, address)

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
