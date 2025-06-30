package gomoon

/*
#cgo !lua_static && !no_link_liblua pkg-config: lua
#cgo lua_static CFLAGS: -D_LUA_STATIC -I${SRCDIR}/internal/lua54

#include "./internal/gomoon_utils.h"
*/
import "C"

import "unsafe"

type _CFunction *[0]byte

type LuaInteger C.lua_Integer

type LuaNumber C.lua_Number

type LuaGoFn func(L *State) (nResults int)

const (
	RegistryIndex        = C.LUA_REGISTRYINDEX
	RegistryIndexGlobals = C.LUA_RIDX_GLOBALS
)

type ThreadStatus uint8

const (
	StatusOk            ThreadStatus = C.LUA_OK
	StatusYield         ThreadStatus = C.LUA_YIELD
	StatusErrRuntime    ThreadStatus = C.LUA_ERRRUN
	StatusErrSyntax     ThreadStatus = C.LUA_ERRSYNTAX
	StatusErrMemory     ThreadStatus = C.LUA_ERRMEM
	StatusErrMsgHandler ThreadStatus = C.LUA_ERRERR
)

type LoadMode uint8

const (
	ModeBinaryText LoadMode = iota
	ModeBinary
	ModeText
)

func (m LoadMode) String() string {
	switch m {
	case ModeBinaryText:
		return "bt"
	case ModeBinary:
		return "b"
	case ModeText:
		return "t"
	}
	panic("unreachable")
}

func UpvalueIndex(index int) int {
	return RegistryIndex - index
}

func asCString(str string) *C.char {
	return (*C.char)(unsafe.Pointer(unsafe.StringData(str + "\x00")))
}
