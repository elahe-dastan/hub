package client

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
	"runtime"
	"strconv"
	"strings"
	"time"
)

type IncomingMessage struct {
	SenderID uint64
	Body     []byte
}

type Client struct {
	conn     net.Conn
	reader   *bufio.Reader
	Who      chan string
	List     chan string
	Incoming chan string
}

func New() *Client {
	return &Client{
		conn:     nil,
		reader:   nil,
		Who:      make(chan string),
		List:     make(chan string),
		Incoming: make(chan string),
	}
}

// Connect to the server using the given address
func (cli *Client) Connect(serverAddr string) error {
	c, err := net.Dial("tcp", serverAddr)
	if err != nil {
		fmt.Println(err)
		return err
	}

	cli.conn = c
	cli.reader = bufio.NewReader(c)

	go cli.HandleIncomingMessages()

	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Print(">> ")
		fmt.Println(runtime.NumGoroutine())

		text, _ := reader.ReadString('\n')
		fmt.Println("After reading console")
		text = strings.TrimSpace(text)

		switch text {
		case "STOP":
			return cli.Close()
		case "Who":
			if _, err2 := cli.WhoAmI(); err2 != nil {
				return err2
			}
		case "List":
			if _, err2 := cli.ListClientIDs(); err2 != nil {
				return err2
			}
		case "Send":
			recipients := make([]uint64, 0)

			for {
				fmt.Println("Enter the next client")
				t, _ := reader.ReadString('\n')
				t = strings.TrimSuffix(t, "\n")

				if t == "END" {
					break
				}

				c, _ := strconv.ParseUint(t, 10, 64)
				recipients = append(recipients, c)

			}

			fmt.Println("Enter the body")

			b, _ := reader.ReadString('\n')

			if err2 := cli.SendMsg(recipients, []byte(b)); err2 != nil {
				return err2
			}
		default:
			time.Sleep(time.Millisecond)
		}
	}
}

// Close the connection to the server
func (cli *Client) Close() error {
	fmt.Println("TCP client exiting...")

	if _, err := fmt.Fprintf(cli.conn, "STOP\r"); err != nil {
		return err
	}

	return nil
}

// Fetch the ID from the server
func (cli *Client) WhoAmI() (uint64, error) {
	if _, err := fmt.Fprintf(cli.conn, "WhoAmI\r"); err != nil {
		return 0, err
	}

	message := <-cli.Who

	fmt.Print("->: " + message)
	u, err2 := strconv.ParseUint(message, 10, 64)

	return u, err2
}

// Fetch the IDs from the server
func (cli *Client) ListClientIDs() ([]uint64, error) {
	if _, err := fmt.Fprintf(cli.conn, "ListClientIDs\r"); err != nil {
		return nil, err
	}

	message := <-cli.List

	fmt.Print("->: " + message)

	arr := strings.Split(message, "\n")
	res := make([]uint64, 0)

	for _, m := range arr {
		r, err2 := strconv.ParseUint(m, 10, 64)

		if err2 != nil {
			return nil, err2
		}

		res = append(res, r)
	}

	return res, nil
}

//  Send the message to the server
func (cli *Client) SendMsg(recipients []uint64, body []byte) error {
	m := "SendMsg\n"

	for _, u := range recipients {
		m += strconv.FormatUint(u, 10)
	}

	m += "\n"
	m += string(body)
	m += "\r"

	if _, err := fmt.Fprintf(cli.conn, m); err != nil {
		return err
	}

	//message, err := bufio.NewReader(cli.conn).ReadString('\r')
	//if err != nil {
	//	return err
	//}

	//fmt.Print("->: " + message)

	return nil
}

// Handle the messages from the server
func (cli *Client) HandleIncomingMessages() {
	for {
		fmt.Println("in incoming")
		m, err := cli.reader.ReadString('\r')

		if err != nil {
			log.Println(err)
		}

		arr := strings.Split(m, ",")
		switch arr[0] {
		case "Who":
			cli.Who <- arr[1]
		}

		fmt.Println("wa")
	}
}
