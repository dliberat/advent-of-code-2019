// Package intcode provides utilities for running intcode programs
package intcode

import (
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
)

const (
	modePosition  = 0
	modeImmediate = 1
	modeRelative  = 2
)

type context struct {
	ic           int // instruction counter
	in           *os.File
	out          *os.File
	inputQueue   []int
	outputBuffer []int
	relativeBase int
}

/*NoInputReadyError indicates that an input instruction
is waiting for input data, and thus the program cannot
continue running.*/
type NoInputReadyError struct{}

func (e *NoInputReadyError) Error() string {
	return "No input ready."
}

type instruction interface {
	execute(memory *[]int, pos *int, ctx *context)
}

type opAdd struct {
	paramModes []int
	params     []int
}

func (o opAdd) execute(memory *[]int, pos *int, ctx *context) {
	var a int
	var b int
	var outpos int
	var sum int

	switch o.paramModes[0] {
	case modePosition:
		a = getValueFromMemory(memory, o.params[0])
	case modeImmediate:
		a = o.params[0]
	case modeRelative:
		a = getValueFromMemory(memory, ctx.relativeBase+o.params[0])
	default:
		panic("Invalid parameter mode")
	}

	switch o.paramModes[1] {
	case modePosition:
		b = getValueFromMemory(memory, o.params[1])
	case modeImmediate:
		b = o.params[1]
	case modeRelative:
		b = getValueFromMemory(memory, ctx.relativeBase+o.params[1])
	default:
		panic("Invalid parameter mode")
	}

	switch o.paramModes[2] {
	case modePosition:
		outpos = o.params[2]
	case modeImmediate:
		panic("Parameters that a function writes to should never be in immediate mode.")
	case modeRelative:
		outpos = ctx.relativeBase + o.params[2]
	default:
		panic("Invalid parameter mode")
	}

	sum = a + b

	writeValueToMemory(memory, outpos, sum)
	*pos = *pos + 4
}

type opMultiply struct {
	paramModes []int
	params     []int
}

func (o opMultiply) execute(memory *[]int, pos *int, ctx *context) {
	var a int
	var b int
	var outpos int
	var prod int

	switch o.paramModes[0] {
	case modePosition:
		a = getValueFromMemory(memory, o.params[0])
	case modeImmediate:
		a = o.params[0]
	case modeRelative:
		a = getValueFromMemory(memory, ctx.relativeBase+o.params[0])
	default:
		panic("Invalid parameter mode.")
	}

	switch o.paramModes[1] {
	case modePosition:
		b = getValueFromMemory(memory, o.params[1])
	case modeImmediate:
		b = o.params[1]
	case modeRelative:
		b = getValueFromMemory(memory, ctx.relativeBase+o.params[1])
	default:
		panic("Invalid parameter mode.")
	}

	switch o.paramModes[2] {
	case modePosition:
		outpos = o.params[2]
	case modeImmediate:
		panic("Output parameters should never be in immediate mode.")
	case modeRelative:
		outpos = ctx.relativeBase + o.params[2]
	}

	prod = a * b

	writeValueToMemory(memory, outpos, prod)
	*pos = *pos + 4
}

type opInput struct {
	paramModes []int
	params     []int
}

func (o opInput) execute(memory *[]int, pos *int, ctx *context) {
	// read user input and put output into the position
	// indicated by the parameter

	var value int = 0

	if ctx.in == os.Stdin {
		fmt.Println("Input: ")
		for true {
			_, err := fmt.Fscanf(ctx.in, "%d", &value)
			if err == nil {
				break
			}
		}
	} else {
		value = ctx.inputQueue[0]
		ctx.inputQueue = ctx.inputQueue[1:]
	}

	var outpos int

	switch o.paramModes[0] {
	case modePosition:
		outpos = o.params[0]
	case modeImmediate:
		panic("Positions that a program writes to should never be in immediate mode.")
	case modeRelative:
		outpos = ctx.relativeBase + o.params[0]
	default:
		panic("Invalid position mode")
	}
	writeValueToMemory(memory, outpos, value)
	*pos = *pos + 2
}

type opOutput struct {
	paramModes []int
	params     []int
}

