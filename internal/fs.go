package internal

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"

	"github.com/ozansz/gls/internal/local"
	"github.com/ozansz/gls/log"
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
	sort          bool
	sizeFormatter SizeFormatter
	sizeThreshold int64
	ignoreChecker *local.IgnoreChecker
}

func NewFileTreeBuilder(path string, opts ...FileTreeBuilderOption) *FileTreeBuilder {
	b := &FileTreeBuilder{
		root:          nil,
		path:          path,
		sort:          false,
		sizeFormatter: NoFormat,
		sizeThreshold: 0,
		ignoreChecker: nil,
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

func WithSortingBySize() FileTreeBuilderOption {
	return func(b *FileTreeBuilder) {
		b.sort = true
	}
}

func WithSizeThreshold(thresh int64) FileTreeBuilderOption {
	return func(b *FileTreeBuilder) {
		b.sizeThreshold = thresh
	}
}

func WithIgnoreChecker(ic *local.IgnoreChecker) FileTreeBuilderOption {
	return func(b *FileTreeBuilder) {
		b.ignoreChecker = ic
	}
}

func (b *FileTreeBuilder) Root() *Node {
	return b.root
}

func (b *FileTreeBuilder) Build() error {
	var err error
	b.root, err = listDir(b.path, b.sizeThreshold, b.ignoreChecker)
	if err != nil {
		return err
	}
	if b.sort {
		b.root.SortChildrenBySize()
	}
	return nil
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

func listDir(path string, sizeThreshold int64, ic *local.IgnoreChecker) (*Node, error) {
	finfo, err := os.Lstat(path)
	if err != nil {
		return nil, fmt.Errorf("failed to lstat %s: %v", path, err)
	}
	root := &Node{
		Name:             filepath.Base(path),
		Mode:             finfo.Mode(),
		IsDir:            true,
		Children:         make([]*Node, 0),
		Parent:           nil,
		LastModification: finfo.ModTime(),
	}
	checkIgnore := func(f os.FileInfo) bool {
		if ic == nil {
			return true
		}
		return ic.ShouldIgnore(f.Name(), f.IsDir())
	}
	myWc := make(chan struct{})
	var wg sync.WaitGroup
	wg.Add(1)
	go walkDir(path, &wg, root, myWc, sizeThreshold, checkIgnore)
	wg.Wait()
	return root, nil
}

func walkDir(dir string, wg *sync.WaitGroup,
	root *Node, wc chan struct{},
	sizeThreshold int64, checkIgnore func(os.FileInfo) bool) {
	defer func() {
		wg.Done()
		wc <- struct{}{}
	}()
	visit := func(path string, f os.FileInfo, err error) error {
		if checkIgnore(f) {
			// Ignore this folder/file
			log.Debugf("Ignoring %s/%s", dir, f.Name())
			if f.IsDir() {
				return filepath.SkipDir
			}
			return nil
		}
		if f.IsDir() && path != dir {
			this := &Node{
				Name:             f.Name(),
				Mode:             f.Mode(),
				IsDir:            true,
				Children:         []*Node{},
				Parent:           root,
				LastModification: f.ModTime(),
			}
			myWc := make(chan struct{})
			wg.Add(1)
			go walkDir(path, wg, this, myWc, sizeThreshold, checkIgnore)
			<-myWc
			root.IncrementSize(this.Size)
			if this.Size >= sizeThreshold {
				root.AddChild(this)
			}
			return filepath.SkipDir
		}
		if f.Mode().IsRegular() {
			size := f.Size()
			root.IncrementSize(size)
			if size >= sizeThreshold {
				root.AddChild(&Node{
					Name:             f.Name(),
					Mode:             f.Mode(),
					Size:             size,
					IsDir:            false,
					Parent:           root,
					LastModification: f.ModTime(),
				})
			}
		}
		return nil
	}
	filepath.Walk(dir, visit)
}
