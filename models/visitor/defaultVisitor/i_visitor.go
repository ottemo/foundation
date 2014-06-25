package defaultVisitor

import "github.com/ottemo/foundation/models"

func (it *DefaultVisitor) GetEmail() string     { return it.Email }
func (it *DefaultVisitor) GetFullName() string  { return it.FirstName + " " + it.LastName }
func (it *DefaultVisitor) GetFirstName() string { return it.FirstName }
func (it *DefaultVisitor) GetLastName() string  { return it.LastName }

func (it *DefaultVisitor) GetShippingAddress() IVisitorAddress {
	return it.ShippingAddress
}

func (it *DefaultVisitor) SetShippingAddress(address IVisitorAddress) error {
	it.ShippingAddress = address
	return nil
}

func (it *DefaultVisitor) GetBillingAddress() IVisitorAddress {
	return it.BillingAddress
}

func (it *DefaultVisitor) SetBillingAddress(address IVisitorAddress) error {
	it.BillingAddress = address
	return nil
}

// Internal usage function which returns IVisitorAddress model filled with values from DB
// or blank structure if no id found in DB
func (it *DefaultVisitor) getVisitorAddressById(addressId string) IVisitorAddress {
	address_model, err := models.GetModel("VisitorAddress")
	if err != nil {
		return nil
	}

	if address_model, ok := address_model.(IVisitorAddress); ok {
		if addressId != "" {
			address_model.Load(addressId)
		}

		return address_model
	}

	return nil
}
