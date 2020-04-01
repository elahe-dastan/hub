package client

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
	"strconv"
	"strings"

	"github.com/elahe-dastan/applifier/config"
	"github.com/elahe-dastan/applifier/message"
)

type Client struct {
	conn     net.Conn
	reader   *bufio.Reader
	writer   *bufio.Writer
	console  *bufio.Reader
	Who      chan string
	List     chan string
	Incoming chan string
}

func New() *Client {
	return &Client{
		Who:      make(chan string, 20),
		List:     make(chan string, 20),
		Incoming: make(chan string, 20),
	}
}

// Connect to the server using the given address
func (cli *Client) Connect(cc config.ClientConfig) error {
	serverAddr := cc.IP + ":" + cc.Port
	c, err := net.Dial("tcp", serverAddr)
	if err != nil {
		return err
	}

	cli.conn = c
	cli.reader = bufio.NewReader(c)
	cli.writer = bufio.NewWriter(c)
	cli.console = bufio.NewReader(os.Stdin)

	go cli.HandleIncomingMessages()

	for {
		fmt.Print(">> ")

		req, _ := cli.console.ReadString('\n')
		req = strings.TrimSpace(req)

		cli.sendReq(req)
	}
}

// Close the connection to the server
func (cli *Client) Close() {
	fmt.Println("TCP client exiting...")

	cli.flushBuffer("STOP\n")
}

// Fetch the ID from the server
func (cli *Client) WhoAmI() {
	cli.flushBuffer(message.WhoAmI + "\n")

	message := <-cli.Who

	fmt.Print("->: " + message)
}

// Fetch the IDs from the server
func (cli *Client) ListClientIDs() {
	cli.flushBuffer(message.ListClientIDs + "\n")

	m := <-cli.List

	IDs := strings.Split(m, "-")

	for _, id := range IDs {
		fmt.Println(id)
	}
}

//  Send the message to the server
func (cli *Client) SendMsg() {
	req := message.SendMsg + ","

	for {
		fmt.Println("Enter the next client or END")
		t, _ := cli.console.ReadString('\n')
		t = strings.TrimSuffix(t, "\n")

		if t == "END" {
			break
		}

		if _, err := strconv.Atoi(t); err != nil {
			fmt.Println("Enter a number or END")
		} else {
			req = req + t + "-"
		}
	}

	fmt.Println("Enter the body")

	b, _ := cli.console.ReadString('\n')

	req = req + "," + b + "\n"

	cli.flushBuffer(req)
}

// Handle the messages from the server
func (cli *Client) HandleIncomingMessages() {
	for {
		m, err := cli.reader.ReadString('\n')

		if err != nil {
			log.Println(err)
		}

		arr := strings.Split(m, ",")
		switch arr[0] {
		case "Who":
			cli.Who <- arr[1]
		case "List":
			cli.List <- arr[1]
		case "Send":
			fmt.Println(arr[1])
			cli.Incoming <- arr[1]
		}
	}
}

func (cli *Client) sendReq(req string) {
	switch req {
	case "STOP":
		cli.Close()
	case "Who":
		cli.WhoAmI()
	case "List":
		cli.ListClientIDs()
	case "Send":
		cli.SendMsg()
	}
}

func (cli Client) flushBuffer(m string) {
	if _, err := fmt.Fprintf(cli.writer, m); err != nil {
		log.Println(err)
	}

	if err := cli.writer.Flush(); err != nil {
		log.Println(err)
	}
}
