package braintree

import (
	"strings"

	"github.com/lionelbarrow/braintree-go"

	"github.com/ottemo/foundation/app/models/order"
	"github.com/ottemo/foundation/app/models/visitor"
	"github.com/ottemo/foundation/env"
	"github.com/ottemo/foundation/utils"
)

// getCustomerIDByVisitorID returns 3rd party customer ID by visitor registered ID
func getCustomerIDByVisitorID(visitorID string) string {
	var absentIDValue = ""
	var customerIDAttribute = "customer_id"

	if visitorID == "" {
		_ = env.ErrorDispatch(env.ErrorNew(ConstErrorModule, ConstErrorLevel, "0f6c678f-66a3-470e-8a80-5cc2ff619058", "empty visitor ID passed to look up customer token"))
		return absentIDValue
	}

	model, _ := visitor.GetVisitorCardCollectionModel()
	if err := model.ListFilterAdd("visitor_id", "=", visitorID); err != nil {
		_ = env.ErrorNew(ConstErrorModule, ConstErrorLevel, "472ae679-c6f4-4e32-864f-6bbb88e11d3c", err.Error())
	}
	if err := model.ListFilterAdd("payment", "=", constCCMethodCode); err != nil {
		_ = env.ErrorNew(ConstErrorModule, ConstErrorLevel, "d45b2df6-9850-479e-821e-77941d6ade7b", err.Error())
	}

	// 3rd party customer identifier, used by braintree
	err := model.ListAddExtraAttribute(customerIDAttribute)
	if err != nil {
		_ = env.ErrorDispatch(err)
	}

	visitorCards, err := model.List()
	if err != nil {
		_ = env.ErrorDispatch(err)
	}

	for _, visitorCard := range visitorCards {
		return utils.InterfaceToString(visitorCard.Extra[customerIDAttribute])
	}

	return absentIDValue
}

// braintreeCardFormatExpirationDate returns braintree CreditCard expiration date formatted as MMYY
func braintreeCardFormatExpirationDate(card braintree.CreditCard) (string, error) {
	var expirationDate = utils.InterfaceToString(card.ExpirationMonth)

	if len(card.ExpirationMonth) < 1 {
		return "", env.ErrorNew(ConstErrorModule, ConstErrorLevel, "f861d39f-516d-4bdd-8316-b7c4d34e3531", "unexpected month value coming back from braintree: "+card.ExpirationMonth)
	}

	// pad with a zero
	if len(card.ExpirationMonth) < 2 {
		expirationDate = "0" + expirationDate
	}

	// append the last two year digits
	year := utils.InterfaceToString(card.ExpirationYear)
	if len(year) == 4 {
		expirationDate = expirationDate + year[2:]
	} else {
		return "", env.ErrorNew(ConstErrorModule, ConstErrorLevel, "950aea13-16e8-4d20-9ad0-f5cee26c03c2", "unexpected year length coming back from braintree: "+year)
	}

	return expirationDate, nil
}

// braintreeAddressFromVisitorAddress populates braintree Address by InterfaceVisitorAddress
func braintreeAddressFromVisitorAddress(visitorAddress visitor.InterfaceVisitorAddress) *braintree.Address {
	return &braintree.Address{
		FirstName:       visitorAddress.GetFirstName(),
		LastName:        visitorAddress.GetLastName(),
		Company:         visitorAddress.GetCompany(),
		StreetAddress:   visitorAddress.GetAddressLine1(),
		ExtendedAddress: visitorAddress.GetAddressLine2(),

		CountryCodeAlpha2: visitorAddress.GetCountry(),
		Locality:          visitorAddress.GetCity(),
		Region:            visitorAddress.GetState(),
		PostalCode:        visitorAddress.GetZipCode(),
	}
}

