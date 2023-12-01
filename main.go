// Package main is the entry point for the kissbom CLI.
package main

import (
	"os"

	"github.com/devops-kung-fu/kissbom/cmd"
)

func main() {
	defer os.Exit(0)
	cmd.Execute()
}
