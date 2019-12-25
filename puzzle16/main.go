package main

import (
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"

	"github.com/dliberat/advent-of-code-2019/matmul"
)

func absoluteLastDigit(num int) int {
	if num < 0 {
		num *= -1
	}
	return num % 10
}

func generateEncodingMatrix(offset int, msglen int) matmul.CSRMatrix {
	M := matmul.CSRMatrix{}

	// do this once for each element in the message
	for i := offset; i < offset+msglen; i++ {

		pattern := make([]int, msglen)

		// starting from column i, i+1 columns are filled with 1s
		// then, skip i+1 columns (because they are zeroes)
		// then, i+1 columns are filled with -1s
		// then, skip i+1 columns
		// repeat until the column number equals the length of the message

		j := i
		for j < len(pattern) {

			for n := 0; n < i+1 && j < len(pattern); n++ {
				pattern[j] = 1
				j++
			}
			j += i + 1
			for n := 0; n < i+1 && j < len(pattern); n++ {
				pattern[j] = -1
				j++
			}
			j += i + 1
		}

		M.AddRow(pattern)

	}

	return M
}

func part1(signal []int) {

	/*
		The ith row of the encoding matrix represents the pattern
		corresponding to the ith entry in the signal.
	*/
	M := generateEncodingMatrix(0, len(signal))

	for i := 0; i < 100; i++ {
		signal = matmul.VectorMultiplyAbsUnits(&signal, &M)
	}

	fmt.Println("[Part1] The first 8 digits after 100 iterations are:", signal[0:8])
}

func part2(signal []int, offset int) {

	// This matrix is on the order of 500,000 x 500,000.
	// Even in CSR format it's too big to hold in memory.
	// M := generateEncodingMatrix(offset, len(signal))

	/* Instead we can take advantage of the fact that the entire message length is
	   about 6,500,000 numbers, but the message offset is on the order of 6,000,000.
		 That's well beyond the halfway mark, and we can note that the lower half
		 of the encoding matrix is a simple triangular matrix. Here is an 8x8 example:

		1  0 -1  0  1  0 -1  0
		0  1  1  0  0 -1 -1  0
		0  0  1  1  1  0  0  0
		0  0  0  1  1  1  1  0
		0  0  0  0  1  1  1  1  <- halfway down
		0  0  0  0  0  1  1  1
		0  0  0  0  0  0  1  1
		0  0  0  0  0  0  0  1
	*/

	for i := 0; i < 100; i++ {

		prev := 0
		newsignal := make([]int, len(signal))

		for j := len(signal) - 1; j >= 0; j-- {
			newsignal[j] = absoluteLastDigit(prev + signal[j])
			prev = newsignal[j]
		}

		signal = newsignal
	}

	fmt.Println("[Part2] Your message is:", signal[:8])
}

func main() {
	data, err := ioutil.ReadFile("input.txt")
	if err != nil {
		panic("Cannot read input data.")
	}
	nums := string(data)
	arr := strings.Split(nums, "")
	signal := make([]int, 0)

	for _, n := range arr {
		val, _ := strconv.Atoi(n)
		signal = append(signal, val)
	}

	part1(signal)

	// for part 2, the signal is the original signal repeated 10000 times
	signal = nil

	for i := 0; i < 10000; i++ {
		for _, n := range arr {
			val, _ := strconv.Atoi(n)
			signal = append(signal, val)
		}
	}

	offsetStr := fmt.Sprintf("%d%d%d%d%d%d%d", signal[0], signal[1], signal[2], signal[3], signal[4], signal[5], signal[6])
	offset, _ := strconv.Atoi(offsetStr)
	fmt.Println("Message offset is:", offset)
	fmt.Println("Original signal is ", len(signal), " items long.")

	trimmedSignal := make([]int, len(signal)-offset)
	for i := 0; i < len(trimmedSignal); i++ {
		trimmedSignal[i] = signal[offset+i]
	}
	part2(trimmedSignal, offset)

}
