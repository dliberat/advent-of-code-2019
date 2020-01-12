package main

import (
	"fmt"
	"io/ioutil"
	"strings"
)

/**hasBit returns true if the bit at position pos in the
binary representation of the number n is set to 1.*/
func hasBit(n int, pos int) bool {
	val := n & (1 << pos)
	return (val > 0)
}

/**clearBit returns the number n minus 2^pos*/
func clearBit(n int, pos int) int {
	mask := ^(1 << pos)
	n &= mask
	return n
}

/**txt2state converts a string representation of a Game of Life
 * grid into an integer representation.
 *
 * A 5x5 grid would be mapped as follows: Starting from the top-left
 * corner of the grid and working rightward along each row,
 * the corresponding bit of a binary number will be set to 1 if
 * that space has a bug, or 0 if the space does not contain a bug.
 *
 * Example: A 5x5 grid with each position index.
 *	 0  1  2  3  4
 *	 5  6  7  8  9
 *	10 11 12 13 14
 *	15 16 17 18 19
 *	20 21 22 23 24
 *
 * If a grid contains a bug in the upper left corner,
 * the 2^0 bit would be set to 1.
 *
 * In the case of a grid that ONLY contains a single bug
 * at index 4, 4, the binary representation would be
 * 10000 00000 00000 00000 00000, or decimal 16777216.
 */
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

/**getNextState calculates the state of the board after a single
 * iteration of the following algorithm:
 *
 * - A bug dies (becoming an empty space) unless there is exactly
 * one bug adjacent to it.
 *
 * - An empty space becomes infested with a bug
 * if exactly one or two bugs are adjacent to it.
 */
