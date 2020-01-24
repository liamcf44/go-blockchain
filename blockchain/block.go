// Package blockchain handles the construction of all parts of the blockchain, the chain itself, blocks, transactions etc.
package blockchain

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"log"
)

// Block stores all the parts of a block, including the hash of the previous block
type Block struct {
	Hash         []byte
	Data         []byte
	PreviousHash []byte
	Counter      int
}

// CreateBlock takes some data and a previous hash and returns a new block
func CreateBlock(d string, ph []byte) *Block {
	// Create a new instance of Block
	b := &Block{[]byte{}, []byte(d), ph, 0}

	// Create a new proof of work for the block
	pow := NewProof(b)

	// Run the proof of work
	c, h := pow.Run()

	// Set the hash and the counter for the block
	b.Hash = h[:]
	b.Counter = c

	return b
}

// CreateInitialBlock makes a first block in a chain
func CreateInitialBlock() *Block {
	// Create the intial block
	return CreateBlock("Initial Block", []byte{})
}

// Serialise is a method on the Block struct that serialises the block's data
func (b *Block) Serialise() []byte {
	// Create the data buffer variable
	var d bytes.Buffer

	// Create a new encoder with the data buffer
	e := gob.NewEncoder(&d)

	// Encode the block with the encoder, handling any errors
	err := e.Encode(b)
	HandleError(err)

	return d.Bytes()
}

// Deserialise takes some data and returns it in the form of a block
func Deserialise(d []byte) *Block {
	// Create a storage variable for the block
	var b Block

	// Create a new encoder with the given data
	dc := gob.NewDecoder(bytes.NewReader(d))

	// Decode the data, handling any errors
	err := dc.Decode(&b)
	HandleError(err)

	return &b
}

// HandleError is a generic function to handle any errors in the process
func HandleError(err error) {
	if err != nil {
		fmt.Println("Error occurred: ")

		log.Panic(err)
	}
}
