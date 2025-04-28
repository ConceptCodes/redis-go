package resp

import "fmt"

type ErrInvalidDataType struct {
	Data byte
}

func (e *ErrInvalidDataType) Error() string {
	return fmt.Sprintf("invalid RESP data type: %c", e.Data)
}

type ErrUnexpectedEOF struct{}

func (e *ErrUnexpectedEOF) Error() string {
	return "unexpected EOF in RESP data"
}

type ErrInvalidFormat struct {
	Reason string
}

func (e *ErrInvalidFormat) Error() string {
	return fmt.Sprintf("-ERR invalid format: %s\r\n", e.Reason)
}

type GenericError struct {
	Message string
}

func (e *GenericError) Error() string {
	return fmt.Sprintf("-ERR : %s\r\n", e.Message)
}
