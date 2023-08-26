package main

import (
	"errors"
	"fmt"
	"github.com/PPTide/gojdk/parse"
	"log"
	"os/exec"
)

func main() {
	//defer profile.Start().Stop()
	cmd := exec.Command("/Library/Java/JavaVirtualMachines/zulu-8.jdk/Contents/Home/bin/javac", `main.java`)
	if err := cmd.Run(); err != nil {
		out, _ := cmd.CombinedOutput()
		log.Fatal(errors.Join(err, errors.New(string(out))))
	}

	res, err := parse.Parse("Main.class")
	//pp.Println(res)
	if err != nil {
		//niceShow(res)
		panic(fmt.Sprintf("Error: %+v", err))
	}

	if res.MajorVersion > 52 {
		panic("This version of the interpreter only recognizes class files until version 52.0")
	}
	//niceShow(res)

	err = execute(res)
	if err != nil {
		panic(err)
	}
}
