package client

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
	"strconv"
	"strings"

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
	bufferSize := 20

	return &Client{
		Who:      make(chan string, bufferSize),
		List:     make(chan string, bufferSize),
		Incoming: make(chan string, bufferSize),
	}
}

// Connects to the server with the address specified in the config file
// and spawns another goroutine which is waiting for the server messages
// all the time and the function itself starts an infinite for loop which
// reads from the console and sends the command to the server
func (cli *Client) Connect(serverAddr string) error {
	c, err := net.Dial("tcp", serverAddr)

	if err != nil {
		return err
	}

	cli.conn = c
	cli.reader = bufio.NewReader(c)
	cli.writer = bufio.NewWriter(c)
	cli.console = bufio.NewReader(os.Stdin)

	go cli.HandleIncomingMessages()
	go cli.privateMessage()

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

	m := <-cli.Who

	fmt.Print("->: " + m)
}

// Fetch the IDs of the clients currently connected to the server
func (cli *Client) ListClientIDs() {
	cli.flushBuffer(message.ListClientIDs + "\n")

	m := <-cli.List

	m = strings.TrimSuffix(m, "-\n")
	IDs := strings.Split(m, "-")

	for _, id := range IDs {
		fmt.Println(id)
	}
}

// Gets the IDs of the clients user wants to send a message to
// and the message itself then sends a req to the server to do this
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

	req = req + "," + b

	cli.flushBuffer(req)
}

// Starts an infinite for loop which is repeatedly waiting for
// messages from the server
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
			cli.Incoming <- arr[1]
		}
	}
}

// Based on the type of the command sends a proper
// message using the right protocol to the server
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

// Writes the message that should be sent to the server in the
// buffer and the flushes it
func (cli Client) flushBuffer(m string) {
	if _, err := cli.writer.WriteString(m); err != nil {
		log.Println(err)
	}

	if err := cli.writer.Flush(); err != nil {
		log.Println(err)
	}
}

func (cli Client) privateMessage() {
	fmt.Print(<-cli.Incoming)
}
