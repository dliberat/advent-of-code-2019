// Package intcode provides utilities for running intcode programs
package intcode

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

const (
	modePosition  = 0
	modeImmediate = 1
)

type instruction interface {
	execute(memory []int, pos *int)
}

type opMultiply struct {
	paramModes []int
	params     []int
}

func (o opMultiply) execute(memory []int, pos *int) {
	var a int
	var b int
	var prod int

	switch o.paramModes[0] {
	case modeImmediate:
		a = o.params[0]
	case modePosition:
		a = memory[o.params[0]]
	default:
		// this should really be an error
		a = memory[o.params[0]]
	}

	switch o.paramModes[1] {
	case modeImmediate:
		b = o.params[1]
	case modePosition:
		b = memory[o.params[1]]
	default:
		b = memory[o.params[1]]
	}

	prod = a * b

	memory[o.params[2]] = prod
	*pos = *pos + 4
}

type opAdd struct {
	paramModes []int
	params     []int
}

func (o opAdd) execute(memory []int, pos *int) {
	var a int
	var b int
	var sum int

	switch o.paramModes[0] {
	case modeImmediate:
		a = o.params[0]
	case modePosition:
		a = memory[o.params[0]]
	default:
		// this should really be an error
		a = memory[o.params[0]]
	}

	switch o.paramModes[1] {
	case modeImmediate:
		b = o.params[1]
	case modePosition:
		b = memory[o.params[1]]
	default:
		b = memory[o.params[1]]
	}

	sum = a + b

	memory[o.params[2]] = sum
	*pos = *pos + 4
}

type opInput struct {
	paramModes []int
	params     []int
}

func (o opInput) execute(memory []int, pos *int) {

	reader := bufio.NewReader(os.Stdin)
	var value int = 0
	var err error

	for true {
		fmt.Print("Enter your input: ")
		text, _ := reader.ReadString('\n')
		value, err = strconv.Atoi(strings.Trim(text, "\n"))
		if err != nil {
			fmt.Println("Bad input. Please provide an integer.", text)
		} else {
			break
		}
	}

	memory[o.params[0]] = value
	*pos = *pos + 2
}

type opOutput struct {
	paramModes []int
	params     []int
}

func (o opOutput) execute(memory []int, pos *int) {
	switch o.paramModes[0] {
	case 0:
		fmt.Println(memory[o.params[0]])
	case 1:
		fmt.Println(o.params[0])
	}

	*pos = *pos + 2
}

type opJumpIfTrue struct {
	paramModes []int
	params     []int
}

func (o opJumpIfTrue) execute(memory []int, pos *int) {
	var a int
	var b int

	switch o.paramModes[0] {
	case modeImmediate:
		a = o.params[0]
	case modePosition:
		a = memory[o.params[0]]
	default:
		a = memory[o.params[0]]
	}

	switch o.paramModes[1] {
	case modeImmediate:
		b = o.params[1]
	case modePosition:
		b = memory[o.params[1]]
	default:
		b = memory[o.params[1]]
	}

	// if the instruction modifies the instruction pointer,
	// that value is used and the instruction pointer is not
	// automatically increased.
	// Presumably this means that if the instruction pointer is
	// NOT modified by an operation (even though that operation could,
	// in theory perform a modification), then the value is incremented
	// as it would be normally
	if a != 0 {
		*pos = b
	} else {
		*pos = *pos + 3
	}
}

type opJumpIfFalse struct {
	paramModes []int
	params     []int
}

func (o opJumpIfFalse) execute(memory []int, pos *int) {
	var a int
	var b int

	switch o.paramModes[0] {
	case modeImmediate:
		a = o.params[0]
	case modePosition:
		a = memory[o.params[0]]
	default:
		a = memory[o.params[0]]
	}

	switch o.paramModes[1] {
	case modeImmediate:
		b = o.params[1]
	case modePosition:
		b = memory[o.params[1]]
	default:
		b = memory[o.params[1]]
	}

	if a == 0 {
		*pos = b
	} else {
		*pos = *pos + 3
	}
}

type opLessThan struct {
	paramModes []int
	params     []int
}

func (o opLessThan) execute(memory []int, pos *int) {
	var a int
	var b int

	switch o.paramModes[0] {
	case modeImmediate:
		a = o.params[0]
	case modePosition:
		a = memory[o.params[0]]
	default:
		a = memory[o.params[0]]
	}

	switch o.paramModes[1] {
	case modeImmediate:
		b = o.params[1]
	case modePosition:
		b = memory[o.params[1]]
	default:
		b = memory[o.params[1]]
	}

	if a < b {
		memory[o.params[2]] = 1
	} else {
		memory[o.params[2]] = 0
	}
	*pos = *pos + 4
}

