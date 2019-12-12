package intcode

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"testing"
)

func TestParseOpcode_1002(t *testing.T) {
	code, params := ParseOpCode("1002")

	if code != 2 {
		t.Errorf("Expected code == 2, got: %d", code)
	}

	if len(params) != 3 {
		t.Errorf("Expected 3 params, got: %d", len(params))
	}

	if params[0] != 0 || params[1] != 1 || params[2] != 0 {
		t.Errorf("Expected params [0 1 0]. Got: %v", params)
	}
}

func TestParseOpCode_99(t *testing.T) {
	code, params := ParseOpCode("99")
	if code != 99 {
		t.Errorf("Expected code == 99, got %d", code)
	}
	if len(params) != 0 {
		t.Errorf("Expected zero-length params. Got %v", params)
	}
}

func TestParseOpCode_1(t *testing.T) {
	code, params := ParseOpCode("1")
	if code != 1 {
		t.Errorf("Expected code 1, got %d", code)
	}
	if len(params) != 3 || params[0] != 0 || params[1] != 0 || params[2] != 0 {
		t.Errorf("Expected params [0 0 0] but got %v", params)
	}
}

func TestParseOpCode_leadingZero(t *testing.T) {
	code, _ := ParseOpCode("02")
	if code != 2 {
		t.Errorf("Expected code 2, got %d", code)
	}
}

func TestOpTerminate_99(t *testing.T) {
	cpu := MakeComputer("99", nil, nil)
	result := cpu.Run()
	if result != 99 {
		t.Errorf("Expected 99, got %d", result)
	}
}

func TestOpAdd_1_1_1_10_99(t *testing.T) {
	cpu := MakeComputer("1,1,1,0,99", nil, nil)
	result := cpu.Run()
	if result != 2 {
		t.Errorf("Expected 1+1=2 but got %d", result)
	}
}

func TestOpAdd_PositionMode(t *testing.T) {
	cpu := MakeComputer("1,5,7,0,99,5,0,7", nil, nil)
	result := cpu.Run()
	if result != 12 {
		cpu.PrintState()
		t.Errorf("Expected 5+7=12 but got '%d'", result)
	}
}

func TestOpAdd_GetBeyondMemoryBounds(t *testing.T) {
	cpu := MakeComputer("1,43,4,0,99", nil, nil)
	result := cpu.Run()
	if result != 99 {
		cpu.PrintState()
		t.Errorf("Expected 0+99 to equal 99 but got '%d'", result)
	}
}

func TestOpAddImmediateMode_1101_20_22_0_99(t *testing.T) {
	cpu := MakeComputer("1101,20,22,0,99", nil, nil)
	result := cpu.Run()
	if result != 42 {
		t.Errorf("Expected 20+22=42 but got %d", result)
	}
}

func TestOpAddRelativeMode(t *testing.T) {
	cpu := MakeComputer("109,2,22201,-1,0,-2,99", nil, nil)
	result := cpu.Run()
	if result != 22203 {
		t.Errorf("Expected 22203 but got '%d'", result)
	}
}

func TestOpMultiply_02_0_0_1_99(t *testing.T) {
	cpu := MakeComputer("02,0,0,0,99", nil, nil)
	result := cpu.Run()
	if result != 4 {
		t.Errorf("Expected 2x2=4 but got %d", result)
	}
}

func TestOpMultiplyImmediateMode_102_5_0_0_99(t *testing.T) {
	cpu := MakeComputer("102,5,0,0,99", nil, nil)
	result := cpu.Run()
	if result != 510 {
		t.Errorf("Expected 102x5=510 but got %d", result)
	}
}

func TestOpMultiplyRelativeMode_zeroShift(t *testing.T) {
	cpu := MakeComputer("1202,0,2,0,99", nil, nil)
	result := cpu.Run()
	if result != 2404 {
		t.Errorf("Expected 1202x2=2404 but got '%d'", result)
	}
}

func TestOpMultiplyRelativeMode_withShift(t *testing.T) {
	cpu := MakeComputer("109,-1,22202,5,3,1,99", nil, nil)
	result := cpu.Run()
	if result != 66606 {
		cpu.PrintState()
		t.Errorf("Expected 66606 but got '%d'", result)
	}
}

