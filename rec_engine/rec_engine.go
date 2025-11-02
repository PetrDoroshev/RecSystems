package rec_engine

import (
	"fmt"
	"sort"

	. "github.com/PetrDoroshev/RS/matrix"
)

type User struct {
	Id int
}

func (u User) String() string {
	return fmt.Sprintf("U%d", u.Id)
}

type Item struct {
	Id   int
	Name string
}

func (it Item) String() string {
	return fmt.Sprintf("P%d", it.Id)
}

type ItemRating struct {
	Item   Item
	Rating float64
}

type Key interface {
	User | Item
}

type similarityStrategy[T Key] interface {
	BuildSimilarityMatrix(objects_to_comp []T, preferenceMatrix *KeyedMatrix[float64, Item, User]) *KeyedMatrix[float64, T, T]
	PredictRating(recEngine *RecEngine[T], target_user User, target_item Item, output bool) float64
}

type RecEngine[T Key] struct {
	PreferenceMatrix KeyedMatrix[float64, Item, User]
	Strategy         similarityStrategy[T]
}

func NewRecEngine[T Key](preferenceMatrix KeyedMatrix[float64, Item, User], strategy similarityStrategy[T]) *RecEngine[T] {

	return &RecEngine[T]{PreferenceMatrix: preferenceMatrix, Strategy: strategy}
}

func (re *RecEngine[T]) AvgItemRating(item Item) float64 {

	sum := 0.0
	n := 0

	item_index := re.PreferenceMatrix.RowKeyToIndex[item]

	for col_n := range re.PreferenceMatrix.ColsN() {

		rating := re.PreferenceMatrix.Get(item_index, col_n)
		if rating > 0 {
			sum += rating
			n++
		}
	}

	return sum / float64(n)

}

func (re *RecEngine[T]) AvgUserRating(user User) float64 {

	n := 0
	sum := 0.0

	user_index := re.PreferenceMatrix.ColKeyToIndex[user]

	for row_n := range re.PreferenceMatrix.RowsN() {

		rating := re.PreferenceMatrix.Get(row_n, user_index)

		if rating != 0 {
			sum += rating
			n++
		}
	}

	if n == 0 {
		return 0.0
	}

	return sum / float64(n)

}

func (re *RecEngine[T]) PredictRating(target_user User, target_item Item, output bool) float64 {

	return re.Strategy.PredictRating(re, target_user, target_item, output)
	/*
		var rating float64

		target_user_index := re.PreferenceMatrix.ColKeyToIndex[target_user]
		target_item_index := re.PreferenceMatrix.RowKeyToIndex[target_item]

		users_to_comp := make([]User, 0, re.PreferenceMatrix.ColsN())

		for col_n := range re.PreferenceMatrix.ColsN() {

			if re.PreferenceMatrix.Get(target_item_index, col_n) != 0 || col_n == target_user_index {

				users_to_comp = append(users_to_comp, re.PreferenceMatrix.ColKeys[col_n])
			}
		}

		similarityMatrix := re.buildSimilarityMatrix(users_to_comp)

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

		target_user_avg_rating := re.AvgUserRating(target_user)
		sum_of_dist := 0.0
		sum_of_rating_diff := 0.0

		for _, u := range nearest_neighbours {

			user_avg_rating := re.AvgUserRating(u)

			sum_of_rating_diff += (re.PreferenceMatrix.GetByKey(target_item, u) - user_avg_rating) * similarityMatrix.GetByKey(target_user, u)
			sum_of_dist += math.Abs(similarityMatrix.GetByKey(target_user, u))

		}

		rating = target_user_avg_rating + (sum_of_rating_diff / sum_of_dist)

		return rating
	*/
}

func (re *RecEngine[T]) getItemPredictedRatings(user User) []ItemRating {

	recommendations := make([]ItemRating, 0, re.PreferenceMatrix.RowsN())

	for _, item := range re.PreferenceMatrix.RowKeys {

		rating := re.PreferenceMatrix.GetByKey(item, user)

		if rating == 0 {

			predicted_rating := re.PredictRating(user, item, false)
			recommendations = append(recommendations, ItemRating{Item: item, Rating: predicted_rating})
		}
	}

	return recommendations
}

func (re *RecEngine[T]) MakeRecommendationTHD(user User, threshold float64) []ItemRating {

	var recommendations []ItemRating

	if re.AvgUserRating(user) == 0 {

		fmt.Println(re.AvgUserRating(user))

		for _, item := range re.PreferenceMatrix.RowKeys {
			recommendations = append(recommendations, ItemRating{Item: item, Rating: re.AvgItemRating(item)})
		}
		fmt.Println(recommendations)

	} else {
		recommendations = re.getItemPredictedRatings(user)
	}

	sort.Slice(recommendations, func(i, j int) bool {
		return recommendations[i].Rating > recommendations[j].Rating
	})

	n := 0

	for n < len(recommendations) && recommendations[n].Rating >= threshold {
		n++
	}

	return recommendations[:n]

}

func (re *RecEngine[T]) MakeRecommendationTopN(user User, N int) []ItemRating {

	var recommendations []ItemRating

	if re.AvgUserRating(user) == 0 {

		for _, item := range re.PreferenceMatrix.RowKeys {
			recommendations = append(recommendations, ItemRating{Item: item, Rating: re.AvgItemRating(item)})
		}

	} else {
		recommendations = re.getItemPredictedRatings(user)
	}

	sort.Slice(recommendations, func(i, j int) bool {
		return recommendations[i].Rating > recommendations[j].Rating
	})

	return recommendations[:min(N, len(recommendations))]
}

func PrintPreferenceMatrix[T Numeric](preferenceMatrix *KeyedMatrix[T, Item, User]) {

	fmt.Print("\t")
	for col_index := range preferenceMatrix.ColsN() {
		fmt.Printf("U%d\t", preferenceMatrix.ColKeys[col_index].Id)
	}

	fmt.Print("\n")

	for row_index := range preferenceMatrix.RowsN() {

		fmt.Printf("P%d\t", preferenceMatrix.RowKeys[row_index].Id)

		for col_index := range preferenceMatrix.ColsN() {
			fmt.Printf("%.2v\t", preferenceMatrix.Get(row_index, col_index))
		}
		fmt.Println()
	}
}

func PrintSimilarityMatrix[T Numeric](similarityMatrix *KeyedMatrix[T, User, User]) {

	fmt.Print("\t")
	for col_index := range similarityMatrix.ColsN() {
		fmt.Printf("U%d\t", similarityMatrix.ColKeys[col_index].Id)
	}

	fmt.Print("\n")

	for row_index := range similarityMatrix.RowsN() {

		fmt.Printf("U%d\t", similarityMatrix.RowKeys[row_index].Id)

		for col_index := range similarityMatrix.ColsN() {
			fmt.Printf("%.2v\t", similarityMatrix.Get(row_index, col_index))
		}
		fmt.Println()
	}
}