func (o opOutput) execute(memory *[]int, pos *int, ctx *context) {
	// output the contents of the parameter

	var data int

	switch o.paramModes[0] {
	case modePosition:
		data = getValueFromMemory(memory, o.params[0])
	case modeImmediate:
		data = o.params[0]
	case modeRelative:
		data = getValueFromMemory(memory, ctx.relativeBase+o.params[0])
	default:
		panic("Invalid parameter mode.")
	}

	if ctx.out == nil {
		ctx.outputBuffer = append(ctx.outputBuffer, data)
	} else {
		val := fmt.Sprintf("%v\n", data)
		io.WriteString(ctx.out, val)
	}

	*pos = *pos + 2
}

type opJumpIfTrue struct {
	paramModes []int
	params     []int
}

func (o opJumpIfTrue) execute(memory *[]int, pos *int, ctx *context) {
	var a int
	var b int

	switch o.paramModes[0] {
	case modePosition:
		a = getValueFromMemory(memory, o.params[0])
	case modeImmediate:
		a = o.params[0]
	case modeRelative:
		a = getValueFromMemory(memory, ctx.relativeBase+o.params[0])
	default:
		panic("Invalid parameter mode.")
	}

	switch o.paramModes[1] {
	case modePosition:
		b = getValueFromMemory(memory, o.params[1])
	case modeImmediate:
		b = o.params[1]
	case modeRelative:
		b = getValueFromMemory(memory, ctx.relativeBase+o.params[1])
	default:
		panic("Invalid parameter mode.")
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

func (o opJumpIfFalse) execute(memory *[]int, pos *int, ctx *context) {
	var a int
	var b int

	switch o.paramModes[0] {
	case modePosition:
		a = getValueFromMemory(memory, o.params[0])
	case modeImmediate:
		a = o.params[0]
	case modeRelative:
		a = getValueFromMemory(memory, ctx.relativeBase+o.params[0])
	default:
		panic("Invalid parameter mode.")
	}

	switch o.paramModes[1] {
	case modePosition:
		b = getValueFromMemory(memory, o.params[1])
	case modeImmediate:
		b = o.params[1]
	case modeRelative:
		b = getValueFromMemory(memory, ctx.relativeBase+o.params[1])
	default:
		panic("Invalid parameter mode.")
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

func (o opLessThan) execute(memory *[]int, pos *int, ctx *context) {
	var a int
	var b int
	var outpos int

	switch o.paramModes[0] {
	case modePosition:
		a = getValueFromMemory(memory, o.params[0])
	case modeImmediate:
		a = o.params[0]
	case modeRelative:
		a = getValueFromMemory(memory, ctx.relativeBase+o.params[0])
	default:
		panic("Invalid parameter mode")
	}

	switch o.paramModes[1] {
	case modePosition:
		b = getValueFromMemory(memory, o.params[1])
	case modeImmediate:
		b = o.params[1]
	case modeRelative:
		b = getValueFromMemory(memory, ctx.relativeBase+o.params[1])
	default:
		panic("Invalid parameter mode")
	}

	switch o.paramModes[2] {
	case modePosition:
		outpos = o.params[2]
	case modeImmediate:
		panic("Positions that a function writes to should never be in immediate mode.")
	case modeRelative:
		outpos = ctx.relativeBase + o.params[2]
	}

	if a < b {
		writeValueToMemory(memory, outpos, 1)
	} else {
		writeValueToMemory(memory, outpos, 0)
	}
	*pos = *pos + 4
}

type opEquals struct {
	paramModes []int
	params     []int
}

func (o opEquals) execute(memory *[]int, pos *int, ctx *context) {
	var a int
	var b int
	var outpos int

	switch o.paramModes[0] {
	case modePosition:
		a = getValueFromMemory(memory, o.params[0])
	case modeImmediate:
		a = o.params[0]
	case modeRelative:
		a = getValueFromMemory(memory, ctx.relativeBase+o.params[0])
	default:
		panic("invalid parameter mode.")
	}

	switch o.paramModes[1] {
	case modePosition:
		b = getValueFromMemory(memory, o.params[1])
	case modeImmediate:
		b = o.params[1]
	case modeRelative:
		b = getValueFromMemory(memory, ctx.relativeBase+o.params[1])
	default:
		panic("Invalid parameter mode.")
	}

	switch o.paramModes[2] {
	case modePosition:
		outpos = o.params[2]
	case modeImmediate:
		panic("Positions that a function writes to should never be in immediate mode.")
	case modeRelative:
		outpos = ctx.relativeBase + o.params[2]
	}

	if a == b {
		writeValueToMemory(memory, outpos, 1)
	} else {
		writeValueToMemory(memory, outpos, 0)
	}
	*pos = *pos + 4
}

type opRelativeBase struct {
	paramModes []int
	params     []int
}

func (o opRelativeBase) execute(memory *[]int, pos *int, ctx *context) {
	var a int
	switch o.paramModes[0] {
	case modePosition:
		a = getValueFromMemory(memory, o.params[0])
	case modeImmediate:
		a = o.params[0]
	case modeRelative:
		a = getValueFromMemory(memory, ctx.relativeBase+o.params[0])
	}

	ctx.relativeBase = ctx.relativeBase + a
	*pos = *pos + 2
}

type opTerminate struct {
}

func (opTerminate) execute(memory *[]int, pos *int, ctx *context) {
	*pos = len(*memory)
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
	case 9:
		return 1
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

func getValueFromMemory(memory *[]int, address int) int {
	if address < len(*memory) {
		return (*memory)[address]
	}
	return 0
}

func writeValueToMemory(memory *[]int, address int, value int) {
	if address >= cap(*memory) {
		newMem := make([]int, address+1)
		copy(newMem, *memory)
		*memory = newMem
	}

	(*memory)[address] = value
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
	case 9:
		return opRelativeBase{paramModes: paramModes, params: params}
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

/*Computer represents an instance of an intcode processor.*/
type Computer struct {
	ctx    context
	memory []int
}

/*Run executes the program provided and returns the
result at the 0th index of memory when the program is complete.*/
func (cpu *Computer) Run() (int, error) {

	for true {
		currentCode := cpu.memory[cpu.ctx.ic]

		opcode, paramModes := ParseOpCode(strconv.Itoa(currentCode))

		if opcode == 99 {
			break
		}

		// input instruction but there is no input  ready
		if opcode == 3 && cpu.ctx.in != os.Stdin && len(cpu.ctx.inputQueue) == 0 {
			return 0, &NoInputReadyError{}
		}

		argv := GetArgumentCount(opcode)
		m := cpu.ctx.ic + 1
		instr := makeInstruction(opcode, paramModes, cpu.memory[m:m+argv])
		instr.execute(&cpu.memory, &cpu.ctx.ic, &cpu.ctx)
	}

	return cpu.memory[0], nil
}

/*QueueInput queues a variable number of inputs to be used
as arguments into input instructions when required by the
intcode program. Queued inputs are ignored if the input
source is set to os.Stdin*/
func (cpu *Computer) QueueInput(values ...int) {
	cpu.ctx.inputQueue = append(cpu.ctx.inputQueue, values...)
}

/*FlushOutput returns a copy of all the outputs the
computer has buffered and clears the buffer.*/
func (cpu *Computer) FlushOutput() []int {
	cpy := make([]int, 0)

	for _, v := range cpu.ctx.outputBuffer {
		cpy = append(cpy, v)
	}

	cpu.ctx.outputBuffer = nil

	return cpy
}

/*PrintState shows some of the internal state of the computer.*/
func (cpu *Computer) PrintState() {
	fmt.Println("Relative base: ", cpu.ctx.relativeBase)
	fmt.Println("Memory: ", cpu.memory)
}

/*Clone produces a replica of the computer with all
its internal state intact.*/
func (cpu *Computer) Clone() Computer {
	mem := make([]int, len(cpu.memory))
	copy(mem, cpu.memory)

	iq := make([]int, len(cpu.ctx.inputQueue))
	copy(iq, cpu.ctx.inputQueue)

	oq := make([]int, len(cpu.ctx.outputBuffer))
	copy(oq, cpu.ctx.outputBuffer)

	ctx := context{ic: cpu.ctx.ic, in: cpu.ctx.in, out: cpu.ctx.out, inputQueue: iq, outputBuffer: oq, relativeBase: cpu.ctx.relativeBase}
	return Computer{memory: mem, ctx: ctx}
}

/*MakeComputer creates a computer object that can be used to
process intcode programs.*/
func MakeComputer(memory string, in *os.File, out *os.File) Computer {

	ctx := context{in: in, out: out, ic: 0}
	cpu := Computer{ctx: ctx}

	txt := strings.Split(memory, ",")
	cpu.memory = make([]int, len(txt))
	for i, t := range txt {
		cpu.memory[i], _ = strconv.Atoi(t)
	}

	return cpu
}
