package rts

import (
	"strings"
	"time"

	"github.com/ottemo/foundation/api"
	"github.com/ottemo/foundation/app"
	"github.com/ottemo/foundation/app/models/cart"
	"github.com/ottemo/foundation/db"
	"github.com/ottemo/foundation/env"
	"github.com/ottemo/foundation/utils"
)

func referrerHandler(event string, eventData map[string]interface{}) bool {

	if _, present := eventData["context"]; present {
		if context, ok := eventData["context"].(api.InterfaceApplicationContext); ok && context != nil {

			xReferrer := utils.InterfaceToString(context.GetRequestSetting("X-Referer"))
			if xReferrer == "" {
				return true
			}

			referrer, err := GetReferrer(xReferrer)
			if err != nil {
				return true
			}

			// excluding itself (i.e. "storefront" requests)
			if strings.Contains(app.GetStorefrontURL(""), referrer) {
				return true
			}

			if _, present := referrers[referrer]; !present {
				updateSync.Lock()
				referrers[referrer] = 0
				updateSync.Unlock()
			}
			referrers[referrer]++

			if err := saveNewReferrer(referrer); err != nil {
				env.LogError(err)
			}
		}
	}

	return true
}

func visitsHandler(event string, eventData map[string]interface{}) bool {

	if sessionInstance, ok := eventData["session"].(api.InterfaceSession); ok && sessionInstance != nil {
		if sessionID := sessionInstance.GetID(); sessionID != "" {
			currentHour := time.Now().Truncate(time.Hour).Unix()
			CheckHourUpdateForStatistic()
			if _, present := statistic[currentHour]; present && statistic[currentHour] != nil {
				statistic[currentHour].TotalVisits++
			}

			if _, present := visitState[sessionID]; !present {
				visitState[sessionID] = false
				if _, present := statistic[currentHour]; present && statistic[currentHour] != nil {
					statistic[currentHour].Visit++
				}

				err := SaveStatisticsData()
				if err != nil {
					env.LogError(err)
				}
			}
		}
	}

	return true
}

func addToCartHandler(event string, eventData map[string]interface{}) bool {

	if sessionInstance, ok := eventData["session"].(api.InterfaceSession); ok && sessionInstance != nil {
		if sessionID := sessionInstance.GetID(); sessionID != "" {

			currentHour := time.Now().Truncate(time.Hour).Unix()
			CheckHourUpdateForStatistic()

			// Add cart counter if it's a visitor that work in this hour
			if haveCard, present := visitState[sessionID]; present {
				if !haveCard {
					visitState[sessionID] = true

					if _, present := statistic[currentHour]; present && statistic[currentHour] != nil {
						statistic[currentHour].Cart++
					}

					err := SaveStatisticsData()
					if err != nil {
						env.LogError(err)
					}
				}

				// Add cart and visit counter if it's a visitor that work for a past hour
			} else {
				visitState[sessionID] = true

				if _, present := statistic[currentHour]; present && statistic[currentHour] != nil {
					statistic[currentHour].Visit++
					statistic[currentHour].TotalVisits++
					statistic[currentHour].Cart++
				}

				err := SaveStatisticsData()
				if err != nil {
					env.LogError(err)
				}
			}
		}
	}

	return true
}

func visitCheckoutHandler(event string, eventData map[string]interface{}) bool {

	currentHour := time.Now().Truncate(time.Hour).Unix()
	CheckHourUpdateForStatistic()

	if _, present := statistic[currentHour]; present && statistic[currentHour] != nil {
		statistic[currentHour].VisitCheckout++
	}

	err := SaveStatisticsData()
	if err != nil {
		env.LogError(err)
	}

	return true
}

func setPaymentHandler(event string, eventData map[string]interface{}) bool {

	currentHour := time.Now().Truncate(time.Hour).Unix()
	CheckHourUpdateForStatistic()

	if _, present := statistic[currentHour]; present && statistic[currentHour] != nil {
		statistic[currentHour].SetPayment++
	}

	err := SaveStatisticsData()
	if err != nil {
		env.LogError(err)
	}

	return true
}

func purchasedHandler(event string, eventData map[string]interface{}) bool {

	if sessionInstance, ok := eventData["session"].(api.InterfaceSession); ok && sessionInstance != nil {
		if sessionID := sessionInstance.GetID(); sessionID != "" {
			saleAmount := float64(0)
			if cartInstance, ok := eventData["cart"].(cart.InterfaceCart); ok {
				saleAmount = cartInstance.GetSubtotal()
			}

			currentHour := time.Now().Truncate(time.Hour).Unix()
			CheckHourUpdateForStatistic()

			if _, present := visitState[sessionID]; !present {
				// Increasing sales, cart and visit counters for visitor of a purchase
				if _, present := statistic[currentHour]; present && statistic[currentHour] != nil {
					statistic[currentHour].Visit++
					statistic[currentHour].TotalVisits++
					statistic[currentHour].Cart++
					statistic[currentHour].VisitCheckout++
					statistic[currentHour].SetPayment++
				}
			}

			visitState[sessionID] = false
			if _, present := statistic[currentHour]; present && statistic[currentHour] != nil {
				statistic[currentHour].Sales++
				statistic[currentHour].SalesAmount += saleAmount
			}

			err := SaveStatisticsData()
			if err != nil {
				env.LogError(err)
			}
		}
	}

	return true
}