type opEquals struct {
	paramModes []int
	params     []int
}

func (o opEquals) execute(memory []int, pos *int) {
	var a int
	var b int

	switch o.paramModes[0] {
	case modeImmediate:
		a = o.params[0]
	case modePosition:
		a = memory[o.params[0]]
	default:
		a = memory[o.params[0]]
	}

	switch o.paramModes[1] {
	case modeImmediate:
		b = o.params[1]
	case modePosition:
		b = memory[o.params[1]]
	default:
		b = memory[o.params[1]]
	}

	if a == b {
		memory[o.params[2]] = 1
	} else {
		memory[o.params[2]] = 0
	}
	*pos = *pos + 4
}

type opTerminate struct {
}

func (opTerminate) execute(memory []int, pos *int) {
	*pos = len(memory)
}

/*GetArgumentCount returns the number of arguments
that should be provided for a given opcode*/
func GetArgumentCount(opcode int) int {
	switch opcode {
	case 1:
		return 3
	case 2:
		return 3
	case 3:
		return 1
	case 4:
		return 1
	case 5:
		return 2
	case 6:
		return 2
	case 7:
		return 3
	case 8:
		return 3
	case 99:
		return 0
	}
	return -1
}

/*ParseOpCode splits an opcode into its component
opcode and parameter mode codes. */
func ParseOpCode(opcode string) (int, []int) {
	// the final two digits are the opcode.
	// However, if there's no leading zero and all
	// parameter mode values are zeroes, the opcode
	// could be represented as a single digit
	modesLen := 0
	if len(opcode)-2 > 0 {
		modesLen = len(opcode) - 2
	}

	operation, _ := strconv.Atoi(opcode[modesLen:])

	modesLen = GetArgumentCount(operation)
	paramModes := make([]int, modesLen)
	var j int = 0

	for i := len(opcode) - 3; i >= 0; i-- {
		paramModes[j], _ = strconv.Atoi(opcode[i : i+1])
		j++
	}

	return operation, paramModes
}

func makeInstruction(opcode int, paramModes []int, params []int) instruction {
	switch opcode {
	case 99:
		return opTerminate{}
	case 1:
		return opAdd{paramModes: paramModes, params: params}
	case 2:
		return opMultiply{paramModes: paramModes, params: params}
	case 3:
		return opInput{paramModes: paramModes, params: params}
	case 4:
		return opOutput{paramModes: paramModes, params: params}
	case 5:
		return opJumpIfTrue{paramModes: paramModes, params: params}
	case 6:
		return opJumpIfFalse{paramModes: paramModes, params: params}
	case 7:
		return opLessThan{paramModes: paramModes, params: params}
	case 8:
		return opEquals{paramModes: paramModes, params: params}
	}

	return nil
}

func parseProgram(memory []int) []instruction {
	var currentCode int
	var argv int

	instructions := make([]instruction, len(memory))
	i := 0 // index of next element to to add to instructions array
	m := 0 // index of next element to process in memory array

	for m < len(memory) {
		currentCode = memory[m]
		opcode, paramModes := ParseOpCode(strconv.Itoa(currentCode))
		argv = GetArgumentCount(opcode)

		instructions[i] = makeInstruction(opcode, paramModes, memory[m+1:m+1+argv])
		i++
		m += argv + 1
	}

	return instructions[0:i]
}

/*Run executes the program provided and returns the
result at the 0th index of memory when the program is complete.*/
func Run(memory string) int {
	txt := strings.Split(memory, ",")
	nums := make([]int, len(txt))
	for i, t := range txt {
		nums[i], _ = strconv.Atoi(t)
	}

	var m int = 0 // index into the head of the current instruction

	for m < len(nums) {
		currentCode := nums[m]
		opcode, paramModes := ParseOpCode(strconv.Itoa(currentCode))

		if opcode == 99 {
			break
		}

		argv := GetArgumentCount(opcode)
		instr := makeInstruction(opcode, paramModes, nums[m+1:m+1+argv])
		instr.execute(nums, &m)
	}

	return nums[0]
}

/*CountInstructions counts the number
of instructions in a program*/
func CountInstructions(memory string) int {
	txt := strings.Split(memory, ",")
	nums := make([]int, len(txt))

	for i, t := range txt {
		nums[i], _ = strconv.Atoi(t)
	}

	instructions := parseProgram(nums)
	return len(instructions)
}
