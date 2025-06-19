package gomoon

// #cgo !gomoon-static pkg-config: lua
// #cgo gomoon-static CFLAGS: -DGOMOON_STATIC
// #cgo gomoon-static LDFLAGS: -lm ${SRCDIR}/minilua.c
// #cgo gomoon-static && gomoon-apicheck CFLAGS: -DLUA_USE_APICHECK
// #cgo gomoon-static && !gomoon-no-luacompat CFLAGS: -DLUA_COMPAT_5_3
// #cgo gomoon-static && windows CFLAGS: -DLUA_USE_WINDOWS
// #cgo gomoon-static && !windows CFLAGS: -DLUA_USE_POSIX -DLUA_USE_DLOPEN
// #cgo gomoon-static && !windows LDFLAGS: -ldl
/*
#ifdef GOMOON_STATIC
#include "minilua.h"
#else // GOMOON_STATIC
#include <lua.h>
#include <lauxlib.h>
#include <lualib.h>
#endif // GOMOON_STATIC

#include "../internal/gomoon_utils.h"
*/
import "C"

import (
	"fmt"
	"sync"
)

// An opaque type that represents the Lua state. This is the only way to interact and access the Lua state.
//
// Functions that interacts with the stack are documented like this:
//
//	\[-n, +n\]
//
//	 - The first field, `-n` represents how many elements the function pops from the stack.
//	 - The second field, `+n` represents how many elements the function pushes to the stack.
//	 - `n` in first and second field can be represented as follows:
//	   - `(x|y)`, which means the function can push/pop `x` or `y` elements, depending on the situation.
//	   - `?`, which means that we don't know exactly `n` elements does the function pushes/pops.
//
// This format is largely based on the format found in the C API docs, but it has been modified to fit with
// the bindings' implementation.
type State struct {
	mu    sync.Mutex
	state *C.lua_State
}

// An error structure type that is returned by [NewState] when an error has occurred.
type InitError struct {
	Reason string
}

func (e InitError) Error() string {
	return fmt.Sprintf("cannot initialize Lua state: %s", e.Reason)
}

// An integer type, reflects on `lua_Integer` type.
type Integer C.lua_Integer

type Type int

const (
	None          Type = C.LUA_TNONE
	Nil           Type = C.LUA_TNIL
	Number        Type = C.LUA_TNUMBER
	Boolean       Type = C.LUA_TBOOLEAN
	String        Type = C.LUA_TSTRING
	Table         Type = C.LUA_TTABLE
	Function      Type = C.LUA_TFUNCTION
	Userdata      Type = C.LUA_TUSERDATA
	Thread        Type = C.LUA_TTHREAD
	LightUserdata Type = C.LUA_TLIGHTUSERDATA
)

func (t Type) String() string {
	switch t {
	case None:
		return "no value"
	case Nil:
		return "nil"
	case Number:
		return "number"
	case Boolean:
		return "boolean"
	case String:
		return "string"
	case Table:
		return "table"
	case Function:
		return "function"
	case Userdata:
		return "userdata"
	case Thread:
		return "thread"
	case LightUserdata:
		return "userdata"
	}
	return "unknown"
}

func NewState() (State, error) {
	raw := C.luaL_newstate()
	if raw == nil {
		return State{}, InitError{"OOM"}
	}
	return State{
		state: raw,
	}, nil
}

func (L *State) IsClosed() bool {
	return L.state == nil
}

func (L *State) Close() {
	if L.IsClosed() {
		return
	}
	C.lua_close(L.state)
	L.state = nil
}

func (L *State) GetTop() int {
	return int(C.lua_gettop(L.state))
}

func (L *State) SetTop(index int) {
	L.mu.Lock()
	defer L.mu.Unlock()
	C.lua_settop(L.state, C.int(index))
}

func (L *State) Type(index int) Type {
	return Type(C.lua_type(L.state, C.int(index)))
}

func (L *State) Pop(n int) {
	L.SetTop(-n - 1)
}

func (L *State) PushNil() {
	L.mu.Lock()
	defer L.mu.Unlock()
	C.lua_pushnil(L.state)
}

func (L *State) PushInteger(num Integer) {
	L.mu.Lock()
	defer L.mu.Unlock()
	C.lua_pushinteger(L.state, C.lua_Integer(num))
}
