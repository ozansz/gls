package size

import (
	"github.com/ozansz/gls/internal"
	"io/fs"
)

func (fsInfo FsInfo) GetSize(fInfo fs.FileInfo) (int64, error) {
	return fInfo.Size(), nil
}

func (fsInfo FsInfo) GetSizeOnDisk(fInfo fs.FileInfo) (int64, error) {
	size := fInfo.Size()
	sizeOnDisk := (size + internal.NSClusterSize - 1) / internal.NSClusterSize * internal.NSClusterSize
	return sizeOnDisk, nil
}
