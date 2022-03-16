package int

func ByteToKB(b int) int {
	return b / 1024
}

func ByteToMB(b int64) int64 {
	return b / 1024 / 1024
}

func KBToByte(b int) int {
	return b * 1024
}

func MBToByte(b int) int {
	return b * 1024 * 1024
}
