package blockchain

import (
	"fmt"

	"github.com/dgraph-io/badger"
)

const (
	dbPath = "./tmp/blocks"
)

type BlockChain struct {
	LastHash []byte
	Database *badger.DB
}

type BlockChainIterator struct {
	CurrentHash []byte
	Database    *badger.DB
}

func (bc *BlockChain) AppendBlock(d string) {
	var lh []byte

	err := bc.Database.View(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte("lh"))
		HandleError(err)

		err = item.Value(func(val []byte) error {
			lh = append([]byte{}, val...)

			return nil
		})

		return err
	})

	HandleError(err)

	nb := CreateBlock(d, lh)

	err = bc.Database.Update(func(txn *badger.Txn) error {
		err := txn.Set(nb.Hash, nb.Serialise())
		HandleError(err)

		err = txn.Set([]byte("lh"), nb.Hash)

		bc.LastHash = nb.Hash

		return err

	})
	HandleError(err)

}

func InitialiseBlockChain() *BlockChain {
	var lh []byte

	o := badger.DefaultOptions("")
	o.Dir = dbPath
	o.ValueDir = dbPath

	db, err := badger.Open(o)
	HandleError(err)

	err = db.Update(func(txn *badger.Txn) error {
		if _, err := txn.Get([]byte("lh")); err == badger.ErrKeyNotFound {
			fmt.Println("No existing blockchain found")

			ib := CreateInitialBlock()

			fmt.Println("Initial block created and proved")

			err = txn.Set(ib.Hash, ib.Serialise())
			HandleError(err)

			err = txn.Set([]byte("lh"), ib.Hash)

			lh = ib.Hash

			return err
		} else {
			item, err := txn.Get([]byte("lh"))

			HandleError(err)

			err = item.Value(func(val []byte) error {
				lh = append([]byte{}, val...)

				return nil
			})

			return err
		}
	})

	HandleError(err)

	bc := BlockChain{lh, db}
	return &bc
}

func (bc *BlockChain) CreateIterator() *BlockChainIterator {
	it := &BlockChainIterator{bc.LastHash, bc.Database}

	return it
}

func (it *BlockChainIterator) Next() *Block {
	var eb []byte
	var b *Block

	err := it.Database.View(func(txn *badger.Txn) error {
		item, err := txn.Get(it.CurrentHash)
		HandleError(err)

		err = item.Value(func(val []byte) error {
			eb = append([]byte{}, val...)

			return nil
		})

		b = Deserialise(eb)

		return err
	})

	HandleError(err)

	it.CurrentHash = b.PreviousHash

	return b
}
