package litexml

import (
	"bytes"
	"strings"
)

// Marshal writes a liteXML struct to a liteXML document.
func Marshal(v interface{}) ([]byte, error) {
	buf := bytes.Buffer{}
	err := NewEncoder(&buf).Encode(v)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

// Unmarshal parses an liteXML document into a liteXML struct.
func Unmarshal(data []byte, v interface{}) error {
	return NewDecoder(strings.NewReader(string(data))).Decode(v)
}
