package server

import (
	"github.com/ShowYL/TCPtoHTTP/internal/request"
	"github.com/ShowYL/TCPtoHTTP/internal/response"
)

type Handler func(w *response.Writer, res *request.Request)
