# RecSystems

A lightweight recommendation engine written in Go. The project demonstrates two classic collaborative-filtering approaches based on cosine similarity:

- **User-based filtering** — finds users with similar preferences and uses their ratings to predict missing ones.
- **Item-based filtering** — finds similar items and uses the target user's ratings to build recommendations.

The repository also contains a generic matrix package with labeled matrices and sparse representations used by the recommendation engine.

## Features

- Generic numeric matrices implemented with Go generics.
- Key-based access to preference matrices through `User` and `Item` labels.
- Cosine similarity for users and items.
- Rating prediction for a selected user-item pair.
- Top-N recommendations.
- Threshold-based recommendations.
- Conversion of dense matrices to COO, CSR, and ELLPACK formats.

## Project structure

```text
.
├── matrix/                 # Generic and labeled matrix types
├── rec_engine/             # Recommendation engine and filtering strategies
├── utils/                  # Similarity and vector utilities
├── task_4/                 # User-based filtering example
├── task_5/                 # Item-based filtering example
└── go.mod
```

## Requirements

- Go 1.23 or newer

## Getting started

Clone the repository and run the examples:

```bash
git clone https://github.com/PetrDoroshev/RecSystems.git
cd RecSystems

go run ./task_4  # user-based recommendations
go run ./task_5  # item-based recommendations
```

The examples print the preference matrix, predicted ratings, and generated recommendations to the console.

To verify all packages compile:

```bash
go test ./...
```

## Usage

The engine expects a preference matrix where rows represent items, columns represent users, and `0` means that no rating is available.

```go
package main

import (
	"fmt"

	"github.com/PetrDoroshev/RS/matrix"
	"github.com/PetrDoroshev/RS/rec_engine"
)

func main() {
	users := []rec_engine.User{{Id: 1}, {Id: 2}}
	items := []rec_engine.Item{{Id: 1}, {Id: 2}, {Id: 3}}
	ratings := [][]float64{
		{5, 4},
		{4, 0},
		{0, 5},
	}

	preferenceMatrix, err := matrix.NewKeyedMatrix(
		*matrix.NewMatrix(ratings),
		items,
		users,
	)
	if err != nil {
		panic(err)
	}

	engine := rec_engine.NewRecEngine(
		*preferenceMatrix,
		rec_engine.UserBasedStrategy{},
	)

	recommendations := engine.MakeRecommendationTopN(users[0], 5)
	for _, recommendation := range recommendations {
		fmt.Printf("recommend %s with predicted rating %.2f\n",
			recommendation.Item,
			recommendation.Rating,
		)
	}
}
```

For item-based filtering, use `rec_engine.ItemBasedStrategy{}` when creating the engine:

```go
engine := rec_engine.NewRecEngine(
	*preferenceMatrix,
	rec_engine.ItemBasedStrategy{},
)
```

## Recommendation methods

```go
// Return up to N unrated items sorted by predicted rating.
recommendations := engine.MakeRecommendationTopN(user, 5)

// Return unrated items whose predicted rating reaches the threshold.
recommendations := engine.MakeRecommendationTHD(user, 4.0)
```

Individual predictions can be calculated with:

```go
predictedRating := engine.PredictRating(user, item, false)
```

Set the third argument to `true` to print the similarity matrix and nearest neighbours used for the prediction.

## Notes

This is an educational prototype intended to demonstrate the core mechanics of collaborative filtering. It currently works with in-memory rating matrices and does not include data loading, persistence, model training, or a production API.

