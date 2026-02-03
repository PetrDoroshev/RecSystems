package matrix

import (
	"errors"
)

type KeyedMatrix[T Numeric, K1 comparable, K2 comparable] struct {
	matrix Matrix[T]

	RowKeyToIndex map[K1]int
	ColKeyToIndex map[K2]int

	RowKeys []K1
	ColKeys []K2
}

func NewKeyedMatrix[T Numeric, K1 comparable, K2 comparable](matrix Matrix[T], rowKeys []K1, colKeys []K2) (*KeyedMatrix[T, K1, K2], error) {

	if matrix.Rows != len(rowKeys) {
		return nil, errors.New("Matrix rows amount does't equal to row labels length")
	}

	if matrix.Cols != len(colKeys) {
		return nil, errors.New("Matrix columns amount does't equal to column labels length")
	}

	km := &KeyedMatrix[T, K1, K2]{

		matrix:        matrix,
		RowKeyToIndex: make(map[K1]int),
		ColKeyToIndex: make(map[K2]int),
		RowKeys:       make([]K1, matrix.Rows),
		ColKeys:       make([]K2, matrix.Cols)}

	for i, rl := range rowKeys {
		km.RowKeyToIndex[rl] = i
		km.RowKeys[i] = rl
	}

	for i, cl := range colKeys {
		km.ColKeyToIndex[cl] = i
		km.ColKeys[i] = cl
	}

	return km, nil
}

func (im *KeyedMatrix[T, K1, K2]) RowsN() int {
	return im.matrix.Rows
}

func (im *KeyedMatrix[T, K1, K2]) ColsN() int {
	return im.matrix.Cols
}

func (lm *KeyedMatrix[T, K1, K2]) Get(row_n int, col_n int) T {

	return lm.matrix.Get(row_n, col_n)
}

func (lm *KeyedMatrix[T, K1, K2]) GetByKey(row_key K1, col_key K2) T {
	return lm.matrix.Get(lm.RowKeyToIndex[row_key], lm.ColKeyToIndex[col_key])
}

func (lm *KeyedMatrix[T, K1, K2]) Set(row_n int, col_n int, val T) {

	lm.matrix.Set(row_n, col_n, val)
}

func (lm *KeyedMatrix[T, K1, K2]) SetByKey(row_key K1, col_key K2, val T) {
	lm.matrix.Set(lm.RowKeyToIndex[row_key], lm.ColKeyToIndex[col_key], val)
}

func (lm *KeyedMatrix[T, K1, K2]) GetRow(row_n int) []T {
	return lm.matrix.GetRow(row_n)
}

func (lm *KeyedMatrix[T, K1, K2]) GetRowByKey(row_key K1) []T {
	return lm.matrix.GetRow(lm.RowKeyToIndex[row_key])
}

func (lm *KeyedMatrix[T, K1, K2]) GetCol(col_n int) []T {
	return lm.matrix.GetCol(col_n)
}

func (lm *KeyedMatrix[T, K1, K2]) GetColByKey(col_key K2) []T {
	return lm.matrix.GetCol(lm.ColKeyToIndex[col_key])
}
