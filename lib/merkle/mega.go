package merkle

import (
	"context"
	"crypto/sha256"
	"fmt"
	"log"

	grpc "rpctest/lib/grpc"
	grpc_dag "rpctest/lib/grpc"
	grpc_merkle "rpctest/lib/grpc/merkle"
	merkle_dag "rpctest/lib/merkle/dag"
	mt "rpctest/lib/merkle/tree"

	merkle_tree "github.com/txaty/go-merkletree"

	"github.com/multiformats/go-multibase"
)

type MegaTree struct {
	Root []byte
	Dag  *merkle_dag.MerkleDag
	Tree *merkle_tree.MerkleTree
}

func CreateMegaTree(dag *merkle_dag.MerkleDag) (*MegaTree, error) {
	mega := &MegaTree{}
	mega.Dag = dag

	builder := mt.CreateTree()
	for _, node := range dag.Nodes {
		if len(node.Data) > 0 {
			builder.AddBlock(node.Data)
		}
	}

	tree, err := builder.Build()
	if err != nil {
		return nil, err
	}
	mega.Tree = tree

	//fmt.Println("Merkle Dag Root Hash")
	//fmt.Println(encoder.Encode(dag.Root))
	//fmt.Println("Merkle Tree Root Hash")
	//fmt.Println(encoder.Encode(mega.Tree.Root))

	combinedRoots := append(dag.Root, tree.Root...)
	newHash := sha256.Sum256(combinedRoots)
	mega.Root = newHash[:]

	//fmt.Println("Mega Tree Root Hash")
	//fmt.Println(encoder.Encode(mega.Root))

	return mega, nil
}

func SendMegaTree(megaTree *MegaTree, encoding ...multibase.Encoding) error {
	var e multibase.Encoding
	if len(encoding) > 0 {
		e = encoding[0]
	} else {
		e = multibase.Base64
	}
	encoder := multibase.MustNewEncoder(e)

	grpcMegaRoot := &grpc_dag.MerkleRoot{
		Root:     encoder.Encode(megaTree.Root),
		TreeRoot: encoder.Encode(megaTree.Tree.Root),
		DagRoot:  encoder.Encode(megaTree.Dag.Root),
	}

	conn, err := grpc.Open()
	if err != nil {
		log.Printf("Failed to open grpc connection: %v", err)
		return err
	}

	pClient, err := grpc_merkle.CreateService(conn)
	if err != nil {
		log.Printf("Failed to open grpc connection: %v", err)
		return err
	}

	client := *pClient

	response, err := client.SendMerkleRoot(context.Background(), grpcMegaRoot)
	if err != nil {
		log.Fatalf("SendMerkleRoot failed: %v", err)
		return err
	}
	fmt.Println(response.Message)

	for _, node := range megaTree.Dag.Nodes {
		grpcMerkleNode := &grpc_dag.MerkleNode{
			Hash:  encoder.Encode(node.Hash),
			Path:  node.Path,
			Type:  string(node.Type),
			Data:  node.Data,
			Links: node.Links,
		}

		response, err := client.SendMerkleNode(context.Background(), grpcMerkleNode)
		if err != nil {
			log.Fatalf("SendMerkleNode failed %v", err)
			return err
		}
		fmt.Println(response.Message)
	}

	reciept, err := client.NotifyCompletion(context.Background(), grpcMegaRoot)
	if err != nil {
		log.Fatalf("NotifyCompletion failed: %v", err)
		return err
	}
	fmt.Println("Reciept recieved: " + reciept.Root.Root)

	signedReciept := &grpc_dag.SignedReciept{
		Root:      reciept,
		Signature: "Signature",
	}

	response, err = client.SendSignedReciept(context.Background(), signedReciept)
	if err != nil {
		log.Fatalf("SendSignedReciept failed: %v", err)
		return err
	}
	fmt.Println(response.Message)

	return nil
}
