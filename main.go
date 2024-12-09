package main

import (
	"fmt"

	"github.com/Dev79844/bitcask/internal"
)

func main() {
	b, _ := internal.Open("hello.txt")
	b.Put("1", []byte("hello"))
	b.Put("2", []byte("world"))

	val1, _ := b.Get("1")
	fmt.Println("val:", string(val1))

	val2, _ := b.Get("2")
	fmt.Println("val:", string(val2))
}