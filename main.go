package main

import (
	"github.com/PPTide/goJDK/validate"
	"github.com/PPTide/gojdk/parse"
	"github.com/k0kubun/pp"
)

func main() {
	res, err := parse.Parse("Square.class")
	pp.Print(res)
	if err != nil {
		panic(err)
	}
	res2, err := validate.Validate(res)
	pp.Print(res2)
	if err != nil {
		panic(err)
	}
}
