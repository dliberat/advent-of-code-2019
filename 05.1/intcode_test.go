package intcode

import "testing"

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

func TestCountInstructions_1_0_0_0_99(t *testing.T) {
	count := CountInstructions("1,0,0,0,99")
	if count != 2 {
		t.Errorf("Expected 2 instructions, got %d", count)
	}
}

func TestCountInstructions_1_0_0_0_2_0_0_0_99(t *testing.T) {
	count := CountInstructions("1,0,0,0,2,0,0,0,99")
	if count != 3 {
		t.Errorf("Expected 3 instructions, got %d", count)
	}
}

func TestOpTerminate_99(t *testing.T) {
	result := Run("99")
	if result != 99 {
		t.Errorf("Expected 99, got %d", result)
	}
}

func TestOpAdd_1_1_1_0_99(t *testing.T) {
	result := Run("1,1,1,0,99")
	if result != 2 {
		t.Errorf("Expected 1+1=2 but got %d", result)
	}
}

func TestOpAddImmediateMode_1101_20_22_0_99(t *testing.T) {
	result := Run("1101,20,22,0,99")
	if result != 42 {
		t.Errorf("Expected 20+22=42 but got %d", result)
	}
}

func TestOpMultiply_02_0_0_0_99(t *testing.T) {
	result := Run("02,0,0,0,99")
	if result != 4 {
		t.Errorf("Expected 2x2=4 but got %d", result)
	}
}

func TestOpMultiplyImmediateMode_102_5_0_0_99(t *testing.T) {
	result := Run("102,5,0,0,99")
	if result != 510 {
		t.Errorf("Expected 102x5=510 but got %d", result)
	}
}

func TestRun_Cyclic(t *testing.T) {
	result := Run("1002,4,3,4,33")
	if result != 1002 {
		t.Errorf("Expected 1002, 33, 3, 4, 99")
	}
}
