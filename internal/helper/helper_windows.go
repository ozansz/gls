package helper

import (
	"github.com/ozansz/gls/internal"
	"io/fs"
)

func (fsInfo FsInfo) GetSize(fInfo fs.FileInfo) (int64, error) {
	return fInfo.Size(), nil
}

func (fsInfo FsInfo) GetSizeOnDisk(fInfo fs.FileInfo) (int64, error) {
	size := fInfo.Size()
	sizeOnDisk := (size + internal.ClusterSize - 1) / internal.ClusterSize * internal.ClusterSize
	return sizeOnDisk, nil
}
