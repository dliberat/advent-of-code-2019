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

func part2() {
	data, err := ioutil.ReadFile("input.txt")
	if err != nil {
		panic("Cannot read intcode program.")
	}

	cpu := intcode.MakeComputer(string(data), nil, nil)

	cpu.QueueASCIIInput("NOT E T\n") // T is 1 if E is a hole
	cpu.QueueASCIIInput("NOT H J\n") // J is 1 if H is a hole
	cpu.QueueASCIIInput("AND T J\n") // J is 1 if both E and H are holes
	cpu.QueueASCIIInput("NOT J J\n") // J is 0 if both E and H are holes

	cpu.QueueASCIIInput("NOT C T\nNOT T T\n") // T is 1 if C is scaffolding
	cpu.QueueASCIIInput("AND B T\n")          // T is 1 if both B and C are scaffolding
	cpu.QueueASCIIInput("NOT T T\n")          // T is 1 if either B or C are holes
	cpu.QueueASCIIInput("AND D T\n")          // T is 1 if B or C is a hole and D is scaffolding

	/* J is 1 if:
	- at least one of E or H is NOT a hole
	- B or C is a hole, and D is scaffolding
	*/
	cpu.QueueASCIIInput("AND T J\n")

	cpu.QueueASCIIInput("NOT A T\n") // if there's a hole at A, jump anyway 'cause otherwise you fall off
	cpu.QueueASCIIInput("OR T J\n")

	cpu.QueueASCIIInput("RUN\n")
	cpu.Run()
	output := cpu.FlushASCIIOutput()

	fmt.Print(output)
}

func main() {

	fmt.Println("[Part 1]")
	part1()

	fmt.Println("\n\n[Part 2]")
	part2()
}
