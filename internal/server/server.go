package server

import (
	"bytes"
	"fmt"
	"io"
	"net"
	"sync/atomic"

<<<<<<< HEAD
=======
	"httpfromtcp/internal/request"
>>>>>>> 43b3835 (test lazygit ui)
	"httpfromtcp/internal/response"
)

type HandlerError struct {
	StatusCode response.StatusCode
	Message    string
}

type Server struct {
	listener net.Listener
	closed   atomic.Bool
	handler  Handler
}

func WriteHandlerError(w io.Writer, handlerErr *HandlerError) error {
	err := response.WriteStatusLine(w, handlerErr.StatusCode)
	if err != nil {
		return err
	}

	headers := response.GetDefaultHeaders(len(handlerErr.Message))

	err = response.WriteHeaders(w, headers)
	if err != nil {
		return err
	}
	_, err = w.Write([]byte(handlerErr.Message))
	return err
}

func Serve(port int, handler Handler) (*Server, error) {
	address := fmt.Sprintf(":%d", port)
	listener, err := net.Listen("tcp", address)
	if err != nil {
		return nil, err
	}

	server := &Server{
		listener: listener,
		handler:  handler,
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

<<<<<<< HEAD
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
=======
	req, err := request.RequestFromReader(conn)
	if err != nil {
		fmt.Println("error parsing request:", err)
		return
	}

	var buf bytes.Buffer

	handlerErr := s.handler(&buf, req)

	if handlerErr != nil {
		err := WriteHandlerError(conn, handlerErr)
		if err != nil {
			fmt.Println("error writing handler error:", err)
		}
		return
	}

	headers := response.GetDefaultHeaders(buf.Len())
	err = response.WriteStatusLine(conn, response.StatusOK)
	if err != nil {
		fmt.Println("error writing status line:", err)
		return
	}

	err = response.WriteHeaders(conn, headers)
	if err != nil {
		fmt.Println("error writing headers:", err)
		return
	}

	_, err = conn.Write(buf.Bytes())
	if err != nil {
		fmt.Println("error writing body:", err)
	}
>>>>>>> 43b3835 (test lazygit ui)
}

type Handler func(w io.Writer, req *request.Request) *HandlerError
