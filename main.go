package main

import (
	"flag"
	"fmt"
	"os"
	"runtime"
	"strconv"

	"github.com/liamcf44/go-blockchain.git/blockchain"
)

// CLI stores a blockchain to allow the Command Line Interface to interact with it 
type CLI struct {
	blockchain *blockchain.BlockChain
}

// Prints out the different CLI options available
func (cli *CLI) printUsage() {
	fmt.Println("/* Usage /*")
	fmt.Println(" add -block BLOCK_DATA - add a block to teh chain")
	fmt.Println(" print - Prints the blocks in the chain")
}

// Validates the given CLI arguments
func (cli *CLI) validateArgs() {
	// If there are less than two arguments, print usage and exit
	if len(os.Args) < 2 {
		cli.printUsage()
		runtime.Goexit()
	}
}

// Handles the 'add --block' CLI option
func (cli *CLI) addBlock(d string) {
	// Pass the given data to the AppendBlock blockchain method
	cli.blockchain.AppendBlock(d)
	fmt.Println("Added Block")
}

// Handles the 'print' CLI option
func (cli *CLI) printChain() {
	// Create a new iterator to cycle through the current blockchain
	it := cli.blockchain.CreateIterator()

	// Set up a loop...
	for {
		// Get the next block in the chain
		b := it.Next()

		// Print out the various parts of the block
		fmt.Printf("Hash ==> %x\n", b.Hash)
		fmt.Printf("Data ==> %s\n", b.Data)
		fmt.Printf("PreviousHash ==> %x\n", b.PreviousHash)

		// Create a Proof of Work for the block and print if it is valid
		pow := blockchain.NewProof(b)

		fmt.Printf("Proof of Work ==> %s\n", strconv.FormatBool(pow.ValidateProof()))
		fmt.Println()

		// If there is no previous block then the end of the chain has been reached, break.
		if len(b.PreviousHash) == 0 {
			break
		}
	}
}

// Function to run the CLI
func (cli *CLI) run() {
	// Make a call to validate the given arguments
	cli.validateArgs()

	// Set the flags for each option
	addBlockCmd := flag.NewFlagSet("add", flag.ExitOnError)
	printChainCmd := flag.NewFlagSet("print", flag.ExitOnError)
	addBlockData := addBlockCmd.String("block", "", "Block data")

	// Check which argument has been provided
	switch os.Args[1] {
	// For add...
	case "add":
		// Parse the arguemnts through addBlockCmd, handling any errors.
		err := addBlockCmd.Parse(os.Args[2:])
		blockchain.HandleError(err)

	// For print...
	case "print":
		// Parse the arguemnts through printChainCmd, handling any errors.
		err := printChainCmd.Parse(os.Args[2:])
		blockchain.HandleError(err)

	// In any other scenario...
	default:
		// Print the chain and exit
		cli.printChain()
		runtime.Goexit()
	}


	// If arguments have been parsed through addBlockCmd do the following...
	if addBlockCmd.Parsed() {

		// Check if the data passed is a blank string, if so print the usage and exit
		if *addBlockData == "" {
			addBlockCmd.Usage()
			runtime.Goexit()
		}

		// Otherwise make a call to addBlock with the data
		cli.addBlock(*addBlockData)
	}

	// If arguments have been parsed through printChainCmd do the following...
	if printChainCmd.Parsed() {
		// Make a call to printChain
		cli.printChain()
	}

}

// main function
func main() {
	// Defer exiting the process
	defer os.Exit(0)

	// Make a call to InitialiseBlockchain and defer the database closing
	bc := blockchain.InitialiseBlockChain()
	defer bc.Database.Close()

	// Use the returned blockchain to create a new CLI struct
	cli := CLI{bc}

	// Run the CLI
	cli.run()
}
