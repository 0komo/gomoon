package gomoon

/*
#cgo !lua_static pkg-config: lua
#cgo lua_static CFLAGS: -D_LUA_STATIC -I${SRCDIR}/internal/lua54

#include "./internal/gomoon_utils.h"
*/
import "C"

import (
	"runtime/cgo"
	"sync"
	"unsafe"
)

const (
	RegistryIndex = C.LUA_REGISTRYINDEX
)

type LuaInteger C.lua_Integer

type LuaNumber C.lua_Number

type LuaGoFunction func(*State) int

// A type that represents the Lua state, this is the only way to interact and access into Lua.
type State struct {
	mu sync.Mutex
	s  *C.lua_State
}

type _CFunction *[0]byte

// Initializes a Lua state, returns `nil` if it failed (such as because of OOM).
func NewState() *State {
	raw := C.luaL_newstate()
	if raw == nil {
		return nil
	}
	return &State{
		s: raw,
	}
}

func FromState(s unsafe.Pointer) *State {
	return &State{
		s: (*C.lua_State)(s),
	}
}

type AllocFn func(ud, ptr unsafe.Pointer, osize, nsize uintptr) unsafe.Pointer

type _AllocData struct {
	fn *cgo.Handle
	ud unsafe.Pointer
}

//export allocFnHandler
func allocFnHandler(ud, ptr unsafe.Pointer, osize, nsize C.size_t) unsafe.Pointer {
	handleData := *(*cgo.Handle)(ud)
	data := handleData.Value().(_AllocData)
	fn := data.fn.Value().(AllocFn)
	return fn(data.ud, ptr, uintptr(osize), uintptr(nsize))
}

// Initializes a Lua state with an allocator function, the function handles every allocations happened inside the state.
// The behavior is same as [NewState], only the difference is it's that this function also accepts an allocation function.
//
// It accepts a function and an opaque pointer to pass into the function.
func NewStateWithAllocFn(fn AllocFn, ud unsafe.Pointer) *State {
	handleFn := cgo.NewHandle(fn)
	handleData := cgo.NewHandle(_AllocData{&handleFn, unsafe.Pointer(ud)})

	raw := C.lua_newstate(_CFunction(C.allocFnHandler), unsafe.Pointer(&handleData))
	if raw == nil {
		return nil
	}
	return &State{
		s: raw,
	}
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
	L.mu.Lock()
	defer L.mu.Unlock()
	return int(C.luaopen_base(L.s))
}

func (L *State) OpenPackagelib() int {
	L.mu.Lock()
	defer L.mu.Unlock()
	return int(C.luaopen_package(L.s))
}

func (L *State) OpenCorolib() int {
	L.mu.Lock()
	defer L.mu.Unlock()
	return int(C.luaopen_coroutine(L.s))
}

func (L *State) OpenTablelib() int {
	L.mu.Lock()
	defer L.mu.Unlock()
	return int(C.luaopen_table(L.s))
}

func (L *State) OpenIolib() int {
	L.mu.Lock()
	defer L.mu.Unlock()
	return int(C.luaopen_io(L.s))
}

func (L *State) OpenOslib() int {
	L.mu.Lock()
	defer L.mu.Unlock()
	return int(C.luaopen_os(L.s))
}

func (L *State) OpenStringlib() int {
	L.mu.Lock()
	defer L.mu.Unlock()
	return int(C.luaopen_string(L.s))
}

func (L *State) OpenMathlib() int {
	L.mu.Lock()
	defer L.mu.Unlock()
	return int(C.luaopen_math(L.s))
}

func (L *State) OpenDebuglib() int {
	L.mu.Lock()
	defer L.mu.Unlock()
	return int(C.luaopen_debug(L.s))
}

func (L *State) GetTop() int {
	return int(C.lua_gettop(L.s))
}

func (L *State) SetTop(index int) {
	L.mu.Lock()
	defer L.mu.Unlock()
	C.lua_settop(L.s, C.int(index))
}

func (L *State) Pop(n int) {
	L.SetTop(-n - 1)
}

func (L *State) PushNil() {
	L.mu.Lock()
	defer L.mu.Unlock()
	C.lua_pushnil(L.s)
}

