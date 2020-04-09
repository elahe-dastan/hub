package server

import (
	"fmt"
	"log"

	"github.com/elahe-dastan/hub/config"
	"github.com/elahe-dastan/hub/internal/server"
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
				c := config.Read()

				if err := s.Start(c); err != nil {
					log.Println(err)
				}
			},
		},
	)
}
