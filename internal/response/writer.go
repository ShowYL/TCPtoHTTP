package response

import (
	"bytes"
	"fmt"
	"io"
	"log"

	"github.com/ShowYL/TCPtoHTTP/internal/request"
	"github.com/ShowYL/TCPtoHTTP/internal/request/headers"
)

type Writer struct {
	writer     io.Writer
	statusLine []byte
	headers    []byte
	body       []byte
	headersOjb headers.Headers
	StatusCode StatusCode
}

func NewWriter(conn io.Writer) *Writer {
	return &Writer{
		writer: conn,
	}
}

func (w *Writer) WriteStatusLine(statusCode StatusCode) {
	switch statusCode {
	case 200:
		w.statusLine = ([]byte("HTTP/1.1 200 OK\r\n"))
		w.StatusCode = statusCode
	case 400:
		w.statusLine = ([]byte("HTTP/1.1 400 Bad Request\r\n"))
		w.StatusCode = statusCode
	case 500:
		w.statusLine = ([]byte("HTTP/1.1 500 Internal Server Error\r\n"))
		w.StatusCode = statusCode
	}
}

func (w *Writer) WriteHeaders(headers headers.Headers) {
	w.headersOjb = headers
	buf := bytes.NewBuffer([]byte(headers.ToByte()))
	buf.Write([]byte("\r\n"))
	w.headers = buf.Bytes()
}

func (w *Writer) WriteBody(m string) {
	w.body = []byte(m)
}

func (w *Writer) Build(req *request.Request) {
	w.headersOjb.Replace("content-length", fmt.Sprint(len(w.body)))
	w.WriteHeaders(w.headersOjb)
	w.writer.Write(w.statusLine)
	w.writer.Write(w.headers)
	w.writer.Write(w.body)
	log.Printf("%s %s %d\n", req.RequestLine.Method, req.RequestLine.RequestTarget, w.StatusCode)
}
