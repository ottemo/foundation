package review

import (
	"time"

	"github.com/ottemo/foundation/api"
	"github.com/ottemo/foundation/app/models/product"
	"github.com/ottemo/foundation/app/models/visitor"
	"github.com/ottemo/foundation/db"
	"github.com/ottemo/foundation/env"
	"github.com/ottemo/foundation/utils"
)

// setupAPI setups package related API endpoint routines
func setupAPI() error {

	var err error

	err = api.GetRestService().RegisterAPI("product/:productID/reviews", api.GET, APIListProductReviews)
	if err != nil {
		return env.ErrorDispatch(err)
	}
	err = api.GetRestService().RegisterAPI("product/:productID/review", api.POST, APICreateProductReview)
	if err != nil {
		return env.ErrorDispatch(err)
	}
	err = api.GetRestService().RegisterAPI("product/:productID/ratedreview/:stars", api.POST, APICreateProductReview)
	if err != nil {
		return env.ErrorDispatch(err)
	}
	err = api.GetRestService().RegisterAPI("product/:productID/review/:reviewID", api.DELETE, APIDeleteProductReview)
	if err != nil {
		return env.ErrorDispatch(err)
	}

	err = api.GetRestService().RegisterAPI("product/:productID/rating", api.GET, APIGetProductRating)
	if err != nil {
		return env.ErrorDispatch(err)
	}

	return nil
}

// APIListProductReviews returns a list of reviews for specified products
//   - product id should be specified in "productID" argument
func APIListProductReviews(context api.InterfaceApplicationContext) (interface{}, error) {

	productObject, err := product.LoadProductByID(context.GetRequestArgument("productID"))
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	collection, err := db.GetCollection(ConstReviewCollectionName)
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	collection.AddFilter("product_id", "=", productObject.GetID())
	collection.AddFilter("review", "!=", "")

	records, err := collection.Load()
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	return records, nil
}

// APICreateProductReview creates a new review for a specified product
//   - product id should be specified in "productID" argument
//   - if "stars" argument specified and is not blank rating mark will be also created
func APICreateProductReview(context api.InterfaceApplicationContext) (interface{}, error) {

	visitorObject, err := visitor.GetCurrentVisitor(context)
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	if visitorObject.IsGuest() {
		return nil, env.ErrorNew(ConstErrorModule, ConstErrorLevel, "2e671776-659b-4c1d-8590-a61f00a9d969", "guest visitor is no allowed to add review")
	}

	productObject, err := product.LoadProductByID(context.GetRequestArgument("productID"))
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	reviewCollection, err := db.GetCollection(ConstReviewCollectionName)
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	// rating update if was set
	//-------------------------
	ratingValue := 0
	if starsValue := context.GetRequestArgument("stars"); starsValue != "" {

		starsNum := utils.InterfaceToInt(starsValue)
		if starsNum <= 0 || starsNum > 5 {
			return nil, env.ErrorNew(ConstErrorModule, env.ConstErrorLevelAPI, "1a7d1a3d-aa79-4722-b02b-c030bffb7557", "stars should be value integer beetween 1 and 5")
		}

		reviewCollection.AddFilter("product_id", "=", productObject.GetID())
		reviewCollection.AddFilter("visitor_id", "=", visitorObject.GetID())
		reviewCollection.AddFilter("rating", ">", 0)

		records, err := reviewCollection.Count()
		if err != nil {
			return nil, env.ErrorDispatch(err)
		}

		if records > 0 {
			return nil, env.ErrorNew(ConstErrorModule, env.ConstErrorLevelAPI, "93032088-8575-4754-92cb-0146b9c4fa97", "you have already vote for that product")
		}

		ratingValue = starsNum

		ratingCollection, err := db.GetCollection(ConstRatingCollectionName)
		if err != nil {
			return nil, env.ErrorDispatch(err)
		}

		ratingCollection.AddFilter("product_id", "=", productObject.GetID())
		ratingRecords, err := ratingCollection.Load()
		if err != nil {
			return nil, env.ErrorDispatch(err)
		}

		recordAttribute := "stars_" + utils.InterfaceToString(ratingValue)
		var ratingRecord map[string]interface{}

		if len(ratingRecords) > 0 {
			ratingRecord = ratingRecords[0]

			ratingRecord[recordAttribute] = utils.InterfaceToInt(ratingRecord[recordAttribute]) + 1
		} else {
			ratingRecord = map[string]interface{}{
				"product_id": productObject.GetID(),
				"stars_1":    0,
				"stars_2":    0,
				"stars_3":    0,
				"stars_4":    0,
				"stars_5":    0,
			}

			ratingRecord[recordAttribute] = 1
		}
		ratingCollection.Save(ratingRecord)
	}

	// review add new record
	//----------------------
	storingValues := map[string]interface{}{
		"product_id": productObject.GetID(),
		"visitor_id": visitorObject.GetID(),
		"username":   visitorObject.GetFullName(),
		"rating":     ratingValue,
		"review":     context.GetRequestContent(),
		"created_at": time.Now(),
	}

	newID, err := reviewCollection.Save(storingValues)
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	storingValues["_id"] = newID

	return storingValues, nil
}

// APIDeleteProductReview  deletes existing review
//   - review if should be specified in "reviewID" argiment
func APIDeleteProductReview(context api.InterfaceApplicationContext) (interface{}, error) {

	reviewID := context.GetRequestArgument("reviewID")

	visitorObject, err := visitor.GetCurrentVisitor(context)
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	if visitorObject.IsGuest() {
		return nil, env.ErrorNew(ConstErrorModule, ConstErrorLevel, "2e671776-659b-4c1d-8590-a61f00a9d969", "guest visitor is no allowed to edit review")
	}

	collection, err := db.GetCollection(ConstReviewCollectionName)
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	reviewRecord, err := collection.LoadByID(reviewID)
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	if visitorID, present := reviewRecord["visitor_id"]; present {

		// check rights
		if err := api.ValidateAdminRights(context); err != nil {
			if visitorID != visitorObject.GetID() {
				return nil, env.ErrorDispatch(err)
			}
		}

		// rating update if was set
		//-------------------------
		reviewRating := utils.InterfaceToInt(reviewRecord["rating"])

		if reviewRating > 0 {
			ratingCollection, err := db.GetCollection(ConstRatingCollectionName)
			if err != nil {
				return nil, env.ErrorDispatch(err)
			}

			ratingCollection.AddFilter("product_id", "=", reviewRecord["product_id"])
			ratingRecords, err := ratingCollection.Load()
			if err != nil {
				return nil, env.ErrorDispatch(err)
			}

			var ratingRecord map[string]interface{}

			if len(ratingRecords) > 0 {
				ratingRecord = ratingRecords[0]

				recordAttribute := "stars_" + utils.InterfaceToString(reviewRating)
				ratingRecord[recordAttribute] = utils.InterfaceToInt(ratingRecord[recordAttribute]) - 1
				ratingCollection.Save(ratingRecord)
			}
		}

		// review remove
		//--------------
		collection.DeleteByID(reviewID)
	}

	return "ok", nil
}

// APIGetProductRating returns rating info for a specified product
//   - product id should be specified in "productID" argument
func APIGetProductRating(context api.InterfaceApplicationContext) (interface{}, error) {

	productObject, err := product.LoadProductByID(context.GetRequestArgument("productID"))
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	ratingCollection, err := db.GetCollection(ConstRatingCollectionName)
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	ratingCollection.AddFilter("product_id", "=", productObject.GetID())
	ratingRecords, err := ratingCollection.Load()
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	return ratingRecords, nil
}
