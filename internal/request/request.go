package request

import (
	"bytes"
	"errors"
	"io"
	"strings"
	"unicode"

	"httpfromtcp/internal/headers"
)

type Request struct {
	RequestLine RequestLine
	Headers     headers.Headers
	Body        []byte
	state       parserState
}

type parserState int

const (
	requestStateInitialized parserState = iota
	requestStateParsingHeaders
	requestStateDone
)

type RequestLine struct {
	HttpVersion   string
	RequestTarget string
	Method        string
}

func RequestFromReader(reader io.Reader) (*Request, error) {
	data, err := io.ReadAll(reader)
	if err != nil {
		return nil, err
	}

	r := &Request{
		Headers: headers.NewHeaders(),
		state:   requestStateInitialized,
	}

	_, err = r.parse(data)
	if err != nil {
		return nil, err
	}

	if r.state != requestStateDone {
		return nil, errors.New("incomplete request")
	}
	return r, nil
}

func parseRequestLine(data []byte) (RequestLine, int, error) {
	idx := bytes.Index(data, []byte("\r\n"))

	if idx == -1 {
		return RequestLine{}, 0, nil
	}

	line := string(data[:idx])
	consumed := idx + 2

	parts := strings.Split(line, " ")
	if len(parts) != 3 {
		return RequestLine{}, consumed, errors.New("invalid request line")
	}
	method := parts[0]
	target := parts[1]
	versionPart := parts[2]

	// validate method
	for _, ch := range method {
		if !unicode.IsLetter(ch) || !unicode.IsUpper(ch) {
			return RequestLine{}, consumed, errors.New("invalid method")
		}
	}
	// validate target

	// validate version

	if !strings.HasPrefix(versionPart, "HTTP/") {
		return RequestLine{}, consumed, errors.New("invalid method")
	}
	version := strings.TrimPrefix(versionPart, "HTTP/")

	if version != "1.1" {
		return RequestLine{}, consumed, errors.New("invalid HTTP version")
	}
	return RequestLine{
		Method:        method,
		RequestTarget: target,
		HttpVersion:   version,
	}, consumed, nil
}

func (r *Request) parse(data []byte) (int, error) {
	totalBytesParsed := 0

	for r.state != requestStateDone {
		n, err := r.parseSingle(data[totalBytesParsed:])
		if err != nil {
			return totalBytesParsed, err
		}

		// n == 0 means there isn't enough data yet
		if n == 0 {
			break
		}

		totalBytesParsed += n
	}

	return totalBytesParsed, nil
}

func (r *Request) parseSingle(data []byte) (int, error) {
	switch r.state {
	case requestStateInitialized:
		requestLine, consumed, err := parseRequestLine(data)
		if err != nil {
			return 0, err
		}

		if consumed == 0 {
			return 0, nil
		}

		r.RequestLine = requestLine
		r.state = requestStateParsingHeaders

		return consumed, nil

	case requestStateParsingHeaders:
		consumed, done, err := r.Headers.Parse(data)
		if err != nil {
			return consumed, err
		}

		if done {
			r.state = requestStateDone
		}

		return consumed, nil

	case requestStateDone:
		return 0, errors.New("cannot parse request in done state")

	default:
		return 0, errors.New("unknown parser state")
	}
}
