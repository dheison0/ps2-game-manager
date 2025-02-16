package utils

import "fmt"

var crctab [256]uint32

// init Populates CRC table
func init() {
	var crc, table, count int32
	for table = 0; table < 256; table++ {
		crc = table << 24
		for count = 8; count > 0; count-- {
			if crc < 0 {
				crc = crc << 1
			} else {
				crc = (crc << 1) ^ 0x04C11DB7
			}
		}
		crctab[255-table] = uint32(crc)
	}
}

// Crc32 Turns data into a Crc32 hash as 8-characters hexadecimal
func Crc32(data []byte) string {
	var crc int32
	data = append(data, byte(0))
	for _, b := range data {
		crc = int32(crctab[int32(b)^((crc>>24)&0xFF)] ^ ((uint32(crc) << 8) & 0xFFFFFF00))
		if b == 0x0 {
			break
		}
	}
	return fmt.Sprintf("%08X", uint32(crc))
}
