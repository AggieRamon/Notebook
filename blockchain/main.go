package main

import (
	"encoding/json"
	"log"
	"net/http"

	"aggieramon.com/blockchain/blockchain"
)

type BlockchainHandler struct {
	Blockchain *blockchain.Blockchain
}

func (b BlockchainHandler) mineblock() ([]byte, int) {
	previousBlock := b.Blockchain.GetPreviousBlock()
	previousProof := previousBlock.GetProof()
	proof := b.Blockchain.ProofOfWork(previousProof)
	previousHash := b.Blockchain.Hash(*previousBlock)
	block := b.Blockchain.CreateBlock(proof, previousHash)
	body := struct {
		Message      string `json:"message"`
		Index        int    `json:"index"`
		Timestamp    int64  `json:"timestamp"`
		Proof        int    `json:"proof"`
		PreviousHash string `json:"previousHash"`
	}{
		Message:      "Congrats you just mined a block",
		Index:        block.GetIndex(),
		Timestamp:    block.GetTimestamp(),
		Proof:        block.GetProof(),
		PreviousHash: block.GetPreviousHash(),
	}

	res, err := json.Marshal(body)
	if err != nil {
		return nil, 500
	}

	return res, 200
}

func (b BlockchainHandler) getChain() ([]byte, int) {
	body := struct {
		Chain  []blockchain.Block `json:"chain"`
		Length int                `json:"length"`
	}{
		Chain:  b.Blockchain.GetChain(),
		Length: len(b.Blockchain.GetChain()),
	}

	res, err := json.Marshal(body)
	if err != nil {
		return nil, 500
	}

	return res, 200
}

func (b BlockchainHandler) isValid() ([]byte, int) {
	valid := b.Blockchain.IsChainValid()
	if valid {
		return []byte("All good"), 200
	}

	return []byte("Chain is not valid"), 200
}

func (b BlockchainHandler) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	switch {
	case req.Method == "GET" && req.RequestURI == "/mineblock":
		body, status := b.mineblock()
		res.WriteHeader(status)
		res.Write(body)
	case req.Method == "GET" && req.RequestURI == "/getchain":
		body, status := b.getChain()
		res.WriteHeader(status)
		res.Write(body)
	case req.Method == "GET" && req.RequestURI == "/valid":
		body, status := b.isValid()
		res.WriteHeader(status)
		res.Write(body)
	default:
		res.Write([]byte("This route does not exist"))
	}
}

func main() {
	master_blockchain := blockchain.Blockchain{}
	master_blockchain.Init()
	handler := BlockchainHandler{
		Blockchain: &master_blockchain,
	}
	http.Handle("/", handler)

	log.Fatal(http.ListenAndServe(":4200", nil))
}
