package cp

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.sazak.io/gls/log"
)

func TestMain(m *testing.M) {
	m.Run()

	// Repair file content
	err := os.WriteFile("./testdata/file/test_file_1_dst/src.txt", []byte("test file 1 dst src.txt"), 0666)
	log.Error(err)

	// Remove folders
	err = os.RemoveAll("./testdata/folder/test_folder_0_dst")
	log.Error(err)
	err = os.RemoveAll("./testdata/folder/test_folder_2_dst")
	log.Error(err)

	// Repair file contents
	err = os.WriteFile("./testdata/folder/test_folder_1_dst/dst_folder/src.txt", []byte("test folder 1 dst src.txt"), 0666)
	log.Error(err)
	err = os.WriteFile("./testdata/folder/test_folder_3_dst/sub_folder_1/sub_folder_2/sub_folder_3/src.txt", []byte("test folder 3 dst sub folder 1 2 3 src.txt"), 0666)
	log.Error(err)
}

func TestFile(t *testing.T) {
	testCases := []struct {
		name     string
		src      string
		dst      string
		fileName string
		expText  []byte
	}{
		{
			name:     "Copy a file to non-exist-file location",
			src:      "./testdata/file/test_file_0_src/src.txt",
			dst:      "./testdata/file/test_file_0_dst/",
			fileName: "src.txt",
			expText:  []byte("test file 0 src.txt"),
		},
		{
			name:     "Copy a file to exist-file location",
			src:      "./testdata/file/test_file_1_src/src.txt",
			dst:      "./testdata/file/test_file_1_dst/",
			fileName: "src.txt",
			expText:  []byte("test file 1 src.txt"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Call File function.
			err := File(filepath.Join(tc.dst, tc.fileName), tc.src)
			assert.Nil(t, err)

			// Read file content and check.
			b, err := os.ReadFile(filepath.Join(tc.dst, tc.fileName))
			assert.Nil(t, err)
			assert.Equal(t, tc.expText, b)
		})
	}
}

func TestFolder(t *testing.T) {
	testCases := []struct {
		name    string
		src     string
		dst     string
		dstPath string
		content string
	}{
		{
			name:    "Copy src_folder to non-exist-folder location",
			src:     "./testdata/folder/test_folder_0_src/src_folder",
			dst:     "./testdata/folder/test_folder_0_dst/src_folder",
			dstPath: "./testdata/folder/test_folder_0_dst/src_folder/src.txt",
			content: "test folder 0 src_folder/src.txt",
		},
		{
			name:    "Copy src_folder to exist-file location - changes the file content",
			src:     "./testdata/folder/test_folder_1_src/src_folder",
			dst:     "./testdata/folder/test_folder_1_dst/dst_folder",
			dstPath: "./testdata/folder/test_folder_1_dst/dst_folder/src.txt",
			content: "test folder 1 src src.txt",
		},
		{
			name:    "Copy src_folder which contains sub-folders",
			src:     "./testdata/folder/test_folder_2_src/src_folder",
			dst:     "./testdata/folder/test_folder_2_dst/src_folder",
			dstPath: "./testdata/folder/test_folder_2_dst/src_folder/src_sub_folder_1/src_sub_folder_2/src_sub_folder_3/src.txt",
			content: "src folder sub-folder 1 2 3 src.txt",
		},
		{
			name:    "Copy sub_folder_1 which contains sub-folders and sub-files - changes the file content",
			src:     "./testdata/folder/test_folder_3_src/sub_folder_1",
			dst:     "./testdata/folder/test_folder_3_dst/sub_folder_1",
			dstPath: "./testdata/folder/test_folder_3_dst/sub_folder_1/sub_folder_2/sub_folder_3/src.txt",
			content: "test folder 3 sub folder 1 2 3 src.txt",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := Folder(tc.dst, tc.src)
			assert.Nil(t, err)
			ok, act, err := checkFileContent(tc.dstPath, tc.content)
			assert.True(t, ok)
			assert.Equal(t, "", act)
			assert.Nil(t, err)
		})
	}
}

// checkFileContent checks the content of the test files.
func checkFileContent(path string, content string) (bool, string, error) {
	b, err := os.ReadFile(path)
	if err != nil {
		return false, "", err
	}
	if string(b) != content {
		return false, string(b), nil
	}
	return true, "", nil
}
