package main

import (
	"errors"
	"github.com/PPTide/gojdk/parse"
	"log"
	"os/exec"
)

func main() {
	cmd := exec.Command("javac", `main.java`)
	if err := cmd.Run(); err != nil {
		out, _ := cmd.CombinedOutput()
		log.Fatal(errors.Join(err, errors.New(string(out))))
	}

	res, err := parse.Parse("Main.class")
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