func getNextState(state int) int {
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

/**
 * Calculates the index of the bit that is associated
 * with a given position in a 5x5 grid (see txt2state)
 */
func pos2bitIndex(x, y int) int {
	return 5*y + x
}

/**getNextStateRecursive calculates the state of the board after a single
 * iteration of the following algorithm:
 *
 * - A bug dies (becoming an empty space) unless there is exactly
 * one bug adjacent to it.
 *
 * - An empty space becomes infested with a bug
 * if exactly one or two bugs are adjacent to it.
 *
 * The center cell inside the grid represented by `state` contains
 * within it another grid (`inner`). Similarly, `state` itself is
 * also housed inside another grid (`outer`). This recursive structure
 * can be visualized as follows:
 *
 * 		     |     |         |     |
 *		  1  |  2  |    3    |  4  |  5
 *		     |     |         |     |
 *		-----+-----+---------+-----+-----
 *		     |     |         |     |
 *		  6  |  7  |    8    |  9  |  10
 *		     |     |         |     |
 *		-----+-----+---------+-----+-----
 *		     |     |A|B|C|D|E|     |
 *		     |     |-+-+-+-+-|     |
 *		     |     |F|G|H|I|J|     |
 *		     |     |-+-+-+-+-|     |
 *		 11  | 12  |K|L|?|N|O|  14 |  15
 *		     |     |-+-+-+-+-|     |
 *		     |     |P|Q|R|S|T|     |
 *		     |     |-+-+-+-+-|     |
 *		     |     |U|V|W|X|Y|     |
 *		-----+-----+---------+-----+-----
 *		     |     |         |     |
 *		 16  | 17  |    18   |  19 |  20
 *		     |     |         |     |
 *		-----+-----+---------+-----+-----
 *		     |     |         |     |
 *		 21  | 22  |    23   |  24 |  25
 *		     |     |         |     |
 *
 * Some examples of how neighboring cells are calculated in this space:
 * - Cell 8 has five neighbors (5, 3, 9, A, B, C, D, E).
 * - Cell K has four neighbors (F, L, P, 12)
 * - Cell E has four neighbors (D, J, 8, 14)
 * - Cell T has four neighbors (O, S, Y, 14)
 */
func getNextStateRecursive(state int, outer int, inner int) int {

	width := 5
	n := state // next state

	for y := 0; y < width; y++ {
		for x := 0; x < width; x++ {

			if y == 2 && x == 2 {
				// center cell is home to the recursive inner grid
				// so it will never hold a bug
				continue
			}

			neighbors := 0

			/* left */
			if x == 0 {
				// Cells A, F, K, P, U have cell 12 as their left-hand neighbor
				if hasBit(outer, pos2bitIndex(1, 2)) {
					neighbors++
				}

			} else if x == 3 && y == 2 {
				// Cell 14's left-hand neighbors are E, J, O, T, Y
				for i := 0; i < width; i++ {
					if hasBit(inner, pos2bitIndex(4, i)) {
						neighbors++
					}
				}

			} else { // all other squares are calculated normally
				if hasBit(state, y*width+(x-1)) {
					neighbors++
				}
			}

			/* right */
			if x == 4 {
				// cells E, J, O, T, Y all have cell 14 as their right-hand neighbor
				if hasBit(outer, pos2bitIndex(3, 2)) {
					neighbors++
				}

			} else if x == 1 && y == 2 {
				// cell 12 has cells A, F, K, P, U as its right-hand neighbors
				for i := 0; i < width; i++ {
					if hasBit(inner, pos2bitIndex(0, i)) {
						neighbors++
					}
				}

			} else { // all other squares are calculated normally
				if hasBit(state, y*width+(x+1)) {
					neighbors++
				}
			}

			/* up */
			if y == 0 {
				// Cells A, B, C, D, and E all have Cell 8 as their upper neighbor
				if hasBit(outer, pos2bitIndex(2, 1)) {
					neighbors++
				}
			} else if y == 3 && x == 2 {
				// Cell 18 has Cells U, V, W, X, Y as its upper neighbors
				for i := 0; i < width; i++ {
					if hasBit(inner, pos2bitIndex(i, 4)) {
						neighbors++
					}
				}
			} else {
				// all other squares are calculated normally
				if hasBit(state, (y-1)*width+x) {
					neighbors++
				}
			}

			/* down */
			if y == 4 {
				// Cells U, V, W, X, Y all have Cell 18 as their lower neighbor
				if hasBit(outer, pos2bitIndex(2, 3)) {
					neighbors++
				}
			} else if y == 1 && x == 2 {
				// Cell 8 has A, B, C, D, E as its lower neighbors
				for i := 0; i < width; i++ {
					if hasBit(inner, pos2bitIndex(i, 0)) {
						neighbors++
					}
				}
			} else {
				// all other squares are calculated normally
				if hasBit(state, (y+1)*width+x) {
					neighbors++
				}
			}

			pos := pos2bitIndex(x, y)

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
		if i == 12 {
			fmt.Print("?")
		} else if hasBit(state, i) {
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

func countBugs(state int) int {
	count := 0
	for i := 0; i < 25; i++ {
		if hasBit(state, i) {
			count++
		}
	}
	return count
}

func part1(initialState int) {
	state := initialState
	memo := make(map[int]bool)

	for true {
		if memo[state] {
			fmt.Println("[Part 1] The biodiversity rating of the first state to be repeated is:", state)
			break
		} else {
			memo[state] = true
		}

		state = getNextState(state)
	}
}

func part2(initialState int) {

	// number of iterations
	n := 200

	// every two iterations we will reach a new inner grid, and a new outer grid.
	// Therefore the amount of space needed in the array is:
	// = 1 for the original, starting grid
	//   + n/2 outer recursive grids
	//   + n/2 inner recursive grids
	length := n + 1
	var gridsA []int = make([]int, length)
	var gridsB []int = make([]int, length) // store intermediate values here
	gridsA[length/2] = initialState

	a := &gridsA
	b := &gridsB

	var outer int
	var inner int

	for timestep := 0; timestep < n; timestep++ {
		for i := 0; i < length; i++ {
			if i == 0 {
				outer = 0
			} else {
				outer = (*a)[i-1]
			}
			if i == length-1 {
				inner = 0
			} else {
				inner = (*a)[i+1]
			}

			(*b)[i] = getNextStateRecursive((*a)[i], outer, inner)
		}

		// swap pointers
		a, b = b, a
	}

	ttl := 0
	for _, s := range *a {
		if s != 0 {
			ttl += countBugs(s)
		}
	}

	fmt.Println("[Part 2] The total number of bugs after", n, "iterations is", ttl)
}

func main() {
	data, err := ioutil.ReadFile("input.txt")
	if err != nil {
		panic("Can't read input file.")
	}

	state := txt2state(string(data))

	part1(state)

	part2(state)

}
