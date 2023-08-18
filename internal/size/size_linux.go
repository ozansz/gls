package size

import (
	"fmt"
	"io/fs"
	"syscall"

	"go.sazak.io/gls/internal"
)

func (fsInfo FsInfo) GetSize(fInfo fs.FileInfo) (int64, error) {
	st, ok := fInfo.Sys().(*syscall.Stat_t)
	if !ok {
		return 0, fmt.Errorf("could not cast %T to syscall.Stat_t", fInfo.Sys())
	}
	return st.Size, nil
}

func (fsInfo FsInfo) GetSizeOnDisk(fInfo fs.FileInfo) (int64, error) {
	st, ok := fInfo.Sys().(*syscall.Stat_t)
	if !ok {
		return 0, fmt.Errorf("could not cast %T to syscall.Stat_t", fInfo.Sys())
	}
	return (st.Blocks * internal.UNIXSizeOfBlock), nil
}
