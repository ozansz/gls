package internal

type ByteSize int64

const (
	B  ByteSize = 1
	KB ByteSize = 1 << (10 * iota)
	MB
	GB
	TB
	PB
)

const (
	UNIXSizeOfBlock = 512
	NSClusterSize   = 4096
)