func TestOpInput_3_0_99(t *testing.T) {
	in, err := ioutil.TempFile("", "")
	if err != nil {
		t.Fatal(err)
	}
	defer in.Close()

	_, err = in.Seek(0, os.SEEK_SET)
	if err != nil {
		os.Remove(in.Name())
		t.Fatal(err)
	}
	io.WriteString(in, "42\n")
	if err != nil {
		os.Remove(in.Name())
		t.Fatal(err)
	}

	_, err = in.Seek(0, os.SEEK_SET)
	if err != nil {
		os.Remove(in.Name())
		t.Fatal(err)
	}

	cpu := MakeComputer("3,0,99", in, nil)
	result := cpu.Run()

	os.Remove(in.Name())

	if result != 42 {
		t.Errorf("Expected 42, got: %d", result)
	}
}

func TestOpInput_Two_Inputs(t *testing.T) {
	in, err := ioutil.TempFile("", "")
	if err != nil {
		t.Fatal(err)
	}
	defer in.Close()

	_, err = in.Seek(0, os.SEEK_SET)
	if err != nil {
		os.Remove(in.Name())
		t.Fatal(err)
	}
	io.WriteString(in, "3\n42\n")
	if err != nil {
		os.Remove(in.Name())
		t.Fatal(err)
	}

	_, err = in.Seek(0, os.SEEK_SET)
	if err != nil {
		os.Remove(in.Name())
		t.Fatal(err)
	}

	cpu := MakeComputer("3,2,1,0,99", in, nil)
	result := cpu.Run()

	os.Remove(in.Name())

	if result != 42 {
		t.Errorf("Expected 42, got: %d", result)
	}
}

func TestOpInput_RelativeMode(t *testing.T) {
	in, err := ioutil.TempFile("", "")
	if err != nil {
		t.Fatal(err)
	}
	defer in.Close()

	_, err = in.Seek(0, os.SEEK_SET)
	if err != nil {
		os.Remove(in.Name())
		t.Fatal(err)
	}
	io.WriteString(in, "42\n")
	if err != nil {
		os.Remove(in.Name())
		t.Fatal(err)
	}

	_, err = in.Seek(0, os.SEEK_SET)
	if err != nil {
		os.Remove(in.Name())
		t.Fatal(err)
	}

	cpu := MakeComputer("109,50,203,-50,99,0", in, nil)
	result := cpu.Run()

	os.Remove(in.Name())

	if result != 42 {
		t.Errorf("Expected 42, got: %d", result)
	}
}

func TestOpOutput_99(t *testing.T) {
	f, err := ioutil.TempFile("", "")
	if err != nil {
		t.Fatal(err)
	}

	cpu := MakeComputer("4,2,99", nil, f)
	result := cpu.Run()

	_, err = f.Seek(0, os.SEEK_SET)
	if err != nil {
		os.Remove(f.Name())
		t.Fatal(err)
	}

	bytes, err := ioutil.ReadFile(f.Name())
	str := string(bytes)

	os.Remove(f.Name())

	if str != "99\n" || result != 4 {
		t.Errorf("Expected 99 with result 4 but got: '%s' with result '%d'", str, result)
	}

}

func TestOpOutput_ImmediateMode(t *testing.T) {

	out, err := ioutil.TempFile("", "")
	if err != nil {
		t.Fatal(err)
	}

	cpu := MakeComputer("104,42,99", nil, out)

	result := cpu.Run()

	_, err = out.Seek(0, os.SEEK_SET)
	if err != nil {
		os.Remove(out.Name())
		t.Fatal(err)
	}

	bytes, err := ioutil.ReadFile(out.Name())
	str := string(bytes)

	os.Remove(out.Name())

	if str != "42\n" || result != 104 {
		t.Errorf("Expected 42 with result 4 but got: '%s' with result '%d'", str, result)
	}

}