func (L *State) PushInteger(n LuaInteger) {
	L.mu.Lock()
	defer L.mu.Unlock()
	C.lua_pushinteger(L.s, C.lua_Integer(n))
}

func (L *State) PushLightUserdata(ptr unsafe.Pointer) {
	L.mu.Lock()
	defer L.mu.Unlock()
	C.lua_pushlightuserdata(L.s, unsafe.Pointer(ptr))
}

func (L *State) PushString(str string) {
	L.mu.Lock()
	defer L.mu.Unlock()
	C.lua_pushlstring(L.s, (*C.char)(unsafe.Pointer(unsafe.StringData(str))), C.size_t(len(str)))
}

func (L *State) PushNumber(n LuaNumber) {
	L.mu.Lock()
	defer L.mu.Unlock()
	C.lua_pushnumber(L.s, C.lua_Number(n))
}

func (L *State) PushBool(b bool) {
	L.mu.Unlock()
	defer L.mu.Unlock()

	var n C.int
	if b {
		n = 1
	} else {
		n = 0
	}
	C.lua_pushboolean(L.s, n)
}

func (L *State) PushThread(tL *State) {
	L.mu.Lock()
	defer L.mu.Unlock()
	C.lua_pushthread(tL.s)
}

func (L *State) PushValue(index int) {
	L.mu.Lock()
	defer L.mu.Unlock()
	C.lua_pushvalue(L.s, C.int(index))
}

//export goFuncHandler
func goFuncHandler(raw *C.lua_State) C.int {
	L := FromState(unsafe.Pointer(raw))
	L.RawGetPtr(RegistryIndex, unsafe.Pointer(C.goFuncHandler))
	h := (*cgo.Handle)(L.ToUserdata(-1))
	L.Pop(1)
	fn := h.Value().(LuaGoFunction)
	return C.int(fn(L))
}

func (L *State) PushGoClosure(fn LuaGoFunction, nUpvalues int) {
	// hacky workaround
	h := cgo.NewHandle(fn)
	L.PushLightUserdata(unsafe.Pointer(&h))
	L.RawSetPtr(RegistryIndex, unsafe.Pointer(C.goFuncHandler))
	func() {
		L.mu.Lock()
		defer L.mu.Unlock()
		C.lua_pushcclosure(L.s, _CFunction(unsafe.Pointer(C.goFuncHandler)), C.int(nUpvalues))
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
	return FromState(C.lua_tothread(L.s, C.int(index)))
}

func (L *State) ToUserdata(index int) unsafe.Pointer {
	return C.lua_touserdata(L.s, C.int(index))
}

func (L *State) SetTable(index int) {
	L.mu.Lock()
	defer L.mu.Unlock()
	C.lua_settable(L.s, C.int(index))
}

func (L *State) GetTable(index int) {
	L.mu.Lock()
	defer L.mu.Unlock()
	C.lua_gettable(L.s, C.int(index))
}

func (L *State) RawSet(index int) {
	L.mu.Lock()
	defer L.mu.Unlock()
	C.lua_rawset(L.s, C.int(index))
}

func (L *State) RawSetElem(index int, n LuaInteger) {
	L.mu.Lock()
	defer L.mu.Unlock()
	C.lua_rawseti(L.s, C.int(index), C.lua_Integer(n))
}

func (L *State) RawSetPtr(index int, ptr unsafe.Pointer) {
	L.mu.Lock()
	defer L.mu.Unlock()
	C.lua_rawsetp(L.s, C.int(index), ptr)
}

func (L *State) RawGet(index int) {
	L.mu.Lock()
	defer L.mu.Unlock()
	C.lua_rawget(L.s, C.int(index))
}

func (L *State) RawGetElem(index int, n LuaInteger) {
	L.mu.Lock()
	defer L.mu.Unlock()
	C.lua_rawgeti(L.s, C.int(index), C.lua_Integer(n))
}

func (L *State) RawGetPtr(index int, ptr unsafe.Pointer) {
	L.mu.Lock()
	defer L.mu.Unlock()
	C.lua_rawgetp(L.s, C.int(index), ptr)
}
