package client

import (
	"fmt"
	"net"

	"github.com/elahe-dastan/applifier/config"
	"github.com/elahe-dastan/applifier/internal/client"
	"github.com/spf13/cobra"
)

func Register(root *cobra.Command) {
	root.AddCommand(
		&cobra.Command{
			Use:                        "client",
			Short:                      "Runs the client",
			Run: func(cmd *cobra.Command, args []string) {
				fmt.Println("Hello from client")
				client := client.New()
				c := config.ReadClient()
				ladder := net.TCPAddr{
					IP:   net.IP{byte(c.First),byte(c.Second),byte(c.Third),byte(c.Fourth)},
					Port: c.Port,
				}
				client.Connect(&ladder)
				client.WhoAmI()
				client.ListClientIDs()
			},
		})
}