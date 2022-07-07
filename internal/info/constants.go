package info

import "fmt"

const (
	Project = "gls"
	Version = "0.1.0"
)

func ProjectNameWithVersion() string {
	return fmt.Sprintf("%s v%s", Project, Version)
}
