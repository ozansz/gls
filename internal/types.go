package internal

import (
	"fmt"
	"os"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/ozansz/gls/internal/analyzer"
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

func (n *Node) clone() (*Node, error) {
	n.mu.Lock()
	defer n.mu.Unlock()
	if n.Parent != nil {
		return nil, fmt.Errorf("can only clone root node")
	}
	root := &Node{
		Name:             n.Name,
		Mode:             n.Mode,
		Size:             n.Size,
		IsDir:            n.IsDir,
		LastModification: n.LastModification,
		Parent:           nil,
	}
	for _, child := range n.Children {
		root.AddChild(child.cloneWithParent(root))
	}
	return root, nil
}

func (n *Node) cloneWithParent(root *Node) *Node {
	n.mu.Lock()
	defer n.mu.Unlock()
	clone := &Node{
		Name:             n.Name,
		Mode:             n.Mode,
		Size:             n.Size,
		IsDir:            n.IsDir,
		LastModification: n.LastModification,
		Parent:           root,
	}
	for _, child := range n.Children {
		clone.AddChild(child.cloneWithParent(clone))
	}
	return clone
}

func (n *Node) ConstructSearchTreeWithSearchString(substring string) (*Node, error) {
	tree, err := n.clone()
	if err != nil {
		return nil, err
	}
	weights := make(map[*Node]int)
	weights[tree] = getSearchTreeWeight(tree, weights, substring)
	removeZeroWeightsFromSearchTree(tree, weights)
	return tree, nil
}

func removeZeroWeightsFromSearchTree(n *Node, weights map[*Node]int) {
	n.mu.Lock()
	defer n.mu.Unlock()
	newChildren := make([]*Node, 0)
	for _, child := range n.Children {
		if weights[child] > 0 {
			newChildren = append(newChildren, child)
		}
	}
	n.Children = newChildren
	for _, child := range n.Children {
		removeZeroWeightsFromSearchTree(child, weights)
	}
}

func getSearchTreeWeight(n *Node, weights map[*Node]int, substring string) int {
	n.mu.Lock()
	defer n.mu.Unlock()
	if n.IsDir {
		weight := 0
		for _, child := range n.Children {
			weight += getSearchTreeWeight(child, weights, substring)
		}
		weights[n] = weight
		return weight
	} else {
		if strings.Contains(n.Name, substring) {
			weights[n] = 1
			return 1
		}
		weights[n] = 0
		return 0
	}
}

func (n *Node) GetFileType(parentPath string) (string, error) {
	if n.IsDir {
		return "directory", nil
	}
	typ, err := analyzer.AnalyzeFileType(n.RelativePath(parentPath))
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%s (%s)", typ.MIME.Value, typ.Extension), nil
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
