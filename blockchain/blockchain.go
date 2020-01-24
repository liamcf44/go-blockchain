// Package blockchain handles the construction of all parts of the blockchain, the chain itself, blocks, transactions etc.
package blockchain

import (
	"fmt"

	"github.com/dgraph-io/badger"
)

// Path for the database to write to
const (
	dbPath = "./tmp/blocks"
)

// BlockChain holds the last hash and a pointer to the database
type BlockChain struct {
	LatestHash []byte
	Database *badger.DB
}

// BlockChainIterator holds the current hash and a pointer to the database
type BlockChainIterator struct {
	CurrentHash []byte
	Database    *badger.DB
}

// AppendBlock is a method on the BlockChain struct which adds a block the the chain
func (bc *BlockChain) AppendBlock(d string) {
	// Storage variable for the latest hash in the chain
	var lh []byte

	// Make a call to the database...
	err := bc.Database.View(func(txn *badger.Txn) error {
		// Get the block stored under the latest hash, handling any error
		item, err := txn.Get([]byte("lh"))
		HandleError(err)

		// Get the value of the item retrieved 
		err = item.Value(func(val []byte) error {
			// Append the value to the latest hash variable
			lh = append([]byte{}, val...)

			return nil
		})

		return err
	})

	// Handle any errors that occurred
	HandleError(err)

	// Create a new block with the given data and the latest hash
	nb := CreateBlock(d, lh)

	// Make a call to the database to update it...
	err = bc.Database.Update(func(txn *badger.Txn) error {
		// Set the hash of the new block and the serialised data, handling any errors
		err := txn.Set(nb.Hash, nb.Serialise())
		HandleError(err)

		// Set the hash of the new block to the latest hash for future use
		err = txn.Set([]byte("lh"), nb.Hash)

		// Set the latest hash on the blockchain to the hash of the new block
		bc.LatestHash = nb.Hash

		return err

	})

	// Handle any errors that occurred
	HandleError(err)

}

// InitialiseBlockChain is a function which returns a new BlockChain
func InitialiseBlockChain() *BlockChain {
	// Storage variable for the latest hash in the chain
	var lh []byte

	// Create an instance of the database options and set the path to the constant above
	// Both the directory and value directory live on the same path
	o := badger.DefaultOptions("")
	o.Dir = dbPath
	o.ValueDir = dbPath

	// Open the connection to the database with the options, creating the database variable, handling any errors
	db, err := badger.Open(o)
	HandleError(err)

	// Open a transaction with the database to update it...
	err = db.Update(func(txn *badger.Txn) error {

		// Use the transaction to check if there is a latest hash, i.e. a blockchain has already been initiated
		if _, err := txn.Get([]byte("lh")); err == badger.ErrKeyNotFound {
			// If there isn't a latest hash then no blockchain exists and do the following...
			fmt.Println("No existing blockchain found")

			// Create an initial block
			ib := CreateInitialBlock()

			fmt.Println("Initial block created and proved")

			// Set the initial blocks hash and serialised data, handling any errors
			err = txn.Set(ib.Hash, ib.Serialise())
			HandleError(err)

			// Set the initialblocks hash as the latest hash in the database
			err = txn.Set([]byte("lh"), ib.Hash)

			// Set the latest hash variable as the initial blocks hash
			lh = ib.Hash

			return err
		} else {
			// If there is already a blockchain initialised do the following...

			// Use the transaction to get the item stored under lh, handling any errors
			item, err := txn.Get([]byte("lh"))
			HandleError(err)

			// Append the value to the latest hash variable
			err = item.Value(func(val []byte) error {
				lh = append([]byte{}, val...)

				return nil
			})

			return err
		}
	})

	// Handle any errors
	HandleError(err)

	// Create the blockchain with the latest hash and the database and return it
	bc := BlockChain{lh, db}
	return &bc
}

// CreateIterator is a method on the BlockChain struct that creates a new BlockChainIterator
func (bc *BlockChain) CreateIterator() *BlockChainIterator {
	// Create the iterator with the current latest hash and the database pointer
	it := &BlockChainIterator{bc.LatestHash, bc.Database}

	// return the iterator
	return it
}

// Next is a method on the BlockChainIterator struct that returns the next block in the chain
func (it *BlockChainIterator) Next() *Block {
	// Storage variable for the raw block returned from the database and the eventual block returned from the function
	var rb []byte
	var b *Block

	// Create a transaction with the database...
	err := it.Database.View(func(txn *badger.Txn) error {
		// Get the item stored under the iterators current hash, handling any errors
		item, err := txn.Get(it.CurrentHash)
		HandleError(err)

		// Fetch the value of the item and append this to the raw block storage variable
		err = item.Value(func(val []byte) error {
			rb = append([]byte{}, val...)

			return nil
		})

		// Deserialise the raw block data into the block storage variable
		b = Deserialise(rb)

		return err
	})

	// Handle any erros that have occurred
	HandleError(err)

	// Store the previous hash of the block as the iterators current hash for future use
	it.CurrentHash = b.PreviousHash

	return b
}
