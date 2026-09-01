package main

import (
	"fmt"
	"os"

	"github.com/kritama/memovee-cli/internal/version"
)

const exitUsage = 2

func main() {
	if len(os.Args) == 2 && os.Args[1] == "version" {
		fmt.Println(version.String())
		return
	}

	fmt.Fprintln(os.Stderr, "usage: memovee version")
	os.Exit(exitUsage)
}
