package medelagateway

import (
	"net/http"
)

type responseWriter struct {
	body       []byte
	statusCode int
	header     http.Header
}

func newResponseWriter() *responseWriter {
	return &responseWriter{
		header: make(http.Header),
	}
}

func (respw *responseWriter) Header() http.Header {
	return respw.header
}

func (respw *responseWriter) Write(b []byte) (int, error) {
	respw.body = append(respw.body, b...)
	return len(respw.body), nil
}

func (respw *responseWriter) WriteHeader(statusCode int) {
	respw.statusCode = statusCode
}
