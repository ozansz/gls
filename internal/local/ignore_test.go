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
			desc:      "simple file path",
			str:       "simple/file/path",
			wantRe:    `simple/file/path$`,
			wantIsDir: false,
		},
		{
			desc:      "simple directory path",
			str:       "simple/dir/path/",
			wantRe:    `simple/dir/path$`,
			wantIsDir: true,
		},
		{
			desc:      "simpe file path with dot",
			str:       "simple/file/path/with.dot",
			wantRe:    `simple/file/path/with\.dot$`,
			wantIsDir: false,
		},
		{
			desc:      "simple directory path with dot",
			str:       "simple/dir/path/with/.dot/",
			wantRe:    `simple/dir/path/with/\.dot$`,
			wantIsDir: true,
		},
		{
			desc:      "only file name with wildcard",
			str:       "*.mp4",
			wantRe:    `.*\.mp4$`,
			wantIsDir: false,
		},
		{
			desc:      "directory path with special characters",
			str:       `~/too-secret\ directory/.dir/`,
			wantRe:    `~/too-secret\\ directory/\.dir$`,
			wantIsDir: true,
		},
		{
			desc:      "directory path with special characters and wildcard",
			str:       "~/too-secret/.dir_*/",
			wantRe:    `~/too-secret/\.dir_.*$`,
			wantIsDir: true,
		},
		/// With re: and re:dir: prefixes
		{
			desc:      "prefixed simple file path",
			str:       `re:simple/file/path$`,
			wantRe:    `simple/file/path$`,
			wantIsDir: false,
		},
		{
			desc:      "prefixed simple directory path",
			str:       `re:dir:simple/dir/path$`,
			wantRe:    `simple/dir/path$`,
			wantIsDir: true,
		},
		{
			desc:      "prefixed simpe file path with dot",
			str:       `re:simple/file/path/with\.dot$`,
			wantRe:    `simple/file/path/with\.dot$`,
			wantIsDir: false,
		},
		{
			desc:      "prefixed simple directory path with dot",
			str:       `re:dir:simple/dir/path/with/\.dot$`,
			wantRe:    `simple/dir/path/with/\.dot$`,
			wantIsDir: true,
		},
		{
			desc:      "prefixed only file name with wildcard",
			str:       `re:^.*\.mp4$`,
			wantRe:    `^.*\.mp4$`,
			wantIsDir: false,
		},
		{
			desc:      "prefixed directory path with special characters",
			str:       `re:dir:~/too-secret\\ directory/\.dir$`,
			wantRe:    `~/too-secret\\ directory/\.dir$`,
			wantIsDir: true,
		},
		{
			desc:      "prefixed directory path with special characters and wildcard",
			str:       `re:dir:~/too-secret/\.dir_.*$`,
			wantRe:    `~/too-secret/\.dir_.*$`,
			wantIsDir: true,
		},
	}
	for _, c := range cases {
		c := c
		t.Run(c.desc, func(t *testing.T) {
			gotRe, gotIsDir := convertIgnorePatternToRegex(c.str)
			if gotRe != c.wantRe {
				t.Errorf("regex: wanted: '%s', got: '%s'", c.wantRe, gotRe)
			}
			if gotIsDir != c.wantIsDir {
				t.Errorf("isDir: wanted: %t, got: %t", c.wantIsDir, gotIsDir)
			}
		})
	}
}
