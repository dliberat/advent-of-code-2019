package main

import (
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"
)

func abs(n int) int {
	if n < 0 {
		n *= -1
	}
	return n
}

func floorMod(n, mod int) int {
	r := n % mod
	if (mod > 0 && r < 0) || (mod < 0 && r > 0) {
		return r + mod
	}
	return r
}

/*Deck represents a deck of space cards*/
type Deck struct {
	offset int
	factor int
	size   int
}

func (d *Deck) cut(N int) {
	d.offset = floorMod(d.offset-N, d.size)
}

func (d *Deck) dealWithIncrement(N int) {
	d.factor = floorMod(d.factor*N, d.size)
	d.offset = floorMod(d.offset*N, d.size)
}

func (d *Deck) cutNandDealWithIncM(N, M int) {
	d.factor = floorMod(d.factor*M, d.size)
	d.offset = floorMod((d.offset-N)*M, d.size)
}

func (d *Deck) dealIntoNewStack() {
	d.offset = floorMod((d.size-1)-d.offset, d.size)
	d.factor = floorMod(d.factor*-1, d.size)
}

func (d *Deck) cardPos(N int) int {
	// (A * B) mod C = (A mod C * B mod C) mod C
	// (A + B) mod C = (A mod C + B mod C) mod C
	mx := floorMod(d.factor*N, d.size)
	return floorMod(abs(mx+d.offset), d.size)
}

func (d *Deck) cardAtPos(N int) int {
	// there must be a more efficient way of doing this
	for i := 0; i < d.size; i++ {
		if d.cardPos(i) == N {
			return i
		}
	}
	return -1
}

func (d *Deck) toSlice() []int {
	repr := make([]int, d.size)
	for i := 0; i < d.size; i++ {
		pos := d.cardPos(i)
		repr[pos] = i
	}
	return repr
}

/*MakeDeck creates a new deck of size "size" in factory order.*/
func MakeDeck(size int) Deck {
	return Deck{size: size, offset: 0, factor: 1}
}

func trackCardPosition(instruction string, cardpos int, decksize int) int {
	if instruction == "deal into new stack" {
		return abs(decksize - cardpos - 1)
	}
	if instruction[:4] == "cut " {
		val, err := strconv.Atoi(instruction[4:])
		if err != nil {
			panic("Can't parse number")
		}

		newpos := ((decksize - val) + cardpos) % decksize
		return abs(newpos)
	}
	if instruction[:20] == "deal with increment " {
		val, err := strconv.Atoi(instruction[20:])
		if err != nil {
			panic("Can't parse number.")
		}

		newpos := (cardpos * val) % decksize
		return abs(newpos)
	}

	return -1
}

func part1() {

	instr := getInput("input.txt")
	decksize := 10007
	card := 2019

	for _, i := range instr {
		card = trackCardPosition(i, card, decksize)
	}

	fmt.Println("[Part 1] Card 2019 is at position", card)
}

func runInstructionSet(d *Deck, instr []string) {
	for _, instruction := range instr {
		if instruction == "deal into new stack" {
			d.dealIntoNewStack()
			continue
		}

		if instruction[:4] == "cut " {
			val, err := strconv.Atoi(instruction[4:])
			if err != nil {
				panic("Can't parse number")
			}

			d.cut(val)
			continue
		}

		if instruction[:20] == "deal with increment " {
			val, err := strconv.Atoi(instruction[20:])
			if err != nil {
				panic("Can't parse number.")
			}

			d.dealWithIncrement(val)
			continue
		}
	}
}

func part2() {
	instr := getInput("input.txt")

	d := MakeDeck(119315717514047)

	runInstructionSet(&d, instr)
	fmt.Println(d)

	runInstructionSet(&d, instr)
	fmt.Println(d)

	runInstructionSet(&d, instr)
	fmt.Println(d)

	runInstructionSet(&d, instr)
	fmt.Println(d)
}

func main() {
	part1()
	part2()
}

func getInput(path string) []string {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		panic("Can't read input file.")
	}

	txt := string(data)
	instructions := strings.Split(txt, "\n")
	return instructions
}
