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

// Options holds configuration for the SSE writer.
type Options struct {
	ResponseStatus int
	Encoding       string
}

// NewResponseWriter creates a new Writer for Server-Sent Events.
func NewResponseWriter(w http.ResponseWriter, opts Options) Writer {
	rw := &responseWriter{
		writer:  w,
		nonce:   0,
		options: opts,
	}
	rw.sendHeaders()
	return rw
}

type responseWriter struct {
	writer  http.ResponseWriter
	nonce   uint64
	options Options
}

// Write sends a message to the client.
func (rw *responseWriter) Write(event string, data interface{}) error {
	rw.nonce = (rw.nonce + 1) % NonceMax

	output := fmt.Sprintf("id: %d\n", rw.nonce)
	if event != "" {
		output += fmt.Sprintf("event: %s\n", event)
	}
	if data != nil {
		encodedData, err := json.Marshal(data)
		if err != nil {
			return err
		}
		output += fmt.Sprintf("data: %s\n", encodedData)
	}
	output += "\n"

	encodedOutput, err := encode(rw.options.Encoding, output)
	if err != nil {
		return err
	}

	if _, err := rw.writer.Write(encodedOutput); err != nil {
		return err
	}
	return rw.flush()
}

// sendHeaders sends the headers for Server-Sent Events.
func (rw *responseWriter) sendHeaders() {
	headers := rw.writer.Header()
	headers.Set("Content-Type", "text/event-stream")
	headers.Set("Cache-Control", "no-store")
	headers.Set("Connection", "keep-alive")

	if rw.options.Encoding != EncodeNone {
		headers.Set("Content-Encoding", rw.options.Encoding)
	}

	status := rw.options.ResponseStatus
	if status == 0 {
		status = http.StatusOK
	}
	rw.writer.WriteHeader(status)
	//nolint:staticcheck
	if err := rw.flush(); err != nil {
		// Intentionally ignored: cannot recover from flush error after headers are sent
	}
}

// flush flushes the response.
func (rw *responseWriter) flush() error {
	if flusher, ok := rw.writer.(http.Flusher); ok {
		flusher.Flush()
		return nil
	}
	return fmt.Errorf("ResponseWriter is not a Flusher")
}
