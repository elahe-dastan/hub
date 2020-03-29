package server

import (
	"fmt"
	"log"
	"net"

	"github.com/elahe-dastan/applifier/config"
	"github.com/elahe-dastan/applifier/internal/server"
	"github.com/spf13/cobra"
)

func Register(root *cobra.Command) {
	root.AddCommand(
		&cobra.Command{
			Use:   "server",
			Short: "Run server to serve the requests",
			Run: func(cmd *cobra.Command, args []string) {
				fmt.Println("Hello from server")
				s := server.New()
				c := config.ReadServer()
				ladder := net.TCPAddr{
					Port: c.Port,
				}
				if err := s.Start(&ladder); err != nil {
					log.Println(err)
				}
			},
		},
	)
}
