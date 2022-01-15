package blockchain

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"math"
	"strconv"
	"time"
)

type Block struct {
	index        int
	timestamp    int64
	proof        int
	previousHash string
}

func (b Block) GetPreviousHash() string {
	return b.previousHash
}

func (b Block) GetProof() int {
	return b.proof
}

func (b Block) GetIndex() int {
	return b.index
}

func (b Block) GetTimestamp() int64 {
	return b.timestamp
}

func (b Block) MarshalJSON() ([]byte, error) {
	j := struct {
		Index        int
		Timestamp    int64
		Proof        int
		PreviousHash string
	}{
		Index:        b.index,
		Timestamp:    b.timestamp,
		Proof:        b.proof,
		PreviousHash: b.previousHash,
	}

	m, err := json.Marshal(j)
	if err != nil {
		return nil, err
	}
	return m, nil
}

type Blockchain struct {
	chain []Block
}

func (b *Blockchain) Init() {
	b.CreateBlock(1, "0")
}

func (b *Blockchain) CreateBlock(proof int, previousHash string) Block {
	newBlock := Block{
		index:        len(b.chain) + 1,
		timestamp:    time.Now().Unix(),
		proof:        proof,
		previousHash: previousHash,
	}
	b.chain = append(b.chain, newBlock)
	return newBlock
}

func (b *Blockchain) GetChain() []Block {
	return b.chain
}

func (b *Blockchain) GetPreviousBlock() *Block {
	return &b.chain[len(b.chain)-1]
}

func (b *Blockchain) ProofOfWork(previousProof int) int {
	newProof := 1
	checkProof := false
	for !checkProof {
		h := sha256.New()
		newProofSquared := int(math.Pow(float64(newProof), 2.0))
		previousProofSquared := int(math.Pow(float64(b.GetPreviousBlock().GetProof()), 2.0))
		calcString := strconv.Itoa(newProofSquared - previousProofSquared)
		h.Write([]byte(calcString))
		hash_operation := hex.EncodeToString(h.Sum(nil))
		if hash_operation[0:4] == "0000" {
			checkProof = true
		} else {
			newProof += 1
		}

	}
	return newProof
}

func (b *Blockchain) Hash(block Block) string {
	h := sha256.New()
	blockString := fmt.Sprintf("%v", block)
	h.Write([]byte(blockString))
	return hex.EncodeToString(h.Sum(nil))
}

func (b *Blockchain) IsChainValid() bool {
	chain := b.GetChain()
	previousBlock := chain[0]
	for i := 1; i < len(chain); i++ {
		block := chain[i]
		previousHash := block.GetPreviousHash()
		if previousHash != b.Hash(previousBlock) {
			return false
		} else {
			previousProof := previousBlock.GetProof()
			proof := block.GetProof()
			h := sha256.New()
			proofSquared := int(math.Pow(float64(proof), 2.0))
			previousProofSquared := int(math.Pow(float64(previousProof), 2.0))
			calcString := strconv.Itoa(proofSquared - previousProofSquared)
			h.Write([]byte(calcString))
			hash_operation := hex.EncodeToString(h.Sum(nil))
			if hash_operation[0:4] != "0000" {
				return false
			} else {
				previousBlock = block
			}
		}
	}
	return true
}
