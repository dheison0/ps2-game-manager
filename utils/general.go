package utils

import "bytes"

func BytesToString(data []byte) string {
	n := bytes.IndexByte(data, 0)
	if n < 0 {
		n = len(data) - 1
	}
	return string(data[:n])
}
