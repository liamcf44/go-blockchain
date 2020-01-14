package blockchain

import (
	"bytes"
	"crypto/sha256"
	"encoding/binary"
	"fmt"
	"math"
	"math/big"
)

const Difficulty = 20

type ProofOfWork struct {
	Block  *Block
	Target *big.Int
}

func NewProof(b *Block) *ProofOfWork {
	t := big.NewInt(1)
	t.Lsh(t, uint(256-Difficulty))

	pow := &ProofOfWork{b, t}

	return pow
}

func (pow *ProofOfWork) InitialiseData(c int) []byte {
	d := bytes.Join(
		[][]byte{
			pow.Block.PreviousHash,
			pow.Block.Data,
			ToHex(int64(c)),
			ToHex(int64(Difficulty)),
		},
		[]byte{},
	)

	return d
}

func ToHex(n int64) []byte {
	b := new(bytes.Buffer)

	err := binary.Write(b, binary.BigEndian, n)
	HandleError(err)

	return b.Bytes()
}

func (pow *ProofOfWork) ValidateProof() bool {
	var ih big.Int

	d := pow.InitialiseData(pow.Block.Counter)
	h := sha256.Sum256(d)

	ih.SetBytes(h[:])

	return ih.Cmp(pow.Target) == -1
}

func (pow *ProofOfWork) Run() (int, []byte) {
	var ih big.Int
	var h [32]byte

	c := 0

	for c < math.MaxInt64 {
		d := pow.InitialiseData(c)
		h = sha256.Sum256(d)

		fmt.Printf("\r%x", h)

		ih.SetBytes(h[:])

		if ih.Cmp(pow.Target) == -1 {
			break
		} else {
			c++
		}
	}

	fmt.Println()

	return c, h[:]
}
