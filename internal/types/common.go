package types

import "fmt"

type SizeFormatter func(int64) string

func NoFormat(size int64) string {
	return fmt.Sprint(size)
}

func SizeFormatterBytes(size int64) string {
	const (
		KB = 1024
		MB = 1024 * KB
		GB = 1024 * MB
		TB = 1024 * GB
	)
	if size < KB {
		return fmt.Sprintf("%d B", size)
	}
	if size < MB {
		return fmt.Sprintf("%.2f KB", float64(size)/float64(KB))
	}
	if size < GB {
		return fmt.Sprintf("%.2f MB", float64(size)/float64(MB))
	}
	if size < TB {
		return fmt.Sprintf("%.2f GB", float64(size)/float64(GB))
	}
	return fmt.Sprintf("%.2f TB", float64(size)/float64(TB))
}

func SizeFormatterPow10(size int64) string {
	const (
		KiB = 1000
		MiB = 1000 * KiB
		GiB = 1000 * MiB
		TiB = 1000 * GiB
	)
	if size < KiB {
		return fmt.Sprintf("%d b", size)
	}
	if size < MiB {
		return fmt.Sprintf("%.2f KiB", float64(size)/float64(KiB))
	}
	if size < GiB {
		return fmt.Sprintf("%.2f MiB", float64(size)/float64(MiB))
	}
	if size < TiB {
		return fmt.Sprintf("%.2f GiB", float64(size)/float64(GiB))
	}
	return fmt.Sprintf("%.2f TiB", float64(size)/float64(TiB))
}
