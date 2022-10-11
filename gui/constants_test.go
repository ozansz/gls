package gui

import (
	"github.com/gdamore/tcell/v2"
	"testing"
)

func TestReadGLSRCFile(t *testing.T) {
	testCases := []struct {
		name    string
		path    string
		expVars map[any]any
		expErr  bool
	}{
		{
			name: "Non-exist .glsrc file",
			path: "./testdata/non_exist_glsrc",
			expVars: map[any]any{
				GridTitleColor:       tcell.ColorRed,
				TreeViewTitleColor:   tcell.ColorGreen,
				FileInfoTitleColor:   tcell.ColorOrange,
				DirectoryColor:       tcell.ColorTurquoise,
				BorderColor:          tcell.ColorLightGray,
				FileInfoAttrColor:    tcell.ColorPowderBlue,
				FileInfoValueColor:   tcell.ColorLightSkyBlue,
				SearchFormTitleColor: tcell.ColorLightSkyBlue,
				UnmarkedFileColor:    tcell.ColorWhite,
				MarkedFileColor:      tcell.ColorRed,
				FileInfoTabAttrWidth: 20,
			},
			expErr: true,
		},
		{
			name: "Customize variables",
			path: "./testdata/.glsrc_test1",
			expVars: map[any]any{
				GridTitleColor:       tcell.ColorBlue,
				TreeViewTitleColor:   tcell.ColorYellow,
				FileInfoTitleColor:   tcell.ColorLightGreen,
				DirectoryColor:       tcell.ColorRed,
				BorderColor:          tcell.ColorWhite,
				FileInfoAttrColor:    tcell.ColorOrange,
				FileInfoValueColor:   tcell.ColorPink,
				SearchFormTitleColor: tcell.ColorBrown,
				UnmarkedFileColor:    tcell.ColorDeepPink,
				MarkedFileColor:      tcell.ColorGray,
				FileInfoTabAttrWidth: 30,
			},
		},
		{
			name: "Some empty lines and mixed-case keys",
			path: "./testdata/.glsrc_test2",
			expVars: map[any]any{
				GridTitleColor:       tcell.ColorBlue,
				TreeViewTitleColor:   tcell.ColorPink,
				FileInfoTitleColor:   tcell.ColorOrange,
				DirectoryColor:       tcell.ColorTurquoise,
				BorderColor:          tcell.ColorLightGray,
				FileInfoAttrColor:    tcell.ColorPowderBlue,
				FileInfoValueColor:   tcell.ColorLightSkyBlue,
				SearchFormTitleColor: tcell.ColorLightSkyBlue,
				UnmarkedFileColor:    tcell.ColorWhite,
				MarkedFileColor:      tcell.ColorRed,
				FileInfoTabAttrWidth: 20,
			},
		},
		{
			name: "Unexpected key and values",
			path: "./testdata/.glsrc_test3",
			expVars: map[any]any{
				GridTitleColor:       tcell.ColorRed,
				TreeViewTitleColor:   tcell.ColorGreen,
				FileInfoTitleColor:   tcell.ColorOrange,
				DirectoryColor:       tcell.ColorTurquoise,
				BorderColor:          tcell.ColorLightGray,
				FileInfoAttrColor:    tcell.ColorPowderBlue,
				FileInfoValueColor:   tcell.ColorLightSkyBlue,
				SearchFormTitleColor: tcell.ColorLightSkyBlue,
				UnmarkedFileColor:    tcell.ColorWhite,
				MarkedFileColor:      tcell.ColorRed,
				FileInfoTabAttrWidth: 20,
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := readGLSRCFile(tc.path)
			if tc.expErr {
				if err == nil {
					t.Fatal("error expected, but error is nil")
				}
				t.Log("error:", err)
			} else {
				if err != nil {
					t.Fatal(err)
				}
			}
			for key, val := range tc.expVars {
				if ok, v := checkVal(key, val); !ok {
					t.Fatalf("expected %v, got %v", val, v)
				}
			}
		})
	}
}

// checkVal checks the given val whether it is expected val.
func checkVal(key any, val any) (bool, any) {
	switch key {
	case "GridTitleColor":
		if GridTitleColor != val {
			return false, GridTitleColor
		}
	case "TreeViewTitleColor":
		if TreeViewTitleColor != val {
			return false, TreeViewTitleColor
		}
	case "FileInfoTitleColor":
		if FileInfoTitleColor != val {
			return false, FileInfoTitleColor
		}
	case "DirectoryColor":
		if DirectoryColor != val {
			return false, DirectoryColor
		}
	case "BorderColor":
		if BorderColor != val {
			return false, BorderColor
		}
	case "FileInfoAttrColor":
		if FileInfoAttrColor != val {
			return false, FileInfoAttrColor
		}
	case "FileInfoValueColor":
		if FileInfoValueColor != val {
			return false, FileInfoValueColor
		}
	case "SearchFormTitleColor":
		if SearchFormTitleColor != val {
			return false, SearchFormTitleColor
		}
	case "UnmarkedFileColor":
		if UnmarkedFileColor != val {
			return false, UnmarkedFileColor
		}
	case "MarkedFileColor":
		if MarkedFileColor != val {
			return false, MarkedFileColor
		}
	case "FileInfoTabAttrWidth":
		if FileInfoTabAttrWidth != val {
			return false, FileInfoTabAttrWidth
		}
	}
	return true, nil
}
