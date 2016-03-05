package shareasale

import (
	"fmt"
	"net/http"

	"github.com/ottemo/foundation/app/models/order"
	"github.com/ottemo/foundation/env"
	"github.com/ottemo/foundation/utils"
)

func checkoutSuccessHandler(event string, eventData map[string]interface{}) bool {

	// if Share A Sale is not enabled ignore this handler
	if enabled := utils.InterfaceToBool(env.ConfigGetValue(ConstConfigPathShareASaleEnabled)); !enabled {
		return true
	}

	// inpsect the order
	var checkoutOrder order.InterfaceOrder
	if eventItem, present := eventData["order"]; present {
		if typedItem, ok := eventItem.(order.InterfaceOrder); ok {
			checkoutOrder = typedItem
		}
	}

	if checkoutOrder != nil {
		aSale, _ := createAffiliateSale(checkoutOrder)
		go sendSale(aSale)
	}

	return true
}

func createAffiliateSale(order order.InterfaceOrder) (*AffiliateSale, error) {

	var merchantID string
	var grandTotal, taxes, shipping, discount float64
	var aSale AffiliateSale

	// retrieve merchant id
	if merchantID = utils.InterfaceToString(env.ConfigGetValue(ConstConfigPathShareASaleMerchantID)); merchantID == "" {
		return nil, env.ErrorNew(ConstErrorModule, env.ConstErrorLevelAPI, "247fec44-cb47-494c-ae5f-adfd6f28eb2d", "Share a Sale Merchant ID may not be empty.")
	}
	aSale.MerchantID = merchantID

	// calculate true subtotal
	grandTotal = utils.InterfaceToFloat64(order.Get("grand_total"))
	taxes = utils.InterfaceToFloat64(order.Get("tax_amount"))
	shipping = utils.InterfaceToFloat64(order.Get("shipping_amount"))
	discount = utils.InterfaceToFloat64(order.Get("discount"))

	aSale.SubTotal = grandTotal - (taxes + shipping + discount)

	// order number
	aSale.OrderNo = utils.InterfaceToString(order.Get("_id"))

	return &aSale, nil
}

// sendSale is to send tracking information to Share A Sale
func sendSale(sale *AffiliateSale) error {

	var err error
	var url string

	// construct the url
	url = fmt.Sprintf("https://shareasale.com/sale.cfm?amount=%.2f&tracking=%s&transtype=SALE&merchantID=%s", sale.SubTotal, sale.OrderNo, sale.MerchantID)
	fmt.Println("ShareASale Tracking url: " + url)

	// send tracking info
	response, err := http.Get(url)
	fmt.Printf("Response Code: %v\n", response.StatusCode)
	fmt.Printf("Response Header: %+v\n", response.Header)
	if err != nil {
		return env.ErrorDispatch(err)
	}

	// check the status code
	if response.StatusCode != http.StatusOK {
		return env.ErrorNew(ConstErrorModule, env.ConstErrorLevelAPI, "6a3346e3-115e-4c72-bc0f-ee14d7f15c43", "ShareASale.com is not responding when attempting to send affiliate sales.")
	}

	defer response.Body.Close()

	return nil
}
