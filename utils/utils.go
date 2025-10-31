package utils

import (
	"errors"
	"math"
)

func GetVectorLength(vector []float64) float64 {

	sum := 0.0

	for _, n := range vector {
		sum += n * n
	}

	return math.Sqrt(sum)
}

func DotProduct(vector_1 []float64, vector_2 []float64) float64 {

	vector_len := min(len(vector_1), len(vector_2))

	dot := 0.0

	for i := 0; i < vector_len; i++ {
		dot += vector_1[i] * vector_2[i]
	}

	return dot

}

func CosSimilarity(v1 []float64, v2 []float64) (float64, error) {

	v1_len := GetVectorLength(v1)
	v2_len := GetVectorLength(v2)

	if v1_len == 0 || v2_len == 0 {
		return math.NaN(), errors.New("vector's length cannot be equals 0")
	}

	return DotProduct(v1, v2) / (v1_len * v2_len), nil
}
