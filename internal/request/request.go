// Package request
package request

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"slices"
	"strconv"

	"github.com/ShowYL/TCPtoHTTP/internal/request/headers"
)

type parserState string

const (
	ParsingRequestLine parserState = "parsing the request line"
	ParsingHeaders     parserState = "parsing the headers"
	ParsingBody        parserState = "parsing the body"
	ParsingDone        parserState = "parsing done"
)

type Request struct {
	RequestLine *RequestLine
	Headers     *headers.Headers
	Body        string
	state       parserState
}

func newRequest() *Request {
	return &Request{
		state:   ParsingRequestLine,
		Headers: headers.NewHeaders(),
		Body:    "",
	}
}

func (r *Request) ToString() string {
	return r.RequestLine.ToString() + r.Headers.ToString() + r.Body + "\n"
}

type RequestLine struct {
	HTTPVersion   string
	RequestTarget string
	Method        string
}

func (rl *RequestLine) isValidHTTP() bool {
	return rl.HTTPVersion == "1.1"
}

func (rl *RequestLine) isValidMethod() bool {
	return slices.Contains(Methods, rl.Method)
}

func (rl *RequestLine) ToString() string {
	res := "Request Line :\n"

	res += fmt.Sprintf(" - Method: %s\n", rl.Method)

	res += fmt.Sprintf(" - Target: %s\n", rl.RequestTarget)

	res += fmt.Sprintf(" - Version: %s\n", rl.HTTPVersion)

	return res
}

var (
	ErrWrongStartLine     = errors.New("wrong start line")
	ErrWrongHTTPVersion   = errors.New("wrong http version, only 1.1 is supported")
	ErrWrongVersionFormat = errors.New("wrong format for the version of http")
	ErrWrongMethod        = errors.New("wrong method used")
)

var Methods = []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS", "HEAD", "CONNECT", "TRACE"}

var SEPARATOR = []byte("\r\n")

var bufferSize = 1024

func RequestFromReader(reader io.Reader) (*Request, error) {
	buf := make([]byte, bufferSize)
	bufIndex := 0

	request := newRequest()

	for request.state != ParsingDone {
		n, err := reader.Read(buf[bufIndex:])
		if err != nil {
			return nil, err
		}

		bufIndex += n

		readN, err := request.parse(buf[:bufIndex])
		if err != nil {
			return nil, err
		}

		copy(buf, buf[readN:bufIndex])
		bufIndex -= readN
	}

	return request, nil
}

func getInt(headers headers.Headers, key string) (int, error) {
	val := headers.Get(key)
	if val == "" {
		return 0, nil
	}

	res, err := strconv.Atoi(val)
	if err != nil {
		return 0, err
	}

	return res, nil
}

func (r *Request) parse(data []byte) (int, error) {
	read := 0

outer:
	for {
		currentData := data[read:]
		switch r.state {
		case ParsingRequestLine:
			rl, n, err := parseRequestLine(currentData)
			if err != nil {
				return 0, err
			}

			if n == 0 {
				break outer
			}

			r.RequestLine = rl
			r.state = ParsingHeaders
			read += n

		case ParsingHeaders:
			n, done, err := r.Headers.Parse(currentData)
			if err != nil {
				return 0, err
			}

			if n == 0 {
				break outer
			}

			if done {
				r.state = ParsingBody
			}

			read += n

		case ParsingBody:
			length, err := getInt(*r.Headers, "content-length")
			if err != nil {
				return 0, err
			}

			if length == 0 {
				r.state = ParsingDone
				break
			}

			remaining := min(length-len(r.Body), len(currentData))
			r.Body += string(currentData[:remaining])
			read += remaining

			if len(r.Body) == length {
				r.state = ParsingDone
			}

			if len(currentData) == 0 {
				break outer
			}

		case ParsingDone:
			break outer
		}
	}

	return read, nil
}

func parseRequestLine(b []byte) (*RequestLine, int, error) {
	idx := bytes.Index(b, SEPARATOR)
	if idx == -1 {
		return nil, 0, nil
	}
	startLine := b[:idx]
	read := idx + len(SEPARATOR)

	parts := bytes.Split(startLine, []byte(" "))

	if len(parts) != 3 {
		return nil, 0, ErrWrongStartLine
	}

	versionParts := bytes.Split(parts[2], []byte("/"))
	if len(versionParts) != 2 || string(versionParts[0]) != "HTTP" {
		return nil, 0, ErrWrongVersionFormat
	}

	rl := &RequestLine{
		Method:        string(parts[0]),
		RequestTarget: string(parts[1]),
		HTTPVersion:   string(versionParts[1]),
	}

	if !rl.isValidHTTP() {
		return nil, 0, ErrWrongHTTPVersion
	}

	if !rl.isValidMethod() {
		return nil, 0, ErrWrongMethod
	}

	return rl, read, nil
}
