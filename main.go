package main

import (
	"fmt"
	"log"
	"net"
)

type Server struct {
	listenAddr string
	ln         net.Listener
	quitchan   chan struct{} // empty struct for not taking memory
}

func NewServer(listenAddr string) *Server {
	return &Server{
		listenAddr: listenAddr,
		quitchan:   make(chan struct{}),
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
	buf := make([]byte, 2048)

	for {
		n, err := conn.Read(buf)
		if err != nil {
			log.Println("error reading from connection:", err)
			continue
		}
		msg := buf[:n]
		fmt.Println(string(msg))
	}
}

func main() {
	server := NewServer(":8000")
	log.Fatal(server.Start())
}