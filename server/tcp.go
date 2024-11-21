package server

import (
	"fmt"
	"io"
	"log"
	"net"
	"strings"

	"github.com/suman7383/lru-cache/cache"
)

type Tcp_Server struct {
	listenAddr string
	ln         net.Listener
	cache      cache.Cacher
}

func New(listenAddr string, cache cache.Cacher) *Tcp_Server {
	return &Tcp_Server{
		listenAddr: listenAddr,
		cache:      cache,
	}
}

func (t *Tcp_Server) Start() {
	ln, err := net.Listen("tcp", t.listenAddr)

	if err != nil {
		log.Fatal("could not start server", err)
	}

	log.Print("server started : ", ln.Addr().String())

	t.ln = ln

	// Accept incomming connections
	for {
		conn, err := t.ln.Accept()

		if err != nil {
			log.Print("error handling incomming connection", err)
			continue
		}

		log.Print("incomming connection: ", conn.RemoteAddr().String())

		go t.handleConn(conn)
	}
}

func (t *Tcp_Server) handleConn(conn net.Conn) {
	// close after done
	defer conn.Close()

	for {
		buf := make([]byte, 1024)

		n, err := conn.Read(buf)

		if err != nil {
			if err == io.EOF {
				log.Print("connection closed", conn.RemoteAddr().String())
				break
			}

			log.Print("error reading from connection", err)
		}

		in := string(buf[:n])
		println("received: ", in)

		// Remove trailing newlines(\n)
		in = strings.TrimSuffix(in, "\n")

		// parse the input
		parsedIn := strings.Split(in, " ")

		switch parsedIn[0] {
		case "PUT":
			// check for valid arguments
			if len(parsedIn) != 3 {
				t.writeToConn(conn, []byte("invalid args. wanted 3 args\n"))
				break
			}

			res, err := t.cache.Put(parsedIn[1], parsedIn[2])

			if err != nil {
				log.Printf("error [PUT]: %q, [VALUE]: %q", parsedIn[1], parsedIn[2])
				break
			}

			t.writeToConn(conn, []byte(res+"\n"))

		case "GET":
			// check for valid arguments
			if len(parsedIn) != 2 {
				t.writeToConn(conn, []byte("invalid args. wanted 2 args\n"))
				break
			}

			res, err := t.cache.Get(parsedIn[1])

			if err != nil {
				log.Printf("error [GET]: %q", parsedIn[1])
				t.writeToConn(conn, []byte(err.Error()+"\n"))
				break
			}

			t.writeToConn(conn, []byte(res+"\n"))

		case "DEL":
			// check for valid arguments
			if len(parsedIn) != 2 {
				t.writeToConn(conn, []byte("invalid args. wanted 2 args\n"))
				break
			}

			res, err := t.cache.Del(parsedIn[1])

			if err != nil {
				log.Printf("error [DEL]: %q", parsedIn[1])
				t.writeToConn(conn, []byte(err.Error()+"\n"))
				break
			}

			t.writeToConn(conn, []byte(res+"\n"))

		case "SIZE":
			res := t.cache.Size()

			t.writeToConn(conn, []byte(fmt.Sprint(res)+"\n"))

		default:
			t.writeToConn(conn, []byte("invalid command\n"))
			log.Print("invalid command")
		}

	}
}

func (t *Tcp_Server) writeToConn(conn net.Conn, msg []byte) {
	_, err := conn.Write(msg)

	if err != nil {
		log.Print("error writing to connection")
		return
	}
}
