package main

import (
	"math/big"
	"testing"
)

func checkOperation(f lcf, endOrder []int64, t *testing.T) {
	for i := 0; i < len(endOrder); i++ {
		x := big.NewInt(endOrder[i])
		expected := big.NewInt(int64(i))

		actual := f.operate(x)
		if actual.Cmp(expected) != 0 {
			t.Errorf("Expected card %v to be at position %v but got %v", x, expected, actual)
		}
	}
}

func TestLcfOperateDealNewStack(t *testing.T) {
	endOrder := []int64{9, 8, 7, 6, 5, 4, 3, 2, 1, 0}
	deckSize := big.NewInt(int64(len(endOrder)))
	f := instr2Lcf("deal into new stack", deckSize)
	checkOperation(f, endOrder, t)
}

func TestLcfOperateCut3Cards(t *testing.T) {
	endOrder := []int64{3, 4, 5, 6, 7, 8, 9, 0, 1, 2}
	deckSize := big.NewInt(int64(len(endOrder)))
	f := instr2Lcf("cut 3", deckSize)
	checkOperation(f, endOrder, t)
}

func TestLcfOperateCut6Cards(t *testing.T) {
	endOrder := []int64{6, 7, 8, 9, 0, 1, 2, 3, 4, 5}
	deckSize := big.NewInt(int64(len(endOrder)))
	f := instr2Lcf("cut 6", deckSize)
	checkOperation(f, endOrder, t)
}

func TestLcfOperateCutNCardsNegative(t *testing.T) {
	endOrder := []int64{6, 7, 8, 9, 0, 1, 2, 3, 4, 5}
	deckSize := big.NewInt(int64(len(endOrder)))
	f := instr2Lcf("cut -4", deckSize)
	checkOperation(f, endOrder, t)
}

func TestLcfOperateDealWithIncrement3(t *testing.T) {
	endOrder := []int64{0, 7, 4, 1, 8, 5, 2, 9, 6, 3}
	deckSize := big.NewInt(int64(len(endOrder)))
	f := instr2Lcf("deal with increment 3", deckSize)
	checkOperation(f, endOrder, t)
}

func TestLcfOperateDealWithIncrement7(t *testing.T) {
	endOrder := []int64{0, 3, 6, 9, 2, 5, 8, 1, 4, 7}
	deckSize := big.NewInt(int64(len(endOrder)))
	f := instr2Lcf("deal with increment 7", deckSize)
	checkOperation(f, endOrder, t)

}

func TestLcfCompose01(t *testing.T) {
	instr := []string{
		"deal with increment 7",
		"deal into new stack",
		"deal into new stack",
	}
	endOrder := []int64{0, 3, 6, 9, 2, 5, 8, 1, 4, 7}
	deckSize := big.NewInt(int64(len(endOrder)))
	lcfs := convertInstrToLcf(instr, deckSize)
	f := compose(lcfs)
	checkOperation(f, endOrder, t)
}

func TestLcfCompose02(t *testing.T) {
	instr := []string{
		"cut 6",
		"deal with increment 7",
	}
	endOrder := []int64{6, 9, 2, 5, 8, 1, 4, 7, 0, 3}
	deckSize := big.NewInt(int64(len(endOrder)))
	lcfs := convertInstrToLcf(instr, deckSize)
	f := compose(lcfs)
	checkOperation(f, endOrder, t)
}

func TestLcfCompose03(t *testing.T) {
	instr := []string{
		"cut 6",
		"deal with increment 7",
		"deal into new stack",
	}
	endOrder := []int64{3, 0, 7, 4, 1, 8, 5, 2, 9, 6}
	deckSize := big.NewInt(int64(len(endOrder)))
	lcfs := convertInstrToLcf(instr, deckSize)
	f := compose(lcfs)
	checkOperation(f, endOrder, t)
}

func TestLcfCompose04(t *testing.T) {
	instr := []string{
		"deal with increment 7",
		"deal with increment 9",
		"cut -2",
	}
	endOrder := []int64{6, 3, 0, 7, 4, 1, 8, 5, 2, 9}
	deckSize := big.NewInt(int64(len(endOrder)))
	lcfs := convertInstrToLcf(instr, deckSize)
	f := compose(lcfs)
	checkOperation(f, endOrder, t)
}

func TestLcfCompose05(t *testing.T) {
	instr := []string{
		"deal into new stack",
		"cut -2",
		"deal with increment 7",
		"cut 8",
		"cut -4",
		"deal with increment 7",
		"cut 3",
		"deal with increment 9",
		"deal with increment 3",
		"cut -1",
	}
	endOrder := []int64{9, 2, 5, 8, 1, 4, 7, 0, 3, 6}
	deckSize := big.NewInt(int64(len(endOrder)))
	lcfs := convertInstrToLcf(instr, deckSize)
	f := compose(lcfs)
	checkOperation(f, endOrder, t)
}