// braintreeCardToAuthorizeResult populates authorize result by creadit card and customer info
func braintreeCardToAuthorizeResult(card braintree.CreditCard, customerID string) (map[string]interface{}, error) {
	expirationDate, err := braintreeCardFormatExpirationDate(card)
	if err != nil {
		return nil, env.ErrorNew(ConstErrorModule, ConstErrorLevel, "7aa0ea8e-679e-4ac3-b84b-40aad71ead5f", "unable to format expiration date: "+err.Error())
	}

	var result = map[string]interface{}{
		"transactionID":      card.Token,     // token_id
		"creditCardLastFour": card.Last4,     // number
		"creditCardType":     card.CardType,  // type
		"creditCardExp":      expirationDate, // expiration_date
		"customerID":         customerID,     // customer_id
	}

	return result, nil
}

// braintreeCustomerParamsByVisitorID populates braintree customer by visitor ID
func braintreeCustomerParamsByVisitorID(visitorID string) (*braintree.Customer, error) {
	visitorData, err := visitor.LoadVisitorByID(visitorID)
	if err != nil {
		return nil, env.ErrorNew(ConstErrorModule, ConstErrorLevel, "09ec64dd-d5c7-4179-aad3-a019c0cd857f", "internal error: unable to load visitor by ID.")
	}

	var customerParamsPtr = &braintree.Customer{
		FirstName: visitorData.GetFirstName(),
		LastName:  visitorData.GetLastName(),
		Email:     visitorData.GetEmail(),
	}

	return customerParamsPtr, nil
}

// braintreeCustomerParamsByVisitorData populates braintree customer by visitor info
func braintreeCustomerParamsByVisitorData(visitorInfo map[string]interface{}) (*braintree.Customer, error) {
	if !utils.KeysInMapAndNotBlank(visitorInfo, "billing_name", "email") {
		return nil, env.ErrorNew(ConstErrorModule, ConstErrorLevel, "5a6abb6a-9e6a-486e-b296-f4a7e4a196df", "visitor data incomplete: billing_name or email")
	}

	var nameParts = strings.SplitN(utils.InterfaceToString(visitorInfo["billing_name"])+" ", " ", 2)
	var firstName = strings.TrimSpace(nameParts[0])
	var lastName = strings.TrimSpace(nameParts[1])

	var customerParamsPtr = &braintree.Customer{
		FirstName: firstName,
		LastName:  lastName,
		Email:     utils.InterfaceToString(visitorInfo["email"]),
	}

	return customerParamsPtr, nil
}

// braintreeCreateCustomer creates 3rd party customer
func braintreeCreateCustomer(visitorInfo map[string]interface{}) (*braintree.Customer, error) {
	visitorIDValue, _ := visitorInfo["visitor_id"]
	visitorID := utils.InterfaceToString(visitorIDValue)

	var customerParamsPtr *braintree.Customer

	if visitorID == "" {
		paramsPtr, err := braintreeCustomerParamsByVisitorData(visitorInfo)
		if err != nil {
			return nil, env.ErrorNew(ConstErrorModule, ConstErrorLevel, "e67e5617-115f-409e-81c1-271bbc3eaa3f", "unable to create customer params: "+err.Error())
		}

		customerParamsPtr = paramsPtr
	} else {
		paramsPtr, err := braintreeCustomerParamsByVisitorID(visitorID)
		if err != nil {
			return nil, env.ErrorNew(ConstErrorModule, ConstErrorLevel, "6d8155c9-969d-4f03-a7c3-e61005e661e1", "unable to create customer params: "+err.Error())
		}

		customerParamsPtr = paramsPtr
	}

	braintreeInstance, err := getBraintreeInstance()
	if err != nil {
		return nil, env.ErrorNew(ConstErrorModule, ConstErrorLevel, "71055f8d-2236-4c5c-891b-8b41c26e77a2", "internal error: unable to initialize braintree: "+err.Error())
	}

	return braintreeInstance.Customer().Create(customerParamsPtr)
}

