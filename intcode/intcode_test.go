package intcode

import (
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

func TestOpAdd_1_1_1_0_99(t *testing.T) {
	cpu := MakeComputer("1,1,1,0,99", nil, nil)
	result := cpu.Run()
	if result != 2 {
		t.Errorf("Expected 1+1=2 but got %d", result)
	}
}

func TestOpAddImmediateMode_1101_20_22_0_99(t *testing.T) {
	cpu := MakeComputer("1101,20,22,0,99", nil, nil)
	result := cpu.Run()
	if result != 42 {
		t.Errorf("Expected 20+22=42 but got %d", result)
	}
}

func TestOpMultiply_02_0_0_0_99(t *testing.T) {
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
	SetInFile(in)

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
	SetInFile(in)

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

func TestOpOutput_99(t *testing.T) {
	f, err := ioutil.TempFile("", "")
	if err != nil {
		t.Fatal(err)
	}

	cpu := MakeComputer("4,2,99", nil, f)
	result := cpu.Run()

	_, err = f.Seek(0, os.SEEK_SET)
	if err != nil {
		os.Remove(out.Name())
		t.Fatal(err)
	}

	bytes, err := ioutil.ReadFile(f.Name())
	str := string(bytes)

	os.Remove(out.Name())

	if str != "99\n" || result != 4 {
		t.Errorf("Expected 99 with result 4 but got: '%s' with result '%d'", str, result)
	}

}

func TestOpOutput_ImmediateMode(t *testing.T) {

	f, err := ioutil.TempFile("", "")
	if err != nil {
		t.Fatal(err)
	}

	cpu := MakeComputer("104,42,99", nil, f)

	result := cpu.Run()

	_, err = f.Seek(0, os.SEEK_SET)
	if err != nil {
		os.Remove(out.Name())
		t.Fatal(err)
	}

	bytes, err := ioutil.ReadFile(f.Name())
	str := string(bytes)

	os.Remove(out.Name())

	if str != "42\n" || result != 104 {
		t.Errorf("Expected 42 with result 4 but got: '%s' with result '%d'", str, result)
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

	f, err := ioutil.TempFile("", "")
	if err != nil {
		t.Fatal(err)
	}

	prog := "3,23,3,24,1002,24,10,24,1002,23,-1,23,101,5,23,23,1,24,23,23,4,23,99,0,0"
	cpu := MakeComputer(prog, in, f)
	result := cpu.Run()

	_, err = f.Seek(0, os.SEEK_SET)
	if err != nil {
		os.Remove(out.Name())
		t.Fatal(err)
	}

	bytes, err := ioutil.ReadFile(f.Name())
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

func runCPU(cpu *Computer) {
	cpu.Run()
}

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
