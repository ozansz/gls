package internal

import (
	"fmt"
	"os/exec"
	"runtime"
	"strconv"
)

type ByteSize int64

const (
	B  ByteSize = 1
	KB ByteSize = 1 << (10 * iota)
	MB
	GB
	TB
	PB
)

func ParseByteSize(s string) (ByteSize, int64, error) {
	if i, err := strconv.ParseInt(s, 10, 64); err == nil {
		return B, i, nil
	}
	if s[len(s)-1] == 'b' || s[len(s)-1] == 'B' {
		if i, err := strconv.ParseInt(s[:len(s)-1], 10, 64); err == nil {
			return B, i, nil
		}
	}
	i, err := strconv.ParseInt(s[:len(s)-2], 10, 64)
	if err != nil {
		return 0, 0, fmt.Errorf("invalid formatting %q: %v", s, err)
	}
	suffix := s[len(s)-2:]
	switch suffix {
	case "KB", "Kb", "kb", "kB":
		return KB, i, nil
	case "MB", "Mb", "mb", "mB":
		return MB, i, nil
	case "GB", "Gb", "gb", "gB":
		return GB, i, nil
	case "TB", "Tb", "tb", "tB":
		return TB, i, nil
	case "PB", "Pb", "pb", "pB":
		return PB, i, nil
	}
	return 0, 0, fmt.Errorf("invalid suffix: %s", suffix)
}

func OpenFile(path string) error {
	switch runtime.GOOS {
	case "windows":
		return exec.Command("cmd", "/C", "start", path).Run()
	case "darwin":
		return exec.Command("open", path).Run()
	default:
		return exec.Command("xdg-open", path).Run()
	}
}

func OpenFileWithProgram(path, program string) error {
	return exec.Command(program, path).Run()
}
