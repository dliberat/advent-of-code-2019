package main

import (
	"fmt"
	"io/ioutil"
	"strings"
)

func hasBit(n int, pos int) bool {
	val := n & (1 << pos)
	return (val > 0)
}

func clearBit(n int, pos int) int {
	mask := ^(1 << pos)
	n &= mask
	return n
}

func txt2state(data string) int {
	n := 0
	data = strings.ReplaceAll(data, "\n", "")

	for i, char := range data {
		if string(char) == "#" {
			n |= (1 << i)
		}
	}

	return n
}

func getNextState(state int) int {
	/*
		  01234
		0 .....
		1 .....
		2 .....
		3 .....
		4 .....
	*/

	width := 5
	n := state // next state

	for y := 0; y < width; y++ {
		for x := 0; x < width; x++ {

			neighbors := 0

			/* left */
			if x > 0 { // leftmost row has no left neighbors
				if hasBit(state, y*width+(x-1)) {
					neighbors++
				}
			}

			/* right */
			if x < 4 { // rightmost row has no right neighbors
				if hasBit(state, y*width+(x+1)) {
					neighbors++
				}
			}

			/* up */
			if y > 0 { // top row has no top neighbors
				if hasBit(state, (y-1)*width+x) {
					neighbors++
				}
			}

			/* down */
			if y < 4 { // bottom row has no neighbors below it
				if hasBit(state, (y+1)*width+x) {
					neighbors++
				}
			}

			pos := y*width + x

			if hasBit(state, pos) && neighbors != 1 {
				/*A bug dies (becoming an empty space) unless
				there is exactly one bug adjacent to it.*/
				n = clearBit(n, pos)

			} else if !hasBit(state, pos) && (neighbors == 1 || neighbors == 2) {

				/*An empty space becomes infested with a bug
				if exactly one or two bugs are adjacent to it.*/
				n |= (1 << pos)
			}

		}
	}
	return n
}

func printState(state int) {

	width := 5
	for i := 0; i < width*width; i++ {
		if hasBit(state, i) {
			fmt.Print("#")
		} else {
			fmt.Print(".")
		}
		if i%width == width-1 {
			fmt.Print("\n")
		}

	}
	fmt.Print("\n")
}

func main() {
	data, err := ioutil.ReadFile("input.txt")
	if err != nil {
		panic("Can't read input file.")
	}

	state := txt2state(string(data))

	memo := make(map[int]bool)

	for true {
		if memo[state] {
			fmt.Println("State", state, "repeated.")
			printState(state)
			break
		} else {
			memo[state] = true
		}

		state = getNextState(state)
	}

}
