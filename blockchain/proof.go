// Package blockchain handles the construction of all parts of the blockchain, the chain itself, blocks, transactions etc.
package blockchain

import (
	"bytes"
	"crypto/sha256"
	"encoding/binary"
	"fmt"
	"math"
	"math/big"
)

// Difficulty is a constant that controls the difficult of hash generation.
// The higher the difficulty the more computing power needed per hash
const Difficulty = 18

// ProofOfWork is a struct to hold a block and a target intiger for the hash
type ProofOfWork struct {
	Block  *Block
	Target *big.Int
}

// NewProof is a method on the Block struct that generates a proof of work
func NewProof(b *Block) *ProofOfWork {
	// Creates a target intiger for the hash
	t := big.NewInt(1)

	// Use Lsh to retrun t << 256 - Difficulty constant
	t.Lsh(t, uint(256-Difficulty))

	// Create a new ProofOfWork with the block and the target
	pow := &ProofOfWork{b, t}

	return pow
}

// InitialiseData is a method on the ProofOfWork struct that takes the current counter and returns data
func (pow *ProofOfWork) InitialiseData(c int) []byte {
	// Create the data by concatting the elements below into a new byte slice
	d := bytes.Join(
		[][]byte{
			pow.Block.PreviousHash,
			pow.Block.HashTransactions(),
			ToHex(int64(c)),
			ToHex(int64(Difficulty)),
		},
		[]byte{},
	)

	return d
}

// ToHex takes an int and returns a hex as a slice of bytes
func ToHex(n int64) []byte {
	// Create a new buffer
	b := new(bytes.Buffer)

	// Write the buffer, along with the int, handling any errors
	err := binary.Write(b, binary.BigEndian, n)
	HandleError(err)

	return b.Bytes()
}

// ValidateProof is a method on the ProofOfWork struct that checks if the proof of work for block
func (pow *ProofOfWork) ValidateProof() bool {
	// Storage variable for initial hash
	var ih big.Int

	// Pass the current counter on the block to InitialiseData
	d := pow.InitialiseData(pow.Block.Counter)

	// Create a hash with sha256
	h := sha256.Sum256(d)

	// Set the bytes of the hash to the initial hash variable
	ih.SetBytes(h[:])

	// Compare the hash with the target to check if it has been validated
	return ih.Cmp(pow.Target) == -1
}

// Run is a method on the ProofOfWork struct to run a proof of work process
func (pow *ProofOfWork) Run() (int, []byte) {
	// Create storage variables for the initial hash and the has to be returned
	var ih big.Int
	var h [32]byte

	// Set the counter to 0
	c := 0

	// Whilst the counter is smaller than a MaxInt do the following...
	for c < math.MaxInt64 {
		// Create the data through InitialiseData
		d := pow.InitialiseData(c)

		// Create a sha256 hash
		h = sha256.Sum256(d)

		fmt.Printf("\r%x", h)

		// Set the bytes of the hash to the initial hash
		ih.SetBytes(h[:])

		// Check if the hash and the target match..
		if ih.Cmp(pow.Target) == -1 {
			// If they do then break
			break
		} else {
			// Otherwise up the counter
			c++
		}
	}

	// Return the counter and the bytes of the hash
	return c, h[:]
}
