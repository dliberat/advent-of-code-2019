package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"strings"

	"github.com/dliberat/advent-of-code-2019/intcode"
)

const (
	up    = iota
	down  = iota
	left  = iota
	right = iota
)

type coord struct {
	x int
	y int
}

type robot struct {
	x   int
	y   int
	dir int
	cpu intcode.Computer
	out *os.File
}

func (bot *robot) turn(dir int) {
	if dir == 0 {
		switch bot.dir {
		case up:
			bot.dir = left
		case left:
			bot.dir = down
		case down:
			bot.dir = right
		case right:
			bot.dir = up
		}
	} else if dir == 1 {
		switch bot.dir {
		case up:
			bot.dir = right
		case right:
			bot.dir = down
		case down:
			bot.dir = left
		case left:
			bot.dir = up
		}
	}
}

func (bot *robot) move() {
	switch bot.dir {
	case up:
		bot.y--
	case right:
		bot.x++
	case down:
		bot.y++
	case left:
		bot.x--
	}
}

func (bot *robot) paint(panels map[coord]int) (int, error) {
	currentColor := panels[coord{x: bot.x, y: bot.y}]
	bot.cpu.QueueInput(currentColor)
	return bot.cpu.Run()
}

func check(e error) {
	if e != nil {
		panic(e)
	}
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

func drawPanels(panels map[coord]int) {

	var maxX int
	var maxY int
	for panel := range panels {
		if panel.x > maxX {
			maxX = panel.x
		}

		if panel.y > maxY {
			maxY = panel.y
		}
	}

	var line string

	for y := 0; y <= maxY; y++ {
		line = ""

		for x := 0; x <= maxX; x++ {
			for panel, color := range panels {
				if panel.x == x && panel.y == y {
					c := "■"
					if color == 0 {
						c = " "
					}
					line = fmt.Sprintf("%s%s", line, c)
				}
			}
		}
		fmt.Println(line)
	}
	fmt.Println("")
}

func main() {
	data, err := ioutil.ReadFile("input.txt")
	check(err)
	memory := string(data)

	out, err := ioutil.TempFile("", "")
	if err != nil {
		panic("Can't create temp file.")
	}

	bot := robot{x: 0, y: 0, dir: up, out: out}
	bot.cpu = intcode.MakeComputer(memory, nil, out)

	panels := make(map[coord]int)

	// set 0 (black) for part 1
	// set 1 (white) for part 2
	panels[coord{x: 0, y: 0}] = 1

	for true {
		_, needsinput := bot.paint(panels)

		if needsinput == nil {
			break
		} else {
			data := readPipe(out)
			panels[coord{x: bot.x, y: bot.y}] = data[0]
			bot.turn(data[1])
			bot.move()
		}
	}

	drawPanels(panels)
	fmt.Println(len(panels))

	os.Remove(out.Name())

}