func salesHandler(event string, eventData map[string]interface{}) bool {

	if cartInstance, ok := eventData["cart"].(cart.InterfaceCart); ok && cartInstance != nil {
		cartProducts := cartInstance.GetItems()

		if len(cartProducts) == 0 {
			return true
		}

		productQty := make(map[string]int)
		for _, cartItem := range cartProducts {
			productQty[cartItem.GetProductID()] += cartItem.GetQty()
		}

		salesHistoryCollection, err := db.GetCollection(ConstCollectionNameRTSSalesHistory)
		if err != nil {
			env.LogError(err)
			return true
		}
		currentDate := time.Now().Truncate(time.Hour).Add(time.Hour)
		for productID, count := range productQty {

			salesHistoryRecord := make(map[string]interface{})

			salesHistoryCollection.ClearFilters()
			salesHistoryCollection.AddFilter("created_at", "=", currentDate)
			salesHistoryCollection.AddFilter("product_id", "=", productID)
			dbSaleRow, err := salesHistoryCollection.Load()
			if err != nil {
				env.LogError(err)
				return true
			}

			//	rewrite existing record if we have one in database
			newCount := utils.InterfaceToInt(count)
			if len(dbSaleRow) > 0 {
				salesHistoryRecord["_id"] = dbSaleRow[0]["_id"]
				newCount = newCount + utils.InterfaceToInt(dbSaleRow[0]["count"])
			}

			// saving new history record
			if newCount > 0 {
				salesHistoryRecord["product_id"] = productID
				salesHistoryRecord["created_at"] = currentDate
				salesHistoryRecord["count"] = newCount
				_, err = salesHistoryCollection.Save(salesHistoryRecord)
				if err != nil {
					env.LogError(err)
					return true
				}
			}
		}
	}

	return true
}

func registerVisitorAsOnlineHandler(event string, eventData map[string]interface{}) bool {

	if sessionInstance, ok := eventData["session"].(api.InterfaceSession); ok {
		sessionID := sessionInstance.GetID()

		referrerType := ConstReferrerTypeDirect
		referrer := ""

		if event == "api.rts.visit" {
			if context, ok := eventData["context"].(api.InterfaceApplicationContext); ok && context != nil {
				xRreferrer := context.GetResponseSetting("X-Referer")
				referrer = utils.InterfaceToString(xRreferrer)
			}
		}

		if event == "api.request" {
			if referrerValue, present := eventData["referrer"]; present {
				referrer = utils.InterfaceToString(referrerValue)
			}
		}

		if referrer != "" {
			referrer, err := GetReferrer(referrer)
			if err != nil {
				return true
			}

			isSearchEngine := false
			for index := 0; index < len(knownSearchEngines); index++ {
				if strings.Contains(referrer, knownSearchEngines[index]) {
					isSearchEngine = true
				}
			}

			if isSearchEngine {
				referrerType = ConstReferrerTypeSearch
			} else {
				referrerType = ConstReferrerTypeSite
			}
		}

		if _, present := OnlineSessions[sessionID]; !present || OnlineSessions[sessionID] == nil {
			updateSync.Lock()
			OnlineSessions[sessionID] = &OnlineReferrer{}
			updateSync.Unlock()

			IncreaseOnline(referrerType)
			if OnlineSessionsCount := len(OnlineSessions); OnlineSessionsCount > OnlineSessionsMax {
				OnlineSessionsMax = OnlineSessionsCount
			}
		} else {
			if OnlineSessions[sessionID].referrerType != referrerType {
				DecreaseOnline(OnlineSessions[sessionID].referrerType)
				IncreaseOnline(referrerType)
			}
		}

		if _, present := OnlineSessions[sessionID]; present && OnlineSessions[sessionID] == nil {
			OnlineSessions[sessionID].time = time.Now()
			OnlineSessions[sessionID].referrerType = referrerType
		}
	}

	return true
}

func visitorOnlineActionHandler(event string, eventData map[string]interface{}) bool {

	if sessionInstance, ok := eventData["session"].(api.InterfaceSession); ok && sessionInstance != nil {
		if sessionID := sessionInstance.GetID(); sessionID != "" {
			if _, present := OnlineSessions[sessionID]; present && OnlineSessions[sessionID] != nil {
				OnlineSessions[sessionID].time = time.Now()
			}
		}
	}

	return true
}
