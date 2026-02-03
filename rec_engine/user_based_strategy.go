package rec_engine

import (
	"fmt"
	"math"

	. "github.com/PetrDoroshev/RS/matrix"
	"github.com/PetrDoroshev/RS/utils"
)

type UserBasedStrategy struct{}

func (s UserBasedStrategy) BuildSimilarityMatrix(objects_to_comp []User, preferenceMatrix *KeyedMatrix[float64, Item, User]) *KeyedMatrix[float64, User, User] {

	similarityMatrix, _ := NewKeyedMatrix(*NewZeroMatrix[float64](len(objects_to_comp), len(objects_to_comp)),
		objects_to_comp,
		objects_to_comp,
	)

	for i, user_1 := range objects_to_comp {

		for k := i + 1; k < len(objects_to_comp); k++ {

			user_2 := objects_to_comp[k]

			similarity, _ := utils.CosSimilarity(preferenceMatrix.GetColByKey(user_1),
				preferenceMatrix.GetColByKey(user_2))

			similarityMatrix.Set(i, k, similarity)
			similarityMatrix.Set(k, i, similarity)
		}
	}

	return similarityMatrix
}

func (s UserBasedStrategy) PredictRating(recEngine *RecEngine[User], target_user User, target_item Item, output bool) float64 {

	var rating float64

	target_user_index := recEngine.PreferenceMatrix.ColKeyToIndex[target_user]
	target_item_index := recEngine.PreferenceMatrix.RowKeyToIndex[target_item]

	users_to_comp := make([]User, 0, recEngine.PreferenceMatrix.ColsN())

	for col_n := range recEngine.PreferenceMatrix.ColsN() {

		if recEngine.PreferenceMatrix.Get(target_item_index, col_n) != 0 || col_n == target_user_index {

			users_to_comp = append(users_to_comp, recEngine.PreferenceMatrix.ColKeys[col_n])
		}
	}

	similarityMatrix := s.BuildSimilarityMatrix(users_to_comp, &recEngine.PreferenceMatrix)

	if output {
		fmt.Println("\nМатрица подобия:")
		PrintSimilarityMatrix(similarityMatrix)
	}

	nearest_neighbours := []User{}
	similarity_threshold := 0.65

	for i, dist := range similarityMatrix.GetRowByKey(target_user) {

		u := similarityMatrix.RowKeys[i]

		if dist >= similarity_threshold && u != target_user {
			nearest_neighbours = append(nearest_neighbours, u)
		}

	}

	if output {
		fmt.Println("\nБлижайшие соседи:")
		fmt.Println(nearest_neighbours)
	}

	target_user_avg_rating := recEngine.AvgUserRating(target_user)
	sum_of_dist := 0.0
	sum_of_rating_diff := 0.0

	for _, u := range nearest_neighbours {

		user_avg_rating := recEngine.AvgUserRating(u)

		sum_of_rating_diff += (recEngine.PreferenceMatrix.GetByKey(target_item, u) - user_avg_rating) * similarityMatrix.GetByKey(target_user, u)
		sum_of_dist += math.Abs(similarityMatrix.GetByKey(target_user, u))

	}

	rating = target_user_avg_rating + (sum_of_rating_diff / sum_of_dist)

	return rating
}
