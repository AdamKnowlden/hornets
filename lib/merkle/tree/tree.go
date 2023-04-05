package tree

import (
	"crypto/sha256"
	"encoding/base64"
	"fmt"

	mt "github.com/txaty/go-merkletree"
)

type TreeContent struct {
	blocks []mt.DataBlock
}

type Block struct {
	data string
}

func (b *Block) Serialize() ([]byte, error) {
	return []byte(b.data), nil
}

func CreateHash(data []byte) string {
	hash := sha256.Sum256(data)
	hashString := base64.StdEncoding.EncodeToString(hash[:])

	return string(hashString)
}

func CreateTree() *TreeContent {
	tree := TreeContent{
		[]mt.DataBlock{},
	}

	return &tree
}

func (tc *TreeContent) AddBlock(data string) {
	block := &Block{data}

	tc.blocks = append(tc.blocks, block)
}

func (tc *TreeContent) Build() (*mt.MerkleTree, error) {
	tree, err := mt.New(nil, tc.blocks)
	if err != nil {
		return nil, err
	}

	result := VerifyTree(tree, tc.blocks)
	if result {
		fmt.Println("Merkle Tree Verified: tree")
	}

	result = VerifyRoot(tree.Root, tree.Proofs, tc.blocks)
	if result {
		fmt.Println("Merkle Tree Verified: root")
	}

	return tree, err
}

func VerifyTree(tree *mt.MerkleTree, blocks []mt.DataBlock) bool {
	result := true

	proofs := tree.Proofs
	// verify the proofs
	for i := 0; i < len(proofs); i++ {
		ok, err := tree.Verify(blocks[i], proofs[i])
		if err != nil {
			fmt.Println("Verification failed for block")
		}

		if ok == false {
			result = ok
		}
	}

	return result
}

func VerifyRoot(root []byte, proofs []*mt.Proof, blocks []mt.DataBlock) bool {
	result := true

	for i := 0; i < len(blocks); i++ {
		// if hashFunc is nil, use SHA256 by default
		ok, err := mt.Verify(blocks[i], proofs[i], root, nil)
		if err != nil {
			fmt.Println("Verification failed for block")
		}

		if ok == false {
			result = ok
		}
	}

	return result
}
