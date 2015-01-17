package tests

import (
	"errors"
	"fmt"
	"math/rand"

	randomdata "github.com/Pallinder/go-randomdata"
	golorem "github.com/drhodes/golorem"

	"github.com/ottemo/foundation/api"
	"github.com/ottemo/foundation/utils"

	"github.com/ottemo/foundation/app/models/cart"
	"github.com/ottemo/foundation/app/models/checkout"
	"github.com/ottemo/foundation/app/models/product"
	"github.com/ottemo/foundation/app/models/visitor"
)

// GetRandomVisitor returns visitor object with randomly filled data
func GetRandomVisitor() (visitor.InterfaceVisitor, error) {
	randomVisitor, err := visitor.GetVisitorModel()
	if err != nil {
		return nil, err
	}

	visitorEmail := randomdata.SillyName() + "@ottemo-test.com"
	err = randomVisitor.LoadByEmail(visitorEmail)
	if err != nil {
		err := randomVisitor.Set("email", visitorEmail)
		if err != nil {
			return nil, err
		}

		err = randomVisitor.Set("first_name", randomdata.FirstName(randomdata.RandomGender))
		if err != nil {
			return nil, err
		}

		err = randomVisitor.Set("last_name", randomdata.LastName())
		if err != nil {
			return nil, err
		}

		err = randomVisitor.Set("validate", "")
		if err != nil {
			return nil, err
		}

		err = randomVisitor.Set("password", "123")
		if err != nil {
			return nil, err
		}

		err = randomVisitor.Save()
		if err != nil {
			return nil, err
		}
	}
	return randomVisitor, nil
}

// GetNewCheckout returns new checkout object with assigned new session, and cart to it
func GetNewCheckout(checkoutVisitor visitor.InterfaceVisitor) (checkout.InterfaceCheckout, error) {
	newSession, err := api.NewSession()
	if err != nil {
		return nil, err
	}

	newCheckout, err := checkout.GetCheckoutModel()
	if err != nil {
		return nil, err
	}
	err = newCheckout.SetSession(newSession)
	if err != nil {
		return nil, err
	}

	err = newCheckout.SetVisitor(checkoutVisitor)
	if err != nil {
		return nil, err
	}

	newCart, err := cart.GetCartModel()
	if err != nil {
		return nil, err
	}

	err = newCart.MakeCartForVisitor(checkoutVisitor.GetID())
	if err != nil {
		return nil, err
	}

	err = newCheckout.SetCart(newCart)
	if err != nil {
		return nil, err
	}

	return newCheckout, nil
}

// AddRandomProductsToCart adds n count products to checkout cart
func AddRandomProductsToCart(currentCheckout checkout.InterfaceCheckout, n int) error {
	if n <= 0 {
		return nil
	}

	// randomizing products
	productCollection, err := product.GetProductCollectionModel()
	if err != nil {
		return err
	}

	productsCount, err := productCollection.GetDBCollection().Count()
	if err != nil {
		return err
	}

	err = productCollection.ListLimit(rand.Intn(productsCount), n)
	if err != nil {
		return err
	}

	// filling cart with randomized products
	currentCart := currentCheckout.GetCart()
	for _, productModel := range productCollection.ListProducts() {
		options := make(map[string]interface{})

		// checking for required option for product
		for optionName, productOption := range productModel.GetOptions() {
			if typedProductOption, ok := productOption.(map[string]interface{}); ok {
				if isRequired, present := typedProductOption["required"]; present && utils.InterfaceToBool(isRequired) {
					options[optionName] = "something"

					if optionValues, present := typedProductOption["options"]; present {
						if typedOptionValues, ok := optionValues.(map[string]interface{}); ok {
							for optionValueName := range typedOptionValues {
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
		_, err := currentCart.AddItem(productModel.GetID(), rand.Intn(3)+1, options)
		if err != nil {
			return err
		}
	}
	err = currentCart.Save()
	if err != nil {
		return err
	}

	return nil
}

// RandomizeShippingAndBillingAddresses sets shipping and billing addresses for checkout object
func RandomizeShippingAndBillingAddresses(currentCheckout checkout.InterfaceCheckout) error {
	currentVisitor := currentCheckout.GetVisitor()
	if currentVisitor == nil {
		return errors.New("visitor for checkout is not set")
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
	err = addressModel.Set("phone", randomdata.StringNumber(5, ""))
	if err != nil {
		return err
	}
	err = addressModel.Set("zip", randomdata.PostalCode("US"))
	if err != nil {
		return err
	}

	err = addressModel.Save()
	if err != nil {
		return err
	}

	err = currentCheckout.SetShippingAddress(addressModel)
	if err != nil {
		return err
	}
	err = currentCheckout.SetBillingAddress(addressModel)
	if err != nil {
		return err
	}

	return nil
}

// UpdateShippingAndPaymentMethods sets check money order payment method and flat rate shipping method to checkout
func UpdateShippingAndPaymentMethods(currentCheckout checkout.InterfaceCheckout) error {
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
	return nil
}

// FullCheckout is a routine to emulate full checkout process at once
func FullCheckout() error {

	currentVisitor, err := GetRandomVisitor()
	if err != nil {
		return err
	}

	currentCheckout, err := GetNewCheckout(currentVisitor)
	if err != nil {
		return err
	}

	err = AddRandomProductsToCart(currentCheckout, rand.Intn(4)+1)
	if err != nil {
		return err
	}

	err = RandomizeShippingAndBillingAddresses(currentCheckout)
	if err != nil {
		return err
	}

	err = UpdateShippingAndPaymentMethods(currentCheckout)
	if err != nil {
		return err
	}

	_, err = currentCheckout.Submit()
	if err != nil {
		return err
	}

	return nil
}

// MakeSureProductsCount checks products count in DB and adds missing if needed
func MakeSureProductsCount(countShouldBe int) error {

	// getting database products count
	productCollection, err := product.GetProductCollectionModel()
	if err != nil {
		return err
	}

	productDBCollection := productCollection.GetDBCollection()
	if productDBCollection == nil {
		return errors.New("can't get db collection")
	}

	count, err := productDBCollection.Count()
	if err != nil {
		return err
	}

	// adding missed products
	for i := count; i <= countShouldBe; i++ {
		productModel, err := product.GetProductModel()
		if err != nil {
			return err
		}

		productModel.Set("enabled", true)
		productModel.Set("sku", fmt.Sprintf("test-%d", i))
		productModel.Set("name", fmt.Sprintf("Test Product %d", i))

		productModel.Set("short_description", golorem.Paragraph(1, 5))
		productModel.Set("description", golorem.Paragraph(5, 10))

		productModel.Set("default_image", "")
		productModel.Set("price", utils.RoundPrice(rand.Float64()*10))
		productModel.Set("weight", utils.RoundPrice(rand.Float64()*10))

		// productModel.Set("options", "")

		err = productModel.Save()
		if err != nil {
			return err
		}
	}

	return nil
}
