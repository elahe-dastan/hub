package client

import (
	"fmt"
	"net"

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
				c := client.New()
				ladder := net.TCPAddr{
					IP:   net.IP{127,0,0,1},
					Port: 8080,
				}
				c.Connect(&ladder)
				c.WhoAmI()
				c.ListClientIDs()
			},
		})
}