package size

import (
	"io/fs"

	"go.sazak.io/gls/internal"
)

func (fsInfo FsInfo) GetSize(fInfo fs.FileInfo) (int64, error) {
	return fInfo.Size(), nil
}

func (fsInfo FsInfo) GetSizeOnDisk(fInfo fs.FileInfo) (int64, error) {
	size := fInfo.Size()
	sizeOnDisk := (size + internal.NSClusterSize - 1) / internal.NSClusterSize * internal.NSClusterSize
	return sizeOnDisk, nil
}
