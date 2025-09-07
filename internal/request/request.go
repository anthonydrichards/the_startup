package request

import (
	"fmt"
	"io"
	"strings"
	"unicode"
)

type RequestStatus int

const (
	Initialized RequestStatus = iota
	Done
)

type Request struct {
	RequestLine RequestLine
	Status      RequestStatus
}

type RequestLine struct {
	HttpVersion   string
	RequestTarget string
	Method        string
}

func (r *Request) parse(data []byte) (int, error) {
	if r.Status == Initialized {
		requestLine, bytesConsumed, err := parseRequestLine(string(data))
		if err != nil {
			return 0, err
		}
		if bytesConsumed == 0 {
			return 0, nil
		}
		r.RequestLine = requestLine
		r.Status = Done
		return bytesConsumed, nil
	}
	return 0, nil
}

func RequestFromReader(reader io.Reader) (*Request, error) {
	httpMessageBytes := make([]byte, 8)
	totalBytesRead := 0
	totalBytesParesed := 0
	var request Request
	for {
		newHttpMessageBytes := make([]byte, 8)
		newBytesRead := 0
		newBytesRead, err := reader.Read(newHttpMessageBytes)
		if err != nil {
			if err == io.EOF {
				break
			}
			return &request, err
		}
		httpMessageBytes = append(httpMessageBytes[:totalBytesRead-totalBytesParesed], newHttpMessageBytes[:newBytesRead]...)
		totalBytesRead += newBytesRead
		bytesParsed, err := request.parse(httpMessageBytes)
		if err != nil {
			return &request, err
		}
		totalBytesParesed += bytesParsed
		if bytesParsed != 0 {
			if bytesParsed == len(httpMessageBytes) {
				httpMessageBytes = make([]byte, 8)
			} else {
				httpMessageBytes = httpMessageBytes[bytesParsed:]
			}
		}
	}
	return &request, nil
}

func parseRequestLine(httpMessage string) (RequestLine, int, error) {
	if !strings.Contains(httpMessage, "\r\n") {
		return RequestLine{}, 0, nil
	}
	lines := strings.Split(httpMessage, "\r\n")
	requestLineElements := strings.Split(lines[0], " ")
	if len(requestLineElements) != 3 {
		return RequestLine{}, 0, fmt.Errorf("invalid number of parts in request line: %s", lines[0])
	}

	method := requestLineElements[0]
	requestTarget := requestLineElements[1]
	httpVersion := strings.TrimPrefix(requestLineElements[2], "HTTP/")

	for _, r := range method {
		if !unicode.IsLetter(r) {
			return RequestLine{}, 0, fmt.Errorf("invalid character in method: %s", method)
		}
		if !unicode.IsUpper(r) {
			return RequestLine{}, 0, fmt.Errorf("method must be uppercase: %s", method)
		}
	}
	if httpVersion != "1.1" {
		return RequestLine{}, 0, fmt.Errorf("unsupported HTTP version: %s", httpVersion)
	}

	return RequestLine{
		Method:        method,
		RequestTarget: requestTarget,
		HttpVersion:   httpVersion,
	}, len(lines[0]) + 2, nil
}
