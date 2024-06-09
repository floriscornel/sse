# SSE - Server-Sent Events Writer

[![GoDoc](https://godoc.org/github.com/floriscornel/sse?status.svg)](https://godoc.org/github.com/floriscornel/sse)
[![Go Report Card](https://goreportcard.com/badge/github.com/floriscornel/sse)](https://goreportcard.com/report/github.com/floriscornel/sse)
[![Codecov](https://img.shields.io/codecov/c/github/floriscornel/sse.svg)](https://codecov.io/gh/floriscornel/sse)
[![License](https://img.shields.io/github/license/floriscornel/sse.svg)](https://github.com/floriscornel/sse/blob/main/LICENSE)

`sse` is a lightweight Go library for writing [Server-Sent Events (SSE)](https://developer.mozilla.org/en-US/docs/Web/API/Server-sent_events). This library allows easy integration of SSE into web servers, providing real-time updates to connected clients.

## Installation

To install the `sse` package, use the following command:

```bash
go get github.com/floriscornel/sse
```

## Usage

### Importing the Package

```go
import (
    "net/http"
    "github.com/floriscornel/sse"
)
```

### Creating a Server-Sent Events Writer

To create an SSE writer, use the `NewResponseWriter` function. This function requires an `http.ResponseWriter` and options to configure the SSE writer.

```go
opts := sse.Options{
    ResponseStatus: http.StatusOK,
    Encoding:       sse.EncodeNone,
}

func handler(w http.ResponseWriter, r *http.Request) {
    sseWriter := sse.NewResponseWriter(w, opts)

    for {
        err := sseWriter.Write("message", map[string]string{"hello": "world"})
        if err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
            return
        }
        time.Sleep(2 * time.Second)
    }
}

http.HandleFunc("/events", handler)
log.Fatal(http.ListenAndServe(":8080", nil))
```

### Encoding Options

The `sse` package supports various encoding options to compress the data sent to clients. Available encoding options include:

- `EncodeNone`
- `EncodeDeflate`
- `EncodeCompress`
- `EncodeGzip`
- `EncodeBrotli`
- `EncodeZstd`

Specify the desired encoding in the `Options` struct:

```go
opts := sse.Options{
    ResponseStatus: http.StatusOK,
    Encoding:       sse.EncodeGzip,
}
```

## Examples

You can find more examples in the `examples` directory. To run an example, navigate to the respective directory and execute the following command:

```bash
go run main.go
```

### Example Structure

- `ping`: A simple example of sending periodic ping messages to the client.
- `incremental-updates`: Demonstrates sending incremental updates to the client with different event types.
- `number-of-listeners`: Implements a counter to track the number of connected clients.

## Contributing

Contributions are welcome! Feel free to open issues or submit pull requests. For major changes, please open an issue first to discuss what you would like to change.

## License

This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.

## Contact

If you have any questions, feel free to reach out:

- **Author:** Floris Cornel
- **Email:** floris@cornel.email
- **GitHub:** [github.com/floriscornel/sse](https://github.com/floriscornel/sse)