func TestOpOutput_RelativeMode(t *testing.T) {

	out, err := ioutil.TempFile("", "")
	if err != nil {
		t.Fatal(err)
	}

	cpu := MakeComputer("109,42,204,-42,99", nil, out)

	result := cpu.Run()

	_, err = out.Seek(0, os.SEEK_SET)
	if err != nil {
		os.Remove(out.Name())
		t.Fatal(err)
	}

	bytes, err := ioutil.ReadFile(out.Name())
	str := string(bytes)

	os.Remove(out.Name())

	if str != "109\n" || result != 109 {
		t.Errorf("Expected 109 with result 109 but got: '%s' with result '%d'", str, result)
	}

}

func TestJumpIfTrue_False(t *testing.T) {
	cpu := MakeComputer("109,30,105,0,9,1101,21,21,0,99,1", nil, nil)
	result := cpu.Run()
	if result != 42 {
		cpu.PrintState()
		t.Errorf("Expected 42 but got '%d'", result)
	}
}

func TestJumpIfTrue_True(t *testing.T) {
	cpu := MakeComputer("109,30,205,-29,9,1101,21,21,0,99,1", nil, nil)
	result := cpu.Run()
	if result != 109 {
		cpu.PrintState()
		t.Errorf("Expected 109 but got '%d'", result)
	}
}

func TestJumpIfFalse_False(t *testing.T) {
	cpu := MakeComputer("109,30,106,0,9,1101,21,21,0,99,1", nil, nil)
	result := cpu.Run()
	if result != 109 {
		cpu.PrintState()
		t.Errorf("Expected 109 but got '%d'", result)
	}
}

func TestJumpIfFalse_True(t *testing.T) {
	cpu := MakeComputer("109,30,206,-29,9,1101,21,21,0,99,1", nil, nil)
	result := cpu.Run()
	if result != 42 {
		cpu.PrintState()
		t.Errorf("Expected 42 but got '%d'", result)
	}
}

func TestLessThan_gt(t *testing.T) {
	cpu := MakeComputer("109,4,21007,1,2,-4,99", nil, nil)
	result := cpu.Run()
	if result != 0 {
		cpu.PrintState()
		t.Errorf("Expected 2 to be less than 4.")
	}
}

func TestLessThan_lt(t *testing.T) {
	cpu := MakeComputer("109,4,01207,3,2,0,99", nil, nil)
	result := cpu.Run()
	if result != 1 {
		cpu.PrintState()
		t.Errorf("Expected 0 to be less than 2.")
	}
}

func TestOpEqual_isEqual(t *testing.T) {
	cpu := MakeComputer("9,0,21008,1,0,-9,99", nil, nil)
	result := cpu.Run()
	if result != 1 {
		cpu.PrintState()
		t.Errorf("Expected 0 to be equal to 0.")
	}
}

func TestOpEqual_isNotEqual(t *testing.T) {
	cpu := MakeComputer("9,3,02108,-1,3,0,99,0", nil, nil)
	result := cpu.Run()
	if result != 0 {
		cpu.PrintState()
		t.Errorf("Expected -1 to NOT equal 2108.")
	}
}

func TestRun(t *testing.T) {
	in, err := ioutil.TempFile("", "")
	if err != nil {
		t.Fatal(err)
	}
	defer in.Close()

	_, err = in.Seek(0, os.SEEK_SET)
	if err != nil {
		os.Remove(in.Name())
		t.Fatal(err)
	}
	io.WriteString(in, "1\n1523\n")
	if err != nil {
		os.Remove(in.Name())
		t.Fatal(err)
	}

	_, err = in.Seek(0, os.SEEK_SET)
	if err != nil {
		os.Remove(in.Name())
		t.Fatal(err)
	}

	out, err := ioutil.TempFile("", "")
	if err != nil {
		t.Fatal(err)
	}

	prog := "3,23,3,24,1002,24,10,24,1002,23,-1,23,101,5,23,23,1,24,23,23,4,23,99,0,0"
	cpu := MakeComputer(prog, in, out)
	result := cpu.Run()

	_, err = out.Seek(0, os.SEEK_SET)
	if err != nil {
		os.Remove(out.Name())
		t.Fatal(err)
	}

	bytes, err := ioutil.ReadFile(out.Name())
	str := string(bytes)

	os.Remove(in.Name())
	os.Remove(out.Name())

	if str != "15234\n" || result != 3 {
		t.Errorf("Expected 15234 with result 3 but got: '%s' with result '%d'", str, result)
	}

}

