package client

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strconv"
	"strings"
)

type IncomingMessage struct {
	SenderID uint64
	Body     []byte
}

type Client struct {
	conn net.Conn
}

func New() *Client {
	return &Client{conn: nil}
}

// Connect to the server using the given address
func (cli *Client) Connect(serverAddr string) error {
	c, err := net.Dial("tcp", serverAddr)
	if err != nil {
		fmt.Println(err)
		return err
	}

	cli.conn = c

	for {
		reader := bufio.NewReader(os.Stdin)

		fmt.Print(">> ")

		text, _ := reader.ReadString('\n')

		if _, err := fmt.Fprintf(c, text+"\n"); err != nil {
			return err
		}

		message, _ := bufio.NewReader(c).ReadString('\n')
		fmt.Print("->: " + message)

		if strings.TrimSpace(text) == "STOP" {
			fmt.Println("TCP client exiting...")

			return nil
		}
	}
}

func (cli *Client) Close() error {
	fmt.Println("TODO: Close the connection to the server")
	return nil
}

// Fetch the ID from the server
func (cli *Client) WhoAmI() (uint64, error) {
	if _, err := fmt.Fprintf(cli.conn, "WhoAmI"+"\n"); err != nil {
		return 0, err
	}

	message, err := bufio.NewReader(cli.conn).ReadString('\n')
	if err != nil {
		return 0, err
	}

	fmt.Print("->: " + message)
	u, err2 := strconv.ParseUint(message, 10, 64)

	return u, err2
}

func (cli *Client) ListClientIDs() ([]uint64, error) {
	if _, err := fmt.Fprintf(cli.conn, "ListClientIDs"+"\n"); err != nil {
		return nil, err
	}

	//buff := make([]uint64, 50)
	//c := bufio.NewReader(cli.conn)
	//
	//for {
	//	// read a single byte which contains the message length
	//	size, err := c.ReadByte()
	//	if err != nil {
	//		return buff,err
	//	}
	//
	//	// read the full message, or return an error
	//	_, err = io.ReadFull(c, buff[:int(size)])
	//	if err != nil {
	//		return err
	//	}
	//
	//	fmt.Printf("received %x\n", buff[:int(size)])
	//}

	message, err := bufio.NewReader(cli.conn).ReadString('\n')
	fmt.Print("->: " + message)
	//u,_ := strconv.ParseUint(message, 10, 64)
	//return u, nil

	fmt.Println("TODO: Fetch the IDs from the server")

	return nil, err
}

func (cli *Client) SendMsg(recipients []uint64, body []byte) error {
	fmt.Println("TODO: Send the message to the server")
	return nil
}

func (cli *Client) HandleIncomingMessages(writeCh chan<- IncomingMessage) {
	fmt.Println("TODO: Handle the messages from the server")
}
