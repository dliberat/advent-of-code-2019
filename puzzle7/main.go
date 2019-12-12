/*https://adventofcode.com/2019/day/7

This code appears to have some bugs in it.
With the specific combination of:
memory = 3,23,3,24,1002,24,10,24,1002,23,-1,23,101,5,23,23,1,24,23,23,4,23,99,0,0
phase = 1
input = 4032

The program becomes unable to read the temporary output file at line 50.
This problem does not occur when working in the VS Code debugger.
*/
package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"math"
	"os"
	"strconv"

	"github.com/dliberat/advent-of-code-2019/intcode"
)

type amplifier struct {
	inputFile *os.File
	outputFile *os.File
	memory string
	phase int
}

func (amp *amplifier) run() int {
	
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func amplify(memory string, phase int, input int) int {

	inputs := fmt.Sprintf("%d\n%d\n", phase, input)
	// fmt.Println("Inputs: ", inputs)

	f, _ := ioutil.TempFile("", "")
	f.Seek(0, os.SEEK_SET)
	io.WriteString(f, inputs)
	f.Seek(0, os.SEEK_SET)
	intcode.SetInFile(f)

	o, _ := ioutil.TempFile("", "outfile.txt")
	o.Seek(0, os.SEEK_SET)
	intcode.SetOutFile(o)

	intcode.Run(memory)

	o.Seek(0, os.SEEK_SET)

	bytes, _ := ioutil.ReadFile(o.Name())
	// fmt.Println("For input ", phase, input, ", Got bytes: ", bytes)
	str := string(bytes[0 : len(bytes)-1])

	output, _ := strconv.Atoi(str)

	// fmt.Println("Got output: ", output)

	os.Remove(f.Name())
	os.Remove(o.Name())

	return output
}

func serialAmplification(memory string, phases []int) int {
	var output int = 0

	for _, phase := range phases {
		output = amplify(memory, phase, output)
	}

	return output
}

func permutations(seq []int, remaining []int, allSequences *[][]int) {

	if len(remaining) == 1 {
		seq = append(seq, remaining[0])
		*allSequences = append(*allSequences, seq)
		return
	}

	for i, num := range remaining {
		cpy := make([]int, len(seq))
		copy(cpy, seq)
		cpy = append(cpy, num)

		remcpy := make([]int, 0)
		for j := 0; j < len(remaining); j++ {
			if j != i {
				remcpy = append(remcpy, remaining[j])
			}
		}

		permutations(cpy, remcpy, allSequences)
	}
}

func main() {
	data, err := ioutil.ReadFile("input.txt")
	check(err)

	// get all permutations of input phases
	phases := []int{0, 1, 2, 3, 4}
	var sequences [][]int
	permutations([]int{}, phases, &sequences)

	// intcode program
	memory := string(data)

	var maxOutput int = math.MinInt64

	for _, p := range sequences {
		fmt.Println("Testing phase sequence ", p)
		output := serialAmplification(memory, p)
		if output > maxOutput {
			maxOutput = output
		}
	}

	fmt.Println("[Part 1] Final processed output: ", maxOutput)
}
