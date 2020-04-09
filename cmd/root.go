package cmd

import (
	"fmt"
	"os"

	"github.com/elahe-dastan/hub/cmd/client"
	"github.com/elahe-dastan/hub/cmd/server"
	"github.com/spf13/cobra"
)

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	// rootCmd represents the base command when called without any subcommands
	var rootCmd = &cobra.Command{
		Use:   "hub",
		Short: "Message deliver",
	}

	server.Register(rootCmd)
	client.Register(rootCmd)

	exitFailure := 1

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(exitFailure)
	}
}
