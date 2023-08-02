package main

import (
	"github.com/PPTide/gojdk/parse"
)

func main() {
	res, err := parse.Parse("Square.class")
	//pp.Println(res)
	if err != nil {
		//niceShow(res)
		panic(err)
	}

	niceShow(res)
}
