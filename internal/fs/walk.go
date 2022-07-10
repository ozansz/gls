package fs

import (
	"fmt"
	"os"
	"sort"
	"syscall"

	"github.com/ozansz/gls/internal/local"
	"github.com/ozansz/gls/internal/types"
	"github.com/ozansz/gls/log"
)

const (
	SizeOfBlock = 512
)

type WalkOptions struct {
	IgnoreChecker *local.IgnoreChecker
	SizeThreshold int64
}

func Walk(path string, opts *WalkOptions) (*types.Node, error) {
	f, err := os.Lstat(path)
	if err != nil {
		log.Warningf("%s: %v", path, err)
		return nil, nil
	}
	st, ok := f.Sys().(*syscall.Stat_t)
	if !ok {
		return nil, fmt.Errorf("could not cast %T to syscall.Stat_t", f.Sys())
	}
	root := &types.Node{
		Name:             f.Name(),
		Mode:             f.Mode(),
		Size:             st.Size,
		SizeOnDisk:       st.Blocks * SizeOfBlock,
		IsDir:            f.IsDir(),
		LastModification: f.ModTime(),
	}
	if root.IsDir {
		names, err := readDirNames(path)
		if err != nil {
			log.Warningf("%s: %v", path, err)
			return root, nil
		}
		for _, name := range names {
			child, err := Walk(path+"/"+name, opts)
			if err != nil {
				return nil, err
			}
			child.Parent = root
			root.Size += child.Size
			root.SizeOnDisk += child.SizeOnDisk
			threshOK := child.Size >= opts.SizeThreshold
			ignoreOK := true
			if opts.IgnoreChecker != nil && opts.IgnoreChecker.ShouldIgnore(child.Name, child.IsDir) {
				ignoreOK = false
			}
			if !ignoreOK {
				log.Debugf("ignore: %s", path+"/"+child.Name)
			}
			if threshOK && ignoreOK {
				root.Children = append(root.Children, child)
			}
		}
	}
	return root, nil
}

// Taken directly from the Go stdlib
func readDirNames(dirname string) ([]string, error) {
	f, err := os.Open(dirname)
	if err != nil {
		return nil, err
	}
	names, err := f.Readdirnames(-1)
	f.Close()
	if err != nil {
		return nil, err
	}
	sort.Strings(names)
	return names, nil
}
