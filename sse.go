package sse

import (
	"encoding/json"
	"fmt"
	"net/http"
)

const (
	// NonceMax is the maximum value of the `id` field in SSE before it resets to 0.
	NonceMax = 1<<63 - 1
)

// Writer is the interface for writing Server-Sent Events.
type Writer interface {
	Write(event string, data interface{}) error
}

type Options struct {
	ResponseStatus int // defaults to http.StatusOK
	Encoding       int // defaults to Compress_None
}

// NewResponseWriter creates a new Writer for Server-Sent Events.
func NewResponseWriter(
	writer http.ResponseWriter,
	options Options,
) Writer {
	s := responseWriter{
		writer:  writer,
		nonce:   0,
		Options: options,
	}
	s.sendHeaders()
	return &s
}

// responseWriter is an implementation of Writer for Server-Sent Events.
type responseWriter struct {
	writer  http.ResponseWriter
	nonce   uint64
	Options Options
}

// Write sends a message to the client.
func (s *responseWriter) Write(
	event string,
	data interface{},
) error {
	s.nonce = (s.nonce + 1) % NonceMax

	output := fmt.Sprintf("id: %d\n", s.nonce)
	if event != "" {
		output += fmt.Sprintf("event: %s\n", event)
	}
	if data != nil {
		encoded, err := json.Marshal(data)
		if err != nil {
			return err
		}
		output += fmt.Sprintf("data: %s\n", encoded)
	}
	output += "\n"

	encoded, err := s.encode(output)
	if err != nil {
		return err
	}
	_, err = s.writer.Write(encoded)
	if err != nil {
		return err
	}
	s.flush()
	return nil
}

// sendHeaders sends the headers for Server-Sent Events.
func (s *responseWriter) sendHeaders() {
	s.writer.Header().Set("Content-Type", "text/event-stream")
	s.writer.Header().Set("Cache-Control", "no-store")
	s.writer.Header().Set("Connection", "keep-alive")

	if s.Options.Encoding == Encode_Gzip {
		s.writer.Header().Set("Content-Encoding", "gzip")
	} else if s.Options.Encoding == Encode_Brotli {
		s.writer.Header().Set("Content-Encoding", "br")
	}

	if s.Options.ResponseStatus == 0 {
		s.writer.WriteHeader(http.StatusOK)
	} else {
		s.writer.WriteHeader(s.Options.ResponseStatus)
	}

	s.flush()
}

// flush flushes the response.
func (s *responseWriter) flush() error {
	if f, ok := s.writer.(http.Flusher); ok {
		f.Flush()
		return nil
	}
	return fmt.Errorf("ResponseWriter is not a Flusher")
}
