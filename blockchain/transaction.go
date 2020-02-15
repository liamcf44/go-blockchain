package blockchain

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"
	"encoding/hex"
	"fmt"
	"log"
)

// Transaction stores the relevant parts of a blockchain transaction, containing multiple inputs and outputs
type Transaction struct {
	ID      []byte
	Inputs  []TxInput
	Outputs []TxOutput
}

// SetID creates a sha256 hash ID for a transaction
func (t *Transaction) SetID() {
	// Create some data as a buffer and a hash variable
	var d bytes.Buffer
	var h [32]byte

	// Create a new encoder, passing the data to it
	var e = gob.NewEncoder(&d)

	// Encode the transaction, handling any errors
	err := e.Encode(t)
	HandleError(err)

	// Create a hash with the datas bytes and assign to the transaction
	h = sha256.Sum256(d.Bytes())
	t.ID = h[:]

}

// CoinbaseTx handles the coinbase (the original transaction)
func CoinbaseTx(r, d string) *Transaction {
	// If the data is empty then assign data to default string
	if d == "" {
		d = fmt.Sprintf("Coins to %s", r)
	}

	// Create a transaction input and output with the given data and recepient
	tIn := TxInput{[]byte{}, -1, d}
	tOut := TxOutput{100, r}

	// Use the above to construct a new transaction
	t := Transaction{nil, []TxInput{tIn}, []TxOutput{tOut}}

	// Call the SetID method
	t.SetID()

	// Return the transaction
	return &t
}

// IsCoinbase checks whether a transaction instance is a coinbase (the original transaction)
func (t *Transaction) IsCoinbase() bool {
	// Check to see there is just 1 input and that it is not linked to any other transactions
	return len(t.Inputs) == 1 && len(t.Inputs[0].ID) == 0 && t.Inputs[0].Out == -1
}

// NewTransaction takes a from and to address, an amount and a block chain and makes a transaction to return
func NewTransaction(f, t string, a int, bc *BlockChain) *Transaction {
	// Create two holding variables for the inputs and outputs
	var i []TxInput
	var o []TxOutput

	// Get the accumulated value and the unspent outputs for the from address, up to the specified amount
	acc, uo := bc.GetSpendableOutputs(f, a)

	// If the accumulator does not reach the amount then the account does not have enough funds
	if acc < a {
		log.Panic("Error : Not enough funds!")
	}

	// If the funds are available, loop through the unspent outputs
	for id, outs := range uo {
		// Decode the transaction ID, handling any errors
		tID, err := hex.DecodeString(id)
		HandleError(err)

		// Loop through the outputs
		for _, out := range outs {
			// Create a new transcation input from the ID, the output and the from address
			in := TxInput{tID, out, f}

			// Append the input to the holding variable
			i = append(i, in)
		}
	}

	// Append a new transaction output, with the given amount and the to address
	o = append(o, TxOutput{a, t})

	// If the accumulated ammount is more than the given ammount then trim the ouput
	if acc > a {
		// Append a new transaction output with some money sent back to the from address
		o = append(o, TxOutput{acc - a, f})
	}

	// Create a new transaction with the inputs and outputs and set its ID
	tx := Transaction{nil, i, o}
	tx.SetID()

	// Return the transaction
	return &tx
}
