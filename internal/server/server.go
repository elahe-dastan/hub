package server

import (
	"bufio"
	"fmt"
	"net"
	"strconv"
	"strings"
)

type Server struct {
	seq  int
	conn map[string]string
}

func New() *Server {
	return &Server{
		seq:  0,
		conn: map[string]string{},
	}
}

// Start handling client connections and messages
func (server *Server) Start(ladder *net.TCPAddr) error {
	PORT := ":" + strconv.Itoa(ladder.Port)
	l, err := net.Listen("tcp4", PORT)
	if err != nil {
		fmt.Println(err)
		return err
	}
	defer l.Close()

	for {
		c, err := l.Accept()
		if err != nil {
			fmt.Println(err)
			return err
		}
		server.assignID(c)
		go server.handleConnection(c)
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

func (server *Server) handleConnection(c net.Conn) {
	fmt.Printf("Serving %s\n", c.RemoteAddr().String())
	r := bufio.NewReader(c)
	w := bufio.NewWriter(c)

	for {
		netData, err := r.ReadString('\n')
		if err != nil {
			fmt.Println(err)
			return
		}

		temp := strings.TrimSpace(netData)
		if temp == "STOP" {
			break
		}

		if temp == "WhoAmI" {
			w.WriteString(server.conn[c.RemoteAddr().String()] + "\n")
		}

		if temp == "ListClientIDs" {
			all := server.ListClientIDs()

			for _,element := range all{
				w.WriteString(strconv.FormatUint(element, 10))
			}

			w.WriteString("\n")
		}
		w.Flush()
	}
	c.Close()
}

func (server *Server) ListClientIDs() []uint64 {
	var result []uint64
	for _, element := range server.conn {
		r,_ := strconv.ParseUint(element, 10, 64)
		result = append(result, r)
	}
	fmt.Println("TODO: Return the IDs of the connected clients")
	return result
}

func (server *Server) Stop() error {
	fmt.Println("TODO: Stop accepting connections and close the existing ones")
	return nil
}

func (server *Server) assignID(c net.Conn) {
	server.seq = server.seq + 1
	server.conn[c.RemoteAddr().String()] = strconv.Itoa(server.seq)
}
