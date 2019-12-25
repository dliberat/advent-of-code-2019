package matmul

import "fmt"

/*CSRMatrix is a matrix stored in
compressed-sparse-row or "Yale" format.*/
type CSRMatrix struct {
	a  []int
	ia []int
	ja []int
}

/*AddRow adds a row to the matrix.*/
func (m *CSRMatrix) AddRow(row []int) {
	nnzInRow := 0

	// A is of length NNZ and holds all non-zero entries
	// of the matrix in left-to-right, top-to-bottom order

	// JA contains the column indices of each non-zero value
	for column, element := range row {
		if element != 0 {
			m.a = append(m.a, element)
			m.ja = append(m.ja, column)
			nnzInRow++
		}
	}

	if m.ia == nil {
		// The first entry of IA is always 0.
		m.ia = append(m.ia, 0)

		// The final entry is the number of non-zero elements
		// in the matrix (i.e., the length of A.)
		// This makes the for loop easier during multiplication.
		m.ia = append(m.ia, 0)
	}

	// index into A where the entries for the current row are stored
	m.ia[len(m.ia)-1] = m.ia[len(m.ia)-2] + nnzInRow

	m.ia = append(m.ia, len(m.a))
}

/*AddRowData appends data for a row to the matrix without
explicitly needing to provide the entire row. Zeroes can be ommitted.
The values slice contains the nonzero values of the row,
and the corresponding indices in the columns slice indicate
their respective columns.*/
func (m *CSRMatrix) AddRowData(values []int, columns []int) {

	for i := range values {
		m.a = append(m.a, values[i])
		m.ja = append(m.ja, columns[i])
	}

	if m.ia == nil {
		// The first entry of IA is always 0.
		m.ia = append(m.ia, 0)

		// The final entry is the number of non-zero elements
		// in the matrix (i.e., the length of A.)
		// This makes the for loop easier during multiplication.
		m.ia = append(m.ia, 0)
	}

	// index into A where the entries for the current row are stored
	m.ia[len(m.ia)-1] = m.ia[len(m.ia)-2] + len(values)

	m.ia = append(m.ia, len(m.a))
}

/*Print displays the matrix to stdout.*/
func (m *CSRMatrix) Print(colcount int) {
	for row := 0; row < len(m.ia)-2; row++ {

		start := m.ia[row]
		end := m.ia[row+1]

		fullRow := make([]int, colcount)

		for i := start; i < end; i++ {
			col := m.ja[i]
			val := m.a[i]
			fullRow[col] = val
		}

		for _, val := range fullRow {
			if val < 0 {
				fmt.Print("", val)
			} else {
				fmt.Print(" ", val)
			}
		}
		fmt.Println("")
	}
}

/*VectorMultiply applies a transformation matrix
to a given vector.*/
func VectorMultiply(vec *[]int, M *CSRMatrix) []int {

	// this is silly. Can this be done without a copy of the original vector?
	result := make([]int, len(*vec))

	for i := range *vec {

		// result[i] = 0 // most entries will end up being zero

		for k := M.ia[i]; k < M.ia[i+1]; k++ {
			result[i] = result[i] + M.a[k]*(*vec)[M.ja[k]]
		}
	}
	return result
}

/*VectorMultiplyAbsUnits applies a transformation matrix
to a given vector but the result is returned as the units
digit of the multiplication result. */
func VectorMultiplyAbsUnits(vec *[]int, M *CSRMatrix) []int {

	// this is silly. Can this be done without a copy of the original vector?
	result := make([]int, len(*vec))

	for i := range *vec {

		// result[i] = 0 // most entries will end up being zero

		for k := M.ia[i]; k < M.ia[i+1]; k++ {
			result[i] = result[i] + M.a[k]*(*vec)[M.ja[k]]
		}

		result[i] = result[i] % 10
		if result[i] < 0 {
			result[i] *= -1
		}
	}
	return result
}
