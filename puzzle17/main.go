package main

import (
	"fmt"
	"io/ioutil"

	"github.com/dliberat/advent-of-code-2019/intcode"
)

func slice2grid(data []int) [][]int {
	grid := make([][]int, 1)

	i := 0

	grid[i] = make([]int, 0)

	for _, v := range data {
		if v != 10 {
			grid[i] = append(grid[i], v)
		} else {
			i++
			grid = append(grid, make([]int, 0))
		}
	}

	return grid[:len(grid)-2]
}

func countTrues(t bool, b bool, l bool, r bool) int {
	trues := 0
	if t {
		trues++
	}
	if b {
		trues++
	}
	if l {
		trues++
	}
	if r {
		trues++
	}
	return trues
}

func part1() {
	data, err := ioutil.ReadFile("input.txt")
	if err != nil {
		panic("Cannot read intcode program.")
	}

	cpu := intcode.MakeComputer(string(data), nil, nil)
	cpu.Run()
	output := cpu.FlushOutput()

	for _, v := range output {
		fmt.Print(string(v))
	}

	grid := slice2grid(output)

	chcksum := 0

	for row, r := range grid {
		for col := range r {

			// is scaffold (#)
			if grid[row][col] == 35 {

				top := row > 0 && grid[row-1][col] == 35
				bot := row < len(grid)-1 && grid[row+1][col] == 35
				left := col > 0 && grid[row][col-1] == 35
				right := col < len(r)-1 && grid[row][col+1] == 35
				if countTrues(top, bot, left, right) >= 3 {
					// is an intersection
					chcksum += row * col
				}

			}
		}
	}

	fmt.Println("[Part 1] Checksum:", chcksum)

}

func main() {
	part1()
}
