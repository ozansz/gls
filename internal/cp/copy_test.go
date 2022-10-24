package cp

import (
	"github.com/ozansz/gls/log"
	"github.com/stretchr/testify/assert"
	"os"
	"path/filepath"
	"testing"
)

func TestMain(m *testing.M) {
	m.Run()

	// Repair file content
	err := os.WriteFile("./testdata/test_file_1_dst/src.txt", []byte("test file 1 dst src.txt"), 0666)
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
			name:     "Copy a file to non-exist filename location",
			src:      "./testdata/test_file_0_src/src.txt",
			dst:      "./testdata/test_file_0_dst/",
			fileName: "src.txt",
			expText:  []byte("test file 0 src.txt"),
		},
		{
			name:     "Copy a file to exist filename location",
			src:      "./testdata/test_file_1_src/src.txt",
			dst:      "./testdata/test_file_1_dst/",
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
