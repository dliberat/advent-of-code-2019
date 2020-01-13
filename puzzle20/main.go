package main

import (
	"io/ioutil"
	"strings"
)

func main() {
	// Part1("input.txt")
	Part2("input.txt")
}

/*GetInput parses the input file into a two-dimensional array
of characters*/
func GetInput(fname string) [][]string {
	data, err := ioutil.ReadFile(fname)
	if err != nil {
		panic("Can't read input file.")
	}

	txt := string(data)

	grid := make([][]string, 0)

	lines := strings.Split(txt, "\n")

	for _, line := range lines {
		row := strings.Split(line, "")
		grid = append(grid, row)
	}

	return grid
}
