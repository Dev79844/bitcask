package main

import (
	"github.com/Dev79844/bitcask/internal"
)

func main() {
	b, _ := internal.Open("hello.txt")
	b.Put("1", []byte("hello"))

	b.Get("1")
}