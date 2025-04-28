package resp

type SimpleString struct {
	Value string
}

type BulkString struct {
	Value  string
	Length int
}

type Integer struct {
	Value int
}

type Array struct {
	Elements []any
}
