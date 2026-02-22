// Package main
package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/ShowYL/TCPtoHTTP/internal/request"
	"github.com/ShowYL/TCPtoHTTP/internal/response"
	"github.com/ShowYL/TCPtoHTTP/internal/server"
)

const port = 42069

func main() {
	server, err := server.Serve(port, func(w *response.Writer, req *request.Request) {
		if req.RequestLine.RequestTarget == "/yourproblem" {
			w.WriteStatusLine(response.BadRequest)
			w.WriteBody(response.BadRequestBody)
			return
		}

		if req.RequestLine.RequestTarget == "/myproblem" {
			w.WriteStatusLine(response.InternalServerError)
			w.WriteBody(response.InternalErrorBody)
			return
		}

		w.WriteBody(response.OkBody)
	})
	if err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
	defer server.Close()
	log.Println("Server started on port", port)

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan
	log.Println("Server gracefully stopped")
}