func TestRun_Cyclic(t *testing.T) {
	cpu := MakeComputer("1002,4,3,4,33", nil, nil)
	result := cpu.Run()
	if result != 1002 {
		t.Errorf("Expected 1002, 33, 3, 4, 99")
	}
}

// func runCPU(cpu *Computer) {
// 	cpu.Run()
// }

// func TestFeedback(t *testing.T) {
// 	in1, _ := ioutil.TempFile("", "infile")
// 	defer in1.Close()
// 	in1.Seek(0, os.SEEK_SET)
// 	io.WriteString(in1, "4\n")
// 	in1.Seek(0, os.SEEK_SET)

// 	out1, _ := ioutil.TempFile("", "outfile")

// 	cpu1 := MakeComputer("3,2,1,7,3,0,99,40", in1, out1)
// 	cpu2 := MakeComputer("3,6,1,7,8,0,99,0,0", out1, in1)

// 	go runCPU(&cpu2)
// 	result := cpu1.Run()

// 	os.Remove(out.Name())

// 	if result != 42 {
// 		t.Errorf("Expected 42 but got: '%d'", result)
// 	}
// }

func TestOpOffset(t *testing.T) {
	cpu := MakeComputer("109,2,21101,21,21,-2,99", nil, nil)
	result := cpu.Run()
	if result != 42 {
		t.Errorf("Expected 21+21=42 but got '%d'", result)
	}
}

func TestOpOffset_increases(t *testing.T) {
	cpu := MakeComputer("109,1,109,1,2201,-2,-1,0,99", nil, nil)
	result := cpu.Run()
	if result != 110 {
		t.Errorf("Expected 110 but got '%d'", result)
	}
}

func TestAccessExtraMemory(t *testing.T) {
	cpu := MakeComputer("2,20,1,0,99", nil, nil)
	result := cpu.Run()
	if result != 0 {
		t.Errorf("Expected 20 x 0 = 0 but got '%d'", result)
	}

}

func TestDay2(t *testing.T) {
	cpu := MakeComputer("1,9,10,3,2,3,11,0,99,30,40,50", nil, nil)
	result := cpu.Run()
	if result != 3500 {
		t.Errorf("Expected 3500 but got '%d'", result)
	}
}

func TestDay5_less_than_8(t *testing.T) {
	in, err := ioutil.TempFile("", "")
	if err != nil {
		t.Fatal(err)
	}
	defer in.Close()

	_, err = in.Seek(0, os.SEEK_SET)
	if err != nil {
		os.Remove(in.Name())
		t.Fatal(err)
	}
	io.WriteString(in, "7\n")
	if err != nil {
		os.Remove(in.Name())
		t.Fatal(err)
	}

	_, err = in.Seek(0, os.SEEK_SET)
	if err != nil {
		os.Remove(in.Name())
		t.Fatal(err)
	}

	out, err := ioutil.TempFile("", "")
	if err != nil {
		t.Fatal(err)
	}

	cpu := MakeComputer("3,21,1008,21,8,20,1005,20,22,107,8,21,20,1006,20,31,1106,0,36,98,0,0,1002,21,125,20,4,20,1105,1,46,104,999,1105,1,46,1101,1000,1,20,4,20,1105,1,46,98,99", in, out)
	cpu.Run()

	bytes, err := ioutil.ReadFile(out.Name())
	str := string(bytes)

	os.Remove(out.Name())

	if str != "999\n" {
		t.Errorf("Expected 999 but got '%s'", str)
	}
}

