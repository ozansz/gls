package cp

import (
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"math"
	"os"
	"path/filepath"
)

var errFolderExist = errors.New("folder exist")

// File copies the given src file to dst location. It changes the mode with
// existing file's mode.
func File(dst, src string) error {
	srcFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer srcFile.Close()

	dstFile, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer dstFile.Close()

	n, err := io.Copy(dstFile, srcFile)
	if err != nil {
		return err
	}
	srcStat, err := srcFile.Stat()
	if n != srcStat.Size() {
		return fmt.Errorf("src file couldn't copied to destination. %v byte(s) is missing", math.Abs(float64(n-srcStat.Size())))
	}

	return dstFile.Chmod(srcStat.Mode())
}

// Folder copies the given src folder to dst location recursively.
func Folder(dst, src string) error {
	srcStat, err := os.Stat(src)
	if err != nil {
		return err
	}
	err = os.MkdirAll(dst, srcStat.Mode())
	if os.IsExist(err) {
		return errFolderExist
	}

	files, err := ioutil.ReadDir(src)
	if err != nil {
		return err
	}

	for _, f := range files {
		if f.IsDir() {
			err = Folder(dst, filepath.Join(src, f.Name()))
			if err != nil {
				return err
			}
		} else {
			err = File(dst, filepath.Join(src, f.Name()))
			if err != nil {
				return err
			}
		}
	}
	return nil
}
