package sse

import (
	"bytes"
	"compress/flate"
	"compress/gzip"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/andybalholm/brotli"
	"github.com/klauspost/compress/zstd"
)

func TestResponseWriter_Write(t *testing.T) {
	tests := []struct {
		name          string
		event         string
		data          interface{}
		options       Options
		expectedData  string
		expectedError bool
	}{
		{
			name:         "No Encoding",
			event:        "test-event",
			data:         map[string]string{"key": "value"},
			options:      Options{Encoding: EncodeNone},
			expectedData: "id: 1\nevent: test-event\ndata: {\"key\":\"value\"}\n\n",
		},
		{
			name:         "Gzip Encoding",
			event:        "test-event",
			data:         map[string]string{"key": "value"},
			options:      Options{Encoding: EncodeGzip},
			expectedData: "id: 1\nevent: test-event\ndata: {\"key\":\"value\"}\n\n",
		},
		{
			name:         "Brotli Encoding",
			event:        "test-event",
			data:         map[string]string{"key": "value"},
			options:      Options{Encoding: EncodeBrotli},
			expectedData: "id: 1\nevent: test-event\ndata: {\"key\":\"value\"}\n\n",
		},
		{
			name:         "Deflate Encoding",
			event:        "test-event",
			data:         map[string]string{"key": "value"},
			options:      Options{Encoding: EncodeDeflate},
			expectedData: "id: 1\nevent: test-event\ndata: {\"key\":\"value\"}\n\n",
		},
		{
			name:         "Zstandard Encoding",
			event:        "test-event",
			data:         map[string]string{"key": "value"},
			options:      Options{Encoding: EncodeZstd},
			expectedData: "id: 1\nevent: test-event\ndata: {\"key\":\"value\"}\n\n",
		},
		{
			name:          "Invalid Data",
			event:         "test-event",
			data:          make(chan int),
			options:       Options{Encoding: EncodeNone},
			expectedError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rec := httptest.NewRecorder()
			writer := NewResponseWriter(rec, tt.options)

			err := writer.Write(tt.event, tt.data)
			if (err != nil) != tt.expectedError {
				t.Fatalf("expected error: %v, got: %v", tt.expectedError, err)
			}

			if tt.expectedError {
				return
			}

			var output string
			switch tt.options.Encoding {
			case EncodeGzip:
				r, err := gzip.NewReader(rec.Body)
				if err != nil {
					t.Fatalf("failed to create gzip reader: %v", err)
				}
				defer r.Close()
				buf := new(bytes.Buffer)
				_, err = buf.ReadFrom(r)
				if err != nil {
					t.Fatalf("failed to read from gzip reader: %v", err)
				}
				output = buf.String()
			case EncodeBrotli:
				r := brotli.NewReader(rec.Body)
				buf := new(bytes.Buffer)
				_, err = buf.ReadFrom(r)
				if err != nil {
					t.Fatalf("failed to read from brotli reader: %v", err)
				}
				output = buf.String()
			case EncodeDeflate:
				r := flate.NewReader(rec.Body)
				defer r.Close()
				buf := new(bytes.Buffer)
				_, err = buf.ReadFrom(r)
				if err != nil {
					t.Fatalf("failed to read from deflate reader: %v", err)
				}
				output = buf.String()
			case EncodeZstd:
				r, err := zstd.NewReader(rec.Body)
				if err != nil {
					t.Fatalf("failed to create zstd reader: %v", err)
				}
				defer r.Close()
				buf := new(bytes.Buffer)
				_, err = buf.ReadFrom(r)
				if err != nil {
					t.Fatalf("failed to read from zstd reader: %v", err)
				}
				output = buf.String()
			default:
				output = rec.Body.String()
			}

			if output != tt.expectedData {
				t.Fatalf("expected data: %q, got: %q", tt.expectedData, output)
			}
		})
	}
}

func TestResponseWriter_SendHeaders(t *testing.T) {
	rec := httptest.NewRecorder()
	options := Options{Encoding: EncodeGzip, ResponseStatus: http.StatusAccepted}
	_ = NewResponseWriter(rec, options)

	headers := rec.Result().Header
	if headers.Get("Content-Type") != "text/event-stream" {
		t.Errorf(
			"expected Content-Type header to be 'text/event-stream', got %s",
			headers.Get("Content-Type"),
		)
	}
	if headers.Get("Cache-Control") != "no-store" {
		t.Errorf(
			"expected Cache-Control header to be 'no-store', got %s",
			headers.Get("Cache-Control"),
		)
	}
	if headers.Get("Connection") != "keep-alive" {
		t.Errorf("expected Connection header to be 'keep-alive', got %s", headers.Get("Connection"))
	}
	if headers.Get("Content-Encoding") != EncodeGzip {
		t.Errorf(
			"expected Content-Encoding header to be 'gzip', got %s",
			headers.Get("Content-Encoding"),
		)
	}
	if rec.Result().StatusCode != http.StatusAccepted {
		t.Errorf(
			"expected status code to be %d, got %d",
			http.StatusAccepted,
			rec.Result().StatusCode,
		)
	}
}

func TestResponseWriter_Flush(t *testing.T) {
	rec := httptest.NewRecorder()
	writer := NewResponseWriter(rec, Options{}).(*responseWriter)

	err := writer.flush()
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	rec = httptest.NewRecorder()
	writer = &responseWriter{
		writer:  &nonFlusherWriter{ResponseWriter: rec},
		options: Options{},
	}

	err = writer.flush()
	if err == nil {
		t.Fatalf("expected error, got nil")
	}
}

type nonFlusherWriter struct {
	http.ResponseWriter
}
