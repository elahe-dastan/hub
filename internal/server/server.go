package server

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net"
	"strconv"
	"strings"

	"github.com/elahe-dastan/applifier/message"
)

type Server struct {
	seq     int
	conn    map[net.Conn]string
	running int
}

func New() *Server {
	return &Server{
		seq:  0,
		conn: map[net.Conn]string{},
	}
}

// Start handling client connections and messages
func (server *Server) Start(ladder *net.TCPAddr) error {
	server.running = 1
	PORT := ":" + strconv.Itoa(ladder.Port)
	l, err := net.Listen("tcp4", PORT)

	if err != nil {
		fmt.Println(err)
		return err
	}

	defer disconnect(l)

	numConn := 100

	tasks := make(chan net.Conn, numConn)

	for i := 0; i < 3; i++ {
		go server.handleConnWorker(tasks)
	}

	for {

		if server.running == 0 {
			return nil
		}

		c, err := l.Accept()

		if err != nil {
			fmt.Println(err)
			return err
		}

		server.assignID(c)
		tasks <- c
	}
}

//func handleConnection1(conn net.Conn) {
//	// read buffer from client after enter is hit
//	bufferBytes, err := bufio.NewReader(conn).ReadBytes('\n')
//
//	if err != nil {
//		log.Println("client left..")
//		conn.Close()
//
//		// escape recursion
//		return
//	}
//
//	// convert bytes from buffer to string
//	message := string(bufferBytes)
//	// get the remote address of the client
//	clientAddr := conn.RemoteAddr().String()
//	// format a response
//	response := fmt.Sprintf(message + " from " + clientAddr + "\n")
//
//	// have server print out important information
//	log.Println(response)
//
//	// let the client know what happened
//	conn.Write([]byte("you sent: " + response))
//
//	// recursive func to handle io.EOF for random disconnects
//	handleConnection1(conn)
//}

func (server Server) handleConnWorker(tasks <-chan net.Conn) {
	for c := range tasks {
		server.handleConnection(c)
	}
}

func (server *Server) handleConnection(c net.Conn) {
	fmt.Printf("Serving %s\n", c.RemoteAddr().String())
	r := bufio.NewReader(c)
	w := bufio.NewWriter(c)

	for {
		//arr, err := readData(r)
		//if err != nil {
		//	log.Println(err)
		//	return
		//}
		netData, err := r.ReadString('\r')
		if err != nil {
			fmt.Println(err)
			return
		}

		arr := strings.Split(netData, "\n")

		command := strings.TrimSpace(arr[0])

		switch command {
		case message.STOP:
			break
		case message.WhoAmI:
			_, err := w.WriteString(server.conn[c] + "\n")
			if err != nil {
				log.Println(err)
			}
		case message.ListClientIDs:
			all := server.ListClientIDs()

			for _, element := range all {
				clientID := strconv.FormatUint(element, 10)
				if clientID != server.conn[c] {
					_, err := w.WriteString(clientID + "\n")
					if err != nil {
						log.Println(err)
					}
				}
			}
		case message.SendMsg:
			recipients, err := strconv.Atoi(arr[1])
			body := arr[2] + "\r"

			if err != nil {
				log.Println(err)
			}

			recipientArr := destCli(recipients)

			server.broadcast(recipientArr, body)
		}

		//_, err = w.WriteString("\r")
		//if err != nil {
		//	log.Println(err)
		//}

		//if err := w.Flush(); err != nil {
		//	log.Println(err)
		//}
	}

	//if err := c.Close(); err != nil {
	//	log.Println(err)
	//}
}

// Return the IDs of the connected clients
func (server *Server) ListClientIDs() []uint64 {
	result := make([]uint64, 0)

	for _, element := range server.conn {
		r, _ := strconv.ParseUint(element, 10, 64)
		result = append(result, r)
	}

	return result
}

// Stop accepting connections and close the existing ones
func (server *Server) Stop() error {
	server.running = 0

	for conn, _ := range server.conn {
		if err := conn.Close(); err != nil {
			log.Println(err)
		}
	}

	server.conn = map[net.Conn]string{}

	return nil
}

func (server *Server) assignID(c net.Conn) {
	server.seq++
	server.conn[c] = strconv.Itoa(server.seq)
}

func disconnect(l io.Closer) {
	if err := l.Close(); err != nil {
		log.Println(err)
	}
}

func destCli(recipients int) []string {
	recipientArr := make([]string, 0)

	for {
		recipientArr = append(recipientArr, fmt.Sprintf("%d", recipients%10))
		recipients /= 10
		if recipients == 0 {
			break
		}
	}

	return recipientArr
}

func (server Server) broadcast(recipientArr []string, body string) {
	for k, v := range server.conn {
		for _, r := range recipientArr {
			if v == r {
				if _, err := bufio.NewWriter(k).WriteString(body); err != nil {
					log.Println(err)
				}
			}
		}
	}
}

//func readData(r *bufio.Reader) ([]string, error) {
//	netData, err := r.ReadString('\r')
//	arr := strings.Split(netData, "\n")
//
//	return arr, err
//}