func TestDay5_equal_to_8(t *testing.T) {
	in, err := ioutil.TempFile("", "")
	if err != nil {
		t.Fatal(err)
	}
	defer in.Close()

	_, err = in.Seek(0, os.SEEK_SET)
	if err != nil {
		os.Remove(in.Name())
		t.Fatal(err)
	}
	io.WriteString(in, "8\n")
	if err != nil {
		os.Remove(in.Name())
		t.Fatal(err)
	}

	_, err = in.Seek(0, os.SEEK_SET)
	if err != nil {
		os.Remove(in.Name())
		t.Fatal(err)
	}

	out, err := ioutil.TempFile("", "")
	if err != nil {
		t.Fatal(err)
	}

	cpu := MakeComputer("3,21,1008,21,8,20,1005,20,22,107,8,21,20,1006,20,31,1106,0,36,98,0,0,1002,21,125,20,4,20,1105,1,46,104,999,1105,1,46,1101,1000,1,20,4,20,1105,1,46,98,99", in, out)
	cpu.Run()

	bytes, err := ioutil.ReadFile(out.Name())
	str := string(bytes)

	os.Remove(out.Name())

	if str != "1000\n" {
		t.Errorf("Expected 1000 but got '%s'", str)
	}
}

func TestDay5_greater_than_8(t *testing.T) {
	in, err := ioutil.TempFile("", "")
	if err != nil {
		t.Fatal(err)
	}
	defer in.Close()

	_, err = in.Seek(0, os.SEEK_SET)
	if err != nil {
		os.Remove(in.Name())
		t.Fatal(err)
	}
	io.WriteString(in, "9\n")
	if err != nil {
		os.Remove(in.Name())
		t.Fatal(err)
	}

	_, err = in.Seek(0, os.SEEK_SET)
	if err != nil {
		os.Remove(in.Name())
		t.Fatal(err)
	}

	out, err := ioutil.TempFile("", "")
	if err != nil {
		t.Fatal(err)
	}

	cpu := MakeComputer("3,21,1008,21,8,20,1005,20,22,107,8,21,20,1006,20,31,1106,0,36,98,0,0,1002,21,125,20,4,20,1105,1,46,104,999,1105,1,46,1101,1000,1,20,4,20,1105,1,46,98,99", in, out)
	cpu.Run()

	bytes, err := ioutil.ReadFile(out.Name())
	str := string(bytes)

	os.Remove(out.Name())

	if str != "1001\n" {
		t.Errorf("Expected 1001 but got '%s'", str)
	}
}

func TestDay5_jump_zero_position(t *testing.T) {
	in, err := ioutil.TempFile("", "")
	if err != nil {
		t.Fatal(err)
	}
	defer in.Close()

	_, err = in.Seek(0, os.SEEK_SET)
	if err != nil {
		os.Remove(in.Name())
		t.Fatal(err)
	}
	io.WriteString(in, "0\n")
	if err != nil {
		os.Remove(in.Name())
		t.Fatal(err)
	}

	_, err = in.Seek(0, os.SEEK_SET)
	if err != nil {
		os.Remove(in.Name())
		t.Fatal(err)
	}

	out, err := ioutil.TempFile("", "")
	if err != nil {
		t.Fatal(err)
	}

	cpu := MakeComputer("3,12,6,12,15,1,13,14,13,4,13,99,-1,0,1,9", in, out)
	cpu.Run()

	bytes, err := ioutil.ReadFile(out.Name())
	str := string(bytes)

	os.Remove(out.Name())

	if str != "0\n" {
		t.Errorf("Expected 0 but got '%s'", str)
	}
}

func TestDay5_jump_non_zero_position(t *testing.T) {
	in, err := ioutil.TempFile("", "")
	if err != nil {
		t.Fatal(err)
	}
	defer in.Close()

	_, err = in.Seek(0, os.SEEK_SET)
	if err != nil {
		os.Remove(in.Name())
		t.Fatal(err)
	}
	io.WriteString(in, "43\n")
	if err != nil {
		os.Remove(in.Name())
		t.Fatal(err)
	}

	_, err = in.Seek(0, os.SEEK_SET)
	if err != nil {
		os.Remove(in.Name())
		t.Fatal(err)
	}

	out, err := ioutil.TempFile("", "")
	if err != nil {
		t.Fatal(err)
	}

	cpu := MakeComputer("3,12,6,12,15,1,13,14,13,4,13,99,-1,0,1,9", in, out)
	cpu.Run()

	bytes, err := ioutil.ReadFile(out.Name())
	str := string(bytes)

	os.Remove(out.Name())

	if str != "1\n" {
		t.Errorf("Expected 1 but got '%s'", str)
	}
}

