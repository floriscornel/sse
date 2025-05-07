package sse

import (
	"bytes"
	"compress/flate"
	"compress/gzip"
	"compress/zlib"
	"testing"

	"github.com/andybalholm/brotli"
	"github.com/klauspost/compress/zstd"
)

func TestEncode(t *testing.T) {
	tests := []struct {
		name    string
		level   string
		input   string
		wantErr bool
	}{
		{"NoEncoding", EncodeNone, "Hello, World!", false},
		{"GzipEncoding", EncodeGzip, "Hello, World!", false},
		{"BrotliEncoding", EncodeBrotli, "Hello, World!", false},
		{"DeflateEncoding", EncodeDeflate, "Hello, World!", false},
		{"CompressEncoding", EncodeCompress, "Hello, World!", false},
		{"ZstdEncoding", EncodeZstd, "Hello, World!", false},
		{"UnknownEncoding", "unknown", "Hello, World!", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := encode(tt.level, tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("encode() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && len(got) == 0 {
				t.Errorf("encode() returned empty result for %v", tt.name)
			}
		})
	}
}

func TestEncodeOutputs(t *testing.T) {
	input := "Hello, World!"

	t.Run("NoEncodingOutput", func(t *testing.T) {
		got, err := encode(EncodeNone, input)
		if err != nil {
			t.Errorf("encode() error = %v", err)
			return
		}
		if string(got) != input {
			t.Errorf("encode() = %v, want %v", string(got), input)
		}
	})

	t.Run("GzipEncodingOutput", func(t *testing.T) {
		got, err := encode(EncodeGzip, input)
		if err != nil {
			t.Errorf("encode() error = %v", err)
			return
		}
		reader, err := gzip.NewReader(bytes.NewReader(got))
		if err != nil {
			t.Errorf("gzip.NewReader() error = %v", err)
			return
		}
		defer func() {
			if err := reader.Close(); err != nil {
				t.Errorf("reader.Close() error = %v", err)
			}
		}()
		buf := new(bytes.Buffer)
		_, err = buf.ReadFrom(reader)
		if err != nil {
			t.Errorf("buf.ReadFrom() error = %v", err)
			return
		}
		if buf.String() != input {
			t.Errorf("gzip decoding = %v, want %v", buf.String(), input)
		}
	})

	t.Run("BrotliEncodingOutput", func(t *testing.T) {
		got, err := encode(EncodeBrotli, input)
		if err != nil {
			t.Errorf("encode() error = %v", err)
			return
		}
		reader := brotli.NewReader(bytes.NewReader(got))
		buf := new(bytes.Buffer)
		_, err = buf.ReadFrom(reader)
		if err != nil {
			t.Errorf("buf.ReadFrom() error = %v", err)
			return
		}
		if buf.String() != input {
			t.Errorf("brotli decoding = %v, want %v", buf.String(), input)
		}
	})

	t.Run("DeflateEncodingOutput", func(t *testing.T) {
		got, err := encode(EncodeDeflate, input)
		if err != nil {
			t.Errorf("encode() error = %v", err)
			return
		}
		reader := flate.NewReader(bytes.NewReader(got))
		defer func() {
			if err := reader.Close(); err != nil {
				t.Errorf("reader.Close() error = %v", err)
			}
		}()
		buf := new(bytes.Buffer)
		_, err = buf.ReadFrom(reader)
		if err != nil {
			t.Errorf("buf.ReadFrom() error = %v", err)
			return
		}
		if buf.String() != input {
			t.Errorf("deflate decoding = %v, want %v", buf.String(), input)
		}
	})

	t.Run("CompressEncodingOutput", func(t *testing.T) {
		got, err := encode(EncodeCompress, input)
		if err != nil {
			t.Errorf("encode() error = %v", err)
			return
		}
		reader, err := zlib.NewReader(bytes.NewReader(got))
		if err != nil {
			t.Errorf("zlib.NewReader() error = %v", err)
			return
		}
		defer func() {
			if err := reader.Close(); err != nil {
				t.Errorf("reader.Close() error = %v", err)
			}
		}()
		buf := new(bytes.Buffer)
		_, err = buf.ReadFrom(reader)
		if err != nil {
			t.Errorf("buf.ReadFrom() error = %v", err)
			return
		}
		if buf.String() != input {
			t.Errorf("compress decoding = %v, want %v", buf.String(), input)
		}
	})

	t.Run("ZstdEncodingOutput", func(t *testing.T) {
		got, err := encode(EncodeZstd, input)
		if err != nil {
			t.Errorf("encode() error = %v", err)
			return
		}
		reader, err := zstd.NewReader(bytes.NewReader(got))
		if err != nil {
			t.Errorf("zstd.NewReader() error = %v", err)
			return
		}
		defer reader.Close()
		buf := new(bytes.Buffer)
		_, err = buf.ReadFrom(reader)
		if err != nil {
			t.Errorf("buf.ReadFrom() error = %v", err)
			return
		}
		if buf.String() != input {
			t.Errorf("zstd decoding = %v, want %v", buf.String(), input)
		}
	})
}
