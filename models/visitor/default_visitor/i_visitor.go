package default_visitor

import (
	"github.com/ottemo/foundation/models/visitor"
	"github.com/ottemo/foundation/models"
)

func (it *DefaultVisitor) GetEmail() string { return it.Email }
func (it *DefaultVisitor) GetFullName() string { return it.FirstName + " " + it.LastName }
func (it *DefaultVisitor) GetFirstName() string { return it.FirstName }
func (it *DefaultVisitor) GetLastName() string { return it.LastName }

func (it *DefaultVisitor) GetShippingAddress() visitor.I_VisitorAddress {
	return it.ShippingAddress
}

func (it *DefaultVisitor) SetShippingAddress(address visitor.I_VisitorAddress) error {
	it.ShippingAddress = address
	return nil
}


func (it *DefaultVisitor) GetBillingAddress() visitor.I_VisitorAddress {
	return it.BillingAddress
}

func (it *DefaultVisitor) SetBillingAddress(address visitor.I_VisitorAddress) error {
	it.BillingAddress = address
	return nil
}


// Internal usage function which returns I_VisitorAddress model filled with values from DB
// or blank structure if no id found in DB
func (it *DefaultVisitor) getVisitorAddressById(addressId string) visitor.I_VisitorAddress {
	address_model, err := models.GetModel("VisitorAddress")
	if err != nil { return nil }

	if address_model, ok := address_model.(visitor.I_VisitorAddress); ok {
		if addressId != "" { address_model.Load(addressId) }

		return address_model
	}

	return nil
}
