package gui

import (
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestFileExistInDstPath(t *testing.T) {
	testCases := []struct {
		name     string
		path     string
		fileName string
		ok       bool
		expErr   bool
	}{
		{
			name:     "File does not exist in destination path",
			path:     "./testdata/test_files/",
			fileName: "new_file.txt",
			ok:       false,
		},
		{
			name:     "File exist in destination path",
			path:     "./testdata/test_files/",
			fileName: "test0.txt",
			ok:       true,
		},
		{
			name:     "Wrong destination path",
			path:     "./testdata/wrong_path/",
			fileName: "new_file.txt",
			ok:       false,
			expErr:   true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ok, err := fileExistInDstPath(tc.path, tc.fileName)
			assert.Equal(t, tc.ok, ok)
			if tc.expErr {
				assert.NotNil(t, err)
			} else {
				assert.Nil(t, err)
			}
		})
	}
}

func TestDirPerm(t *testing.T) {
	dirPath := "./testdata/test_dir"
	err := os.MkdirAll(dirPath, 0775)
	assert.Nil(t, err)
	defer os.RemoveAll(dirPath)

	fileInfo, err := os.Stat(dirPath)

	testCases := []struct {
		name   string
		path   string
		perm   os.FileMode
		expErr bool
	}{
		{
			name: "Get directory perm successfully",
			path: dirPath,
			perm: fileInfo.Mode(),
		},
		{
			name:   "Wrong path",
			path:   "./testdata/wrong_path",
			perm:   0,
			expErr: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			perm, err := dirPerm(tc.path)
			if tc.expErr {
				assert.NotNil(t, err)
			}
			assert.Equal(t, tc.perm, perm)
		})
	}
}
