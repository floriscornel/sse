package sse

import (
	"bytes"
	"compress/gzip"

	"github.com/andybalholm/brotli"
)

const (
	Encode_None = iota
	Encode_Gzip
	Encode_Brotli
)

type Encoder interface {
	Write(p []byte) (int, error)
	Close() error
}

func (s *responseWriter) encode(m string) ([]byte, error) {
	if s.Options.Encoding == Encode_None {
		return []byte(m), nil
	}

	var buffer bytes.Buffer
	var encoder Encoder

	switch s.Options.Encoding {
	case Encode_Gzip:
		encoder = gzip.NewWriter(&buffer)
	case Encode_Brotli:
		encoder = brotli.NewWriter(&buffer)
	}

	_, err := encoder.Write([]byte(m))
	if err != nil {
		return []byte{}, err
	}
	err = encoder.Close()
	if err != nil {
		return []byte{}, err
	}
	return buffer.Bytes(), nil
}
