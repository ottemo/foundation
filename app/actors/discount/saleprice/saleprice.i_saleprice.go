package saleprice

import "time"

// ---------------------------------------------------------------------------------------------------------------------
// InterfaceModel implementation (package "github.com/ottemo/foundation/app/models/discount/interfaces/InterfaceSalePrice")
// ---------------------------------------------------------------------------------------------------------------------

// SetAmount : amount setter
func (it *DefaultSalePrice) SetAmount(amount float64) error {
	it.amount = amount
	return nil
}

// GetAmount : amount getter
func (it *DefaultSalePrice) GetAmount() float64 {
	return it.amount
}

// SetEndDatetime : endDatetime setter
func (it *DefaultSalePrice) SetEndDatetime(endDatetime time.Time) error {
	it.endDatetime = endDatetime
	return nil
}

// GetEndDatetime : endDatetime getter
func (it *DefaultSalePrice) GetEndDatetime() time.Time {
	return it.endDatetime
}

// SetProductID : productID setter
func (it *DefaultSalePrice) SetProductID(productID string) error {
	it.productID = productID
	return nil
}

// GetProductID : productID getter
func (it *DefaultSalePrice) GetProductID() string {
	return it.productID
}

// SetStartDatetime : startDatetime setter
func (it *DefaultSalePrice) SetStartDatetime(startDatetime time.Time) error {
	it.startDatetime = startDatetime
	return nil
}

// GetStartDatetime : startDatetime getter
func (it *DefaultSalePrice) GetStartDatetime() time.Time {
	return it.startDatetime
}
