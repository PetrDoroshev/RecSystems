package main

import (
	_ "errors"
	"fmt"
	"log"
	_ "math"

	"github.com/PetrDoroshev/RS/matrix"
)

func printPreferenceMatrix(matrix *matrix.Matrix[int], item_labels []string, user_labels []string) {

	fmt.Print("\t")
	for _, user_label := range user_labels {
		fmt.Printf("%s\t", user_label)
	}

	fmt.Print("\n")

	for row_n := range matrix.Rows {

		fmt.Printf("%s\t", item_labels[row_n])

		for col_n := range matrix.Cols {
			fmt.Printf("%d\t", matrix.Get(row_n, col_n))
		}
		fmt.Println()
	}
}

func meanRating[T matrix.Numeric](data []T) float64 {

	var sum T
	non_zero_items := 0

	for _, n := range data {

		if n > 0 {
			sum += n
			non_zero_items++
		}
	}

	if non_zero_items == 0 {
		return 0.0
	}

	return float64(sum) / float64(non_zero_items)

}


func main() {

	m := matrix.NewMatrix([][]int{

		{1, 0, 0, 0, 2, 0},
		{0, 0, 3, 4, 0, 0},
		{0, 0, 0, 0, 5, 9},
		{0, 0, 0, 8, 0, 5},
		{0, 0, 0, 0, 0, 0},
		{0, 7, 1, 0, 0, 6}})

	cl := m.ToCoordinates()
	csr := m.ToCSR()
	ellpack := m.ToELLPACK()

	fmt.Println("Матрица:")
	fmt.Println(m.String())

	fmt.Println()
	fmt.Println(cl)

	fmt.Println()
	fmt.Println(csr)

	fmt.Println()
	fmt.Println(ellpack)

	preference_matrix := matrix.NewMatrix([][]int{
		{5, 4, 5, 0, 5},
		{5, 5, 5, 0, 4},
		{0, 0, 0, 0, 0},
		{5, 5, 5, 0, 5},
		{0, 0, 0, 0, 0},
		{5, 4, 5, 0, 4},
	})

	user_labels := make([]string, preference_matrix.Cols)
	item_labels := make([]string, preference_matrix.Rows)

	for i := 0; i < preference_matrix.Cols; i++ {
		user_labels[i] = fmt.Sprintf("U%d", i+1)
	}

	for i := 0; i < preference_matrix.Rows; i++ {
		item_labels[i] = fmt.Sprintf("P%d", i+1)
	}

	fmt.Println("Матрица предпочтений:")
	printPreferenceMatrix(preference_matrix, item_labels, user_labels)

	items_mean_ratings := make([]float64, preference_matrix.Rows)

	for row_n := range preference_matrix.Rows {

		mean_val := meanRating(preference_matrix.GetRow(row_n))
		items_mean_ratings[row_n] = mean_val
	}

	threshold := 3.9

	fmt.Println("\nСредние оценки товаров:")
	fmt.Println(items_mean_ratings)
	fmt.Println()

	for i := 0; i < len(items_mean_ratings); i++ {

		if items_mean_ratings[i] < threshold {

			err := preference_matrix.DeleteRow(i)

				if (err != nil) {
				log.Fatalf("%d out of range", i)
			}

			items_mean_ratings = append(items_mean_ratings[:i], items_mean_ratings[i+1:]...)
			item_labels = append(item_labels[:i], item_labels[i+1:]...)

			i--
		}
	}

	for col_n := 0; col_n < preference_matrix.Cols; col_n++ {

		var row_n int

		for row_n = range preference_matrix.Rows {

			if preference_matrix.Get(row_n, col_n) != 0 {
				break
			}
		}

		if row_n == preference_matrix.Rows - 1 {

			err := preference_matrix.DeleteColumn(col_n)

			if (err != nil) {
				log.Fatalf("%d out of range", row_n)
			}

			user_labels = append(user_labels[:col_n], user_labels[col_n+1:]...)
			col_n--
		}
	}

	fmt.Println("\nОбновленная матрица предпочтений:")
	printPreferenceMatrix(preference_matrix, item_labels, user_labels)

}
