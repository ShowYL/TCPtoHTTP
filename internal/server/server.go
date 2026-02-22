// Package server
package server

import (
	"fmt"
	"io"
	"net"

	"github.com/ShowYL/TCPtoHTTP/internal/request"
	"github.com/ShowYL/TCPtoHTTP/internal/request/headers"
	"github.com/ShowYL/TCPtoHTTP/internal/response"
)

type Server struct {
	handler  Handler
	listener net.Listener
}

func newServer(handler Handler, listener net.Listener) *Server {
	return &Server{
		handler:  handler,
		listener: listener,
	}
}

func (s *Server) Close() {
	s.listener.Close()
}

func (s *Server) listen() {
	for {
		conn, err := s.listener.Accept()
		if err != nil {
			return
		}

		go s.handle(conn)
	}
}

func (s *Server) handle(conn io.ReadWriteCloser) {
	defer conn.Close()

	headers := *getDefaultHeaders()
	writer := response.NewWriter(conn)

	req, err := request.RequestFromReader(conn)
	if err != nil {
		writer.WriteStatusLine(response.BadRequest)
		writer.WriteHeaders(headers)
		writer.WriteBody(response.BadRequestBody)
		writer.Build(req)
		return
	}

	writer.WriteStatusLine(response.Ok)
	writer.WriteHeaders(headers)

	s.handler(writer, req)

	writer.Build(req)
}

func Serve(port int, handler Handler) (*Server, error) {
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return nil, err
	}

	server := newServer(handler, listener)

	go server.listen()

	return server, nil
}

func getDefaultHeaders() *headers.Headers {
	h := headers.NewHeaders()
	h.Set("Connection", "close")
	h.Set("Content-Type", "text/html")

	return h
}
