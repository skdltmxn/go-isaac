package isaac

import "unsafe"

const (
	uintMax  = 1 << 63
	uintMask = uintMax - 1
)

func unsafeCopy(dst, src unsafe.Pointer, size int) {
	for i := 0; i < size; i++ {
		*(*uint8)(dst) = *(*uint8)(src)
		dst = unsafe.Pointer(uintptr(dst) + unsafe.Sizeof(uint8(0)))
		src = unsafe.Pointer(uintptr(src) + unsafe.Sizeof(uint8(0)))
	}
}
