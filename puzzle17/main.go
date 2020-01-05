package main

import (
	"fmt"
	"io/ioutil"
	"sort"

	"github.com/dliberat/advent-of-code-2019/intcode"
)

type substring struct {
	index int
	s     string
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func abs(a int) int {
	if a < 0 {
		return a * -1
	}
	return a
}

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

/*Returns the longest common prefix of a and b*/
func lcp(a string, b string) string {
	m := min(len(a), len(b))

	i := 0

	for i = 0; i < m; i++ {
		if a[i] != b[i] {
			break
		}
	}
	if i == 0 {
		return ""
	}
	return a[:i]
}

/*Returns the longest repeated substring that does not overlap.
This function is buggy and generally not reliable, but works
sufficiently well for the specific input provided in this problem,
since we do not need a true solution. It's enough to get a
rough idea of what the longest repeating non-overlapping string is,
and then the remainder of the deductive work can be done using
visual inspection.*/
func lrs(s string) string {

	// Create suffix array
	N := len(s)
	suffixes := make([]substring, N)
	for i := 0; i < N; i++ {
		suffixes[i] = substring{s: s[i:N], index: i}
	}

	// After sorting, substrings that have repeated characters will appear next to one another
	sort.Slice(suffixes, func(i, j int) bool { return suffixes[i].s < suffixes[j].s })

	longest := ""
	for i := 0; i < N-1; i++ {
		x := lcp(suffixes[i].s, suffixes[i+1].s)

		// there must be at least as much distance between the
		// candidate suffixes as the length of their common substring
		constraint := abs(suffixes[i].index-suffixes[i+1].index) > len(x)

		if len(x) > len(longest) && constraint {
			longest = x
		}
	}

	return longest
}

/*A helper function used to aid in visual analysis of the problem.
First, the complete scaffold route was mapped out in terms of
Right, Left and Forward motions. This was done visually, and the
result is saved in the text file part2_scaffold.txt.

This function extracts the longest repeated, non-overlapping
substring from that sequence of motions, which provides a
first approximation for how the route can be broken down
into distinct subsequences.

Once the longest repeated, non-overlapping substring is extracted
from the complete sequence, the sequence was easily subdivided
into three subsequences using visual inspection. However, some
reordering was needed in order to ensure that each of the
subsequences conformed to the 20 character limit imposed by the
intcode program.

Although this function reports 8R10R10R4R8R10R12R12R4L12L12 as
the longest repeating, non-overlapping subsequence, the actual longest
sequence used to solve the problem was somewhat shorter: R12R4L12L12
*/
func part2Assist() {
	data, err := ioutil.ReadFile("part2_scaffold.txt")
	if err != nil {
		panic("Cannot read file.")
	}

	longestRepeated := lrs(string(data))
	fmt.Println(longestRepeated)
}

func part2() {
	data, err := ioutil.ReadFile("input_part2.txt")
	if err != nil {
		panic("Cannot read intcode program.")
	}

	cpu := intcode.MakeComputer(string(data), nil, nil)

	cpu.QueueASCIIInput("A,B,A,C,A,B,C,A,B,C\n") // main movement routine
	cpu.QueueASCIIInput("R,8,R,10,R,10\n")       // definition of subroutine A
	cpu.QueueASCIIInput("R,4,R,8,R,10,R,12\n")   // definition of subroutine B
	cpu.QueueASCIIInput("R,12,R,4,L,12,L,12\n")  // definition of subroutine C
	cpu.QueueASCIIInput("n\n")                   // Continuous video feed (y/n)

	cpu.Run()
	output := cpu.FlushOutput()

	for _, v := range output {
		if v < 128 {
			// ASCII range
			fmt.Print(string(v))
		} else {
			fmt.Println("[Part 2] Total space dust collected", v)
		}
	}

}

func main() {
	part1()
	// part2Assist()
	part2()

}
