package main

import (
	"embed"
	"fmt"
	"io"
	"io/fs"
	"math/rand"
)

//go:embed assets/*
var art embed.FS

var images []string

func greet() {

    fs.WalkDir(art, "assets", func(path string, d fs.DirEntry, err error) error {
        if d.IsDir() {
            return nil
        }

        f, err := art.Open(path)
        if err != nil {
            return err
        }

        b, err := io.ReadAll(f)
        _ = b
        if err != nil {
            return err
        }

        images = append(images, string(b))

        return nil
    })

    index := rand.Intn(len(images))
    fmt.Print(images[index])
}

func banner() {
    index := rand.Intn(len(images))
    fmt.Print(images[index])
}
