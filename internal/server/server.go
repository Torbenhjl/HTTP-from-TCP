package server

import (
	"fmt"
	"net"
	"sync/atomic"

	"httpfromtcp/internal/response"
)

type Server struct {
	listener net.Listener
	closed   atomic.Bool
}

func Serve(port int) (*Server, error) {
	address := fmt.Sprintf(":%d", port)

	listener, err := net.Listen("tcp", address)
	if err != nil {
		return nil, err
	}

	server := &Server{
		listener: listener,
	}

	go server.listen()

	return server, nil
}

func (s *Server) Close() error {
	s.closed.Store(true)
	return s.listener.Close()
}

func (s *Server) listen() {
	for {
		conn, err := s.listener.Accept()
		if err != nil {
			if s.closed.Load() {
				return
			}

			fmt.Println("error accepting connection:", err)
			continue
		}

		go s.handle(conn)
	}
}

func (s *Server) handle(conn net.Conn) {
	defer conn.Close()

	// response := "HTTP/1.1 200 OK\r\n" +
	// 	"Content-Type: text/plain\r\n" +
	// 	"Content-Length: 13\r\n" +
	// 	"\r\n" +
	// 	"Hello World!\n"
	//

	response.WriteStatusLine(conn, 200)

	headers := response.GetDefaultHeaders(0)
	response.WriteHeaders(conn, headers)

	// _, err := conn.Write([]byte(resp))
	// if err != nil {
	// 	fmt.Println("error writing response:", err)
	// }
}
