package fs

import (
	"fmt"

	"github.com/ozansz/gls/internal/local"
	"github.com/ozansz/gls/internal/types"
)

type FileTreeBuilderOption func(*FileTreeBuilder)

type FileTreeBuilder struct {
	root          *types.Node
	path          string
	sort          bool
	sizeFormatter types.SizeFormatter
	sizeThreshold int64
	ignoreChecker *local.IgnoreChecker
}

func NewFileTreeBuilder(path string, opts ...FileTreeBuilderOption) *FileTreeBuilder {
	b := &FileTreeBuilder{
		root:          nil,
		path:          path,
		sort:          false,
		sizeFormatter: types.NoFormat,
		sizeThreshold: 0,
		ignoreChecker: nil,
	}
	for _, opt := range opts {
		opt(b)
	}
	return b
}

func WithSizeFormatter(f types.SizeFormatter) FileTreeBuilderOption {
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

func (b *FileTreeBuilder) Root() *types.Node {
	return b.root
}

func (b *FileTreeBuilder) Build() error {
	var err error
	b.root, err = Walk(b.path, &WalkOptions{
		SizeThreshold: b.sizeThreshold,
		IgnoreChecker: b.ignoreChecker,
	})
	if err != nil {
		return err
	}
	if b.root == nil {
		return fmt.Errorf("could not build, root is nil")
	}
	if b.sort {
		b.root.SortChildrenBySizeOnDisk()
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
