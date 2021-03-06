package review

import (
	"time"

	"github.com/ottemo/commerce/api"
	"github.com/ottemo/commerce/db"
	"github.com/ottemo/commerce/env"
	"github.com/ottemo/commerce/utils"

	"github.com/ottemo/commerce/app/models"
	"github.com/ottemo/commerce/app/models/product"
	"github.com/ottemo/commerce/app/models/visitor"
)

// setupAPI setups package related API endpoint routines
func setupAPI() error {

	service := api.GetRestService()

	//--------------------------------------------------------------------------------------------------------------
	// In case of api names change - please, fix test.
	//--------------------------------------------------------------------------------------------------------------

	service.POST("review/:productID", APICreateProductReview)
	service.POST("ratedreview/:productID/:stars", APICreateProductReview)

	service.GET("reviews", APIListReviews)
	service.GET("review/:reviewID", APIGetReview)
	service.GET("rating/:productID", APIGetProductRating)

	service.PUT("review/:reviewID", APIUpdateReview)

	service.DELETE("review/:reviewID", APIDeleteProductReview)

	return nil
}

// APIListReviews returns a list of reviews for specified products
//   - product id could be specified in "productID" parameter
func APIListReviews(context api.InterfaceApplicationContext) (interface{}, error) {

	collection, err := db.GetCollection(ConstReviewCollectionName)
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	// product filter, limit
	if err := models.ApplyFilters(context, collection); err != nil {
		return nil, env.ErrorDispatch(err)
	}

	if !api.IsAdminSession(context) {
		visitorObject, err := visitor.GetCurrentVisitor(context)
		if err != nil {
			return nil, env.ErrorDispatch(err)
		}
		if visitorObject.IsGuest() {
			if err := collection.AddFilter("review", "!=", ""); err != nil {
				_ = env.ErrorNew(ConstErrorModule, ConstErrorLevel, "36b0185e-3eed-4f5d-bc08-91d381b0944e", err.Error())
			}
			if err := collection.AddFilter("approved", "=", true); err != nil {
				_ = env.ErrorNew(ConstErrorModule, ConstErrorLevel, "d441a78d-ec6e-4d70-88df-58112552a9b9", err.Error())
			}
		} else {
			collection.SetupFilterGroup("default", false, "")

			collection.SetupFilterGroup("visible", true, "default")
			if err := collection.SetupFilterGroup("visible", true, "default"); err != nil {
				_ = env.ErrorNew(ConstErrorModule, ConstErrorLevel, "a8b90261-bce2-4506-8e5c-39146dcb142f", err.Error())
			}

			if err := collection.SetupFilterGroup("content", false, "visible"); err != nil {
				_ = env.ErrorNew(ConstErrorModule, ConstErrorLevel, "909d45c0-4b43-4025-b5ec-53d1d50167e1", err.Error())
			}
			if err := collection.AddGroupFilter("content", "review", "!=", ""); err != nil {
				_ = env.ErrorNew(ConstErrorModule, ConstErrorLevel, "10955f8e-c56f-4d88-9f5b-4e66add7cb0a", err.Error())
			}
			if err := collection.AddGroupFilter("content", "approved", "=", true); err != nil {
				_ = env.ErrorNew(ConstErrorModule, ConstErrorLevel, "c3ac19ec-c48b-4d28-a764-c8e05de798dd", err.Error())
			}

			if err := collection.AddGroupFilter("visible", "visitor_id", "=", visitorObject.GetID()); err != nil {
				_ = env.ErrorNew(ConstErrorModule, ConstErrorLevel, "c463cd3e-f2fb-449a-9c2e-d055342334e9", err.Error())
			}
		}
	}

	// check "count" request
	if context.GetRequestArgument(api.ConstRESTActionParameter) == "count" {
		return collection.Count()
	}

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
		return nil, env.ErrorNew(ConstErrorModule, ConstErrorLevel, "8d84ba12-5f04-4322-8cb5-c520189cad97", "guest visitor is no allowed to add review")
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

		if err := reviewCollection.AddFilter("product_id", "=", productObject.GetID()); err != nil {
			_ = env.ErrorNew(ConstErrorModule, ConstErrorLevel, "e8b60041-08c3-4483-a747-542305af04bb", err.Error())
		}
		if err := reviewCollection.AddFilter("visitor_id", "=", visitorObject.GetID()); err != nil {
			_ = env.ErrorNew(ConstErrorModule, ConstErrorLevel, "a35c80e1-6608-4f43-97fe-788f81adaa50", err.Error())
		}
		if err := reviewCollection.AddFilter("rating", ">", 0); err != nil {
			_ = env.ErrorNew(ConstErrorModule, ConstErrorLevel, "fa514a76-5c9d-42e7-8f77-0a0fc106cbe4", err.Error())
		}

		records, err := reviewCollection.Count()
		if err != nil {
			return nil, env.ErrorDispatch(err)
		}

		if records > 0 {
			return nil, env.ErrorNew(ConstErrorModule, env.ConstErrorLevelAPI, "93032088-8575-4754-92cb-0146b9c4fa97", "you have already vote for that product")
		}

		ratingValue = starsNum
	}

	requestData, err := api.GetRequestContentAsMap(context)
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	// review add new record
	//----------------------
	storingValues := map[string]interface{}{
		"product_id": productObject.GetID(),
		"visitor_id": visitorObject.GetID(),
		"username":   visitorObject.GetFullName(),
		"rating":     ratingValue,
		"review":     requestData["review"],
		"created_at": time.Now(),
		"approved":   false,
	}

	newID, err := reviewCollection.Save(storingValues)
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	storingValues["_id"] = newID

	return storingValues, nil
}

