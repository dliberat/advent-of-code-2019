package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetCardPosFactoryOrder0(t *testing.T) {
	d := MakeDeck(5)
	pos := d.cardPos(0)
	assert.Equal(t, 0, pos)
}
func TestGetCardPosFactoryOrder4(t *testing.T) {
	d := MakeDeck(7)
	pos := d.cardPos(4)
	assert.Equal(t, 4, pos)
}

func TestFactoryOrderToSlice(t *testing.T) {
	d := MakeDeck(5)
	expected := []int{0, 1, 2, 3, 4}
	assert.Equal(t, expected, d.toSlice(), "The deck should be in factory order.")
}

func TestCardAtPos(t *testing.T) {
	var card int
	d := MakeDeck(1000)

	for i := 0; i < 1000; i++ {
		card = d.cardAtPos(i)
		assert.Equal(t, i, card)
	}
}

func TestCutDeckForward(t *testing.T) {
	d := MakeDeck(5)
	d.cut(1)
	expected := []int{1, 2, 3, 4, 0}
	assert.Equal(t, expected, d.toSlice())
}

func TestCutDeckBackward(t *testing.T) {
	d := MakeDeck(5)
	d.cut(-3)
	expected := []int{2, 3, 4, 0, 1}
	assert.Equal(t, expected, d.toSlice())
}

func TestCutDeckRepeatedly(t *testing.T) {
	d := MakeDeck(6)
	d.cut(2)
	expected := []int{2, 3, 4, 5, 0, 1}
	assert.Equal(t, expected, d.toSlice())

	d.cut(-3)
	expected = []int{5, 0, 1, 2, 3, 4}
	assert.Equal(t, expected, d.toSlice())
}

func TestDealWithIntcrement3(t *testing.T) {
	d := MakeDeck(10)
	d.dealWithIncrement(3)
	expected := []int{0, 7, 4, 1, 8, 5, 2, 9, 6, 3}
	assert.Equal(t, expected, d.toSlice())
}

func TestExample1(t *testing.T) {
	d := MakeDeck(10)
	d.dealWithIncrement(7)
	d.dealIntoNewStack()
	d.dealIntoNewStack()
	expected := []int{0, 3, 6, 9, 2, 5, 8, 1, 4, 7}
	assert.Equal(t, expected, d.toSlice())
}

func TestExample2(t *testing.T) {
	d := MakeDeck(10)
	d.cut(6)
	d.dealWithIncrement(7)
	d.dealIntoNewStack()
	expected := []int{3, 0, 7, 4, 1, 8, 5, 2, 9, 6}
	assert.Equal(t, expected, d.toSlice())
}

func TestExample2_inline(t *testing.T) {
	d := MakeDeck(10)
	d.cutNandDealWithIncM(6, 7)
	d.dealIntoNewStack()
	expected := []int{3, 0, 7, 4, 1, 8, 5, 2, 9, 6}
	assert.Equal(t, expected, d.toSlice())
}

func TestExample3(t *testing.T) {
	d := MakeDeck(10)
	d.dealWithIncrement(7)
	d.dealWithIncrement(9)
	d.cut(-2)
	expected := []int{6, 3, 0, 7, 4, 1, 8, 5, 2, 9}
	assert.Equal(t, expected, d.toSlice())
}

func TestExample4(t *testing.T) {
	d := MakeDeck(10)
	d.dealIntoNewStack()
	d.cut(-2)
	d.dealWithIncrement(7)
	d.cut(8)
	d.cut(-4)
	d.dealWithIncrement(7)
	d.cut(3)
	d.dealWithIncrement(9)
	d.dealWithIncrement(3)
	d.cut(-1)
	expected := []int{9, 2, 5, 8, 1, 4, 7, 0, 3, 6}
	assert.Equal(t, expected, d.toSlice())
}

func TestCumulativeCuts(t *testing.T) {
	d := MakeDeck(10)
	d.cut(2)
	d.cut(4)
	d.cut(-3)

	e := MakeDeck(10)
	e.cut(2 + 4 - 3)
	assert.Equal(t, d.toSlice(), e.toSlice())
}
