package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/dliberat/advent-of-code-2019/intcode"
)

func part1() {
	data, err := ioutil.ReadFile("input.txt")
	if err != nil {
		panic("Cannot read intcode program.")
	}

	memory := string(data)
	cpu := intcode.MakeComputer(memory, nil, nil)

	reader := bufio.NewReader(os.Stdin)

	for true {
		cpu.Run()

		fmt.Print(cpu.FlushASCIIOutput())

		text, _ := reader.ReadString('\n')
		fmt.Println(text)
		cpu.QueueASCIIInput(text)

	}

	/* Solution:
	Items needed: ornament, astrolabe, shell, sand
	Steps: north, take sand, north, north, take astrolabe, south, south, south
	       west, north, take shell, south, south, west, take ornament, west, south, south
	*/
}

func main() {
	part1()
}
