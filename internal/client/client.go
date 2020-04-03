package client

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strconv"
	"strings"

	"github.com/c-bata/go-prompt"
	"github.com/elahe-dastan/applifier/request"
	"github.com/elahe-dastan/applifier/response"
	log "github.com/sirupsen/logrus"
)

type Client struct {
	conn     net.Conn
	reader   *bufio.Reader
	writer   *bufio.Writer
	console  *bufio.Reader
	Who      chan response.Who
	List     chan response.List
	Incoming chan response.Send
}

func New() *Client {
	bufferSize := 20

	Formatter := new(log.TextFormatter)
	Formatter.TimestampFormat = "02-01-2006 15:04:05"
	Formatter.FullTimestamp = true
	log.SetFormatter(Formatter)

	return &Client{
		Who:      make(chan response.Who, bufferSize),
		List:     make(chan response.List, bufferSize),
		Incoming: make(chan response.Send, bufferSize),
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

	p := prompt.New(cli.sendReq, completer, prompt.OptionPrefix("applifier> "))
	p.Run()

	return nil
}

// Close the connection to the server
func (cli *Client) Close() {
	fmt.Println("TCP client exiting...")

	cli.flushBuffer("STOP\n")

	if err := cli.conn.Close(); err != nil {
		log.Fatal(err)
	}

	os.Exit(0)
}

// Fetch the ID from the server
func (cli *Client) WhoAmI() {
	cli.flushBuffer(request.Who{}.Marshal())

	m := <-cli.Who

	fmt.Print("->: " + m.ID)
}

// Fetch the IDs of the clients currently connected to the server
func (cli *Client) ListClientIDs() {
	cli.flushBuffer(request.List{}.Marshal())

	l := <-cli.List
	m := l.ConcatedIds

	m = strings.TrimSuffix(m, "\n")
	IDs := strings.Split(m, "-")

	for _, id := range IDs {
		fmt.Println(id)
	}
}

// Gets the IDs of the clients user wants to send a message to
// and the message itself then sends a req to the server to do this
func (cli *Client) SendMsg() {
	r := request.Send{IDs: []string{}}

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
			r.IDs = append(r.IDs, t)
		}
	}

	fmt.Println("Enter the body")

	b, _ := cli.console.ReadString('\n')
	r.Body = b

	cli.flushBuffer(r.Marshal())
}

// Starts an infinite for loop which is repeatedly waiting for
// messages from the server
func (cli *Client) HandleIncomingMessages() {
	for {
		m, err := cli.reader.ReadString('\n')

		if err != nil {
			log.Fatal(err)
		}

		s := response.Unmarshal(m)

		switch s.(type) {
		case response.Who:
			w := s.(response.Who)
			cli.Who <- w
		case response.List:
			l := s.(response.List)
			cli.List <- l
		case response.Send:
			se := s.(response.Send)
			cli.Incoming <- se
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
		log.Fatal(err)
	}

	if err := cli.writer.Flush(); err != nil {
		log.Fatal(err)
	}
}

func (cli Client) privateMessage() {
	for {
		s := <-cli.Incoming
		fmt.Print(s.Body)
	}
}

func completer(d prompt.Document) []prompt.Suggest {
	s := []prompt.Suggest{
		{Text: "STOP", Description: "Closes the connection"},
		{Text: "Who", Description: "Returns the ID"},
		{Text: "List", Description: "Returns the IDs of other clients"},
		{Text: "Send", Description: "Send a message"},
	}

	return prompt.FilterHasPrefix(s, d.GetWordBeforeCursor(), true)
}
