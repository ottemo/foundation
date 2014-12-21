// Package review is a set of API functions to provide an ability to make reviews for a particular product
package review

import (
	"github.com/ottemo/foundation/env"
)

// Package global constants
const (
	ConstReviewCollectionName = "review"
	ConstRatingCollectionName = "rating"

	ConstErrorModule = "product/review"
	ConstErrorLevel  = env.ConstErrorLevelActor
)
