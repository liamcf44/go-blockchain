package main

import (
	"flag"
	"fmt"
	"os"
	"runtime"
	"strconv"

	"github.com/liamcf44/go-blockchain.git/blockchain"
)

type CLI struct {
	blockchain *blockchain.BlockChain
}

func (cli *CLI) printUsage() {
	fmt.Println("/* Usage /*")
	fmt.Println(" add -block BLOCK_DATA - add a block to teh chain")
	fmt.Println(" print - Prints the blocks in the chain")
}

func (cli *CLI) validateArgs() {
	if len(os.Args) < 2 {
		cli.printUsage()
		runtime.Goexit()
	}
}

func (cli *CLI) addBlock(d string) {
	cli.blockchain.AppendBlock(d)
	fmt.Println("Added Block")
}

func (cli *CLI) printChain() {
	it := cli.blockchain.CreateIterator()

	for {
		b := it.Next()

		fmt.Printf("Hash ==> %x\n", b.Hash)
		fmt.Printf("Data ==> %s\n", b.Data)
		fmt.Printf("PreviousHash ==> %x\n", b.PreviousHash)

		pow := blockchain.NewProof(b)

		fmt.Printf("Proof of Work ==> %s\n", strconv.FormatBool(pow.ValidateProof()))
		fmt.Println()

		if len(b.PreviousHash) == 0 {
			break
		}
	}
}

func (cli *CLI) run() {
	cli.validateArgs()

	addBlockCmd := flag.NewFlagSet("add", flag.ExitOnError)
	printChainCmd := flag.NewFlagSet("print", flag.ExitOnError)
	addBlockData := addBlockCmd.String("block", "", "Block data")

	switch os.Args[1] {
	case "add":
		err := addBlockCmd.Parse(os.Args[2:])
		blockchain.HandleError(err)

	case "print":
		err := printChainCmd.Parse(os.Args[2:])
		blockchain.HandleError(err)

	default:
		cli.printChain()
		runtime.Goexit()
	}

	if addBlockCmd.Parsed() {
		if *addBlockData == "" {
			addBlockCmd.Usage()
			runtime.Goexit()
		}

		cli.addBlock(*addBlockData)
	}

	if printChainCmd.Parsed() {
		cli.printChain()
	}

}

func main() {
	defer os.Exit(0)

	bc := blockchain.InitialiseBlockChain()
	defer bc.Database.Close()

	cli := CLI{bc}
	cli.run()
}
