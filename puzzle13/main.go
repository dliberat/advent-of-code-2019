package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"strings"

	"github.com/dliberat/advent-of-code-2019/intcode"
)

type position struct {
	x int
	y int
}

func str2int(str string) []int {
	txt := strings.Split(str, "\n")
	retval := make([]int, 0)
	for _, t := range txt {
		n, err := strconv.Atoi(t)
		if err == nil {
			retval = append(retval, n)
		}
	}

	return retval
}

func main() {
	data, err := ioutil.ReadFile("input.txt")
	if err != nil {
		panic("Cannot read input file.")
	}

	// intcode program
	memory := string(data)

	out, err := ioutil.TempFile("", "")
	if err != nil {
		panic("Can't create temp file.")
	}

	cpu := intcode.MakeComputer(memory, nil, out)
	cpu.Run()

	bytes, err := ioutil.ReadFile(out.Name())
	str := string(bytes)
	os.Remove(out.Name())
	outputs := str2int(str)

	screen := make(map[position]int)

	i := 0
	for i < len(outputs) {

		pixel := position{x: outputs[i], y: outputs[i+1]}
		screen[pixel] = outputs[i+2]
		i += 3
	}

	blockCount := 0
	for _, val := range screen {
		if val == 2 {
			blockCount++
		}
	}

	fmt.Println("[Part 1] Number of blocks displayed on screen when game ends: ", blockCount)
}
