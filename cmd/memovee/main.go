package main

import (
	"os"

	"github.com/kritama/memovee-cli/internal/command"
)

func main() {
	exitCode := command.Run(os.Args[1:], command.IO{
		Stdin:  os.Stdin,
		Stdout: os.Stdout,
		Stderr: os.Stderr,
	})
	if exitCode != command.ExitSuccess {
		os.Exit(exitCode)
	}
}
