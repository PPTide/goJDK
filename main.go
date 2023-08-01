package main

import (
	"github.com/PPTide/goJDK/parse"
	"github.com/k0kubun/pp"
)

func main() {
	res, err := parse.Parse("Square.class")
	pp.Print(res)
	if err != nil {
		panic(err)
	}
}