// braintreeCreateCardParams populates braintree credit card by credit card info
func braintreeCreateCardParams(creditCardMap map[string]interface{}) (*braintree.CreditCard, error) {
	if !utils.KeysInMapAndNotBlank(creditCardMap, "cvc", "number", "expire_year", "expire_month") {
		return nil, env.ErrorNew(ConstErrorModule, ConstErrorLevel, "bd0a78bf-065a-462b-92c7-d5a1529797c4", "credit card data incomplete: cvc, number or expiration date")
	}

	creditCardParams := &braintree.CreditCard{
		Number:          utils.InterfaceToString(creditCardMap["number"]),
		ExpirationYear:  utils.InterfaceToString(creditCardMap["expire_year"]),
		ExpirationMonth: utils.InterfaceToString(creditCardMap["expire_month"]),
		CVV:             utils.InterfaceToString(creditCardMap["cvc"]),
	}

	return creditCardParams, nil
}

// braintreeCreateCard creates braintree credit card by credita crad and customer info
func braintreeCreateCard(creditCardMap map[string]interface{}, customerID string) (*braintree.CreditCard, error) {
	creditCardParams, err := braintreeCreateCardParams(creditCardMap)
	if err != nil {
		return nil, env.ErrorNew(ConstErrorModule, ConstErrorLevel, "46ddda47-867d-4d28-979e-2c43f6f631dd", "unable to create card params: "+err.Error())
	}

	creditCardParams.CustomerId = customerID
	creditCardParams.Options = &braintree.CreditCardOptions{
		VerifyCard: true,
	}

	braintreeInstance, err := getBraintreeInstance()
	if err != nil {
		return nil, env.ErrorNew(ConstErrorModule, ConstErrorLevel, "150f6b8c-2fb3-41de-bc1d-d5534e3f2e19", "internal error: unable to initialize braintree: "+err.Error())
	}

	return braintreeInstance.CreditCard().Create(creditCardParams)
}

// braintreeRegisterCardForVisitor registers braintree card for visitor
func braintreeRegisterCardForVisitor(visitorInfo map[string]interface{}, creditCardInfo map[string]interface{}) (*braintree.CreditCard, error) {
	// 1. Get visitor token
	// Skip presence check - visitor could be set as by ID as by billing name
	visitorIDValue, _ := visitorInfo["visitor_id"]

	customerID := getCustomerIDByVisitorID(utils.InterfaceToString(visitorIDValue))

	if customerID == "" {
		// 2. We don't have a braintree client id on file, make a new customer
		customerPtr, err := braintreeCreateCustomer(visitorInfo)
		if err != nil {
			return nil, env.ErrorNew(ConstErrorModule, ConstErrorLevel, "69ab6b26-457c-4a7a-b784-19594f25514d", "unable to create customer: "+err.Error())
		}

		customerID = customerPtr.Id
	}

	// 3. Create a card
	return braintreeCreateCard(creditCardInfo, customerID)
}

// braintreeTransactionParamsByOrder populates braintree transaction params by order info
func braintreeTransactionParamsByOrder(orderInstance order.InterfaceOrder) (*braintree.Transaction, error) {
	transactionParams := &braintree.Transaction{
		Type:    "sale",
		Amount:  braintree.NewDecimal(int64(orderInstance.GetGrandTotal()*100), 2),
		OrderId: orderInstance.GetID(),
		Options: &braintree.TransactionOptions{
			SubmitForSettlement: true,
		},
	}

	transactionParams.BillingAddress = braintreeAddressFromVisitorAddress(orderInstance.GetBillingAddress())
	transactionParams.ShippingAddress = braintreeAddressFromVisitorAddress(orderInstance.GetShippingAddress())

	return transactionParams, nil
}

