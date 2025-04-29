package resp

import (
	"fmt"
	"io"
)

type Parser struct {
	reader io.Reader
}

func NewParser(reader io.Reader) *Parser {
	return &Parser{
		reader: reader,
	}
}

func (p *Parser) Parse() (any, error) {
	firstByte := make([]byte, 1)
	_, err := p.reader.Read(firstByte)
	if err != nil {
		return nil, err
	}

	switch firstByte[0] {
	case '+':
		return p.parseSimpleString()
	case '-':
		return p.parseError()
	case ':':
		return p.parseInteger()
	case '$':
		return p.parseBulkString()
	case '*':
		return p.parseArray()
	default:
		return nil, &ErrInvalidFormat{"invalid RESP type"}
	}
}

func (p *Parser) readLine() ([]byte, error) {
	line := make([]byte, 0, 256)
	for {
		b := make([]byte, 1)
		_, err := p.reader.Read(b)
		if err != nil {
			return nil, err
		}
		if b[0] == '\r' {
			continue
		}
		if b[0] == '\n' {
			break
		}
		line = append(line, b[0])
	}
	return line, nil
}

func (p *Parser) parseSimpleString() (any, error) {
	line, err := p.readLine()
	if err != nil {
		return nil, err
	}
	return &SimpleString{Value: string(line)}, nil
}

func (p *Parser) parseError() (any, error) {
	line, err := p.readLine()
	if err != nil {
		return nil, err
	}
	// TODO: Handle different error types
	return &GenericError{Message: string(line)}, nil
}

func (p *Parser) parseInteger() (any, error) {
	line, err := p.readLine()
	if err != nil {
		return nil, err
	}
	var value int
	_, err = fmt.Sscanf(string(line), "%d", &value)
	if err != nil {
		return nil, err
	}
	return &Integer{Value: value}, nil
}

func (p *Parser) parseBulkString() (any, error) {
	line, err := p.readLine()
	if err != nil {
		return nil, err
	}
	var length int
	_, err = fmt.Sscanf(string(line), "%d", &length)
	if err != nil {
		return nil, err
	}
	if length == -1 {
		return &BulkString{Value: "", Length: -1}, nil
	}
	value := make([]byte, length)
	_, err = p.reader.Read(value)
	if err != nil {
		return nil, err
	}
	p.reader.Read(make([]byte, 2)) // Read the CRLF
	return &BulkString{Value: string(value), Length: length}, nil
}

func (p *Parser) parseArray() (any, error) {
	line, err := p.readLine()
	if err != nil {
		return nil, err
	}
	var length int
	_, err = fmt.Sscanf(string(line), "%d", &length)
	if err != nil {
		return nil, err
	}
	if length == -1 {
		return &Array{Elements: nil}, nil
	}
	elements := make([]any, length)
	for i := range length {
		element, err := p.Parse()
		if err != nil {
			return nil, err
		}
		elements[i] = element
	}
	return &Array{Elements: elements}, nil
}
