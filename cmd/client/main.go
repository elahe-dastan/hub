package client

import (
	"fmt"
	"log"

	"github.com/elahe-dastan/hub/internal/client"
	"github.com/spf13/cobra"
)

func Register(root *cobra.Command) {
	c := cobra.Command{
		Use:   "client",
		Short: "Runs the client",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("Hello from client")

			cli := client.New()

			a, err := cmd.Flags().GetString("server")
			if err != nil {
				log.Fatal(err)
			}

			if err := cli.Connect(a); err != nil {
				log.Fatal(err)
			}
		},
	}

	c.Flags().StringP("server", "s", "127.0.0.1:8080", "server address")

	root.AddCommand(
		&c,
	)
}
