package main

import (
	"fmt"
	"io/ioutil"

	"github.com/dliberat/advent-of-code-2019/intcode"
)

func part1() {
	data, err := ioutil.ReadFile("input.txt")
	if err != nil {
		panic("Cannot read intcode program.")
	}

	cpu := intcode.MakeComputer(string(data), nil, nil)

	cpu.QueueASCIIInput("NOT C T\n") // if there's a hole at C and scaffolding at D, jump to D
	cpu.QueueASCIIInput("AND D T\n")
	cpu.QueueASCIIInput("OR T J\n")

	cpu.QueueASCIIInput("NOT A T\n") // if there's a hole at A, jump anyway 'cause otherwise you fall off
	cpu.QueueASCIIInput("OR T J\n")

	cpu.QueueASCIIInput("WALK\n")
	cpu.Run()
	output := cpu.FlushASCIIOutput()

	fmt.Print(output)
}

func main() {

	part1()
}
