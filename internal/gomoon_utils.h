#ifndef GOMOON_UTILS_H
#define GOMOON_UTILS_H

#ifdef _LUA_STATIC
#include "minilua.h"
#else // _LUA_STATIC
#include <lauxlib.h>
#include <lua.h>
#include <lualib.h>
#endif // _LUA_STATIC

extern void *allocFnHandler(void *, void *, size_t, size_t);
extern int goFuncHandler(lua_State *);

#endif // GOMOON_UTILS_H
