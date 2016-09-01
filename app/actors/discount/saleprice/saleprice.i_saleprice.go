package saleprice

import "time"

// ---------------------------------------------------------------------------------------------------------------------
// InterfaceModel implementation (package "github.com/ottemo/foundation/app/models/discount/interfaces/InterfaceSalePrice")
// ---------------------------------------------------------------------------------------------------------------------

// Amount setter
func (it *DefaultSalePrice) SetAmount(amount float64) error {
	it.amount = amount
	return nil
}

// Amount getter
func (it *DefaultSalePrice) GetAmount() float64 {
	return it.amount
}

// EndDatetime setter
func (it *DefaultSalePrice) SetEndDatetime(endDatetime time.Time) error {
	it.endDatetime = endDatetime
	return nil
}

// EndDatetime getter
func (it *DefaultSalePrice) GetEndDatetime() time.Time {
	return it.endDatetime
}

// ProductID setter
func (it *DefaultSalePrice) SetProductID(product_id string) error {
	it.product_id = product_id
	return nil
}

// ProductID getter
func (it *DefaultSalePrice) GetProductID() string {
	return it.product_id
}

// StartDatetime setter
func (it *DefaultSalePrice) SetStartDatetime(startDatetime time.Time) error {
	it.startDatetime = startDatetime
	return nil
}

// StartDatetime getter
func (it *DefaultSalePrice) GetStartDatetime() time.Time {
	return it.startDatetime
}
