package main

import (
	"os"

	"github.com/liamcf44/go-blockchain.git/cli"
)

// main function
func main() {
	// Defer exiting the process
	defer os.Exit(0)

	// Create the new command line struct
	cmd := cli.CLI{}

	// // Run the CLI
	cmd.Run()
}