// chargeGuestVisitor charges UNregistered visitor
func chargeGuestVisitor(orderInstance order.InterfaceOrder, creditCardInfoMap map[string]interface{}) (*braintree.Transaction, error) {
	transactionParams, err := braintreeTransactionParamsByOrder(orderInstance)
	if err != nil {
		return nil, env.ErrorNew(ConstErrorModule, ConstErrorLevel, "395cd75d-ee5f-4e42-9c1d-22c56f5a28c0", "unable to create transaction params: "+err.Error())
	}

	creditCardParams, err := braintreeCreateCardParams(creditCardInfoMap)
	if err != nil {
		return nil, env.ErrorNew(ConstErrorModule, ConstErrorLevel, "d9e6d4a3-c533-428a-9bd6-9978ca7c0c24", "unable to create card params: "+err.Error())
	}

	transactionParams.CreditCard = creditCardParams

	braintreeInstance, err := getBraintreeInstance()
	if err != nil {
		return nil, env.ErrorNew(ConstErrorModule, ConstErrorLevel, "2cd29ffd-38ab-49eb-ad18-7fc9a948ce06", "internal error: unable to initialize braintree: "+err.Error())
	}

	return braintreeInstance.Transaction().Create(transactionParams)
}

// chargeRegisteredVisitor charges Registered visitor
func chargeRegisteredVisitor(orderInstance order.InterfaceOrder, creditCard visitor.InterfaceVisitorCard) (*braintree.Transaction, error) {
	var err error
	cardToken := creditCard.GetToken()
	customerID := creditCard.GetCustomerID()

	if cardToken == "" || customerID == "" {
		return nil, env.ErrorNew(ConstErrorModule, ConstErrorLevel, "6b43e527-9bc7-48f7-8cdd-320ceb6d77e6", "invalid token or customer id")
	}

	braintreeInstance, err := getBraintreeInstance()
	if err != nil {
		return nil, env.ErrorNew(ConstErrorModule, ConstErrorLevel, "06a80ccb-82e2-42e9-89bf-7e8ee2936728", "internal error: unable to initialize braintree: "+err.Error())
	}

	if _, err := braintreeInstance.CreditCard().Find(cardToken); err != nil {
		return nil, env.ErrorNew(ConstErrorModule, ConstErrorLevel, "bb3748d9-67f6-4c1f-b2da-d1e7a6f0e519", "unable to find credit card: "+err.Error())
	}

	transactionParams, err := braintreeTransactionParamsByOrder(orderInstance)
	if err != nil {
		return nil, env.ErrorNew(ConstErrorModule, ConstErrorLevel, "5a7ca589-2031-41ee-8997-c3824a0178f4", "unable to create transaction params: "+err.Error())
	}
	transactionParams.CustomerID = customerID
	transactionParams.PaymentMethodToken = cardToken
	transactionParams.Options.StoreInVault = true

	return braintreeInstance.Transaction().Create(transactionParams)
}

func getBraintreeInstance() (*braintree.Braintree, error) {
	var environmentValue = utils.InterfaceToString(env.ConfigGetValue(ConstGeneralConfigPathEnvironment))

	// braintree package could panic
	if environmentValue != string(braintree.Development) &&
		environmentValue != string(braintree.Sandbox) &&
		environmentValue != string(braintree.Production) {

		return nil, env.ErrorNew(ConstErrorModule, ConstErrorLevel, "1e2b4af6-0256-4324-97a5-c451957119d4", "internal error: invalid braintree environment ["+environmentValue+"]")
	}

	// we do not check other config values because they points that something configured incorrectly
	var braintreeInstance = braintree.New(
		braintree.Environment(environmentValue),
		utils.InterfaceToString(env.ConfigGetValue(ConstGeneralConfigPathMerchantID)),
		utils.InterfaceToString(env.ConfigGetValue(ConstGeneralConfigPathPublicKey)),
		utils.InterfaceToString(env.ConfigGetValue(ConstGeneralConfigPathPrivateKey)),
	)

	return braintreeInstance, nil
}
