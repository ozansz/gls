package helper

import (
	"os"
	"runtime"
	"testing"
)

func TestFsHelper(t *testing.T) {
	var diskUsage DiskUsage = &FsInfo{}
	path := ""
	if runtime.GOOS == "windows" {
		path = "..\\"
	} else {
		path = "~/"
	}
	f, err := os.Lstat(path)
	if err != nil {
		t.Fatalf("%s not found or couldn't get infos", path)
	}
	size, _ := diskUsage.GetSize(f)
	if err != nil {
		t.Fatalf("getSize methods gives error during execution")
	}
	sizeOnDisk, err := diskUsage.GetSizeOnDisk(f)
	if err != nil {
		t.Fatalf("getSizeOnDisk methods gives error during execution")
	}
	if size == 0 || sizeOnDisk == 0 {
		t.Fatalf("file size or sizeondisk errors, size: %v , sizeondisk: %v", size, sizeOnDisk)
	}
}
