package sse

import (
	"bytes"
	"compress/flate"
	"compress/gzip"
	"compress/zlib"
	"fmt"
	"io"

	"github.com/andybalholm/brotli"
	"github.com/klauspost/compress/zstd"
)

const (
	// No encoding applied
	EncodeNone = ""

	// Compress the data using the DEFLATE algorithm
	// See https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers/Content-Encoding#deflate
	EncodeDeflate = "deflate"

	// Compress the data using the LZW algorithm
	// See https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers/Content-Encoding#compress
	EncodeCompress = "compress"

	// Compress the data using the GZIP algorithm
	// See https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers/Content-Encoding#gzip
	EncodeGzip = "gzip"

	// Compress the data using the Brotli algorithm
	// See https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers/Content-Encoding#br
	EncodeBrotli = "br"

	// Compress the data using the Zstandard algorithm
	// See https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers/Content-Encoding#zstd
	EncodeZstd = "zstd"
)

// encode applies the selected encoding to the input string and returns the encoded bytes.
// It returns an error if the encoding process fails.
func encode(level string, m string) ([]byte, error) {
	buffer := new(bytes.Buffer)

	encoder, err := encoderFactory(level, buffer)
	if err != nil {
		return nil, fmt.Errorf("error creating encoder: %v", err)
	}
	if encoder == nil {
		return []byte(m), nil
	}

	if _, err = encoder.Write([]byte(m)); err != nil {
		return nil, fmt.Errorf("error writing to encoder: %v", err)
	}

	if err = encoder.Close(); err != nil {
		return nil, fmt.Errorf("error closing encoder: %v", err)
	}

	return buffer.Bytes(), nil
}

// encoderFactory creates an encoder based on the provided level.
func encoderFactory(level string, buffer *bytes.Buffer) (io.WriteCloser, error) {
	switch level {
	case EncodeNone:
		return nil, nil
	case EncodeGzip:
		return gzip.NewWriter(buffer), nil
	case EncodeBrotli:
		return brotli.NewWriter(buffer), nil
	case EncodeDeflate:
		return flate.NewWriter(buffer, flate.DefaultCompression)
	case EncodeCompress:
		return zlib.NewWriter(buffer), nil
	case EncodeZstd:
		return zstd.NewWriter(buffer)
	default:
		return nil, fmt.Errorf("unknown encoding level: %s", level)
	}
}
