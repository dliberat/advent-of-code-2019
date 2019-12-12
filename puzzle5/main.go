package main

import (
	"fmt"
	"io/ioutil"

	"github.com/dliberat/advent-of-code-2019/intcode"
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func main() {
	data, err := ioutil.ReadFile("input.txt")
	check(err)

	cpu := intcode.MakeComputer(string(data), nil, nil)
	retval := cpu.Run()
	fmt.Println("Run completed. Output: ", retval)
}
