package fs

import (
	"context"
	"os"
	"sort"
	"strings"
	"sync"

	"github.com/ozansz/gls/internal/helper"
	"github.com/ozansz/gls/internal/local"
	"github.com/ozansz/gls/internal/types"
	"github.com/ozansz/gls/log"
	"golang.org/x/sync/errgroup"
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
	//implements os dependent proper methods.
	var diskUsage helper.DiskUsage = &helper.FsInfo{}
	size, _ := diskUsage.GetSize(f)
	if err != nil {
		return nil, err
	}
	sizeOnDisk, err := diskUsage.GetSizeOnDisk(f)
	if err != nil {
		return nil, err
	}
	root := &types.Node{
		Name:             f.Name(),
		Mode:             f.Mode(),
		Size:             size,
		SizeOnDisk:       sizeOnDisk,
		IsDir:            f.IsDir(),
		LastModification: f.ModTime(),
	}
	if root.IsDir {
		names, err := readDirNames(path)
		if err != nil {
			log.Warningf("%s: %v", path, err)
			return root, nil
		}

		eg, _ := errgroup.WithContext(context.Background())
		rl := &sync.Mutex{}

		for _, name := range names {
			var curr = name
			eg.Go(func() error {
				child, err := Walk(path+"/"+curr, opts)
				if err != nil {
					return err
				}
				child.Parent = root
				rl.Lock()
				defer rl.Unlock()
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
				return nil
			})
		}

		if err := eg.Wait(); err != nil {
			return nil, err
		}

		sort.Slice(root.Children, func(i, j int) bool {
			return strings.Compare(root.Children[i].Name, root.Children[j].Name) == -1
		})
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
	return names, nil
}
