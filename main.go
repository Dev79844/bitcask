package main

import (
	"fmt"

	"github.com/Dev79844/bitcask/internal"
)

func main() {
	b, _ := internal.Open("hello.txt")
	b.Put("1", []byte("hello"))
	val1, _ := b.Get("1")
	fmt.Println("val:", string(val1))
	b.Put("1", []byte("world"))

	val2, _ := b.Get("1")
	fmt.Println("val:", string(val2))

	b.Delete("1")
}