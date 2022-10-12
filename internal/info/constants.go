package info

import "fmt"

const (
	Project      = "gls"
	VersionMajor = 1
	VersionMinor = 1
	VersionPatch = 1
)

func ProjectNameWithVersion() string {
	return fmt.Sprintf("%s v%d.%d.%d", Project, VersionMajor, VersionMinor, VersionPatch)
}
