/*I could not have completed this challenge without the walkthrough provided by Spheniscine:
https://codeforces.com/blog/entry/72593 */
package main

import (
	"fmt"
	"io/ioutil"
	"math/big"
	"strconv"
	"strings"
)

// lcf represents a linear congruential function in the form f(x) = ax + b mod m
type lcf struct {
	a *big.Int
	b *big.Int
	m *big.Int // modulo
}

// compose two functions f(x) and g(x). First apply f(x), then apply g(x). g(f(x)).
func (f *lcf) compose(g lcf) lcf {
	//(a, b) ; (c, d) = (ac mod m, bc+d mod m)

	if f.m.Cmp(g.m) != 0 {
		panic("cannot compose functions with different modulos")
	}

	newA := big.NewInt(1)
	newA.Mul(f.a, g.a)
	newA.Mod(newA, f.m)

	newB := big.NewInt(1)
	newB.Mul(f.b, g.a)
	newB.Add(newB, g.b)
	newB.Mod(newB, f.m)
	return lcf{a: newA, b: newB, m: f.m}
}

// operate returns the position of the card x after the shuffle represented by this lcf.
func (f *lcf) operate(x *big.Int) *big.Int {
	// f(x) = ax + b mod m
	fx := big.NewInt(-1)
	fx.Mul(f.a, x)
	fx.Add(fx, f.b)
	fx.Mod(fx, f.m)
	return fx
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

// instr2Lcf converts an instruction in text form into its lcf representation.
// deal into new stack: f(x) = -x - 1 mod m
// cut N: f(x) = x - n mod m
// deal with increment N: f(x) = nx mod m
func instr2Lcf(instruction string, deckSize *big.Int) lcf {
	if instruction == "deal into new stack" {
		return lcf{a: big.NewInt(-1), b: big.NewInt(-1), m: deckSize}
	}

	if instruction[:4] == "cut " {
		val, err := strconv.Atoi(instruction[4:])
		if err != nil {
			panic("Can't parse number")
		}
		n := int64(val)
		return lcf{a: big.NewInt(1), b: big.NewInt(-n), m: deckSize}
	}

	if instruction[:20] == "deal with increment " {
		val, err := strconv.Atoi(instruction[20:])
		if err != nil {
			panic("Can't parse number.")
		}
		n := int64(val)

		return lcf{a: big.NewInt(n), b: big.NewInt(0), m: deckSize}
	}
	panic("bad instruction")
}

// convertInstrToLcf converts a slice of instructions in raw text form
// into their lcf representations and returns a slice of lcfs
func convertInstrToLcf(instructions []string, deckSize *big.Int) []lcf {
	lcfs := make([]lcf, len(instructions))
	for i, instr := range instructions {
		lcfs[i] = instr2Lcf(instr, deckSize)
	}
	return lcfs
}

// compose a series of shuffle operations into a single lcf
func compose(lcfs []lcf) lcf {
	f := lcfs[0]
	for i := 1; i < len(lcfs); i++ {
		f = f.compose(lcfs[i])
	}
	return f
}

// powCompose composes a function f(x) into itself k times
func powCompose(f lcf, k int64) lcf {
	g := lcf{a: big.NewInt(1), b: big.NewInt(0), m: f.m}
	for k > 0 {
		if k%2 == 1 {
			g = g.compose(f)
		}
		k /= 2
		f = f.compose(f)
	}
	return g
}

func part1(instr []string) {
	deckSize := big.NewInt(10007)
	lcfs := convertInstrToLcf(instr, deckSize)
	f := compose(lcfs) // f(x) represents a single shuffle of the deck
	res := f.operate(big.NewInt(2019))
	fmt.Println("[PART 1]: Card 2019 is at position", res)
}

func part2(instr []string) {
	m := big.NewInt(119315717514047)    // deck size
	lcfs := convertInstrToLcf(instr, m) // shuffle sequence
	k := int64(101741582076661)         // no. of shuffles
	f := compose(lcfs)                  // f(x) represents a single shuffle of the deck
	f = powCompose(f, k)                // f(x) represents k shuffles of the deck

	/* f(x) tells us where card x ends up after k shuffles. But the problem asks us which
	card ends up in position 2020. Thus we need to invert f(x).
	The inverse of f(x) must be a function F(x) such that x = aF(x) + b.
	That is to say, F(x) is a function of x that performs the reverse operation of f(x).
	After rearranging to isolate F(x):
	        x - b
	F(x) = ------- mod m
	          a
	*/
	x := big.NewInt(2020)
	numerator := big.NewInt(0)
	numerator.Sub(x, f.b)
	denominator := f.a

	/*Division in modular arithmetic is done by first finding the modular multiplicative
	inverse of the denominator, and then multiplying the numerator by that value.
	       	p/q mod m = p・q^-1 mod m
	*/
	denominator.ModInverse(denominator, m) // modular multiplicative inverse of a = q^-1
	numerator.Mul(numerator, denominator)  // p・q^-1
	numerator.Mod(numerator, m)            // p・q^-1 mod m
	fmt.Println("[PART 2]: The card at position", x, "ends up in position", numerator, "after", k, "shuffles.")
}

func main() {
	instr := getInput("input.txt")
	part1(instr)
	part2(instr)
}
