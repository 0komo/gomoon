#ifndef GOMOON_UTILS_H
#define GOMOON_UTILS_H

#ifdef GOMOON_STATIC
#include "minilua.h"
#else // GOMOON_STATIC
#include <lauxlib.h>
#include <lua.h>
#include <lualib.h>
#endif // GOMOON_STATIC

inline lua_State *cluaL_newstate() { return luaL_newstate(); }

#endif // GOMOON_UTILS_H
