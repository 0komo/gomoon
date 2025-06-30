#ifndef GOMOON_UTILS_H
#define GOMOON_UTILS_H

#ifdef _LUA_STATIC
#include "minilua.h"
#else // _LUA_STATIC
#include <lauxlib.h>
#include <lua.h>
#include <lualib.h>
#endif // _LUA_STATIC

extern void *gomoonAllocFnLayer(void *, void *, size_t, size_t);
extern int gomoonCFnLayer(lua_State *);
extern char *gomoonReaderFnLayer(lua_State *, void *, size_t *);
extern int gomoonKFuncLayer(lua_State *, int, lua_KContext);

#endif // GOMOON_UTILS_H
