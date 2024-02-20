package oplCRC32

//#cgo CFLAGS: -g -Wall
//#cgo LDFLAGS: -Wl,--allow-multiple-definition
//#include "crc32.h"
import "C"
import (
	"fmt"
)

func Crc32(data string) string {
	crc := C.crc32(C.CString(data))
	return fmt.Sprintf("%08X", crc)
}
