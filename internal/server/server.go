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
		conn:    map[net.Conn]string{},
		writers: map[net.Conn]*bufio.Writer{},
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

		server.broadcast(dest, res)
	}
}

// Return the IDs of the connected clients except the client asking for this
func (server *Server) ListClientIDs(c net.Conn, r *response.List) {
	if len(server.conn) == 1 {
		r.IDs = append(r.IDs, "No other client connected")
		return
	}

	for _, id := range server.conn {
		if id != server.conn[c] {
			r.IDs = append(r.IDs, id)
		}
	}
}

// Stop accepting connections and close the existing ones
func (server *Server) Stop() error {
	server.running = 0

	for conn := range server.conn {
		if err := conn.Close(); err != nil {
			log.Error(err)
		}
	}

	server.conn = map[net.Conn]string{}

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

func (server *Server) destCli(recipientIDs string) []net.Conn {
	recipientIDs = strings.TrimSuffix(recipientIDs, "-")
	recipientArr := strings.Split(recipientIDs, "-")

	recipientConn := make([]net.Conn, 0)

	for k, v := range server.conn {
		for _, r := range recipientArr {
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
	arr := strings.Split(data, ",")
	t := strings.TrimSpace(arr[0])

	des := make([]net.Conn, 0)
	res := ""

	switch t {
	case message.STOP:
		des = append(des, c)

		if err := server.Stop(); err != nil {
			res = err.Error()
		} else {
			res = "Server stopped"
		}
	case message.WhoAmI:
		des = append(des, c)
		r := response.Who{ID: server.conn[c]}
		res = r.Marshal()
	case message.ListClientIDs:
		des = append(des, c)
		r := response.List{IDs: []string{}}
		server.ListClientIDs(c, &r)
		res = r.Marshal()
	case message.SendMsg:
		des = server.destCli(arr[1])
		r := response.Send{Body: strings.Join(arr[2:], ",")}
		res = r.Marshal()
	}

	return des, res
}
