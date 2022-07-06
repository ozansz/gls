package main

import (
	"flag"
	"fmt"

	"github.com/ozansz/gls/internal"
)

var (
	path      = flag.String("path", "", "path to list")
	formatter = flag.String("fmt", "bytes", "formatter to use: bytes, pow10 or none")

	formatters = map[string]internal.SizeFormatter{
		"bytes": internal.SizeFormatterBytes,
		"pow10": internal.SizeFormatterPow10,
		"none":  internal.NoFormat,
	}
)

func main() {
	flag.Parse()
	if *path == "" {
		flag.Usage()
		return
	}
	formatterFunc, ok := formatters[*formatter]
	if !ok {
		fmt.Printf("Unknown formatter: %s\n\n", *formatter)
		flag.Usage()
		return
	}
	opts := []internal.FileTreeBuilderOption{
		internal.WithSizeFormatter(formatterFunc),
	}
	b := internal.NewFileTreeBuilder(*path, opts...)
	b.Build()
	if err := b.Print(); err != nil {
		fmt.Printf("Error: %s\n", err)
	}
}
