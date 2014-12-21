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

	err = api.GetRestService().RegisterAPI("product", "GET", "review/list/:pid", restReviewList)
	if err != nil {
		return env.ErrorDispatch(err)
	}
	err = api.GetRestService().RegisterAPI("product", "POST", "review/add/:pid", restReviewAdd)
	if err != nil {
		return env.ErrorDispatch(err)
	}
	err = api.GetRestService().RegisterAPI("product", "POST", "review/add/:pid/:stars", restReviewAdd)
	if err != nil {
		return env.ErrorDispatch(err)
	}
	err = api.GetRestService().RegisterAPI("product", "DELETE", "review/remove/:reviewID", restReviewRemove)
	if err != nil {
		return env.ErrorDispatch(err)
	}
	err = api.GetRestService().RegisterAPI("product", "GET", "rating/info/:pid", restRatingInfo)
	if err != nil {
		return env.ErrorDispatch(err)
	}

	return nil
}

// WEB REST API function used to get list of reviews for particular product
func restReviewList(params *api.StructAPIHandlerParams) (interface{}, error) {

	productObject, err := product.LoadProductByID(params.RequestURLParams["pid"])
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

// WEB REST API function used to add new review for a product
func restReviewAdd(params *api.StructAPIHandlerParams) (interface{}, error) {

	visitorObject, err := visitor.GetCurrentVisitor(params)
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	if visitorObject.IsGuest() {
		return nil, env.ErrorNew(ConstErrorModule, ConstErrorLevel, "2e671776659b4c1d8590a61f00a9d969", "guest visitor is no allowed to add review")
	}

	productObject, err := product.LoadProductByID(params.RequestURLParams["pid"])
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
	if starsValue, present := params.RequestURLParams["stars"]; present {

		starsNum := utils.InterfaceToInt(starsValue)
		if starsNum <= 0 || starsNum > 5 {
			return nil, env.ErrorNew(ConstErrorModule, env.ConstErrorLevelAPI, "1a7d1a3daa794722b02bc030bffb7557", "stars should be value integer beetween 1 and 5")
		}

		reviewCollection.AddFilter("product_id", "=", productObject.GetID())
		reviewCollection.AddFilter("visitor_id", "=", visitorObject.GetID())
		reviewCollection.AddFilter("rating", ">", 0)

		records, err := reviewCollection.Count()
		if err != nil {
			return nil, env.ErrorDispatch(err)
		}

		if records > 0 {
			return nil, env.ErrorNew(ConstErrorModule, env.ConstErrorLevelAPI, "930320888575475492cb0146b9c4fa97", "you have already vote for that product")
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
		"review":     params.RequestContent,
		"created_at": time.Now(),
	}

	newID, err := reviewCollection.Save(storingValues)
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	storingValues["_id"] = newID

	return storingValues, nil
}

// WEB REST API function used to remove review for a product
func restReviewRemove(params *api.StructAPIHandlerParams) (interface{}, error) {

	reviewID := params.RequestURLParams["reviewID"]

	visitorObject, err := visitor.GetCurrentVisitor(params)
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	if visitorObject.IsGuest() {
		return nil, env.ErrorNew(ConstErrorModule, ConstErrorLevel, "2e671776659b4c1d8590a61f00a9d969", "guest visitor is no allowed to edit review")
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
		if err := api.ValidateAdminRights(params); err != nil {
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

// WEB REST API function used to get product rating info
func restRatingInfo(params *api.StructAPIHandlerParams) (interface{}, error) {

	productObject, err := product.LoadProductByID(params.RequestURLParams["pid"])
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
