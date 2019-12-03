/*  https://adventofcode.com/2019/day/2  */

package main

import "bufio"
import "fmt"
import "os"
import "strings"
import "strconv"

func run_intcode(program []int) {

	a := -1
	b := -1
	c := -1

	i := 0
	opcode := program[i]

	for opcode != 99 {
		a = program[i+1]
		b = program[i+2]
		c = program[i+3]

		if opcode == 1 {
			program[c] = program[a] + program[b]
		} else if opcode == 2 {
			program[c] = program[a] * program[b]
		}

		i += 4
		opcode = program[i]
	}
	fmt.Println(program)

}

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {

		txt := strings.Split(scanner.Text(), ",")
		nums := make([]int, len(txt))

		for i, v  := range txt {
			v, err := strconv.Atoi(v)
			if err != nil {
				fmt.Println("ERROR: Cannot parse input.")
			}
			nums[i] = v
		}


		// before running the program, replace position 1 with the value 12
		// and replace position 2 with the value 2
		nums[1] = 12
		nums[2] = 2

		run_intcode(nums)
	}
}
