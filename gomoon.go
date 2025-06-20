package gomoon

/*
#cgo !lua_static pkg-config: lua
#cgo lua_static CFLAGS: -D_LUA_STATIC -I${SRCDIR}/internal/lua54

#include "./internal/gomoon_utils.h"
*/
import "C"

import (
	"runtime/cgo"
	"unsafe"
)

const (
	RegistryIndex        = C.LUA_REGISTRYINDEX
	RegistryIndexGlobals = C.LUA_RIDX_GLOBALS
)

type LuaInteger C.lua_Integer

type LuaNumber C.lua_Number

type LuaGoFunction func(*State) int

type AllocFn func(ptr unsafe.Pointer, osize, nsize uintptr) unsafe.Pointer

// A type that represents the Lua state, this is the only way to interact and access into Lua.
type State struct {
	s *C.lua_State
}

type _CFunction *[0]byte

type _AllocData struct {
	fn *cgo.Handle
}

func UpvalueIndex(index int) int {
	return RegistryIndex - index
}

// Initializes a Lua state, returns `nil` if it failed (such as because of OOM).
func NewState() *State {
	raw := C.luaL_newstate()
	if raw == nil {
		return nil
	}
	return &State{raw}
}

func FromState(s unsafe.Pointer) *State {
	return &State{(*C.lua_State)(s)}
}

//export allocFnHandler
func allocFnHandler(ud, ptr unsafe.Pointer, osize, nsize C.size_t) unsafe.Pointer {
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
	raw := C.lua_newstate(_CFunction(C.allocFnHandler), unsafe.Pointer(&h))
	if raw == nil {
		return nil
	}
	return &State{raw}
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

func (L *State) OpenBaselib() int {
	return int(C.luaopen_base(L.s))
}

func (L *State) OpenPackagelib() int {
	return int(C.luaopen_package(L.s))
}

func (L *State) OpenCorolib() int {
	return int(C.luaopen_coroutine(L.s))
}

func (L *State) OpenTablelib() int {
	return int(C.luaopen_table(L.s))
}

func (L *State) OpenIolib() int {
	return int(C.luaopen_io(L.s))
}

func (L *State) OpenOslib() int {
	return int(C.luaopen_os(L.s))
}

func (L *State) OpenStringlib() int {
	return int(C.luaopen_string(L.s))
}

func (L *State) OpenMathlib() int {
	return int(C.luaopen_math(L.s))
}

func (L *State) OpenDebuglib() int {
	return int(C.luaopen_debug(L.s))
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

//export goFuncHandler
func goFuncHandler(raw *C.lua_State) C.int {
	L := FromState(unsafe.Pointer(raw))
	h := (*cgo.Handle)(L.ToUserdata(UpvalueIndex(1)))
	fn := h.Value().(LuaGoFunction)
	L.Pop(1)
	return C.int(fn(L))
}

func (L *State) PushGoClosure(fn LuaGoFunction, nUpvalues int) {
	// hacky workaround
	h := cgo.NewHandle(fn)
	L.PushLightUserdata(unsafe.Pointer(&h))
	func() {
		C.lua_pushcclosure(L.s, _CFunction(unsafe.Pointer(C.goFuncHandler)), C.int(nUpvalues+1))
	}()
}

func (L *State) PushGoFunction(fn LuaGoFunction) {
	L.PushGoClosure(fn, 0)
}

func (L *State) ToBool(index int) bool {
	return C.lua_toboolean(L.s, C.int(index)) != 0
}

func (L *State) ToInteger(index int) (LuaInteger, bool) {
	var isNum C.int
	n := C.lua_tointegerx(L.s, C.int(index), &isNum)
	return LuaInteger(n), isNum != 0
}

func (L *State) ToString(index int) (string, bool) {
	var strLen C.size_t
	str := C.lua_tolstring(L.s, C.int(index), &strLen)
	if str == nil {
		return "", false
	}
	return unsafe.String((*byte)(unsafe.Pointer(str)), uintptr(strLen)), true
}

func (L *State) ToNumber(index int) (LuaNumber, bool) {
	var isNum C.int
	n := C.lua_tonumberx(L.s, C.int(index), &isNum)
	return LuaNumber(n), isNum != 0
}

func (L *State) ToPointer(index int) unsafe.Pointer {
	return C.lua_topointer(L.s, C.int(index))
}

func (L *State) ToThread(index int) *State {
	return FromState(unsafe.Pointer(C.lua_tothread(L.s, C.int(index))))
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
	L.RawGetElem(RegistryIndex, RegistryIndexGlobals)
	L.PushString(name)
	L.PushValue(-3)
	L.RawSet(-3)
	L.Pop(1)
}

func (L *State) DoString(str string) bool {
	n := C.luaL_loadstring(L.s, C.CString(str))
	if n != 0 {
		return false
	}
	n = C.lua_pcallk(L.s, 0, 0, 0, 0, nil)
	return n == 0
}
