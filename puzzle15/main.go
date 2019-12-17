package main

import (
	"fmt"
	"io/ioutil"

	"github.com/dliberat/advent-of-code-2019/intcode"
)

var visitedCoords []coord

type coord struct {
	x      int
	y      int
	hasO2  bool
	isWall bool
}

type state struct {
	cpu   intcode.Computer
	depth int
	pos   coord
}

func spread02(O2x int, O2y int) {

	minute := 0

	for m := len(visitedCoords) - 1; m >= 0; m-- {
		// oxygen tank location
		if visitedCoords[m].x == O2x && visitedCoords[m].y == O2y {
			visitedCoords[m].hasO2 = true
		}
	}

	for len(visitedCoords) > 0 {
		toModify := make([]int, 0)

		for i := len(visitedCoords) - 1; i >= 0; i-- {

			if !visitedCoords[i].hasO2 {
				continue
			}

			for j, other := range visitedCoords {

				// vertical neighbors
				if (other.x == visitedCoords[i].x) && (other.y == visitedCoords[i].y+1 || other.y == visitedCoords[i].y-1) {
					toModify = append(toModify, j)
				}
				// horizontal neighbors
				if (other.y == visitedCoords[i].y) && (other.x == visitedCoords[i].x+1 || other.x == visitedCoords[i].x-1) {
					toModify = append(toModify, j)
				}
			}

		}

		minute++

		for _, i := range toModify {
			visitedCoords[i].hasO2 = true
		}

		cnt := 0
		for i := len(visitedCoords) - 1; i >= 0; i-- {
			if visitedCoords[i].hasO2 {
				cnt++
			}
		}
		fmt.Println("At minute", minute, ",", cnt, "areas have oxygen.")
		printMap()
		if cnt == len(visitedCoords) {
			break
		}
	}

	fmt.Println("[Part 2] It will take", minute, "minutes to fill with O2")
}

func isVisited(x int, y int) bool {
	for _, s := range visitedCoords {
		if s.x == x && s.y == y {
			return true
		}
	}
	return false
}

func getNextPos(prev coord, move int) coord {
	if move == 1 {
		return coord{x: prev.x, y: prev.y - 1, hasO2: false}
	} else if move == 2 {
		return coord{x: prev.x, y: prev.y + 1, hasO2: false}
	} else if move == 3 {
		return coord{x: prev.x - 1, y: prev.y, hasO2: false}
	} else {
		return coord{x: prev.x + 1, y: prev.y, hasO2: false}
	}
}

func bfs(bot intcode.Computer) (int, int) {

	var O2x int
	var O2y int
	var currentState state
	moves := []int{1, 2, 3, 4}
	queue := make([]state, 0)
	queue = append(queue, state{cpu: bot, depth: 0, pos: coord{x: 0, y: 0, hasO2: false}})

	for len(queue) > 0 {

		currentState = queue[0]
		queue = queue[1:]

		for _, move := range moves {

			nextPos := getNextPos(currentState.pos, move)
			if isVisited(nextPos.x, nextPos.y) {
				continue
			} else {

			}

			clone := currentState.cpu.Clone()
			clone.QueueInput(move)
			clone.Run()
			res := clone.FlushOutput()
			if res[0] == 2 {
				// found the oxygen tank!
				fmt.Println("Found oxygen tank in ", currentState.depth+1, " moves at ", nextPos.x, nextPos.y)
				O2x = nextPos.x
				O2y = nextPos.y
				visitedCoords = append(visitedCoords, nextPos)
			} else if res[0] == 0 {
				// hit a wall
				// fmt.Println("Position ", nextPos, "is a wall.")
			} else {
				visitedCoords = append(visitedCoords, nextPos)
				queue = append(queue, state{cpu: clone, depth: currentState.depth + 1, pos: nextPos})
			}

		}
	}
	return O2x, O2y
}

func printMap() {
	var maxX int
	var maxY int
	var minX int
	var minY int
	for _, c := range visitedCoords {
		if c.x > maxX {
			maxX = c.x
		}
		if c.x < minX {
			minX = c.x
		}
		if c.y > maxY {
			maxY = c.y
		}
		if c.y < minY {
			minY = c.y
		}
	}

	line := "--\t|"
	for x := minX; x <= maxX; x++ {
		val := x % 10
		if val < 0 {
			val = -1 * val
		}
		line = fmt.Sprintf("%s%d", line, val)
	}
	fmt.Println(line + "|")

	for y := minY; y <= maxY; y++ {
		line := fmt.Sprintf("%d\t|", y)

		for x := minX; x <= maxX; x++ {
			flag := false
			for _, c := range visitedCoords {
				if c.x == x && c.y == y {
					z := " "
					if c.hasO2 {
						z = "O"
					}
					line = fmt.Sprintf("%s%s", line, z)
					flag = true
					break
				}
			}
			if flag == false {
				line = fmt.Sprintf("%s#", line)
			}

		}
		fmt.Println(line + "|")
	}

	line = "--\t|"
	for x := minX; x <= maxX; x++ {
		val := x % 10
		if val < 0 {
			val = -1 * val
		}
		line = fmt.Sprintf("%s%d", line, val)
	}
	fmt.Println(line + "|")
}

func main() {
	data, err := ioutil.ReadFile("input.txt")
	if err != nil {
		panic("Cannot read intcode program.")
	}

	cpu := intcode.MakeComputer(string(data), nil, nil)
	x, y := bfs(cpu)

	printMap()

	spread02(x, y)
}
