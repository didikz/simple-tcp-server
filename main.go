package main

import (
	"fmt"
	"log"
	"net"
)

type Message struct {
	source  string
	payload []byte
}

type Server struct {
	listenAddr string
	ln         net.Listener
	quitchan   chan struct{} // empty struct for not taking memory
	msgchan    chan Message
}

func NewServer(listenAddr string) *Server {
	return &Server{
		listenAddr: listenAddr,
		quitchan:   make(chan struct{}),
		msgchan:    make(chan Message, 10),
	}
}

func (s *Server) Start() error {
	ln, err := net.Listen("tcp", s.listenAddr)
	if err != nil {
		return err
	}

	log.Printf("Server start at %s", s.listenAddr)

	defer ln.Close()

	s.ln = ln

	go s.AcceptConnections()

	<-s.quitchan
	close(s.msgchan)

	return nil
}

func (s *Server) AcceptConnections() {
	for {
		conn, err := s.ln.Accept()
		if err != nil {
			log.Println("error accepting connection:", err)
			continue
		}

		log.Println("accepting connection from", conn.RemoteAddr())

		go s.ReadConnection(conn)
	}
}

func (s *Server) ReadConnection(conn net.Conn) {
	defer conn.Close()

	for {
		buf := make([]byte, 2048)
		n, err := conn.Read(buf)
		if err != nil {
			log.Println("error reading from connection:", err)
			continue
		}

		if n == 0 || string(buf[:n]) == "\r\n" {
			continue
		}

		s.msgchan <- Message{
			source:  conn.RemoteAddr().String(),
			payload: buf[:n],
		}
	}
}

func main() {
	server := NewServer(":8000")

	go func() {
		for msg := range server.msgchan {
			fmt.Printf("received message (%s): %s\n", msg.source, string(msg.payload))
		}
	}()
	log.Fatal(server.Start())
}
