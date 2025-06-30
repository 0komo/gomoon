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

type AllocFn func(ptr unsafe.Pointer, osize, nsize uintptr) (newPtr unsafe.Pointer)

// Initializes a Lua state, returns `nil` if it failed (such as because of OOM).
func NewState() *State {
	raw := C.luaL_newstate()
	if raw == nil {
		return nil
	}
	return &State{raw}
}

func FromPtr(s any) *State {
	return &State{s.(*C.lua_State)}
}

//export gomoonAllocFnLayer
func gomoonAllocFnLayer(ud, ptr unsafe.Pointer, osize, nsize C.size_t) unsafe.Pointer {
	h := (*cgo.Handle)(ud)
	fn := h.Value().(AllocFn)
	return fn(ptr, uintptr(osize), uintptr(nsize))
}

// Initializes a Lua state with an allocator function, the function handles every allocations happened inside the state.
// The behavior is same as [NewState], only the difference is it's that this function also accepts an allocation function.
//
// It accepts a function and an opaque pointer to pass into the function.
func NewStateWithAllocFn(fn AllocFn, ud unsafe.Pointer) *State {
	h := cgo.NewHandle(fn)
	raw := C.lua_newstate(_CFunction(C.gomoonAllocFnLayer), unsafe.Pointer(&h))
	if raw == nil {
		return nil
	}
	return &State{raw}
}
