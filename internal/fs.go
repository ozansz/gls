package internal

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"
)

type FileTreeBuilderOption func(*FileTreeBuilder)
type SizeFormatter func(int64) string

func NoFormat(size int64) string {
	return fmt.Sprint(size)
}

func SizeFormatterBytes(size int64) string {
	const (
		KB = 1024
		MB = 1024 * KB
		GB = 1024 * MB
		TB = 1024 * GB
	)
	if size < KB {
		return fmt.Sprintf("%d B", size)
	}
	if size < MB {
		return fmt.Sprintf("%.2f KB", float64(size)/float64(KB))
	}
	if size < GB {
		return fmt.Sprintf("%.2f MB", float64(size)/float64(MB))
	}
	if size < TB {
		return fmt.Sprintf("%.2f GB", float64(size)/float64(GB))
	}
	return fmt.Sprintf("%.2f TB", float64(size)/float64(TB))
}

func SizeFormatterPow10(size int64) string {
	const (
		KiB = 1000
		MiB = 1000 * KiB
		GiB = 1000 * MiB
		TiB = 1000 * GiB
	)
	if size < KiB {
		return fmt.Sprintf("%d b", size)
	}
	if size < MiB {
		return fmt.Sprintf("%.2f KiB", float64(size)/float64(KiB))
	}
	if size < GiB {
		return fmt.Sprintf("%.2f MiB", float64(size)/float64(MiB))
	}
	if size < TiB {
		return fmt.Sprintf("%.2f GiB", float64(size)/float64(GiB))
	}
	return fmt.Sprintf("%.2f TiB", float64(size)/float64(TiB))
}

type FileTreeBuilder struct {
	root          *Node
	path          string
	sizeFormatter SizeFormatter
}

func NewFileTreeBuilder(path string, opts ...FileTreeBuilderOption) *FileTreeBuilder {
	b := &FileTreeBuilder{
		root:          nil,
		path:          path,
		sizeFormatter: NoFormat,
	}
	for _, opt := range opts {
		opt(b)
	}
	return b
}

func WithSizeFormatter(f SizeFormatter) FileTreeBuilderOption {
	return func(b *FileTreeBuilder) {
		b.sizeFormatter = f
	}
}

func (b *FileTreeBuilder) Build() {
	b.root = listDir(b.path)
}

func (b *FileTreeBuilder) Print() error {
	if b.root == nil {
		return fmt.Errorf("no root node built")
	}
	b.root.PrintWithSizeFormatter(b.sizeFormatter)
	return nil
}

func (b *FileTreeBuilder) RootInfo() error {
	if b.root == nil {
		return fmt.Errorf("no root node built")
	}
	b.root.InfoWithSizeFormatter(b.sizeFormatter)
	return nil
}

func listDir(path string) *Node {
	root := &Node{
		Name:     filepath.Base(path),
		Mode:     os.ModeDir,
		IsDir:    true,
		Children: make([]*Node, 0),
	}
	myWc := make(chan struct{})
	var wg sync.WaitGroup
	wg.Add(1)
	go walkDir(path, &wg, root, myWc)
	wg.Wait()
	return root
}

func walkDir(dir string, wg *sync.WaitGroup, root *Node, wc chan struct{}) {
	defer func() {
		wg.Done()
		wc <- struct{}{}
	}()
	visit := func(path string, f os.FileInfo, err error) error {
		if f.IsDir() && path != dir {
			this := &Node{
				Name:     f.Name(),
				Mode:     f.Mode(),
				IsDir:    true,
				Children: []*Node{},
			}
			myWc := make(chan struct{})
			wg.Add(1)
			go walkDir(path, wg, this, myWc)
			<-myWc
			root.AddChild(this)
			root.IncrementSize(this.Size)
			return filepath.SkipDir
		}
		if f.Mode().IsRegular() {
			size := f.Size()
			root.AddChild(&Node{
				Name:  f.Name(),
				Mode:  f.Mode(),
				Size:  size,
				IsDir: false,
			})
			root.IncrementSize(size)
		}
		return nil
	}
	filepath.Walk(dir, visit)
}
