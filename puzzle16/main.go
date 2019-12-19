package main

import (
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"
)

func getPattern(n int, msglen int) []int {
	base := []int{0, 1, 0, -1}

	pattern := make([]int, 0)

	sz := 0 //vector size

	for _, b := range base {
		for i := 0; i < n; i++ {
			pattern = append(pattern, b)
			if sz == msglen {
				break
			}
		}
	}

	return pattern
}

func applyPattern(input []int, pattern []int) int {
	ttl := 0

	i := 1

	for _, val := range input {
		ttl += val * pattern[i]
		i = (i + 1) % len(pattern)
	}

	if ttl < 0 {
		ttl *= -1
	}

	return ttl % 10
}

func applyPatternFast(input []int, n int) int {
	ttl := 0

	sign := 1

	// start on the first element whose pattern number is not a zero
	i := n - 1

	for i < len(input) {
		for j := 0; j < n && i < len(input); j++ {
			ttl += sign * input[i]
			i++
		}
		i += n
		sign *= -1
	}

	if ttl < 0 {
		ttl *= -1
	}

	return ttl % 10

}

func fft(input []int) []int {
	output := make([]int, len(input))

	for i := range output {
		// pattern := getPattern(i+1, len(input))
		output[i] = applyPatternFast(input, i+1)
	}
	return output
}

func runfft(input []int, output *[]int, from int, to int, flag *chan bool) {
	for i := from; i < to; i++ {
		(*output)[i] = applyPatternFast(input, i+1)
	}
	(*flag) <- true
}

func mtFft(input []int) []int {
	output := make([]int, len(input))
	ol := len(output)

	threadcount := 8
	flags := make([]chan bool, threadcount)
	for i := range flags {
		flags[i] = make(chan bool)

		from := i * ol / threadcount
		to := (i + 1) * ol / threadcount

		fmt.Println("Launching process from ", from, "to", to)
		go runfft(input, &output, from, to, &flags[i])
	}

	for i := range flags {
		<-flags[i]
	}

	return output
}

func part1(signal []int) {
	i := 0
	for i < 100 {
		signal = fft(signal)
		i++
	}
	fmt.Println("[Part1] The first 8 digits after 100 iterations are:", signal[0:8])
}

func part2(signal []int) {
	/* This naive approach is too slow. It will never finish. */
	offsetStr := fmt.Sprintf("%d%d%d%d%d%d%d", signal[0], signal[1], signal[2], signal[3], signal[4], signal[5], signal[6])
	offset, _ := strconv.Atoi(offsetStr)
	fmt.Println("Message offset is:", offset)

	for i := 0; i < 100; i++ {
		fmt.Println("iter: ", i)
		signal = mtFft(signal)
	}

	fmt.Println("[Part2] Your message is:", signal[offset:offset+9])
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

	// // for part 2, the signal is the original signal repeated 10000 times
	signal = nil

	for i := 0; i < 10000; i++ {
		for _, n := range arr {
			val, _ := strconv.Atoi(n)
			signal = append(signal, val)
		}
	}

	part2(signal)

}
