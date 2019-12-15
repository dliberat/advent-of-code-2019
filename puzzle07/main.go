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
	"strings"

	"github.com/dliberat/advent-of-code-2019/intcode"
)

type amplifier struct {
	inputFile  *os.File
	outputFile *os.File
	memory     string
	phase      int
}

func (amp *amplifier) run() int {
	return 0
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func readPipe(pipe *os.File) []int {
	pipe.Seek(0, os.SEEK_SET)
	bytes, err := ioutil.ReadFile(pipe.Name())
	str := string(bytes)
	data := strings.Split(str, "\n")
	values := make([]int, 0)
	for _, datum := range data {
		num, err := strconv.Atoi(datum)
		if err == nil {
			values = append(values, num)
		}
	}
	err = os.Truncate(pipe.Name(), 0)
	check(err)
	return values
}

func amplify(memory string, phase int, input int) int {

	inputs := fmt.Sprintf("%d\n%d\n", phase, input)
	// fmt.Println("Inputs: ", inputs)

	f, _ := ioutil.TempFile("", "")
	f.Seek(0, os.SEEK_SET)
	io.WriteString(f, inputs)
	f.Seek(0, os.SEEK_SET)

	o, _ := ioutil.TempFile("", "outfile.txt")
	o.Seek(0, os.SEEK_SET)

	cpu := intcode.MakeComputer(memory, f, o)
	cpu.Run()

	o.Seek(0, os.SEEK_SET)

	bytes, _ := ioutil.ReadFile(o.Name())
	str := string(bytes[0 : len(bytes)-1])

	output, _ := strconv.Atoi(str)

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

func feedbackAmplification(memory string, phases []int) int {

	var outputSignal int = 0
	var pipeBuffer []int = nil

	pipeAB, err := ioutil.TempFile("", "ab")
	check(err)
	pipeBC, err := ioutil.TempFile("", "bc")
	check(err)
	pipeCD, err := ioutil.TempFile("", "cd")
	check(err)
	pipeDE, err := ioutil.TempFile("", "de")
	check(err)
	pipeEA, err := ioutil.TempFile("", "ea")
	check(err)

	ampA := intcode.MakeComputer(memory, nil, pipeAB)
	ampB := intcode.MakeComputer(memory, nil, pipeBC)
	ampC := intcode.MakeComputer(memory, nil, pipeCD)
	ampD := intcode.MakeComputer(memory, nil, pipeDE)
	ampE := intcode.MakeComputer(memory, nil, pipeEA)

	// Load phase settings
	ampA.QueueInput(phases[0])
	ampB.QueueInput(phases[1])
	ampC.QueueInput(phases[2])
	ampD.QueueInput(phases[3])
	ampE.QueueInput(phases[4])

	// initial signal
	pipeBuffer = append(pipeBuffer, 0)

	for true {

		ampA.QueueInput(pipeBuffer...)
		pipeBuffer = nil
		ampA.Run()
		pipeBuffer = readPipe(pipeAB)

		ampB.QueueInput(pipeBuffer...)
		pipeBuffer = nil
		ampB.Run()
		pipeBuffer = readPipe(pipeBC)

		ampC.QueueInput(pipeBuffer...)
		pipeBuffer = nil
		ampC.Run()
		pipeBuffer = readPipe(pipeCD)

		ampD.QueueInput(pipeBuffer...)
		pipeBuffer = nil
		ampD.Run()
		pipeBuffer = readPipe(pipeDE)

		ampE.QueueInput(pipeBuffer...)
		pipeBuffer = nil
		_, err = ampE.Run()

		pipeBuffer = readPipe(pipeEA)

		// a "no input ready" error indicates that
		// the feedback loop needs to continue
		if err == nil {
			outputSignal = pipeBuffer[0]
			break
		}

	}

	os.Remove(pipeAB.Name())
	os.Remove(pipeBC.Name())
	os.Remove(pipeCD.Name())
	os.Remove(pipeDE.Name())
	os.Remove(pipeEA.Name())

	return outputSignal
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

	var phases []int
	var sequences [][]int
	var memory string = string(data)
	var maxOutput int

	// get all permutations of input phases
	// phases = []int{0, 1, 2, 3, 4}
	// permutations([]int{}, phases, &sequences)

	// // intcode program
	// maxOutput = math.MinInt64

	// for _, p := range sequences {
	// 	fmt.Println("Testing phase sequence ", p)
	// 	output := serialAmplification(memory, p)
	// 	if output > maxOutput {
	// 		maxOutput = output
	// 	}
	// }

	// fmt.Println("[Part 1] Final processed output: ", maxOutput)

	phases = []int{5, 6, 7, 8, 9}
	sequences = nil
	permutations([]int{}, phases, &sequences)

	maxOutput = math.MinInt64
	for _, p := range sequences {
		output := feedbackAmplification(memory, p)
		if output > maxOutput {
			maxOutput = output
		}
	}

	fmt.Println("[Part 2] Final processed output: ", maxOutput)
}
