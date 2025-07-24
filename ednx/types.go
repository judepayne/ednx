package ednx

import (
	"errors"
	"os"
)

// Error variables for the ednx package
var (
	ErrJsonUnmarshal = errors.New("error reading json")
	ErrEdnUnmarshal  = errors.New("error reading edn")
	ErrEdnPrettify   = errors.New("error prettifying edn")
	ErrEdnMarshal    = errors.New("error converting to edn")
	ErrJsonPrettify  = errors.New("error prettifying json")
	ErrJsonMarshal   = errors.New("error converting to json")
)

// EdnConvertOptions configures EDN conversion
type EdnConvertOptions struct {
	KeywordizeKeys bool
	PrettyPrint    bool
	WidthLimit     int
}

// Write implements io.Writer interface for EdnConvertOptions
func (opts *EdnConvertOptions) Write(p []byte) (n int, err error) {
	return os.Stdout.Write(p)
}

// JsonConvertOptions configures JSON conversion
type JsonConvertOptions struct {
	PrettyPrint bool
}

// Write implements io.Writer interface for JsonConvertOptions
func (opts *JsonConvertOptions) Write(p []byte) (n int, err error) {
	return os.Stdout.Write(p)
}