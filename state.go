package gomoon

/*
#cgo !lua_static && !no_link_liblua pkg-config: lua
#cgo lua_static CFLAGS: -D_LUA_STATIC -I${SRCDIR}/internal/lua54

#include "./internal/gomoon_utils.h"
*/
import "C"

import (
	"errors"
	"runtime/cgo"
	"unsafe"
)

// A type that represents the Lua state, this is the only way to interact and access into Lua.
type State struct {
	s *C.lua_State
}

// Deinitializes the Lua state, marking it as "closed".
func (L *State) Close() {
	C.lua_close(L.s)
	L.s = nil
}

// Checks if the Lua state is closed.
func (L *State) IsClosed() bool {
	return L.s == nil
}

func (L *State) GetTop() int {
	return int(C.lua_gettop(L.s))
}

func (L *State) SetTop(index int) {
	C.lua_settop(L.s, C.int(index))
}

func (L *State) Pop(n int) {
	L.SetTop(-n - 1)
}

func (L *State) PushNil() {
	C.lua_pushnil(L.s)
}

func (L *State) PushInteger(n LuaInteger) {
	C.lua_pushinteger(L.s, C.lua_Integer(n))
}

func (L *State) PushLightUserdata(ptr unsafe.Pointer) {
	C.lua_pushlightuserdata(L.s, ptr)
}

func (L *State) PushString(str string) {
	C.lua_pushlstring(L.s, (*C.char)(unsafe.Pointer(unsafe.StringData(str))), C.size_t(len(str)))
}

func (L *State) PushNumber(n LuaNumber) {
	C.lua_pushnumber(L.s, C.lua_Number(n))
}

func (L *State) PushBool(b bool) {
	var n C.int
	if b {
		n = 1
	} else {
		n = 0
	}
	C.lua_pushboolean(L.s, n)
}

func (L *State) PushThread(tL *State) {
	C.lua_pushthread(tL.s)
}

func (L *State) PushValue(index int) {
	C.lua_pushvalue(L.s, C.int(index))
}

//export gomoonCFnLayer
func gomoonCFnLayer(raw *C.lua_State) C.int {
	L := FromPtr(raw)
	h := (*cgo.Handle)(L.ToUserdata(UpvalueIndex(1)))
	fn := h.Value().(LuaGoFn)
	L.Pop(1)
	return C.int(fn(L))
}

func (L *State) PushGoClosure(fn LuaGoFn, nUpvalues int) {
	// hacky workaround
	h := cgo.NewHandle(fn)
	L.PushLightUserdata(unsafe.Pointer(&h))
	C.lua_pushcclosure(L.s, _CFunction(unsafe.Pointer(C.gomoonCFnLayer)), C.int(nUpvalues+1))
}

func (L *State) PushGoFunction(fn LuaGoFn) {
	L.PushGoClosure(fn, 0)
}

func (L *State) ToBool(index int) bool {
	return C.lua_toboolean(L.s, C.int(index)) != 0
}

func (L *State) ToString(index int) (str string, ok bool) {
	var strLen C.size_t
	s := C.lua_tolstring(L.s, C.int(index), &strLen)
	if s == nil {
		return "", false
	}
	return unsafe.String((*byte)(unsafe.Pointer(s)), uintptr(strLen)), true
}

func (L *State) ToInteger(index int) (n LuaInteger, ok bool) {
	var isNum C.int
	num := C.lua_tointegerx(L.s, C.int(index), &isNum)
	return LuaInteger(num), isNum != 0
}

func (L *State) ToNumber(index int) (n LuaNumber, ok bool) {
	var isNum C.int
	num := C.lua_tonumberx(L.s, C.int(index), &isNum)
	return LuaNumber(num), isNum != 0
}

func (L *State) ToPointer(index int) unsafe.Pointer {
	return C.lua_topointer(L.s, C.int(index))
}

func (L *State) ToThread(index int) *State {
	return FromPtr(C.lua_tothread(L.s, C.int(index)))
}

func (L *State) ToUserdata(index int) unsafe.Pointer {
	return C.lua_touserdata(L.s, C.int(index))
}

func (L *State) SetTable(index int) {
	C.lua_settable(L.s, C.int(index))
}

func (L *State) GetTable(index int) {
	C.lua_gettable(L.s, C.int(index))
}

func (L *State) RawSet(index int) {
	C.lua_rawset(L.s, C.int(index))
}

func (L *State) RawSetElem(index int, n LuaInteger) {
	C.lua_rawseti(L.s, C.int(index), C.lua_Integer(n))
}

func (L *State) RawSetPtr(index int, ptr unsafe.Pointer) {
	C.lua_rawsetp(L.s, C.int(index), ptr)
}

func (L *State) RawGet(index int) {
	C.lua_rawget(L.s, C.int(index))
}

func (L *State) RawGetElem(index int, n LuaInteger) {
	C.lua_rawgeti(L.s, C.int(index), C.lua_Integer(n))
}

func (L *State) RawGetPtr(index int, ptr unsafe.Pointer) {
	C.lua_rawgetp(L.s, C.int(index), ptr)
}

func (L *State) SetGlobal(name string) {
	L.RawGetElem(RegistryIndex, RegistryIndexGlobals) /* table */
	L.PushString(name)                                /* table, string */
	L.PushValue(-3)                                   /* table, string, ??? */
	L.RawSet(-3)                                      /* table */
	L.Pop(1)                                          /*  */
}

func (L *State) GetGlobal(name string) {
	L.RawGetElem(RegistryIndex, RegistryIndexGlobals)
	L.PushString(name)
	L.RawGet(-2)
	L.Pop(1)
}

func (L *State) RawLoad(reader ReaderFn, chunkname string, mode LoadMode) ThreadStatus {
	h := cgo.NewHandle(reader)
	status := ThreadStatus(C.lua_load(L.s, _CFunction(C.gomoonReaderFnLayer), unsafe.Pointer(&h), asCString(chunkname), asCString(string(mode))))
	return status
}

func (L *State) Load(reader ReaderFn, chunkname string, mode LoadMode) (ThreadStatus, error) {
	var err error
	status := L.RawLoad(reader, chunkname, mode)
	if status != StatusOk {
		str, _ := L.ToString(-1)
		err = errors.New(str)
		L.Pop(1)
	}
	return status, err
}

func (L *State) LoadBytes(data []byte, mode LoadMode) (ThreadStatus, error) {
	return L.Load(func(_ *State) []byte {
		return data
	}, string(data), mode)
}

func (L *State) SetUpvalue(fnIndex, n int) string {
	return C.GoString(C.lua_setupvalue(L.s, C.int(fnIndex), C.int(n)))
}

func (L *State) GetUpvalue(fnIndex, n int) string {
	return C.GoString(C.lua_getupvalue(L.s, C.int(fnIndex), C.int(n)))
}

func (L *State) UnsafeCallK(nArgs, nResults int, ctx KContext, fn KFunction) ThreadStatus {
	var layerFn _CFunction = nil
	if fn != nil {
		layerFn = _CFunction(C.gomoonKFuncLayer)
	}
	return ThreadStatus(C.lua_callk(L.s, C.int(nArgs), C.int(nResults), C.lua_KContext(ctx), layerFn))
}

func (L *State) UnsafeCall(nArgs, nResults int, ctx KContext, fn KFunction) ThreadStatus {
}
