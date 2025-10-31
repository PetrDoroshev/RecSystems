package matrix

import (
	"errors"
	"fmt"
	"strings"
)

type Numeric interface {
	int | int8 | int16 | int32 | int64 |
		uint | uint8 | uint16 | uint32 | uint64 |
		float32 | float64
}

type IMatrix[T Numeric] interface {
	fmt.Stringer

	Get(row_n int, col_n int) T
	Set(row_n int, col_n int, val T)
	GetRow() []T
	DeleteRow(row_n int) error
	DeleteColumn(col_n int) error
	Transpose() *IMatrix[T]
}

type Matrix[T Numeric] struct {
	data [][]T
	Rows int
	Cols int
}

func (m Matrix[T]) String() string {

	var sb strings.Builder

	for i := 0; i < len(m.data); i++ {

		sb.WriteString("  [")

		for j := 0; j < len(m.data[0]); j++ {
			fmt.Fprintf(&sb, "%5v", m.data[i][j])
		}
		sb.WriteString(" ]\n")
	}
	return sb.String()
}

type CoordinateList[T Numeric] struct {
	Values []T
	Row    []int
	Col    []int
}

func (m CoordinateList[T]) String() string {

	var sb strings.Builder
	sb.WriteString("CoordinateList (COO format):\n\n")
	for i := range m.Values {
		fmt.Fprintf(&sb, "(%d, %d) -> %v\n", m.Row[i], m.Col[i], m.Values[i])
	}
	return sb.String()
}

type CSR[T Numeric] struct {
	Values    []T
	Col       []int
	Row_index []int
}

func (m CSR[T]) String() string {

	var sb strings.Builder

	sb.WriteString("CSR (Compressed Sparse Row) format:\n\n")
	fmt.Fprintf(&sb, "Values:    %v\n", m.Values)
	fmt.Fprintf(&sb, "Columns:   %v\n", m.Col)
	fmt.Fprintf(&sb, "Row_index: %v\n", m.Row_index)

	return sb.String()
}

type ELLPACK[T Numeric] struct {
	Value Matrix[T]
	Index Matrix[uint]
}

func (e ELLPACK[T]) String() string {

	var sb strings.Builder

	sb.WriteString("ELLPACK (ELL) format:\n\n")

	sb.WriteString("Values:\n")
	sb.WriteString(e.Value.String())
	sb.WriteString("\nIndices:\n")
	sb.WriteString(e.Index.String())

	return sb.String()
}

func NewMatrix[T Numeric](data [][]T) *Matrix[T] {

	return &Matrix[T]{data: data,
		Rows: len(data),
		Cols: len(data[0])}
}

func NewZeroMatrix[T Numeric](row_n int, col_n int) *Matrix[T] {

	data := make([][]T, row_n)

	for i := range row_n {
		data[i] = make([]T, col_n)
	}

	return &Matrix[T]{data: data, Rows: row_n, Cols: col_n}
}

func (m *Matrix[T]) Get(row_n int, col_n int) T {

	return m.data[row_n][col_n]
}

func (m *Matrix[T]) Set(row_n int, col_n int, val T) {

	m.data[row_n][col_n] = val
}

func (m *Matrix[T]) GetRow(row_n int) []T {
	return m.data[row_n]
}

func (m *Matrix[T]) GetCol(col_n int) []T {

	if m.Rows == 0 {
		return []T {} 
	}

	column := make([]T, 0, m.Rows)	
	
	for row_n := range m.Rows {
		column = append(column, m.data[row_n][col_n])
	}
	
	return column
}

func (m *Matrix[T]) DeleteRow(row_n int) error {

	if row_n > m.Rows-1 || row_n < 0 {
		return errors.New("row number out of range")
	}

	m.data = append(m.data[:row_n], m.data[row_n+1:]...)
	m.Rows--

	return nil
}

func (m *Matrix[T]) DeleteColumn(col_n int) error {

	if col_n > m.Cols-1 || col_n < 0 {
		return errors.New("column number out of range")
	}

	for i := range m.Rows {

		m.data[i] = append(m.data[i][:col_n], m.data[i][col_n+1:]...)
	}
	m.Cols--

	return nil
}

func (m *Matrix[T]) ToCoordinates() CoordinateList[T] {

	cl := CoordinateList[T]{}

	for i, row := range m.data {
		for k, item := range row {

			if item != 0 {

				cl.Values = append(cl.Values, item)
				cl.Row = append(cl.Row, i)
				cl.Col = append(cl.Col, k)

			}
		}
	}

	return cl
}

func (m *Matrix[T]) ToCSR() CSR[T] {

	csr := CSR[T]{}

	for _, row := range m.data {

		csr.Row_index = append(csr.Row_index, len(csr.Values))

		for k, item := range row {

			if item != 0 {
				csr.Values = append(csr.Values, item)
				csr.Col = append(csr.Col, k)
			}
		}
	}
	csr.Row_index = append(csr.Row_index, len(csr.Values))

	return csr
}

func (m *Matrix[T]) ToELLPACK() ELLPACK[T] {

	ellpack := ELLPACK[T]{}
	max_row_len := 0

	for _, row := range m.data {

		row_len := 0

		for _, item := range row {

			if item != 0 {
				row_len++
			}
		}

		if row_len > max_row_len {
			max_row_len = row_len
		}
	}

	ellpack.Value.data = make([][]T, len(m.data))
	ellpack.Index.data = make([][]uint, len(m.data))

	for i := 0; i < len(m.data); i++ {

		ellpack.Value.data[i] = make([]T, max_row_len)
		ellpack.Index.data[i] = make([]uint, max_row_len)
	}

	for i, row := range m.data {

		new_index := 0
		for k, item := range row {

			if item > 0 {

				ellpack.Value.data[i][new_index] = item
				ellpack.Index.data[i][new_index] = uint(k)
				new_index++
			}
		}
	}

	return ellpack
}

func (m *Matrix[T]) Transpose() *Matrix[T] {

	rows := m.Rows
	cols := m.Cols
	transposed_matrix := make([][]T, cols)

	for i := range transposed_matrix {
		transposed_matrix[i] = make([]T, rows)
	}

	for i := 0; i < rows; i++ {
		for j := 0; j < cols; j++ {
			transposed_matrix[j][i] = m.data[i][j]
		}
	}

	return &Matrix[T] {

		data: transposed_matrix,
		Rows: len(transposed_matrix),
		Cols: len(transposed_matrix[0]) }
}
