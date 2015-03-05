package rts

import (
	"strings"
	"time"

	"github.com/ottemo/foundation/api"
	"github.com/ottemo/foundation/app"
	"github.com/ottemo/foundation/app/models/cart"
	"github.com/ottemo/foundation/db"
	"github.com/ottemo/foundation/utils"
	"github.com/ottemo/foundation/env"
)

func referrerHandler(event string, eventData map[string]interface{}) bool {

	if _, present := eventData["context"]; present {
		if context, ok := eventData["context"].(api.InterfaceApplicationContext); ok {
			xReferrer := utils.InterfaceToString(context.GetRequestSetting("X-Referer"))
			if "" == xReferrer {
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

			referrers[referrer]++
		}
	}

	return true
}

func visitsHandler(event string, eventData map[string]interface{}) bool {

	if sessionInstance, ok := eventData["session"].(api.InterfaceSession); ok {
		sessionID := sessionInstance.GetID()

		currentHour := time.Now().Truncate(time.Hour).Unix()
		CheckHourUpdateForStatistic()

		if _, present := visitState[sessionID]; !present {
			visitState[sessionID] = false
			statistic[currentHour].Visit++

			err := SaveStatisticsData()
			if err != nil {
				env.LogError(err)
			}
		}
	}

	return true
}

func addToCartHandler(event string, eventData map[string]interface{}) bool {

	if sessionInstance, ok := eventData["session"].(api.InterfaceSession); ok {
		sessionID := sessionInstance.GetID()

		currentHour := time.Now().Truncate(time.Hour).Unix()
		CheckHourUpdateForStatistic()

		// Add cart counter if it's a visitor that work in this hour
		if haveCard, present := visitState[sessionID]; present {
			if !haveCard {
				visitState[sessionID] = true
				statistic[currentHour].Cart++

				err := SaveStatisticsData()
				if err != nil {
					env.LogError(err)
				}
			}

			// Add cart and visit counter if it's a visitor that work for a past hour
		} else {
			visitState[sessionID] = true
			statistic[currentHour].Visit++
			statistic[currentHour].Cart++

			err := SaveStatisticsData()
			if err != nil {
				env.LogError(err)
			}
		}
	}

	return true
}

func purchasedHandler(event string, eventData map[string]interface{}) bool {

	if sessionInstance, ok := eventData["session"].(api.InterfaceSession); ok {
		sessionID := sessionInstance.GetID()

		currentHour := time.Now().Truncate(time.Hour).Unix()
		CheckHourUpdateForStatistic()

		// Add sales counter if it's a visitor that work in this hour
		if _, present := visitState[sessionID]; present {
			if visitState[sessionID] {
				visitState[sessionID] = false
				statistic[currentHour].Sales++

				err := SaveStatisticsData()
				if err != nil {
					env.LogError(err)
				}
			}

			// Add sales, cart and visit counter if it's a visitor that work for a past hour
		} else {
			visitState[sessionID] = false
			statistic[currentHour].Visit++
			statistic[currentHour].Cart++
			statistic[currentHour].Sales++

			err := SaveStatisticsData()
			if err != nil {
				env.LogError(err)
			}
		}
	}

	return true
}

func salesHandler(event string, eventData map[string]interface{}) bool {

	if cartInstance, ok := eventData["cart"].(cart.InterfaceCart); ok {
		cartProducts := cartInstance.GetItems()

		if len(cartProducts) == 0 {
			return true
		}

		productQtys := make(map[string]interface{})
		for i := range cartProducts {
			productQtys[cartProducts[i].GetProductID()] = cartProducts[i].GetQty()
		}

		salesData := make(map[string]int)

		salesHistoryCollection, err := db.GetCollection(ConstCollectionNameRTSSalesHistory)
		if err != nil {
			return true
		}

		for productID, count := range productQtys {
			currentDate := time.Now().Truncate(time.Hour * 24)

			salesHistoryRecord := make(map[string]interface{})
			salesData[productID] = utils.InterfaceToInt(count)

			salesHistoryCollection.ClearFilters()
			salesHistoryCollection.AddFilter("created_at", "=", currentDate)
			salesHistoryCollection.AddFilter("product_id", "=", productID)
			dbSaleRow, _ := salesHistoryCollection.Load()

			newCount := utils.InterfaceToInt(count)
			if len(dbSaleRow) > 0 {
				salesHistoryRecord["_id"] = utils.InterfaceToString(dbSaleRow[0]["_id"])
				oldCount := utils.InterfaceToInt(dbSaleRow[0]["count"])
				newCount += oldCount
			}

			// saving new history record
			if newCount > 0 {
				salesHistoryRecord["product_id"] = productID
				salesHistoryRecord["created_at"] = currentDate
				salesHistoryRecord["count"] = newCount
				_, err = salesHistoryCollection.Save(salesHistoryRecord)
				if err != nil {
					return true
				}
			}
		}

		SaveSalesData(salesData)
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

		if _, ok := OnlineSessions[sessionID]; !ok {
			OnlineSessions[sessionID] = &OnlineReferrer{}
			IncreaseOnline(referrerType)
			if len(OnlineSessions) > OnlineSessionsMax {
				OnlineSessionsMax = len(OnlineSessions)
			}
		} else {
			if OnlineSessions[sessionID].referrerType != referrerType {
				DecreaseOnline(OnlineSessions[sessionID].referrerType)
				IncreaseOnline(referrerType)
			}
		}

		OnlineSessions[sessionID].time = time.Now()
		OnlineSessions[sessionID].referrerType = referrerType
	}

	return true
}

func visitorOnlineActionHandler(event string, eventData map[string]interface{}) bool {

	if sessionInstance, ok := eventData["session"].(api.InterfaceSession); ok {
		sessionID := sessionInstance.GetID()
		if _, ok := OnlineSessions[sessionID]; ok {
			OnlineSessions[sessionID].time = time.Now()
		}
	}

	return true
}