func TestDay5_jump_zero_immediate(t *testing.T) {
	in, err := ioutil.TempFile("", "")
	if err != nil {
		t.Fatal(err)
	}
	defer in.Close()

	_, err = in.Seek(0, os.SEEK_SET)
	if err != nil {
		os.Remove(in.Name())
		t.Fatal(err)
	}
	io.WriteString(in, "0\n")
	if err != nil {
		os.Remove(in.Name())
		t.Fatal(err)
	}

	_, err = in.Seek(0, os.SEEK_SET)
	if err != nil {
		os.Remove(in.Name())
		t.Fatal(err)
	}

	out, err := ioutil.TempFile("", "")
	if err != nil {
		t.Fatal(err)
	}

	cpu := MakeComputer("3,3,1105,-1,9,1101,0,0,12,4,12,99,1", in, out)
	cpu.Run()

	bytes, err := ioutil.ReadFile(out.Name())
	str := string(bytes)

	os.Remove(out.Name())

	if str != "0\n" {
		t.Errorf("Expected 0 but got '%s'", str)
	}
}

func TestDay5_jump_non_zero_immediate(t *testing.T) {
	in, err := ioutil.TempFile("", "")
	if err != nil {
		t.Fatal(err)
	}
	defer in.Close()

	_, err = in.Seek(0, os.SEEK_SET)
	if err != nil {
		os.Remove(in.Name())
		t.Fatal(err)
	}
	io.WriteString(in, "32\n")
	if err != nil {
		os.Remove(in.Name())
		t.Fatal(err)
	}

	_, err = in.Seek(0, os.SEEK_SET)
	if err != nil {
		os.Remove(in.Name())
		t.Fatal(err)
	}

	out, err := ioutil.TempFile("", "")
	if err != nil {
		t.Fatal(err)
	}

	cpu := MakeComputer("3,3,1105,-1,9,1101,0,0,12,4,12,99,1", in, out)
	cpu.Run()

	bytes, err := ioutil.ReadFile(out.Name())
	str := string(bytes)

	os.Remove(out.Name())

	if str != "1\n" {
		t.Errorf("Expected 1 but got '%s'", str)
	}
}
func TestDay9_quine(t *testing.T) {
	out, err := ioutil.TempFile("", "")
	if err != nil {
		t.Fatal(err)
	}
	cpu := MakeComputer("109,1,204,-1,1001,100,1,100,1008,100,16,101,1006,101,0,99", nil, nil)
	result := cpu.Run()

	bytes, err := ioutil.ReadFile(out.Name())
	str := string(bytes)
	os.Remove(out.Name())

	fmt.Println(result)
	fmt.Println(str)
}

func TestDay9_largeNum(t *testing.T) {

	out, err := ioutil.TempFile("", "")
	if err != nil {
		t.Fatal(err)
	}

	cpu := MakeComputer("104,1125899906842624,99", nil, out)
	cpu.Run()

	bytes, err := ioutil.ReadFile(out.Name())
	str := string(bytes)
	os.Remove(out.Name())

	if str != "1125899906842624\n" {
		t.Errorf("Expected 1125899906842624 but got '%s'", str)
	}
}

func TestDay9_16digitNum(t *testing.T) {

	out, err := ioutil.TempFile("", "")
	if err != nil {
		t.Fatal(err)
	}

	cpu := MakeComputer("1102,34915192,34915192,7,4,7,99,0", nil, out)
	cpu.Run()

	bytes, err := ioutil.ReadFile(out.Name())
	str := string(bytes)
	os.Remove(out.Name())

	if len(str) != 17 {
		t.Errorf("Expected a 16-digit number but but got '%s'", str)
	}
}
