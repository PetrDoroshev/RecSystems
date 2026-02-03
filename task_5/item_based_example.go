package main

import (
	"fmt"

	"github.com/PetrDoroshev/RS/matrix"
	"github.com/PetrDoroshev/RS/rec_engine"
)

func main() {

	users := make([]rec_engine.User, 11)
	items := make([]rec_engine.Item, 10)

	for i := range len(users) {
		users[i] = rec_engine.User{Id: i + 1}
	}

	for i := range len(items) {
		items[i] = rec_engine.Item{Id: i + 1, Name: ""}
	}

	m := [][]float64{

		{5, 4, 0, 0, 4, 2, 5, 4, 5, 3, 0},
		{4, 5, 0, 3, 5, 0, 4, 3, 4, 5, 0},
		{3, 4, 0, 4, 5, 3, 4, 4, 3, 5, 0},
		{4, 5, 4, 5, 0, 4, 5, 5, 4, 5, 0},
		{5, 4, 5, 4, 5, 4, 0, 3, 5, 4, 0},
		{0, 4, 4, 5, 4, 3, 4, 0, 4, 5, 0},
		{4, 5, 5, 4, 4, 5, 4, 5, 0, 4, 0},
		{0, 4, 3, 4, 5, 4, 3, 4, 4, 0, 0},
		{4, 5, 0, 3, 4, 5, 4, 5, 5, 4, 0},
		{0, 4, 5, 4, 3, 4, 5, 4, 4, 5, 0}}

	preferenceMatrix, err := matrix.NewKeyedMatrix(*matrix.NewMatrix(m), items, users)

	if err != nil {
		fmt.Println(err.Error())
	}

	user := users[2]
	//item := items[2]

	rec_engine.PrintPreferenceMatrix(preferenceMatrix)

	re := rec_engine.RecEngine[rec_engine.Item]{PreferenceMatrix: *preferenceMatrix, Strategy: rec_engine.ItemBasedStrategy{}}
	//rating := re.PredictRating(user, item, true)

	//fmt.Printf("\nПредстказанный рейтинг товара %s от пользователя %s: %f\n", item, user, rating)

	recommedations := re.MakeRecommendationTHD(user, 2.0)

	fmt.Printf("\nРекомендации для пользователя %s:\n", user)
	for _, rec := range recommedations {

		fmt.Printf("%s рекомендовать %s, предсказанный рейтинг: %f\n", user, rec.Item, rec.Rating)

	}

	user = users[10]

	recommedations = re.MakeRecommendationTHD(user, 4.0)

	fmt.Printf("\nРекомендации для пользователя %s (товары с наибольшим рейтингом > 4):\n", user)
	for _, rec := range recommedations {

		fmt.Printf("%s рекомендовать %s, предсказанный рейтинг: %f\n", user, rec.Item, rec.Rating)

	}

}
