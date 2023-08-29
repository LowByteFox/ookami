package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/LowByteFox/ookami/lua/modules/fs"
	"github.com/yuin/gopher-lua"
)

func PrintLn(L *lua.LState) int {
    text := L.Get(1)
    println(text.String())
    return 0
}

func colorify(L *lua.LState) int {
    isBg := L.ToBool(1)
    r := L.ToNumber(2)
    g := L.ToNumber(3)
    b := L.ToNumber(4)

    textCode := 38
    if isBg {
        textCode = 48
    }

    colorCode := fmt.Sprintf("%d;2;%d;%d;%d", textCode, r, g, b)
    escapeCode := fmt.Sprintf("\x1b[%sm", colorCode)

    print(escapeCode)
    return 0
}

func resetColor(L *lua.LState) int {
    print("\x1b[0m")
    return 0
}

func startScript(path string) {
    L := lua.NewState()
    defer L.Close()

    for _, pair := range []struct {
        n string
        f lua.LGFunction
    }{
        {lua.LoadLibName, lua.OpenPackage}, // Must be first
        {lua.BaseLibName, lua.OpenBase},
        {lua.TabLibName, lua.OpenTable},
    } {
        if err := L.CallByParam(lua.P{
            Fn:      L.NewFunction(pair.f),
            NRet:    0,
            Protect: true,
        }, lua.LString(pair.n)); err != nil {
            panic(err)
        }
    }

    L.PreloadModule("ookami:fs", fs.Loader)

    table := L.NewTable()

    for _, e := range os.Environ() {
        if i := strings.Index(e, "="); i >= 0 {
            table.RawSetString(e[:i], lua.LString(e[i+1:]))
        }
    }

    L.SetGlobal("colorify", L.NewFunction(colorify))
    L.SetGlobal("resetColor", L.NewFunction(resetColor))
    L.SetGlobal("env", table)
    if err := L.DoFile(path); err != nil {
        panic(err)
    }
}
