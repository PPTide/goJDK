package main

import (
	"github.com/PPTide/gojdk/parse"
	"github.com/k0kubun/pp"
)

func main() {
	res, err := parse.Parse("Square.class")
	pp.Println(res)
	if err != nil {
		panic(err)
	}
	res2, err := parse.Validate(res)
	pp.Println(res2)
	if err != nil {
		panic(err)
	}
}
