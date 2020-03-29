package client

import (
	"fmt"
	"log"
	"net"

	"github.com/elahe-dastan/applifier/config"
	"github.com/elahe-dastan/applifier/internal/client"
	"github.com/spf13/cobra"
)

func Register(root *cobra.Command) {
	root.AddCommand(
		&cobra.Command{
			Use:   "client",
			Short: "Runs the client",
			Run: func(cmd *cobra.Command, args []string) {
				fmt.Println("Hello from client")
				cli := client.New()
				c := config.ReadClient()
				ladder := net.TCPAddr{
					IP:   net.IP{byte(c.First), byte(c.Second), byte(c.Third), byte(c.Fourth)},
					Port: c.Port,
				}

				if err := cli.Connect(ladder.String()); err != nil {
					log.Fatal(err)
				}

				if _, err := cli.WhoAmI(); err != nil {
					log.Println(err)
				}

				if _, err := cli.ListClientIDs(); err != nil {
					log.Println(err)
				}
			},
		})
}
