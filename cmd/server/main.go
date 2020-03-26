package server

import (
	"fmt"
	"net"

	"github.com/elahe-dastan/applifier/internal/server"
	"github.com/spf13/cobra"
)

func Register(root *cobra.Command)  {
	root.AddCommand(
		&cobra.Command{
			Use:   "server",
			Short: "Run server to serve the requests",
			Run: func(cmd *cobra.Command, args []string) {
				fmt.Println("Hello from server")
				s := server.New()
				ladder := net.TCPAddr{
					Port: 8080,
				}
				s.Start(&ladder)
			},
		},
	)
}