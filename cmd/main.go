package main

import (
	"fmt"

	"github.com/Dev79844/bitcask"
)

func main() {
	b, _ := bitcask.Open(bitcask.WithDir("./data"))
	b.Put("1", []byte("hello"))
	b.Put("2", []byte("world"))
	val1, _ := b.Get("1")
	fmt.Println("val:", string(val1))
	b.Put("1", []byte("world"))

	val2, _ := b.Get("1")
	fmt.Println("val:", string(val2))

	keys := b.List_Keys()
	fmt.Println("keys:",keys)

	b.Merge()

	b.Close()
}