// Package headers
package headers

import (
	"bytes"
	"errors"
	"fmt"
	"strings"
)

type Headers struct {
	headers map[string]string
}

var SEPARATOR = "\r\n"

func NewHeaders() *Headers {
	return &Headers{
		headers: map[string]string{},
	}
}

func (h *Headers) Get(key string) string {
	return h.headers[strings.ToLower(key)]
}

func (h *Headers) Replace(key, value string) {
	h.headers[strings.ToLower(key)] = value
}

func (h *Headers) Set(key, value string) {
	stored, exist := h.headers[strings.ToLower(key)]
	if exist {
		h.headers[strings.ToLower(key)] = stored + ", " + value
		return
	}
	h.headers[strings.ToLower(key)] = value
}

func (h *Headers) ToString() string {
	var res strings.Builder
	res.WriteString("Headers:\n")

	for key, value := range h.headers {
		fmt.Fprintf(&res, " - %s: %s\n", key, value)
	}

	return res.String()
}

func (h *Headers) ToByte() []byte {
	var res bytes.Buffer

	for key, value := range h.headers {
		fmt.Fprintf(&res, "%s: %s\r\n", key, value)
	}

	return res.Bytes()
}

func isCorrectKey(s string) bool {
	if len(s) == 0 {
		return false
	}

	for i := 0; i < len(s); i++ {
		c := s[i]

		isAlpha := (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z')
		isDigit := c >= '0' && c <= '9'

		if isAlpha || isDigit {
			continue
		}

		switch c {
		case '!', '#', '$', '%', '&', '\'', '*', '+', '-', '.', '^', '_', '`', '|', '~':
			continue
		default:
			return false
		}
	}

	return true
}

var (
	ErrWrongFormatHeaders = errors.New("wrong format in the headers")
	ErrWrongKeyFormat     = errors.New("wrong symbol used in key")
)

func (h *Headers) Parse(data []byte) (int, bool, error) {
	read := 0
	done := false

	for {
		idx := bytes.Index(data[read:], []byte(SEPARATOR))
		if idx == -1 {
			break
		}

		if idx == 0 {
			read += len(SEPARATOR)
			done = true
			break
		}

		key, value, err := parseHeader(data[read : idx+read])
		if err != nil {
			return 0, false, err
		}

		if !isCorrectKey(key) {
			return 0, false, ErrWrongKeyFormat
		}

		h.Set(key, value)
		read += idx + len(SEPARATOR)
	}

	return read, done, nil
}

func parseHeader(line []byte) (string, string, error) {
	parts := bytes.SplitN(line, []byte(":"), 2)

	if len(parts) != 2 {
		return "", "", ErrWrongFormatHeaders
	}

	if bytes.Contains(parts[0], []byte(" ")) {
		return "", "", ErrWrongFormatHeaders
	}

	return string(parts[0]), string(bytes.TrimSpace(parts[1])), nil
}
