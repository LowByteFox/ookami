package fs

import (
	"os"

	"github.com/yuin/gopher-lua"
)

func Loader(L *lua.LState) int {
    mod := L.SetFuncs(L.NewTable(), exports)
    L.Push(mod)
    return 1
}

var exports = map[string]lua.LGFunction {
    "readFile": readFile,
}

func readFile(L *lua.LState) int {
    path := L.ToString(1)

    content, err := os.ReadFile(path)

    if err != nil {
        panic(err)
    }

    L.Push(lua.LString(string(content)))

    return 1
}
