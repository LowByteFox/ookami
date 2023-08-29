package main

import (
    "os"
    "strings"

    "github.com/yuin/gopher-lua"
)

func startScript(path string) {
    L := lua.NewState(lua.Options{SkipOpenLibs: true})
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

    table := L.NewTable()

    for _, e := range os.Environ() {
        if i := strings.Index(e, "="); i >= 0 {
            table.RawSetString(e[:i], lua.LString(e[i+1:]))
        }
    }

    L.SetGlobal("env", table)
    if err := L.DoFile(path); err != nil {
        panic(err)
    }
}
