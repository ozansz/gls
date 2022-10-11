package helper

import (
	"io/fs"
)

type DiskUsage interface {
	GetSize(f fs.FileInfo) (int64, error)
	GetSizeOnDisk(f fs.FileInfo) (int64, error)
}

type FsInfo struct {
}
