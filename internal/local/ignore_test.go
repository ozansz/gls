package local

import (
	"testing"
)

func TestConvertIgnorePatternToRegex(t *testing.T) {
	cases := []struct {
		desc      string
		str       string
		wantRe    string
		wantIsDir bool
	}{
		{
			str:       "simple/file/path",
			wantRe:    `simple/file/path$`,
			wantIsDir: false,
		},
		{
			str:       "simple/dir/path/",
			wantRe:    `simple/dir/path$`,
			wantIsDir: true,
		},
		{
			str:       "simple/file/path/with.dot",
			wantRe:    `simple/file/path/with\.dot$`,
			wantIsDir: false,
		},
		{
			str:       "simple/dir/path/with/.dot/",
			wantRe:    `simple/dir/path/with/\.dot$`,
			wantIsDir: true,
		},
		{
			str:       "*.mp4",
			wantRe:    `.*\.mp4$`,
			wantIsDir: false,
		},
		{
			str:       "~/too-secret/.dir/",
			wantRe:    `~/too-secret/\.dir$`,
			wantIsDir: true,
		},
		{
			str:       "~/too-secret/.dir_*/",
			wantRe:    `~/too-secret/\.dir_.*$`,
			wantIsDir: true,
		},
	}
	for _, c := range cases {
		c := c
		t.Run(c.desc, func(t *testing.T) {
			gotRe, gotIsDir := convertIgnorePatternToRegex(c.str)
			if gotRe != c.wantRe {
				t.Errorf("regex: wanted: %q, got: %q", c.wantRe, gotRe)
			}
			if gotIsDir != c.wantIsDir {
				t.Errorf("isDir: wanted: %t, got: %t", c.wantIsDir, gotIsDir)
			}
		})
	}
}
