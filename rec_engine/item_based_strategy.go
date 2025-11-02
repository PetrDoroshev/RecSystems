package rec_engine

import (
	. "github.com/PetrDoroshev/RS/matrix"
)

type ItemBasedStrategy struct{}

func (s ItemBasedStrategy) BuildSimilarityMatrix(objects_to_comp []Item, preferenceMatrix *KeyedMatrix[float64, Item, User]) *KeyedMatrix[float64, Item, Item] {

}

func (s ItemBasedStrategy) PredictRating(recEngine *RecEngine[Item], target_user User, target_item Item, output bool) float64 {

}
