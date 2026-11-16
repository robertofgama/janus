package cli

import (
	"os"
	"syscall"
	"unsafe"
)

// isTerminal diz se f é um terminal, pedindo ao sistema os atributos de terminal (TCGETS),
// que só um terminal tem. os.ModeCharDevice foi rejeitado porque /dev/null também é um
// dispositivo de caractere. O projeto é só para Linux, e o syscall da biblioteca padrão
// dispensa golang.org/x/term.
func isTerminal(f *os.File) bool {
	var attrs syscall.Termios
	_, _, errno := syscall.Syscall(syscall.SYS_IOCTL, f.Fd(), syscall.TCGETS, uintptr(unsafe.Pointer(&attrs)))
	return errno == 0
}
