package main

import (
	"fmt"
	"io/ioutil"

	"github.com/dliberat/advent-of-code-2019/intcode"
)

func part1(nums string, display bool) {

	area := make([][]bool, 50)
	for i := 0; i < 50; i++ {
		area[i] = make([]bool, 50)
	}

	totalAffected := 0

	for y := 0; y < 50; y++ {
		for x := 0; x < 50; x++ {

			cpu := intcode.MakeComputer(nums, nil, nil)
			cpu.QueueInput(x)
			cpu.QueueInput(y)
			cpu.Run()
			droneResponse := cpu.FlushOutput()[0]
			if droneResponse == 1 {
				area[y][x] = true
				totalAffected++
			}
		}
	}

	if display {
		for y := range area {
			line := ""
			for x := range area[y] {
				val := "."
				if area[y][x] {
					val = "#"
				}
				line = fmt.Sprintf("%s%s", line, val)
			}
			fmt.Println(line)
		}
	}

	fmt.Println(totalAffected, " squares are affected by the tractor beam")
}

func part2(nums string) {

	shipSize := 100
	maxAreaSz := 10000

	area := make([][]bool, maxAreaSz)
	for i := 0; i < maxAreaSz; i++ {
		area[i] = make([]bool, maxAreaSz)
	}

	totalAffected := 0

	rightMostEdges := make([]int, 0)

	xStart := 0

	for y := 0; y < maxAreaSz; y++ {

		updatedXStart := false
		rightMostEdges = append(rightMostEdges, 0)

		x := xStart

		for x = xStart; x < maxAreaSz; x++ {
			cpu := intcode.MakeComputer(nums, nil, nil)
			cpu.QueueInput(x, y)
			cpu.Run()
			droneResponse := cpu.FlushOutput()[0]

			if droneResponse == 1 {

				rightMostEdges[y] = x

				if !updatedXStart {
					updatedXStart = true
					xStart = x - 4 // on the next iteration, start partway down the row
					if xStart < 0 {
						xStart = 0
					}

					if len(rightMostEdges) >= shipSize && rightMostEdges[y-shipSize+1]-shipSize+1 == x {
						fmt.Println("Santa's ship can fit at ", x, y-shipSize+1)

						return
					}

				}

				area[y][x] = true
				totalAffected++

			} else {
				if updatedXStart {
					break
				}
			}
		}
	}

	for y := range area {
		line := ""
		for x := range area[y] {
			val := "."
			if area[y][x] {
				val = "#"
			}
			line = fmt.Sprintf("%s%s", line, val)
		}
		fmt.Println(line)
	}
	fmt.Println("totalAffected:", totalAffected)
}

func main() {
	data, err := ioutil.ReadFile("input.txt")
	if err != nil {
		panic("Cannot read input data.")
	}
	nums := string(data)

	part1(nums, true)
	part2(nums)
}
