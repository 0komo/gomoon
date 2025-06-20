package tests

import "unsafe"

// #include <stdlib.h>
import "C"

func Free(ptr unsafe.Pointer) {
	C.free(ptr)
}

func Realloc(ptr unsafe.Pointer, nsize uintptr) unsafe.Pointer {
	return C.realloc(ptr, C.size_t(nsize))
}
