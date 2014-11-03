package tests

import (
	"errors"
	"math/rand"
	"strings"
	"testing"

	randomdata "github.com/Pallinder/go-randomdata"

	"github.com/ottemo/foundation/api/session"
	"github.com/ottemo/foundation/utils"

	"github.com/ottemo/foundation/app/models/cart"
	"github.com/ottemo/foundation/app/models/checkout"
	"github.com/ottemo/foundation/app/models/product"
	"github.com/ottemo/foundation/app/models/visitor"
)

// routine to emulate full checkout process at once
func makeCheckout() error {

	// customizing visitor
	//--------------------
	currentVisitor, err := visitor.GetVisitorModel()
	if err != nil {
		return err
	}

	visitorEmail := randomdata.SillyName() + "@ottemo-test.com"
	err = currentVisitor.LoadByEmail(visitorEmail)
	if err != nil {
		err := currentVisitor.Set("email", visitorEmail)
		if err != nil {
			return err
		}

		err = currentVisitor.Set("first_name", randomdata.FirstName(randomdata.RandomGender))
		if err != nil {
			return err
		}

		err = currentVisitor.Set("last_name", randomdata.LastName())
		if err != nil {
			return err
		}

		err = currentVisitor.Set("validate", "")
		if err != nil {
			return err
		}

		err = currentVisitor.Save()
		if err != nil {
			return err
		}
	}

	// making new session and checkout
	//---------------------------------
	currentSession, err := session.NewSession()
	if err != nil {
		return err
	}

	currentCheckout, err := checkout.GetCheckoutModel()
	if err != nil {
		return err
	}
	err = currentCheckout.SetSession(currentSession)
	if err != nil {
		return err
	}

	err = currentCheckout.SetVisitor(currentVisitor)
	if err != nil {
		return err
	}

	currentCart, err := cart.GetCartModel()
	if err != nil {
		return err
	}

	err = currentCart.SetVisitorId(currentVisitor.GetId())
	if err != nil {
		return err
	}

	// filling cart with products
	//----------------------------
	productCollection, err := product.GetProductCollectionModel()
	if err != nil {
		return err
	}

	productsCount, err := productCollection.GetDBCollection().Count()
	if err != nil {
		return err
	}

	err = productCollection.ListLimit(rand.Intn(productsCount), rand.Intn(4)+1)
	if err != nil {
		return err
	}

	for _, productModel := range productCollection.ListProducts() {
		options := make(map[string]interface{})

		// checking for required option for product
		for optionName, productOption := range productModel.GetOptions() {
			if typedProductOption, ok := productOption.(map[string]interface{}); ok {
				if isRequired, present := typedProductOption["required"]; present && utils.InterfaceToBool(isRequired) {
					options[optionName] = "something"

					if optionValues, present := typedProductOption["options"]; present {
						if typedOptionValues, ok := optionValues.(map[string]interface{}); ok {
							for optionValueName, _ := range typedOptionValues {
								if rand.Intn(2) == 1 {
									options[optionName] = optionValueName
									break
								}
							}
						}
					}
				}
			}
		}

		//adding item to cart
		_, err := currentCart.AddItem(productModel.GetId(), rand.Intn(3)+1, options)
		if err != nil {
			return err
		}
	}
	err = currentCart.Save()
	if err != nil {
		return err
	}
	err = currentCheckout.SetCart(currentCart)
	if err != nil {
		return err
	}

	// setting shipping/payment address
	addressModel, err := visitor.GetVisitorAddressModel()
	if err != nil {
		return err
	}
	err = addressModel.Set("first_name", currentVisitor.GetFirstName())
	if err != nil {
		return err
	}
	err = addressModel.Set("last_name", currentVisitor.GetLastName())
	if err != nil {
		return err
	}
	err = addressModel.Set("address_line1", randomdata.Street())
	if err != nil {
		return err
	}
	err = addressModel.Set("country", "US")
	if err != nil {
		return err
	}
	err = addressModel.Set("city", randomdata.City())
	if err != nil {
		return err
	}
	err = addressModel.Set("state", randomdata.State(randomdata.Small))
	if err != nil {
		return err
	}
	err = addressModel.Set("phone", utils.InterfaceToString(randomdata.Number(1000000000, 9999999999)))
	if err != nil {
		return err
	}
	err = addressModel.Set("zip", randomdata.PostalCode("US"))

	err = currentCheckout.SetShippingAddress(addressModel)
	if err != nil {
		return err
	}
	err = currentCheckout.SetBillingAddress(addressModel)
	if err != nil {
		return err
	}

	// setting shipping method
	found := false
	for _, shippingMethod := range checkout.GetRegisteredShippingMethods() {
		if shippingMethod.GetCode() == "flat_rate" {
			currentCheckout.SetShippingMethod(shippingMethod)
			found = true
			break
		}
	}
	if !found {
		errors.New("Shipping method 'flat rate' not registered")
	}

	// setting payment method
	found = false
	for _, paymentMethod := range checkout.GetRegisteredPaymentMethods() {
		if paymentMethod.GetCode() == "checkmo" {
			currentCheckout.SetPaymentMethod(paymentMethod)
			found = true
			break
		}
	}
	if !found {
		errors.New("Payment method 'check money order' not registered")
	}

	// submiting checkout
	_, err = currentCheckout.Submit()
	if err != nil && !strings.Contains(err.Error(), "Authentication required") {
		return err
	}

	return nil
}

// benchmarks order creation from zero
func BenchmarkCheckout(b *testing.B) {
	// starting application and getting product model
	err := StartAppInTestingMode()
	if err != nil {
		b.Error(err)
	}

	err = MakeSureProductsCount(100)
	if err != nil {
		b.Error(err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		err := makeCheckout()
		if err != nil {
			b.Error(err)
		}
	}
}
