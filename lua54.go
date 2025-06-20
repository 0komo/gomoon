//go:build lua54

package gomoon

/*
#cgo !lua_static pkg-config: lua
#cgo lua_static CFLAGS: -D_LUA_STATIC -I${SRCDIR}/internal/lua54
#cgo lua_static LDFLAGS: -lm ${SRCDIR}/internal/lua54/minilua.c
#cgo lua_static && !lua_no_compat CFLAGS: -DLUA_COMPAT_5_3
#cgo lua_static && lua_apicheck CFLAGS: -DLUA_USE_APICHECK
#cgo lua_static && windows CFLAGS: -DLUA_USE_WINDOWS
#cgo lua_static && unix CFLAGS: -DLUA_USE_POSIX
#cgo lua_static && unix LDFLAGS: -ldl

#include "./internal/gomoon_utils.h"
*/
import "C"

const luaVersion = "5.4"

func (L *State) OpenUtf8lib() int {
	return int(C.luaopen_utf8(L.s))
}
