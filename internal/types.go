package internal

import (
	"fmt"
	"os"
	"sort"
	"strings"
	"sync"
	"time"
)

type Node struct {
	mu               sync.Mutex
	Name             string
	Mode             os.FileMode
	Size             int64
	IsDir            bool
	LastModification time.Time
	Children         []*Node
	Parent           *Node
}

func (n *Node) Remove(parentPath string) error {
	n.mu.Lock()
	defer n.mu.Unlock()
	if n.IsDir {
		return fmt.Errorf("cannot remove directory %s", n.Name)
	}
	n.Parent.RemoveChild(n)
	return os.Remove(n.RelativePath(parentPath))
}

func (n *Node) RemoveChild(c *Node) {
	n.mu.Lock()
	defer n.mu.Unlock()
	for i, child := range n.Children {
		if child == c {
			n.Children = append(n.Children[:i], n.Children[i+1:]...)
			return
		}
	}
}

func (n *Node) IncrementSize(size int64) {
	n.mu.Lock()
	defer n.mu.Unlock()
	n.Size += size
}

func (n *Node) AddChild(child *Node) {
	n.mu.Lock()
	defer n.mu.Unlock()
	n.Children = append(n.Children, child)
}

func (n *Node) SortChildrenBySize() {
	for _, c := range n.Children {
		c.SortChildrenBySize()
	}
	sort.Slice(n.Children, func(i, j int) bool {
		return n.Children[i].Size > n.Children[j].Size
	})
}

func (n *Node) Print() {
	n.printWithLevel(0, NoFormat)
}

func (n *Node) PrintWithSizeFormatter(f SizeFormatter) {
	n.printWithLevel(0, f)
}

func (n *Node) printWithLevel(level int, f SizeFormatter) {
	n.mu.Lock()
	defer n.mu.Unlock()
	fmt.Println(n.infoWithLevel(level, f))
	for _, child := range n.Children {
		child.printWithLevel(level+1, f)
	}
}

func (n *Node) Info() string {
	return n.infoWithLevel(0, NoFormat)
}

func (n *Node) InfoWithSizeFormatter(f SizeFormatter) string {
	return n.infoWithLevel(0, f)
}

func (n *Node) infoWithLevel(level int, f SizeFormatter) string {
	return fmt.Sprintf("%s%s [%s] [%s]", strings.Repeat("  ", level), n.Name, n.Mode.String(), f(n.Size))
}

func (n *Node) RelativePath(parent string) string {
	if n.Parent == nil {
		if parent == "" {
			return n.Name
		}
		return parent
	}
	return n.Parent.RelativePath(parent) + "/" + n.Name
}
