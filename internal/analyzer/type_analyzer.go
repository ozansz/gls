package analyzer

import (
	"io"
	"os"

	"github.com/h2non/filetype"
	"github.com/h2non/filetype/types"
)

const (
	lookupBytes = 300
)

var (
	cache = make(map[string]types.Type)
)

func AnalyzeFileType(path string) (types.Type, error) {
	if t, ok := cache[path]; ok {
		return t, nil
	}
	file, err := os.Open(path)
	if err != nil {
		return types.Unknown, err
	}
	defer file.Close()
	buf := make([]byte, lookupBytes)
	_, err = file.Read(buf)
	if err != nil && err != io.EOF {
		return types.Unknown, err
	}
	typ, err := filetype.Match(buf)
	if err != nil {
		return types.Unknown, err
	}
	cache[path] = typ
	return typ, nil
}
