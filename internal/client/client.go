package client

import (
	"bufio"
	"fmt"
	"net"
	"strconv"
)

type IncomingMessage struct {
	SenderID uint64
	Body     []byte
}

type Client struct {
	conn net.Conn
}

func New() *Client {
	return &Client{conn:nil}
}

func (cli *Client) Connect(serverAddr *net.TCPAddr) error {
	c, err := net.Dial("tcp", serverAddr.String())
	if err != nil {
		fmt.Println(err)
		return err
	}

	cli.conn = c
	//for {
	//	reader := bufio.NewReader(os.Stdin)
	//	fmt.Print(">> ")
	//	text, _ := reader.ReadString('\n')
	//	fmt.Fprintf(c, text+"\n")
	//
	//	message, _ := bufio.NewReader(c).ReadString('\n')
	//	fmt.Print("->: " + message)
	//	if strings.TrimSpace(string(text)) == "STOP" {
	//		fmt.Println("TCP client exiting...")
	//		return nil
	//	}
	//}
	fmt.Println("TODO: Connect to the server using the given address")
	return nil
}

func (cli *Client) Close() error {
	fmt.Println("TODO: Close the connection to the server")
	return nil
}

func (cli *Client) WhoAmI() (uint64, error) {
	fmt.Fprintf(cli.conn, "WhoAmI"+"\n")

	message, _ := bufio.NewReader(cli.conn).ReadString('\n')
	fmt.Print("->: " + message)
	fmt.Println("TODO: Fetch the ID from the server")
	u,_ := strconv.ParseUint(message, 10, 64)
	return u, nil
}

func (cli *Client) ListClientIDs() ([]uint64, error) {
	fmt.Fprintf(cli.conn, "ListClientIDs"+"\n")

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

	message, _ := bufio.NewReader(cli.conn).ReadString('\n')
	fmt.Print("->: " + message)
	fmt.Println("TODO: Fetch the ID from the server")
	//u,_ := strconv.ParseUint(message, 10, 64)
	//return u, nil

	fmt.Println("TODO: Fetch the IDs from the server")
	return nil, nil
}

func (cli *Client) SendMsg(recipients []uint64, body []byte) error {
	fmt.Println("TODO: Send the message to the server")
	return nil
}

func (cli *Client) HandleIncomingMessages(writeCh chan<- IncomingMessage) {
	fmt.Println("TODO: Handle the messages from the server")
}
