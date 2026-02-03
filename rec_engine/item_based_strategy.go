package rec_engine

import (
	"fmt"
	"math"

	. "github.com/PetrDoroshev/RS/matrix"
	"github.com/PetrDoroshev/RS/utils"
)

type ItemBasedStrategy struct{}

func (s ItemBasedStrategy) BuildSimilarityMatrix(objects_to_comp []Item, preferenceMatrix *KeyedMatrix[float64, Item, User]) *KeyedMatrix[float64, Item, Item] {

	similarityMatrix, _ := NewKeyedMatrix(*NewZeroMatrix[float64](len(objects_to_comp), len(objects_to_comp)),
		objects_to_comp,
		objects_to_comp,
	)

	for i, item_1 := range objects_to_comp {

		for k := i + 1; k < len(objects_to_comp); k++ {

			item_2 := objects_to_comp[k]

			similarity, _ := utils.CosSimilarity(preferenceMatrix.GetRowByKey(item_1),
				preferenceMatrix.GetRowByKey(item_2))

			similarityMatrix.Set(i, k, similarity)
			similarityMatrix.Set(k, i, similarity)
		}
	}

	return similarityMatrix
}

func (s ItemBasedStrategy) PredictRating(recEngine *RecEngine[Item], target_user User, target_item Item, output bool) float64 {

	similarityMatrix := s.BuildSimilarityMatrix(recEngine.PreferenceMatrix.RowKeys, &recEngine.PreferenceMatrix)

	if output {
		fmt.Println("\nМатрица подобия:")
		PrintSimilarityMatrix(similarityMatrix)
	}

	nearest_neighbours := []Item{}
	similarity_threshold := 0.85

	for i, dist := range similarityMatrix.GetRowByKey(target_item) {

		i := similarityMatrix.RowKeys[i]

		if dist >= similarity_threshold && i != target_item {

			nearest_neighbours = append(nearest_neighbours, i)
		}

	}

	if output {
		fmt.Printf("\nБлижайшие соседи %s (< %.2f):\n", target_item, similarity_threshold)
		fmt.Println(nearest_neighbours)
	}

	n := 0
	for _, i := range nearest_neighbours {
		if recEngine.PreferenceMatrix.GetByKey(i, target_user) != 0 {
			n++
		}
	}

	sum_of_dist := 0.0
	sum_of_rating := 0.0

	if n > 0 {

		for _, i := range nearest_neighbours {

			if recEngine.PreferenceMatrix.GetByKey(i, target_user) != 0 {
				sum_of_rating += recEngine.PreferenceMatrix.GetByKey(i, target_user) * similarityMatrix.GetByKey(target_item, i)
				sum_of_dist += math.Abs(similarityMatrix.GetByKey(target_item, i))
			}

		}
	} else {
		//fmt.Println(2)
		for _, i := range nearest_neighbours {

			users_count := 0

			for _, u := range recEngine.PreferenceMatrix.ColKeys {

				if recEngine.PreferenceMatrix.GetByKey(i, u) != 0 {
					sum_of_rating += recEngine.PreferenceMatrix.GetByKey(i, u) * similarityMatrix.GetByKey(target_item, i)
					users_count++
				}

			}
			sum_of_rating /= float64(users_count)
			sum_of_dist += math.Abs(similarityMatrix.GetByKey(target_item, i))
		}
		//fmt.Println(sum_of_rating)
		//fmt.Println(sum_of_dist)
	}

	return sum_of_rating / sum_of_dist
}
