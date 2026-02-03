package main

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"log"
	"math"
	"os"
	"strconv"
	"strings"
)

func printPreferenceMatrix(matrix [][]int) {

	fmt.Print("\t")
	for i := range len(matrix[0]) {
		fmt.Printf("U%d\t", i+1)
	}

	fmt.Print("\n")

	for i, row := range matrix {

		fmt.Printf("P%d\t", i+1)

		for _, n := range row {
			fmt.Printf("%d\t", n)
		}
		fmt.Println()
	}
}

func transposeMatrix(matrix [][]int) [][]int {

	rows := len(matrix)
	cols := len(matrix[0])
	transposed_matrix := make([][]int, cols)

	for i := range transposed_matrix {
		transposed_matrix[i] = make([]int, rows)
	}

	for i := 0; i < rows; i++ {
		for j := 0; j < cols; j++ {
			transposed_matrix[j][i] = matrix[i][j]
		}
	}
	return transposed_matrix
}

func getLength(vector []int) float64 {

	sum := 0

	for _, n := range vector {
		sum += n * n
	}

	return math.Sqrt(float64(sum))
}

func dotProduct(vector_1 []int, vector_2 []int) float64 {

	vector_len := min(len(vector_1), len(vector_2))

	dot := 0

	for i := 0; i < vector_len; i++ {
		dot += vector_1[i] * vector_2[i]
	}

	return float64(dot)

}

func cosMetric(v1 []int, v2 []int) (float64, error) {

	v1_len := getLength(v1)
	v2_len := getLength(v2)

	if v1_len == 0 || v2_len == 0 {
		return math.NaN(), errors.New("vector's length cannot be equals 0")
	}

	return dotProduct(v1, v2) / (v1_len * v2_len), nil
}

func readMatrixFromFile(filename string) (matrix [][]int) {

	file, err := os.Open(filename)

	if err != nil {
		log.Fatal(err)
	}

	defer file.Close()

	scanner := bufio.NewScanner(file)
	for row := 0; scanner.Scan(); row++ {

		items := strings.Fields(scanner.Text())

		matrix = append(matrix, make([]int, len(items)))

		for i, n := range items {

			int_n, _ := strconv.Atoi(n)
			matrix[row][i] = int_n
		}
	}

	return matrix
}

func getClosestItems(matrix [][]int) (item_1 int, item_2 int, max_dist float64) {

	max_dist = -1.0

	for i := 0; i < len(matrix); i++ {
		for k := i + 1; k < len(matrix); k++ {

			dist, _ := cosMetric(matrix[i], matrix[k])

			//fmt.Printf("%d, %d = %f\n", i+1, k+1, dist)

			if dist > max_dist {

				max_dist = dist
				item_1 = i
				item_2 = k
			}
		}
	}
	return item_1, item_2, max_dist
}

func getClosestUsers(matrix [][]int) (user_1 int, user_2 int, dist float64) {

	transposed_matrix := transposeMatrix(matrix)
	return getClosestItems(transposed_matrix)
}

func main() {

	matrix_filename_ptr := flag.String("matrix_file", "./matrix.txt", "matrix file path arg")

	flag.Parse()

	matrix := readMatrixFromFile(*matrix_filename_ptr)

	fmt.Println("Матрица предпочтений:")
	printPreferenceMatrix(matrix)

	item_1, item_2, dist := getClosestItems(matrix)
	fmt.Printf("\nБлижайшие продукты: P%d, P%d, Расстояние: %f\n", item_1+1, item_2+1, dist)

	user_1, user_2, dist := getClosestUsers(matrix)
	fmt.Printf("Ближайшие пользователи: U%d, U%d, Расстояние: %f\n", user_1+1, user_2+1, dist)

}
