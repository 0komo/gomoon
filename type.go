package gomoon

/*
#cgo !lua_static && !no_link_liblua pkg-config: lua
#cgo lua_static CFLAGS: -D_LUA_STATIC -I${SRCDIR}/internal/lua54

#include "./internal/gomoon_utils.h"
*/
import "C"

const (
	TNone          = C.LUA_TNONE
	TNil           = C.LUA_TNIL
	TNumber        = C.LUA_TNUMBER
	TString        = C.LUA_TSTRING
	TTable         = C.LUA_TTABLE
	TFunction      = C.LUA_TFUNCTION
	TThread        = C.LUA_TTHREAD
	TUserData      = C.LUA_TUSERDATA
	TLightUserData = C.LUA_TLIGHTUSERDATA
)

type LuaType int8

func (t LuaType) String() string {
	switch t {
	case TNone:
		return "no value"
	case TNil:
		return "nil"
	case TNumber:
		return "number"
	case TString:
		return "string"
	case TTable:
		return "table"
	case TFunction:
		return "function"
	case TThread:
		return "thread"
	case TUserData, TLightUserData:
		return "userdata"
	}
	panic("unreachable")
}
