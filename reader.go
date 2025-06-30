package gomoon

/*
#cgo !lua_static && !no_link_liblua pkg-config: lua
#cgo lua_static CFLAGS: -D_LUA_STATIC -I${SRCDIR}/internal/lua54

#include "./internal/gomoon_utils.h"
*/
import "C"

import (
	"runtime/cgo"
	"unsafe"
)

type ReaderFn func(L *State) (partialData []byte)

//export gomoonReaderFnLayer
func gomoonReaderFnLayer(raw *C.lua_State, ud unsafe.Pointer, size *C.size_t) *C.char {
	L := FromPtr(raw)
	h := (*cgo.Handle)(ud)
	fn := h.Value().(ReaderFn)

	data := fn(L)
	if data == nil {
		return nil
	}

	*size = C.size_t(len(data))
	return (*C.char)(asCString(string(data)))
}