// APIDeleteProductReview  deletes existing review
//   - review ID should be specified in "reviewID" argument
func APIDeleteProductReview(context api.InterfaceApplicationContext) (interface{}, error) {

	reviewID := context.GetRequestArgument("reviewID")

	var visitorObject visitor.InterfaceVisitor
	if !api.IsAdminSession(context) {
		visitorObject, err := visitor.GetCurrentVisitor(context)
		if err != nil {
			return nil, env.ErrorDispatch(err)
		}
		if visitorObject.IsGuest() {
			return nil, env.ErrorNew(ConstErrorModule, ConstErrorLevel, "2e671776-659b-4c1d-8590-a61f00a9d969", "guest visitor is no allowed to delete review")
		}
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
		if !api.IsAdminSession(context) {
			if visitorID != visitorObject.GetID() {
				return nil, env.ErrorNew(ConstErrorModule, env.ConstErrorLevelAPI, "b4751e17-309d-4f90-a33a-e986c5f2420a", "Operation not allowed.")
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

			if err := ratingCollection.AddFilter("product_id", "=", reviewRecord["product_id"]); err != nil {
				_ = env.ErrorNew(ConstErrorModule, ConstErrorLevel, "717148ed-bb64-4f70-a324-322ad7b20830", err.Error())
			}
			ratingRecords, err := ratingCollection.Load()
			if err != nil {
				return nil, env.ErrorDispatch(err)
			}

			var ratingRecord map[string]interface{}

			if len(ratingRecords) > 0 {
				ratingRecord = ratingRecords[0]

				recordAttribute := "stars_" + utils.InterfaceToString(reviewRating)
				ratingRecord[recordAttribute] = utils.InterfaceToInt(ratingRecord[recordAttribute]) - 1
				if _, err := ratingCollection.Save(ratingRecord); err != nil {
					return nil, env.ErrorNew(ConstErrorModule, ConstErrorLevel, "2a7b8c27-e46b-4a85-9aa8-7b4324529b8a", err.Error())
				}
			}
		}

		// review remove
		//--------------
		if err := collection.DeleteByID(reviewID); err != nil {
			return nil, env.ErrorNew(ConstErrorModule, ConstErrorLevel, "ca07589f-01da-4a86-add0-b67f0a4157bb", err.Error())
		}
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

	if err := ratingCollection.AddFilter("product_id", "=", productObject.GetID()); err != nil {
		_ = env.ErrorNew(ConstErrorModule, ConstErrorLevel, "5da26dc9-f2ad-46be-b6b4-53fd6f4c06a3", err.Error())
	}
	ratingRecords, err := ratingCollection.Load()
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	return ratingRecords, nil
}

// APIUpdateReview updates an existing review
//   - review ID should be specified in "reviewID" argument
func APIUpdateReview(context api.InterfaceApplicationContext) (interface{}, error) {

	// admin or visitor
	var visitorObject visitor.InterfaceVisitor
	if !api.IsAdminSession(context) {
		var err error
		if visitorObject, err = visitor.GetCurrentVisitor(context); err != nil {
			return nil, env.ErrorDispatch(err)
		}

		if visitorObject.IsGuest() {
			return nil, env.ErrorNew(ConstErrorModule, ConstErrorLevel, "3de19977-8e20-484c-9ed0-7868118d767f", "guest visitor is no allowed to update review")
		}
	}

	reviewID := context.GetRequestArgument("reviewID")

	requestData, err := api.GetRequestContentAsMap(context)
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	reviewCollection, err := db.GetCollection(ConstReviewCollectionName)
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	record, err := reviewCollection.LoadByID(reviewID)
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	if visitorObject != nil && visitorObject.GetID() != "" && visitorObject.GetID() != record["visitor_id"] {
		return nil, env.ErrorNew(ConstErrorModule, ConstErrorLevel, "ba19a94b-088c-4a28-861c-6fe2145f2348", "you not allowed to update review")
	}

	if api.IsAdminSession(context) {
		if record["approved"] != context.GetRequestArgument("approved") {
			ratingValue := utils.InterfaceToInt(record["rating"])

			if 1 <= ratingValue && ratingValue <= 5 {
				var diff = -1
				if context.GetRequestArgument("approved") == "true" {
					diff = 1
				}

				ratingCollection, err := db.GetCollection(ConstRatingCollectionName)
				if err != nil {
					return nil, env.ErrorDispatch(err)
				}

				if err := ratingCollection.AddFilter("product_id", "=", record["product_id"]); err != nil {
					return nil, env.ErrorNew(ConstErrorModule, ConstErrorLevel, "6d0a0ee0-3afb-4f0e-8a15-8a095003e510", err.Error())
				}
				ratingRecords, err := ratingCollection.Load()
				if err != nil {
					return nil, env.ErrorDispatch(err)
				}

				recordAttribute := "stars_" + utils.InterfaceToString(ratingValue)
				var ratingRecord map[string]interface{}

				if len(ratingRecords) > 0 {
					ratingRecord = ratingRecords[0]

					ratingRecord[recordAttribute] = utils.InterfaceToInt(ratingRecord[recordAttribute]) + diff
				} else {
					ratingRecord = map[string]interface{}{
						"product_id": record["product_id"],
						"stars_1":    0,
						"stars_2":    0,
						"stars_3":    0,
						"stars_4":    0,
						"stars_5":    0,
					}

					if diff > 0 {
						ratingRecord[recordAttribute] = 1
					}
				}
				if _, err := ratingCollection.Save(ratingRecord); err != nil {
					return nil, env.ErrorNew(ConstErrorModule, ConstErrorLevel, "956ff15f-681b-415b-9c76-2320df0ebee2", "unable to save rating record")
				}
			}
		}
	} else { // not admin
		record["approved"] = false
	}

	for attrName := range record {
		if value, present := requestData[attrName]; present {
			record[attrName] = value
		}
	}

	if _, err := reviewCollection.Save(record); err != nil {
		return nil, env.ErrorDispatch(err)
	}

	return record, nil
}

// APIGetReview returns an existing review for owner or admin
//   - review ID should be specified in "reviewID" argument
func APIGetReview(context api.InterfaceApplicationContext) (interface{}, error) {
	// admin or visitor
	var visitorObject visitor.InterfaceVisitor
	if !api.IsAdminSession(context) {
		var err error
		if visitorObject, err = visitor.GetCurrentVisitor(context); err != nil {
			return nil, env.ErrorDispatch(err)
		}
	}

	reviewID := context.GetRequestArgument("reviewID")

	reviewCollection, err := db.GetCollection(ConstReviewCollectionName)
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	record, err := reviewCollection.LoadByID(reviewID)
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	if visitorObject != nil {
		if !utils.InterfaceToBool(record["approved"]) ||
			utils.InterfaceToString(record["review"]) == "" {
			if visitorObject.GetID() != utils.InterfaceToString(record["visitor_id"]) {
				return nil, env.ErrorNew(ConstErrorModule, ConstErrorLevel, "69578b10-e0c4-412d-aa54-4bacb7277262", "not allowed to get review")
			}
		}
	}

	return record, nil
}
