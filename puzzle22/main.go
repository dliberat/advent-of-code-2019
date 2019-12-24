package main

import (
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"
)

type deck struct {
	cards []int
}

/*newStack reverses the order of cards in the deck.*/
func (d *deck) newStack() {

	i := 0
	j := len(d.cards) - 1

	for i < j {
		d.cards[i], d.cards[j] = d.cards[j], d.cards[i]
		i++
		j--
	}
}

/*cut cuts the deck at the specified point, shifting the
front n cards to the back of the deck. If n is negative,
the same operation is performed starting from the back
of the deck.*/
func (d *deck) cut(n int) {

	var m int

	if n > 0 {
		m = n
	} else if n < 0 {
		m = len(d.cards) + n // add because n is negative
	} else {
		return
	}

	c := d.cards[:m]
	d.cards = d.cards[m:]
	d.cards = append(d.cards, c...)
}

func (d *deck) dealIncrement(n int) {

	shuffled := make([]int, len(d.cards))
	j := 0

	for i := range d.cards {
		shuffled[j] = d.cards[i]
		j = (j + n) % len(shuffled)
	}

	d.cards = shuffled
}

func (d *deck) findCardPosition(n int) int {
	for i, card := range d.cards {
		if card == n {
			return i
		}
	}
	return -1
}

func (d *deck) display() {
	fmt.Println(d.cards)
}

func makeDeck(size int) deck {
	d := deck{cards: make([]int, size)}
	for i := 0; i < size; i++ {
		d.cards[i] = i
	}
	return d
}

func parseInstruction(instruction string, d *deck) {
	if instruction == "deal into new stack" {
		d.newStack()
		return
	}
	if instruction[:4] == "cut " {
		val, err := strconv.Atoi(instruction[4:])
		if err != nil {
			panic("Can't parse number")
		}
		d.cut(val)
		return
	}
	if instruction[:20] == "deal with increment " {
		val, err := strconv.Atoi(instruction[20:])
		if err != nil {
			panic("Can't parse number.")
		}
		d.dealIncrement(val)
		return
	}

}

func main() {
	d := makeDeck(10007)

	instr := getInput("input.txt")
	for _, i := range instr {
		parseInstruction(i, &d)
	}

	card := d.findCardPosition(2019)
	fmt.Println("[Part 1] Card 2019 is at position", card)
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
