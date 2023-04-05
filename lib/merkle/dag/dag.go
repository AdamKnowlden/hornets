package dag

import (
	"crypto/sha256"
	"encoding/json"
	"io/fs"
	"io/ioutil"
	"log"
	"path/filepath"
	"strconv"

	//cid "github.com/ipfs/go-cid"
	"github.com/multiformats/go-multibase"
)

const ChunkSize = 2048 * 1024 // 256 KiB // 1mb

type NodeType string

const (
	FileNodeType      NodeType = "file"
	ChunkNodeType     NodeType = "chunk"
	DirectoryNodeType NodeType = "directory"
)

type MerkleNode struct {
	Hash  []byte
	Path  string
	Type  NodeType
	Links []string
	Data  string
}

type MerkleDag struct {
	Root  []byte
	Nodes map[string]*MerkleNode
}

type NodeBuilder struct {
	Path     *string
	Links    [][]byte
	NodeType NodeType
	Content  []byte
}

type MerkleDagBuilder struct {
	Nodes map[string]*MerkleNode
}

func (parent *MerkleNode) CreateLink(child *MerkleNode, encoder multibase.Encoder) {
	parent.Links = append(parent.Links, encoder.Encode(child.Hash))
}

func ChunkFile(fileData []byte, chunkSize int) [][]byte {
	var chunks [][]byte
	fileSize := len(fileData)

	for i := 0; i < fileSize; i += chunkSize {
		end := i + chunkSize
		if end > fileSize {
			end = fileSize
		}
		chunks = append(chunks, fileData[i:end])
	}

	return chunks
}

func CreateBuilder(path *string) *NodeBuilder {
	builder := &NodeBuilder{
		Path:  path,
		Links: [][]byte{},
	}

	return builder
}

func (b *NodeBuilder) AddLink(hash []byte) {
	b.Links = append(b.Links, hash)
}

func (b *NodeBuilder) BuildNode(encoder multibase.Encoder) *MerkleNode {
	encodedContent := encoder.Encode(b.Content)

	links := []string{}

	for _, link := range b.Links {
		links = append(links, encoder.Encode(link))
	}

	nodeData := struct {
		Path    string
		Type    NodeType
		Links   []string
		Content string
	}{
		Path:    *b.Path,
		Type:    b.NodeType,
		Links:   links,
		Content: encodedContent,
	}

	serializedNodeData, err := json.Marshal(nodeData)
	if err != nil {
		log.Fatal("Failed to serialize node data: ", err)
	}

	//cid := cid.NewCidV1(cid.Raw, serializedNodeData)
	hash := sha256.Sum256(serializedNodeData)
	return &MerkleNode{
		Hash:  hash[:],
		Path:  *b.Path,
		Type:  b.NodeType,
		Links: links,
		Data:  encodedContent,
	}
}

func processEntry(entry fs.FileInfo, path *string, dag *MerkleDagBuilder, encoder multibase.Encoder) (*MerkleNode, error) {
	entryPath := filepath.Join(*path, entry.Name())
	builder := CreateBuilder(&entryPath)

	if entry.IsDir() {
		builder.NodeType = DirectoryNodeType
		builder.Content = nil

		entries, err := ioutil.ReadDir(entryPath)
		if err != nil {
			return nil, err
		}

		for _, entry := range entries {
			node, err := processEntry(entry, &entryPath, dag, encoder)
			if err != nil {
				return nil, err
			}

			builder.AddLink(node.Hash)
			dag.AddNode(node, encoder)
		}
	} else {
		fileData, err := ioutil.ReadFile(entryPath)
		if err != nil {
			return nil, err
		}

		builder.NodeType = FileNodeType
		builder.Content = nil

		fileChunks := ChunkFile(fileData, ChunkSize)

		if len(fileChunks) == 1 {
			builder.Content = fileChunks[0]
		} else {
			for i, chunk := range fileChunks {
				chunkEntryPath := filepath.Join(entryPath, strconv.Itoa(i))
				chunkBuilder := CreateBuilder(&chunkEntryPath)

				chunkBuilder.Content = chunk
				chunkBuilder.NodeType = ChunkNodeType

				chunkNode := chunkBuilder.BuildNode(encoder)

				builder.AddLink(chunkNode.Hash)
				dag.AddNode(chunkNode, encoder)
			}
		}
	}

	return builder.BuildNode(encoder), nil
}

func CreateMerkleDagBuilder() *MerkleDagBuilder {
	return &MerkleDagBuilder{
		Nodes: make(map[string]*MerkleNode),
	}
}

func (b *MerkleDagBuilder) AddNode(node *MerkleNode, encoder multibase.Encoder) {
	b.Nodes[encoder.Encode(node.Hash)] = node
}

func (b *MerkleDagBuilder) BuildMerkleDag(root []byte) *MerkleDag {
	return &MerkleDag{
		Nodes: b.Nodes,
		Root:  root,
	}
}

func CreateMerkleDag(path string, encoding ...multibase.Encoding) (*MerkleDag, error) {
	var e multibase.Encoding
	if len(encoding) > 0 {
		e = encoding[0]
	} else {
		e = multibase.Base64
	}
	encoder := multibase.MustNewEncoder(e)

	dag := CreateMerkleDagBuilder()

	builder := CreateBuilder(&path)
	builder.NodeType = DirectoryNodeType
	builder.Content = nil

	entries, err := ioutil.ReadDir(path)
	if err != nil {
		return nil, err
	}

	for _, entry := range entries {
		node, err := processEntry(entry, &path, dag, encoder)
		if err != nil {
			return nil, err
		}

		builder.AddLink(node.Hash)
		dag.AddNode(node, encoder)
	}

	node := builder.BuildNode(encoder)
	dag.AddNode(node, encoder)

	return dag.BuildMerkleDag(node.Hash), nil
}
