package helper

import (
	"github.com/ozansz/gls/internal"
	"io/fs"
	"syscall"
)

func (fsInfo FsInfo) GetSize(fInfo fs.FileInfo) (int64, error) {
	st, ok := fInfo.Sys().(*syscall.Stat_t)
	if !ok {
		return nil, fmt.Errorf("could not cast %T to syscall.Stat_t", fInfo.Sys())
	}
	return st.Size, nil
}

func (fsInfo FsInfo) GetSizeOnDisk(fInfo fs.FileInfo) (int64, error) {
	st, ok := fInfo.Sys().(*syscall.Stat_t)
	if !ok {
		return nil, fmt.Errorf("could not cast %T to syscall.Stat_t", fInfo.Sys())
	}
	return (st.Blocks * internal.SizeOfBlock), nil
}