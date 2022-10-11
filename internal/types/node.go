package types

import (
	"fmt"
	"os"
	"regexp"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/ozansz/gls/internal"
	"github.com/ozansz/gls/internal/analyzer"
	"github.com/ozansz/gls/internal/helper"
)

type Node struct {
	mu   sync.Mutex
	Name string
	Mode os.FileMode
	Size int64

	SizeOnDisk int64
	// Blocks           int64

	IsDir            bool
	LastModification time.Time
	Children         []*Node
	Parent           *Node
}

func (n *Node) FileCount() int {
	if n == nil {
		return 0
	}
	n.mu.Lock()
	defer n.mu.Unlock()
	if n.IsDir {
		count := 0
		for _, c := range n.Children {
			count += c.FileCount()
		}
		return count
	}
	return 1
}

type cloneOpts struct {
	discardFiles bool
}

func newNoOpCloneOpts() *cloneOpts {
	return &cloneOpts{
		discardFiles: false,
	}
}

func (n *Node) CloneDirectoryStructure() (*Node, error) {
	return n.clone(&cloneOpts{
		discardFiles: true,
	})
}

func (n *Node) clone(opts *cloneOpts) (*Node, error) {
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
		if c := child.cloneWithParent(root, opts); c != nil {
			root.AddChild(c)
		}
	}
	return root, nil
}

func (n *Node) cloneWithParent(root *Node, opts *cloneOpts) *Node {
	n.mu.Lock()
	defer n.mu.Unlock()
	if opts.discardFiles && !root.IsDir {
		return nil
	}
	clone := &Node{
		Name:             n.Name,
		Mode:             n.Mode,
		Size:             n.Size,
		IsDir:            n.IsDir,
		LastModification: n.LastModification,
		Parent:           root,
	}
	for _, child := range n.Children {
		if c := child.cloneWithParent(clone, opts); c != nil {
			clone.AddChild(c)
		}
	}
	return clone
}

type TreeFilterOptions struct {
	nameContains    string
	re              *regexp.Regexp
	caseInsensitive bool
	invertSelection bool
}

func NewTreeFilterOpts(nameContains, reMatches string, caseInsensitive, invert bool) (*TreeFilterOptions, error) {
	o := &TreeFilterOptions{
		nameContains:    nameContains,
		caseInsensitive: caseInsensitive,
		invertSelection: invert,
	}
	if reMatches == "" {
		o.re = nil
	} else {
		if caseInsensitive {
			reMatches = "(?i)" + reMatches
		}
		var err error
		o.re, err = regexp.Compile(reMatches)
		if err != nil {
			return nil, err
		}
	}
	return o, nil
}

func (o *TreeFilterOptions) IsNoOp() bool {
	if o.nameContains == "" && o.re == nil {
		return true
	}
	return false
}

func (o *TreeFilterOptions) CheckNameContains(original string) bool {
	if o.nameContains == "" {
		return true // pass the test
	}
	name := original
	sub := o.nameContains
	if o.caseInsensitive {
		name = strings.ToLower(name)
		sub = strings.ToLower(sub)
	}
	return strings.Contains(name, sub) != o.invertSelection // return (ok XOR invert)
}

func (o *TreeFilterOptions) CheckRegexMatches(original string) bool {
	if o.re == nil {
		return true // pass the test
	}
	name := original
	if o.caseInsensitive {
		name = strings.ToLower(name)
	}
	return o.re.MatchString(name) != o.invertSelection // return (ok XOR invert)
}

func (n *Node) NewFilteredTree(opts *TreeFilterOptions) (*Node, error) {
	tree, err := n.clone(newNoOpCloneOpts())
	if err != nil {
		return nil, err
	}
	if opts.IsNoOp() {
		return tree, nil
	}
	weights := make(map[*Node]int)
	weights[tree] = getSearchTreeWeight(tree, weights, opts)
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

func getSearchTreeWeight(n *Node, weights map[*Node]int, opts *TreeFilterOptions) int {
	n.mu.Lock()
	defer n.mu.Unlock()
	if n.IsDir {
		weight := 0
		for _, child := range n.Children {
			weight += getSearchTreeWeight(child, weights, opts)
		}
		weights[n] = weight
		return weight
	} else {
		if opts.CheckNameContains(n.Name) && opts.CheckRegexMatches(n.Name) {
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

func (n *Node) CreateChild(fileName, parentPath string) error {
	n.mu.Lock()
	defer n.mu.Unlock()
	if !n.IsDir {
		return fmt.Errorf("cannot create file under a file (%s)", n.Name)
	}
	filePath := n.RelativePath(parentPath) + "/" + fileName
	if _, err := os.Stat(filePath); !os.IsNotExist(err) {
		return fmt.Errorf("file with path %s already exists", filePath)
	}
	f, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("could not create file %s: %v", filePath, err)
	}
	fInfo, err := f.Stat()
	if err != nil {
		return fmt.Errorf("could not stat file %s: %v", filePath, err)
	}
	var diskUsage helper.DiskUsage = &helper.FsInfo{}
	size, err := diskUsage.GetSize(fInfo)
	if err != nil {
		return err
	}
	n.Children = append(n.Children, &Node{
		Name:             fInfo.Name(),
		Mode:             fInfo.Mode(),
		Size:             size,
		SizeOnDisk:       size * internal.UNIXSizeOfBlock,
		IsDir:            fInfo.IsDir(),
		LastModification: fInfo.ModTime(),
		Parent:           n,
	})
	if err = f.Close(); err != nil && err != os.ErrClosed {
		return fmt.Errorf("error while closing the newly created file %s: %v", filePath, err)
	}
	return nil
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

func (n *Node) SortChildrenBySizeOnDisk() {
	for _, c := range n.Children {
		c.SortChildrenBySizeOnDisk()
	}
	sort.Slice(n.Children, func(i, j int) bool {
		return n.Children[i].SizeOnDisk > n.Children[j].SizeOnDisk
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
	return fmt.Sprintf("%s%s [%s]", strings.Repeat("  ", level), n.Name, f(n.SizeOnDisk))
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
