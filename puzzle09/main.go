/*https://adventofcode.com/2019/day/9
The BOOST program will ask for a single input; run it in test mode by
providing it the value 1. It will perform a series of checks on each
opcode, output any opcodes (and the associated parameter modes) that
seem to be functioning incorrectly, and finally output a BOOST keycode.

Once your Intcode computer is fully functional, the BOOST program should
report no malfunctioning opcodes when run in test mode; it should only
output a single value, the BOOST keycode.
*/
package main

import (
	"fmt"
	"io/ioutil"

	"github.com/dliberat/advent-of-code-2019/intcode"
)

func main() {
	data, err := ioutil.ReadFile("input.txt")
	if err != nil {
		panic("Cannot read input file.")
	}

	// intcode program
	memory := string(data)

	cpu := intcode.MakeComputer(memory, nil, nil)
	result := cpu.Run()

	fmt.Println("Final result at position 0: ", result)
}
