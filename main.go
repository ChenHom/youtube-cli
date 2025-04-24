package main

import (
	"os"

	"github.com/ChenHom/ytcli/cmd"
)

func main() {
	if err := cmd.Execute(); err != nil {
		os.Exit(1)
	}
}
