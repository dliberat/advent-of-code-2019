package main

import (
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"
)

func getPattern(n int) []int {
	base := []int{0, 1, 0, -1}

	pattern := make([]int, 0)

	for _, b := range base {
		for i := 0; i < n; i++ {
			pattern = append(pattern, b)
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

func fft(input []int) []int {
	output := make([]int, len(input))

	for i := range output {
		pattern := getPattern(i + 1)
		output[i] = applyPattern(input, pattern)
	}

	return output
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

	i := 0
	for i < 100 {
		signal = fft(signal)
		i++
	}
	fmt.Println("[Part1] The first 8 digits after 100 iterations are:", signal[0:8])

	/* This naive approach is too slow. It will never finish. */
	// for part 2, the signal is the original signal repeated 1000 times
	signal = nil

	for i = 0; i < 10000; i++ {
		for _, n := range arr {
			val, _ := strconv.Atoi(n)
			signal = append(signal, val)
		}
	}

	offsetStr := nums[:7]
	offset, _ := strconv.Atoi(offsetStr)
	fmt.Println("Message offset is:", offset)

	for i = 0; i < 100; i++ {
		signal = fft(signal)
		fmt.Println("iter: ", i)
	}

	fmt.Println("[Part2] Your message is:", signal[offset:offset+9])
}
