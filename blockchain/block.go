package blockchain

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"log"
)

type Block struct {
	Hash         []byte
	Data         []byte
	PreviousHash []byte
	Counter      int
}

func CreateBlock(d string, ph []byte) *Block {
	b := &Block{[]byte{}, []byte(d), ph, 0}
	pow := NewProof(b)
	c, h := pow.Run()

	b.Hash = h[:]
	b.Counter = c

	return b
}

func CreateInitialBlock() *Block {
	return CreateBlock("Initial Block", []byte{})
}

func (b *Block) Serialise() []byte {
	var r bytes.Buffer

	e := gob.NewEncoder(&r)

	err := e.Encode(b)

	HandleError(err)

	return r.Bytes()
}

func Deserialise(d []byte) *Block {
	var b Block

	dc := gob.NewDecoder(bytes.NewReader(d))

	err := dc.Decode(&b)

	HandleError(err)

	return &b
}

func HandleError(err error) {
	if err != nil {
		fmt.Println("Error occurred: ")

		log.Panic(err)
	}
}
