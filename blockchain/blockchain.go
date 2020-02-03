// Package blockchain handles the construction of all parts of the blockchain, the chain itself, blocks, transactions etc.
package blockchain

import (
	"encoding/hex"
	"fmt"
	"os"
	"runtime"

	"github.com/dgraph-io/badger"
)

// Path for the database to write to
const (
	dbPath      = "./tmp/blocks"
	dbFile      = "./tmp/blocks/MANIFEST"
	initialData = "First Transaction from Initialising Chain"
)

// BlockChain holds the last hash and a pointer to the database
type BlockChain struct {
	LatestHash []byte
	Database   *badger.DB
}

// BlockChainIterator holds the current hash and a pointer to the database
type BlockChainIterator struct {
	CurrentHash []byte
	Database    *badger.DB
}

// Checks whether the database exists and is setup
func checkDB() bool {
	// Check if the database file exists and no error is returned
	if _, err := os.Stat(dbFile); os.IsNotExist(err) {
		return false
	}

	return true
}

// AppendBlock is a method on the BlockChain struct which adds a block the the chain
func (bc *BlockChain) AppendBlock(t []*Transaction) {
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
	nb := CreateBlock(t, lh)

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
func InitialiseBlockChain(a string) *BlockChain {
	// Storage variable for the latest hash in the chain
	var lh []byte

	// Check if the database already exists...
	if checkDB() {
		fmt.Println("Blockchain already exists...")
		runtime.Goexit()
	}

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
		// Set up a coinbase (initial) transaction with the address and initial data const
		c := CoinbaseTx(a, initialData)

		// Create an initial block
		ib := CreateInitialBlock(c)

		fmt.Println("Initial block created and proved")

		// Set the initial blocks hash and serialised data, handling any errors
		err = txn.Set(ib.Hash, ib.Serialise())
		HandleError(err)

		// Set the initialblocks hash as the latest hash in the database
		err = txn.Set([]byte("lh"), ib.Hash)

		// Set the latest hash variable as the initial blocks hash
		lh = ib.Hash

		return err
	})

	// Handle any errors
	HandleError(err)

	// Create the blockchain with the latest hash and the database and return it
	bc := BlockChain{lh, db}
	return &bc
}

// ContinueBlockChain is a function which continues off from an existing saved blockchain
func ContinueBlockChain(a string) *BlockChain {
	// Check if the database hasn't been made...
	if checkDB() == false {
		fmt.Println("No existing blockchain found, one needs to be made...")
		runtime.Goexit()
	}

	// Create a storage variable for the latest hash
	var lh []byte

	// Create an instance of the database options and set the path to the constant above
	// Both the directory and value directory live on the same path
	o := badger.DefaultOptions("")
	o.Dir = dbPath
	o.ValueDir = dbPath

	// Open the connection to the database with the options, creating the database variable, handling any errors
	db, err := badger.Open(o)
	HandleError(err)

	// Use the transaction to get the item stored under lh, handling any errors
	err = db.Update(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte("lh"))
		HandleError(err)

		// Append the value to the latest hash variable
		err = item.Value(func(val []byte) error {
			lh = append([]byte{}, val...)

			return nil
		})

		return err
	})

	// Handle any errors
	HandleError(err)

	// Create the chain with the lash hash and the database
	c := BlockChain{lh, db}

	return &c

}

// GetUnspentTransactions is a Blockchain method which returns any unspent transaction for an address
func (bc *BlockChain) GetUnspentTransactions(a string) []Transaction {
	// Create a holding variable for the unspent transactions
	var ut []Transaction

	// Make a map to hold the spent transactions
	st := make(map[string][]int)

	// Create an iterator on the chain to loop through
	i := bc.CreateIterator()

	// Start a for loop for the iterator
	for {
		// Get the next block in the chain
		b := i.Next()

		// For all of those blocks transactions, do the following...
		for _, t := range b.Transactions {

			// Create a transaction ID by encoding the ID
			tID := hex.EncodeToString(t.ID)

			// Create a new labeled loop to go through the outputs
		Outputs:
			for oID, o := range t.Outputs {

				// If the transaction ID does exist in the spent transaction map...
				if st[tID] != nil {

					// Loop through the IDs...
					for _, so := range st[tID] {

						// If the spent output is the same as the outputID then continue
						if so == oID {
							continue Outputs
						}
					}
				}

				// If the output can be unlocked by the address do the following...
				if o.CanBeUnlocked(a) {
					// Append the transaction to the unspent transactions slice
					ut = append(ut, *t)
				}
			}

			// Check if the transaction is a coinbase, if it isn't do the following...
			if t.IsCoinbase() == false {

				// Loop through the inputs for the transaction
				for _, i := range t.Inputs {

					// If the input can unlock the address then do the following...
					if i.CanUnlock(a) {
						// Encode the ID of the input
						iID := hex.EncodeToString(i.ID)

						// Append the inputs output to the spend transactions slice
						st[iID] = append(st[iID], i.Out)
					}
				}
			}

		}

		// Check if the block doesn't have any previous hash, i.e. is the original block, if so break
		if len(b.PreviousHash) == 0 {
			break
		}

	}

	// Return the unspent transactions slice
	return ut

}

// GetUnspentTransactionOutputs is a method on BlockChain which returns the outputs for each unspent transaction
func (bc *BlockChain) GetUnspentTransactionOutputs(a string) []TOutput {
	// Create a holding variable for the unspent transaction outputs
	var uto []TOutput

	// Get the unspent transactions for an address
	ut := bc.GetUnspentTransactions(a)

	// Loop through the unspent transactions
	for _, t := range ut {
		// Loops through each transactions outputs
		for _, o := range t.Outputs {
			// If the output can be unlocked then append
			if o.CanBeUnlocked(a) {
				uto = append(uto, o)
			}
		}
	}

	// Return the unspent transaction outputs
	return uto
}

// GetSpendableOutputs takes an address and a total value to send, it returns spendable outputs for an address, i.e. none coinbase outputs
func (bc *BlockChain) GetSpendableOutputs(a string, v int) (int, map[string][]int) {
	// Make a map to store the unspent outputs
	uo := make(map[string][]int)

	// Go get the unspent transactions for an address
	ut := bc.GetUnspentTransactions(a)

	// Value for accumulated values
	acc := 0

	// Create a labeled loop for each of the unspent transactions
Work:
	for _, t := range ut {

		// Create an ID for the transaction
		tID := hex.EncodeToString(t.ID)

		// Loop through the transactions outputs
		for oID, o := range t.Outputs {

			// Check if the output can be unlocked and that the accumulated value is less than the total given value
			if o.CanBeUnlocked(a) && acc < v {
				// Assign the output value to the accumulated value
				acc += o.Value

				// Assign the output ID to the unspent outputs map
				uo[tID] = append(uo[tID], oID)

				// If the accumulated value gets above the given value then break
				if acc >= v {
					break Work
				}
			}

		}
	}

	// Return the accumulated value and the unspent outputs
	return acc, uo
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
