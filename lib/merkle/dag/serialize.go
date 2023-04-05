package dag

import (
	"encoding/json"
	//grpc_dag "rpctest/lib/grpc"

	cbor "github.com/fxamacker/cbor/v2"
	"github.com/multiformats/go-multibase"
)

type EncodedMerkleNode struct {
	Hash  string
	Path  string
	Type  string
	Links []string
	Data  string
}

type EncodedMerkleDag struct {
	Root  string
	Nodes map[string]*EncodedMerkleNode
}

func ToCBOR(dag *MerkleDag) ([]byte, error) {
	cborData, err := cbor.Marshal(dag)
	if err != nil {
		return nil, err
	}

	return cborData, nil
}

func ToJSONRaw(dag *MerkleDag) ([]byte, error) {
	jsonData, err := json.MarshalIndent(dag, "", "  ")
	if err != nil {
		return nil, err
	}

	return jsonData, nil
}

func ToJSON(dag *MerkleDag, encoding ...multibase.Encoding) ([]byte, error) {
	var e multibase.Encoding
	if len(encoding) > 0 {
		e = encoding[0]
	} else {
		e = multibase.Base64
	}
	encoder := multibase.MustNewEncoder(e)

	nodes := make(map[string]*EncodedMerkleNode, len(dag.Nodes))
	for hash, node := range dag.Nodes {
		nodes[hash] = &EncodedMerkleNode{
			Hash:  encoder.Encode(node.Hash),
			Path:  node.Path,
			Type:  string(node.Type),
			Links: node.Links,
			Data:  node.Data,
		}
	}

	encodedDag := &EncodedMerkleDag{
		Root:  encoder.Encode(dag.Root),
		Nodes: nodes,
	}

	jsonData, err := json.MarshalIndent(encodedDag, "", "  ")
	if err != nil {
		return nil, err
	}

	return jsonData, nil
}

/*
func ToGRPCMerkleDag(dag *MerkleDag) *grpc_dag.MerkleDag {
	root := ToGRPCMerkleNode(dag.Root)
	nodes := make([]*grpc_dag.MerkleNode, 0, len(dag.Nodes))
	for _, node := range dag.Nodes {
		nodes = append(nodes, ToGRPCMerkleNode(node))
	}
	return &grpc_dag.MerkleDag{
		Root:  root,
		Nodes: nodes,
	}
}

func ToGRPCMerkleNode(node *MerkleNode) *grpc_dag.MerkleNode {
	return &grpc_dag.MerkleNode{
		Cid:   node.CID,
		Name:  node.Name,
		Type:  string(node.Type),
		Data:  node.Data,
		Links: node.Links,
	}
}

func FromGRPCMerkleDag(in *grpc_dag.MerkleDag) *MerkleDag {
	// Create a new MerkleDag struct
	dag := &MerkleDag{}

	// Convert the gRPC root node to a MerkleNode struct
	root := FromGRPCMerkleNode(in.Root)
	dag.Root = root

	// Convert each gRPC node to a MerkleNode struct and add it to the MerkleDag
	for _, node := range in.Nodes {
		dagNode := FromGRPCMerkleNode(node)
		dag.Nodes[dagNode.CID] = dagNode
	}

	return dag
}

func FromGRPCMerkleNode(in *grpc_dag.MerkleNode) *MerkleNode {
	// Create a new MerkleNode struct
	node := &MerkleNode{}

	// Assign the CID, name, type, and data fields of the MerkleNode
	node.CID = in.Cid
	node.Name = in.Name
	node.Type = NodeType(in.Type)
	node.Data = in.Data

	// Assign the links of the MerkleNode
	for _, link := range in.Links {
		node.Links = append(node.Links, link)
	}

	return node
}
*/
