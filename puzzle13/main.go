package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/dliberat/advent-of-code-2019/intcode"
)

type position struct {
	x int
	y int
}

func str2int(str string) []int {
	txt := strings.Split(str, "\n")
	retval := make([]int, 0)
	for _, t := range txt {
		n, err := strconv.Atoi(t)
		if err == nil {
			retval = append(retval, n)
		}
	}

	return retval
}

func readPipe(pipe *os.File) []int {
	pipe.Seek(0, os.SEEK_SET)
	bytes, err := ioutil.ReadFile(pipe.Name())
	str := string(bytes)
	data := strings.Split(str, "\n")
	values := make([]int, 0)
	for _, datum := range data {
		num, err := strconv.Atoi(datum)
		if err == nil {
			values = append(values, num)
		}
	}
	err = os.Truncate(pipe.Name(), 0)
	if err != nil {
		panic("Can't erase file.")
	}
	return values
}

func part1() {
	data, err := ioutil.ReadFile("input.txt")
	if err != nil {
		panic("Cannot read input file.")
	}

	// intcode program
	memory := string(data)

	out, err := ioutil.TempFile("", "")
	if err != nil {
		panic("Can't create temp file.")
	}

	cpu := intcode.MakeComputer(memory, nil, out)
	cpu.Run()

	bytes, err := ioutil.ReadFile(out.Name())
	str := string(bytes)
	os.Remove(out.Name())
	outputs := str2int(str)

	screen := make(map[position]int)

	i := 0
	for i < len(outputs) {

		pixel := position{x: outputs[i], y: outputs[i+1]}
		screen[pixel] = outputs[i+2]
		i += 3
	}

	blockCount := 0
	for _, val := range screen {
		if val == 2 {
			blockCount++
		}
	}

	fmt.Println("[Part 1] Number of blocks displayed on screen when game ends: ", blockCount)
}

type pixel struct {
	x int
	y int
	v int
}

func printValue(v int) string {
	switch v {
	case 0: // empty
		return " "
	case 1: // wall
		return "■"
	case 2: // block
		return "="
	case 3: // paddle
		return "-"
	case 4: // ball
		return "○"
	}
	return " "
}

func printScreen(prevScreen []pixel, data []int, ballX *int, paddleX *int) []pixel {
	fmt.Println("\n\n\n\n\n\n\n\n ")

	if len(prevScreen) == 0 {
		var i int = 0
		for i < len(data) {
			p := pixel{x: data[i], y: data[i+1], v: data[i+2]}

			prevScreen = append(prevScreen, p)
			i += 3
		}

	} else {

		var i int = 0
		for i < len(data) {
			p := pixel{x: data[i], y: data[i+1], v: data[i+2]}

			for j := range prevScreen {
				if prevScreen[j].x == p.x && prevScreen[j].y == p.y {
					prevScreen[j].v = p.v
				}
			}
			i += 3
		}
	}

	var maxX int
	var maxY int
	for _, pixel := range prevScreen {
		if pixel.x > maxX {
			maxX = pixel.x
		}

		if pixel.y > maxY {
			maxY = pixel.y
		}
	}

	for _, pixel := range prevScreen {
		if pixel.x == -1 && pixel.y == 0 {
			fmt.Println("     ====== CURRENT SCORE: ", pixel.v, " =======")
			break
		}
	}
	fmt.Println("")

	var line string

	for y := 0; y <= maxY; y++ {
		line = ""

		for x := 0; x <= maxX; x++ {
			for _, pixel := range prevScreen {
				if pixel.x == x && pixel.y == y {
					line = fmt.Sprintf("%s%s", line, printValue(pixel.v))
					if pixel.v == 3 {
						*paddleX = x
					}
					if pixel.v == 4 {
						*ballX = x
					}
				}
			}
		}
		fmt.Println(line)
	}
	fmt.Println("")
	return prevScreen
}

func part2() {
	data, err := ioutil.ReadFile("freeplayinput.txt")
	if err != nil {
		panic("Cannot read intcode program.")
	}

	memory := string(data)
	var ballX int = 0
	var paddleX int = 0
	var prevScreen = make([]pixel, 0)

	out, err := ioutil.TempFile("", "")
	if err != nil {
		panic("Can't create temp file.")
	}

	cpu := intcode.MakeComputer(memory, nil, out)

	for true {
		_, endgame := cpu.Run()

		screen := readPipe(out)
		prevScreen = printScreen(prevScreen, screen, &ballX, &paddleX)

		if ballX > paddleX {
			cpu.QueueInput(1)
		} else if ballX < paddleX {
			cpu.QueueInput(-1)
		} else {
			cpu.QueueInput(0)
		}

		/* To have the CPU play automatically */
		time.Sleep(50 * time.Millisecond)

		/* To play manually */
		// fmt.Print("Input: ")
		// for true {
		// 	_, err := fmt.Fscanf(os.Stdin, "%d", &value)
		// 	if err == nil {

		// 		switch value {
		// 		case 6:
		// 			cpu.QueueInput(1)
		// 		case 4:
		// 			cpu.QueueInput(-1)
		// 		default:
		// 			cpu.QueueInput(0)
		// 		}
		// 		break
		// 	}
		// }

		if endgame == nil {
			return
		}

	}
}

func main() {
	/* Usage: $ ./puzzle13 --part [1|2] */
	partPtr := flag.Int("part", 2, "Which part of the puzzle to solve (1 or 2)")

	flag.Parse()

	if *partPtr == 1 {
		part1()
	} else if *partPtr == 2 {
		part2()
	} else {
		fmt.Println("Invalid option.")
	}
}
