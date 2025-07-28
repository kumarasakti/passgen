package main

import (
	"fmt"
	"os"

	"github.com/kumarasakti/passgen/internal/infrastructure/cli"
)

// Version can be overridden at build time using -ldflags "-X main.Version=x.y.z"
var Version = "v1.1.0"

func main() {
	handler := cli.NewHandler()
	rootCmd := handler.CreateRootCommand(Version)

	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
