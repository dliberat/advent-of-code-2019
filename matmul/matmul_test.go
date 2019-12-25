package matmul

import "testing"

func TestVectorMultiply(t *testing.T) {

	M := CSRMatrix{}
	M.AddRow([]int{11, 12, 0, 14, 0, 0, 0, 0})
	M.AddRow([]int{0, 22, 23, 0, 25, 0, 0, 0})
	M.AddRow([]int{31, 0, 33, 34, 0, 0, 0, 0})
	M.AddRow([]int{0, 42, 0, 0, 45, 46, 0, 0})
	M.AddRow([]int{0, 0, 0, 0, 55, 0, 0, 0})
	M.AddRow([]int{0, 0, 0, 0, 65, 66, 67, 0})
	M.AddRow([]int{0, 0, 0, 0, 75, 0, 77, 78})
	M.AddRow([]int{0, 0, 0, 0, 0, 0, 87, 88})

	vector := []int{1, 0, 0, 0, 0, 0, 0, 0}

	vector = VectorMultiply(&vector, &M)

	expected := []int{11, 0, 31, 0, 0, 0, 0, 0}

	for i := range vector {
		if vector[i] != expected[i] {
			t.Errorf("Expected 11,0,31,0,0,0,0,0 but got: %v", vector)

		}
	}
}

func TestVectorMultiply2(t *testing.T) {

	M := CSRMatrix{}
	M.AddRow([]int{1, 0, -1, 0, 1, 0, -1, 0})
	M.AddRow([]int{0, 1, 1, 0, 0, -1, -1, 0})
	M.AddRow([]int{0, 0, 1, 1, 1, 0, 0, 0})
	M.AddRow([]int{0, 0, 0, 1, 1, 1, 1, 0})
	M.AddRow([]int{0, 0, 0, 0, 1, 1, 1, 1})
	M.AddRow([]int{0, 0, 0, 0, 0, 1, 1, 1})
	M.AddRow([]int{0, 0, 0, 0, 0, 0, 1, 1})
	M.AddRow([]int{0, 0, 0, 0, 0, 0, 0, 1})

	vector := []int{1, 2, 3, 4, 5, 6, 7, 8}

	vector = VectorMultiply(&vector, &M)

	expected := []int{-4, -8, 12, 22, 26, 21, 15, 8}

	for i := range vector {
		if vector[i] != expected[i] {
			t.Errorf("Expected %v but got: %v", expected, vector)

			return
		}
	}
}

func TestVectorMultiplyAbsUnits(t *testing.T) {

	M := CSRMatrix{}
	M.AddRow([]int{1, 0, -1, 0, 1, 0, -1, 0})
	M.AddRow([]int{0, 1, 1, 0, 0, -1, -1, 0})
	M.AddRow([]int{0, 0, 1, 1, 1, 0, 0, 0})
	M.AddRow([]int{0, 0, 0, 1, 1, 1, 1, 0})
	M.AddRow([]int{0, 0, 0, 0, 1, 1, 1, 1})
	M.AddRow([]int{0, 0, 0, 0, 0, 1, 1, 1})
	M.AddRow([]int{0, 0, 0, 0, 0, 0, 1, 1})
	M.AddRow([]int{0, 0, 0, 0, 0, 0, 0, 1})

	vector := []int{1, 2, 3, 4, 5, 6, 7, 8}

	vector = VectorMultiplyAbsUnits(&vector, &M)

	expected := []int{4, 8, 2, 2, 6, 1, 5, 8}

	for i := range vector {
		if vector[i] != expected[i] {
			t.Errorf("Expected %v but got: %v", expected, vector)
			return
		}
	}

	vector = VectorMultiplyAbsUnits(&vector, &M)

	expected = []int{3, 4, 0, 4, 0, 4, 3, 8}

	for i := range vector {
		if vector[i] != expected[i] {
			t.Errorf("Expected %v but got: %v", expected, vector)
			return
		}
	}
}

func TestAddRowData(t *testing.T) {

	M := CSRMatrix{}

	r := []int{1, -1, 1, -1}
	c := []int{0, 2, 4, 6}
	M.AddRowData(r, c)

	r = []int{1, 1, -1, -1}
	c = []int{1, 2, 5, 6}
	M.AddRowData(r, c)

	r = []int{1, 1, 1}
	c = []int{2, 3, 4}
	M.AddRowData(r, c)

	r = []int{1, 1, 1, 1}
	c = []int{3, 4, 5, 6}
	M.AddRowData(r, c)

	M.AddRow([]int{0, 0, 0, 0, 1, 1, 1, 1})
	M.AddRow([]int{0, 0, 0, 0, 0, 1, 1, 1})
	M.AddRow([]int{0, 0, 0, 0, 0, 0, 1, 1})

	r = []int{1}
	c = []int{7}
	M.AddRowData(r, c)

	vector := []int{1, 2, 3, 4, 5, 6, 7, 8}

	vector = VectorMultiply(&vector, &M)

	expected := []int{-4, -8, 12, 22, 26, 21, 15, 8}

	for i := range vector {
		if vector[i] != expected[i] {
			t.Errorf("Expected %v but got: %v", expected, vector)

			return
		}
	}
}
