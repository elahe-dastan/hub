package server

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"strconv"
	"strings"

	"github.com/elahe-dastan/applifier/config"
	"github.com/elahe-dastan/applifier/message"
	"github.com/elahe-dastan/applifier/request"
	"github.com/elahe-dastan/applifier/response"
	log "github.com/sirupsen/logrus"
)

type Server struct {
	seq     int
	conn    map[net.Conn]string
	writers map[net.Conn]*bufio.Writer
	running int
}

func New() *Server {
	Formatter := new(log.TextFormatter)
	Formatter.TimestampFormat = "02-01-2006 15:04:05"
	Formatter.FullTimestamp = true
	log.SetFormatter(Formatter)

	return &Server{
		seq:     0,
		conn:    make(map[net.Conn]string),
		writers: make(map[net.Conn]*bufio.Writer),
	}
}

// Start handling client connections and messages
func (server *Server) Start(c config.ServerConfig) error {
	server.running = 1

	l, err := net.Listen("tcp4", c.Address)

	if err != nil {
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

func (server *Server) handleConnWorker(tasks <-chan net.Conn) {
	for c := range tasks {
		server.handleConnection(c)
	}
}

func (server *Server) handleConnection(c net.Conn) {
	fmt.Printf("Serving %s\n", c.RemoteAddr().String())

	r := bufio.NewReader(c)

	for {
		netData, err := r.ReadString('\n')
		if err != nil {
			log.Error(err)
			return
		}

		dest, res := server.response(netData, c)

		if res == message.STOP {
			break
		}

		server.broadcast(dest, res)
	}
}

// Return the IDs of the connected clients except the client asking for this
func (server *Server) ListClientIDs(c net.Conn, r *response.List) {
	if len(server.conn) == 1 {
		r.ConcatedIds = "No other client connected"
		return
	}

	var ids []string

	for _, id := range server.conn {
		if id != server.conn[c] {
			ids = append(ids, id)
		}
	}

	r.ConcatedIds = strings.Join(ids, "-")
}

// Stop accepting connections and close the existing ones
func (server *Server) Stop() error {
	server.running = 0

	for conn := range server.conn {
		if err := conn.Close(); err != nil {
			log.Error(err)
		}
	}

	server.conn = make(map[net.Conn]string)

	return nil
}

// nolint: interfacer
func (server *Server) assignID(c net.Conn) {
	server.seq++
	server.conn[c] = strconv.Itoa(server.seq)
	server.writers[c] = bufio.NewWriter(c)
}

func disconnect(l io.Closer) {
	if err := l.Close(); err != nil {
		log.Error(err)
	}
}

func (server *Server) destCli(recipientIDs []string) []net.Conn {
	recipientConn := make([]net.Conn, 0)

	for k, v := range server.conn {
		for _, r := range recipientIDs {
			if v == r {
				recipientConn = append(recipientConn, k)
			}
		}
	}

	return recipientConn
}

func (server *Server) broadcast(recipients []net.Conn, res string) {
	for _, c := range recipients {
		w := server.writers[c]
		if _, err := w.WriteString(res); err != nil {
			log.Error(err)
		}

		if err := w.Flush(); err != nil {
			log.Error(err)
		}
	}
}

func (server *Server) response(data string, c net.Conn) ([]net.Conn, string) {
	s := request.Unmarshal(data)

	des := make([]net.Conn, 0)
	res := ""

	switch s.(type) {
	case request.Stop:
		delete(server.conn, c)
		res = response.Stop{}.MarshalRes()
	case request.Who:
		des = append(des, c)
		r := response.Who{ID: server.conn[c]}
		res = r.MarshalRes()
	case request.List:
		des = append(des, c)
		r := response.List{}
		server.ListClientIDs(c, &r)
		res = r.MarshalRes()
	case request.Send:
		se := s.(request.Send)
		des = server.destCli(se.IDs)
		r := response.Send{Body: se.Body}
		res = r.MarshalRes()
	}

	return des, res
}
