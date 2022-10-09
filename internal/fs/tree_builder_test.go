package fs

import (
	"flag"
	"testing"
)

var (
	path = flag.String("path", "", "path to run the benchmark on")
)

func BenchmarkFileTreeBuilder(b *testing.B) {
	for i := 0; i < b.N; i++ {
		bl := NewFileTreeBuilder(*path)
		if err := bl.Build(); err != nil {
			b.Fatalf("error building file tree: %v", err)
		}
	}
}
