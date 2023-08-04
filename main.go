package main

import (
	"github.com/PPTide/gojdk/parse"
	"log"
	"os/exec"
)

func main() {
	cmd := exec.Command("javac", "main.java")
	if err := cmd.Run(); err != nil {
		log.Fatal(err)
	}

	res, err := parse.Parse("Square.class")
	//pp.Println(res)
	if err != nil {
		//niceShow(res)
		panic(err)
	}

	//niceShow(res)

	err = execute(res)
	if err != nil {
		panic(err)
	}
}
