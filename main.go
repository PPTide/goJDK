package main

import (
	"errors"
	"github.com/PPTide/gojdk/parse"
	//"github.com/pkg/profile"
	"log"
	"os/exec"
)

func main() {
	//defer profile.Start().Stop()
	cmd := exec.Command("javac", `main.java`)
	if err := cmd.Run(); err != nil {
		out, _ := cmd.CombinedOutput()
		log.Fatal(errors.Join(err, errors.New(string(out))))
	}

	res, err := parse.Parse("Hello.class")
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
