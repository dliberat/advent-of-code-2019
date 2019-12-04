/*  https://adventofcode.com/2019/day/2  */

package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func runIntcode(memory []int, noun int, verb int) int {

	memory[1] = noun
	memory[2] = verb

	memoryLength := len(memory)
	a := -1
	b := -1
	c := -1

	i := 0
	opcode := memory[i]

	for opcode != 99 {

		a = memory[i+1]
		b = memory[i+2]
		c = memory[i+3]

		if a >= memoryLength || b >= memoryLength || c >= memoryLength {
			return -1
		}

		if opcode == 1 {
			memory[c] = memory[a] + memory[b]
		} else if opcode == 2 {
			memory[c] = memory[a] * memory[b]
		}

		i += 4

		opcode = memory[i]
	}

	return memory[0]
}

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {

		txt := strings.Split(scanner.Text(), ",")
		memory := make([]int, len(txt))

		for i, v := range txt {
			v, err := strconv.Atoi(v)
			if err != nil {
				fmt.Println("ERROR: Cannot parse input.")
			}
			memory[i] = v
		}

		target := 19690720
		var output int
		var verb int
		noun := -1

		for output != target && noun < len(memory) {
			output = 0
			verb = -1
			noun++

			for output != target && verb < len(memory) {

				clone := make([]int, len(memory))
				copy(clone, memory)

				verb++
				output = runIntcode(clone, noun, verb)

			}

		}

		if output != target {
			fmt.Println("No combination of inputs to the supplied program produces the desired output.")
		} else {
			fmt.Println("Noun: ", noun, "; Verb: ", verb, "; Output: ", output)
			fmt.Println("100 x noun + verb = ", (100*noun)+verb)
		}

	}
}
