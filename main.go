package main

import (
	"github.com/k0kubun/pp"
	"goJDK/parse"
)

func main() {
	res, err := parse.Parse("Square.class")
	pp.Print(res)
	if err != nil {
		panic(err)
	}
}
