package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

func printMatrix(matrix [][]float64, labels []string) {

	maxLabelLen := 0

	for _, l := range labels {
		if len(l) > maxLabelLen {
			maxLabelLen = len(l)
		}
	}

	cellWidth := 7

	if maxLabelLen > 7 {
		cellWidth = maxLabelLen
	}

	headerFmt := fmt.Sprintf("%%%ds", cellWidth)
	labelFmt := fmt.Sprintf("%%-%ds", cellWidth)
	numFmt := fmt.Sprintf("%%%d.4f", cellWidth)

	fmt.Printf(headerFmt, "")

	for _, label := range labels {
		fmt.Printf(headerFmt, label)
	}
	fmt.Println()

	for i, row := range matrix {

		fmt.Printf(labelFmt, labels[i])

		for k, n := range row {
			
			if i <= k {
				fmt.Printf(numFmt, n)
			} else {
				fmt.Printf(headerFmt, "")
			}
		}
		fmt.Println()
	}
}

func generateLabels(clusters [][]int) []string {

	labels := make([]string, len(clusters))

	for i, cluster := range clusters {

		parts := make([]string, len(cluster))

		for j, idx := range cluster {
			parts[j] = fmt.Sprintf("U%d", idx+1)
		}

		labels[i] = strings.Join(parts, "+")
	}
	return labels
}

func readDistanceMatrixFromFile(filename string) (matrix [][]float64) {

	file, err := os.Open(filename)

	if err != nil {
		log.Fatal(err)
	}

	defer file.Close()

	scanner := bufio.NewScanner(file)
	for row := 0; scanner.Scan(); row++ {

		items := strings.Fields(scanner.Text())
		matrix = append(matrix, make([]float64, len(items)))

		for i, n := range items {

			float_n, err := strconv.ParseFloat(n, 64)

			if err != nil {
				log.Fatal(err)
			}

			matrix[row][i] = float_n
		}
	}

	return matrix
}

func getMax(matrix [][]float64) (row int, col int, max_val float64) {

	max_val = -1.0

	for i := 0; i < len(matrix); i++ {
		for k := i; k < len(matrix); k++ {

			if matrix[i][k] > max_val {

				max_val = matrix[i][k]
				row = i
				col = k
			}
		}
	}
	return row, col, max_val
}

func rebuildMatrix(matrix [][]float64, item_1 int, item_2 int) (new_matrix [][]float64) {

	new_matrix = make([][]float64, len(matrix)-1)

	for i := range new_matrix {

		new_matrix[i] = make([]float64, len(matrix)-1)

		for k := range new_matrix {
			new_matrix[i][k] = 0
		}
	}

	new_col, new_row := 0, 0
	for i := 0; i < len(matrix)-1; i++ {

		if i == item_2 {
			continue
		}

		new_row = i
		if i > item_2 {
			new_row = i - 1
		}

		for k := i + 1; k < len(matrix); k++ {

			new_col = k
			if k > item_2 {
				new_col = k - 1
			}

			if k == item_2 {
				continue
			}

			if k == item_1 {
				new_matrix[new_row][new_col] = max(matrix[i][item_1], matrix[i][item_2])

			} else if i == item_1 {
				new_matrix[new_row][new_col] = max(matrix[item_1][k], matrix[item_2][k])

			} else {
				new_matrix[new_row][new_col] = matrix[i][k]
			}

			new_matrix[new_col][new_row] = new_matrix[new_row][new_col]
		}
	}
	return new_matrix
}

func MakeClusters(matrix [][]float64, clusters [][]int, cluster_size float64) {

	for step := 0; ; step++ {

		item_1, item_2, max_val := getMax(matrix)

		if max_val < cluster_size {
			break
		}

		labels := generateLabels(clusters)
		fmt.Printf("\n=== Шаг %d ===\n", step)
		fmt.Printf("Объединяются кластеры: [%s] + [%s] (сходство = %.4f)\n\n",
			labels[item_1], labels[item_2], max_val)

		clusters[item_1] = append(clusters[item_1], clusters[item_2]...)
		clusters = append(clusters[:item_2], clusters[item_2+1:]...)

		matrix = rebuildMatrix(matrix, item_1, item_2)

		printMatrix(matrix, generateLabels(clusters))
		fmt.Println()
	}
}

func main() {

	cluster_size_ptr := flag.Float64("cluster_size", 0.85, "cluster size arg")
	matrix_filename_ptr := flag.String("matrix_file", "./similarity_matrix.txt", "matrix file path arg")

	flag.Parse()

	matrix := readDistanceMatrixFromFile(*matrix_filename_ptr)

	clusters := make([][]int, len(matrix))
	for i := range clusters {
		clusters[i] = []int{i}
	}

	printMatrix(matrix, generateLabels(clusters))

	fmt.Println()

	MakeClusters(matrix, clusters, *cluster_size_ptr)

}
