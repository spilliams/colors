package main

import (
	"os"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spilliams/colors/cmd/contrastratio"
)

func main() {
	log.SetLevel(log.DebugLevel)
	log.SetOutput(os.Stderr)

	rootCmd := &cobra.Command{
		Use:   "colors",
		Short: "A tool for playing with colors",
	}

	rootCmd.AddCommand(contrastratio.NewCmd())

	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}
