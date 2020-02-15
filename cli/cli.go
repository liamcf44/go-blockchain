package cli

import (
	"flag"
	"fmt"
	"os"
	"runtime"
	"strconv"

	"github.com/liamcf44/go-blockchain.git/blockchain"
)

// CLI stores a blockchain to allow the Command Line Interface to interact with it
type CLI struct{}

// Prints out the different CLI options available
func (cli *CLI) printUsage() {
	fmt.Println("/* Usage /*")
	fmt.Println(" getbalance -address ADDRESS - get the balance for an address")
	fmt.Println(" createblockchain -address ADDRESS creates a blockchain and sends genesis reward to address")
	fmt.Println(" print - Prints the blocks in the chain")
	fmt.Println(" send -from FROM -to TO -amount AMOUNT - Send amount of coins")
}

// Validates the given CLI arguments
func (cli *CLI) validateArgs() {
	// If there are less than two arguments, print usage and exit
	if len(os.Args) < 2 {
		cli.printUsage()
		runtime.Goexit()
	}
}

// Handles the 'print' CLI option
func (cli *CLI) printChain() {
	// Create a chain with ContinueBlockChain and a blank address
	bc := blockchain.ContinueBlockChain("")

	// Defer the closing of the chain's database
	defer bc.Database.Close()

	// Create a new iterator to cycle through the blockchain
	it := bc.CreateIterator()

	// Set up a loop...
	for {
		// Get the next block in the chain
		b := it.Next()

		// Print out the various parts of the block
		fmt.Printf("Hash ==> %x\n", b.Hash)
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

// createBlockChain creates a new blockchain with a given address
func (cli *CLI) createBlockChain(a string) {
	// Create the new chain with InitialiseBlockChain
	bc := blockchain.InitialiseBlockChain(a)

	// Close the database connection
	bc.Database.Close()

	fmt.Println("New blockchain created!")
}

// getBalance returns the balance for a given address
func (cli *CLI) getBalance(a string) {
	// Create the chain with ContinueBlockChain
	bc := blockchain.ContinueBlockChain(a)

	// Defer the closing of the chains database
	defer bc.Database.Close()

	// Holding variable for the balance
	b := 0

	// Get the unspent transaction outputs for the address
	uto := bc.GetUnspentTransactionOutputs(a)

	// Loop through the unspent ouputs
	for _, o := range uto {
		// Add each outputs value to the balance
		b += o.Value
	}

	// Print out the balance
	fmt.Printf("Balance for %s: %d\n", a, b)

}

// send is a function to send an amount from one address to another
func (cli *CLI) send(f, t string, a int) {
	// Create the blockchain with ContinueBlockChain and the from address
	bc := blockchain.ContinueBlockChain(f)

	// Defer the closing of the database
	defer bc.Database.Close()

	// Create a new transaction with the address, the amount and the chain
	tx := blockchain.NewTransaction(f, t, a, bc)

	// Append the transaction to the chain
	bc.AppendBlock([]*blockchain.Transaction{tx})

	fmt.Printf("Successfully sent %d, from %s to %s\n", a, f, t)
}

// Run is the function to run the CLI process
func (cli *CLI) Run() {
	// Make a call to validate the given arguments
	cli.validateArgs()

	// Set the flags for each option
	getBalanceCmd := flag.NewFlagSet("getbalance", flag.ExitOnError)
	createBlockchainCmd := flag.NewFlagSet("createblockchain", flag.ExitOnError)
	sendCmd := flag.NewFlagSet("send", flag.ExitOnError)
	printCmd := flag.NewFlagSet("print", flag.ExitOnError)

	// Extract the information for each command
	getBalanceAddress := getBalanceCmd.String("address", "", "The address to get balance for")
	createBlockchainAddress := createBlockchainCmd.String("address", "", "The address to send genesis block reward to")
	sendFrom := sendCmd.String("from", "", "Source wallet address")
	sendTo := sendCmd.String("to", "", "Destination wallet address")
	sendAmount := sendCmd.Int("amount", 0, "Amount to send")

	// Check which argument has been provided
	switch os.Args[1] {
	// For getbalance...
	case "getbalance":
		// Parse the arguemnts through addBlockCmd, handling any errors.
		err := getBalanceCmd.Parse(os.Args[2:])
		blockchain.HandleError(err)

	// For createblockchain...
	case "createblockchain":
		// Parse the arguemnts through printChainCmd, handling any errors.
		err := createBlockchainCmd.Parse(os.Args[2:])
		blockchain.HandleError(err)

	// For send...
	case "send":
		// Parse the arguemnts through printChainCmd, handling any errors.
		err := sendCmd.Parse(os.Args[2:])
		blockchain.HandleError(err)

	// For print...
	case "print":
		// Parse the arguemnts through printChainCmd, handling any errors.
		err := printCmd.Parse(os.Args[2:])
		blockchain.HandleError(err)

	// In any other scenario...
	default:
		// Print the chain and exit
		cli.printChain()
		runtime.Goexit()
	}

	// If arguments have been parsed through getBalanceCmd do the following...
	if getBalanceCmd.Parsed() {

		// Check if the address passed is a blank string, if so print the usage and exit
		if *getBalanceAddress == "" {
			getBalanceCmd.Usage()
			runtime.Goexit()
		}

		// Otherwise make a call to getBalance with the address
		cli.getBalance(*getBalanceAddress)
	}

	// If arguments have been parsed through createBlockchainCmd do the following...
	if createBlockchainCmd.Parsed() {

		// Check if the address passed is a blank string, if so print the usage and exit
		if *createBlockchainAddress == "" {
			createBlockchainCmd.Usage()
			runtime.Goexit()
		}

		// Otherwise make a call to getBalance with the address
		cli.createBlockChain(*createBlockchainAddress)
	}

	// If arguments have been parsed through sendCmd do the following...
	if sendCmd.Parsed() {
		// Check if any of the given address are blank, or if there i no amount
		if *sendFrom == "" || *sendTo == "" || *sendAmount <= 0 {
			sendCmd.Usage()
			runtime.Goexit()
		}

		// Otherwise make a call to send with the details
		cli.send(*sendFrom, *sendTo, *sendAmount)
	}

	// If arguments have been parsed through printCmd do the following...
	if printCmd.Parsed() {
		// Make a call to printChain
		cli.printChain()
	}

}
