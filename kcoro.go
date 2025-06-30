package gomoon

/*
#cgo !lua_static && !no_link_liblua pkg-config: lua
#cgo lua_static CFLAGS: -D_LUA_STATIC -I${SRCDIR}/internal/lua54

#include "./internal/gomoon_utils.h"
*/
import "C"

import "runtime/cgo"

type KContext C.lua_KContext

type KFunction func(L *State, status ThreadStatus, ctx KContext) (nResults int)

//export gomoonKFuncLayer
func gomoonKFuncLayer(raw *C.lua_State, status C.int, ctx C.lua_KContext) C.int {
	L := FromPtr(raw)
	h := (*cgo.Handle)(L.ToUserdata(UpvalueIndex(1)))
	fn := h.Value().(KFunction)
	return C.int(fn(L, ThreadStatus(status), KContext(ctx)))
}
