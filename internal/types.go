package internal

import (
	"fmt"
	"os"
	"strings"
	"sync"
)

type Node struct {
	Name     string
	Mode     os.FileMode
	Size     int64
	IsDir    bool
	Children []*Node
	mu       sync.Mutex
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

func (n *Node) Print() {
	n.printWithLevel(0, NoFormat)
}

func (n *Node) PrintWithSizeFormatter(f SizeFormatter) {
	n.printWithLevel(0, f)
}

func (n *Node) printWithLevel(level int, f SizeFormatter) {
	n.mu.Lock()
	defer n.mu.Unlock()
	n.infoWithLevel(level, f)
	for _, child := range n.Children {
		child.printWithLevel(level+1, f)
	}
}

func (n *Node) Info() {
	n.infoWithLevel(0, NoFormat)
}

func (n *Node) InfoWithSizeFormatter(f SizeFormatter) {
	n.infoWithLevel(0, f)
}

func (n *Node) infoWithLevel(level int, f SizeFormatter) {
	fmt.Printf("%s%s [%s] [%s]\n", strings.Repeat("  ", level), n.Name, n.Mode.String(), f(n.Size))
}
