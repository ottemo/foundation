package checkout

import (
	"fmt"
	"github.com/ottemo/foundation/app/models/checkout"
	"testing"
	"github.com/ottemo/foundation/utils"
)

func TestPriceAdjustments(t *testing.T) {
	// PA - price adjustment - object that used as a basic element for checkout calculation
	// in this test priority is not used as it called only in full calculation process
	priceAdjustments := []checkout.StructPriceAdjustment{
		// subtotal holds base cost for items
		checkout.StructPriceAdjustment{
			Code:     checkout.ConstLabelSubtotal,
			Name:     checkout.ConstLabelSubtotal,
			Priority: checkout.ConstCalculateTargetSubtotal,
			Labels:   []string{checkout.ConstLabelSubtotal}, // any amount is applied to grand total
			// and additionally to any label (key) provided here
			PerItem: map[string]float64{
				"1": 100, // product1 1x100
				"2": 300, // product2 2x150
			},
		},

		// sale price per item
		checkout.StructPriceAdjustment{
			Code:     "sale_price",
			Name:     "Sale Price",
			Amount:   0,
			Priority: 1.1,
			Labels:   []string{checkout.ConstLabelSalePriceAdjustment},
			PerItem: map[string]float64{
				"1": -1, // 1$ on product1
				"2": -2, // 2$ on product2
			},
		},

		checkout.StructPriceAdjustment{
			Code:     "default",
			Name:     "Flat Rate",
			Amount:   100,
			Priority: checkout.ConstCalculateTargetShipping,
			Labels:   []string{checkout.ConstLabelShipping},
			// not necessary to provide PerItem key if it applies amount to full cart
		},

		// discount before tax
		// buy one get one - 150 amount for one product2
		checkout.StructPriceAdjustment{
			Code:     "Free item",
			Name:     "one free item",
			Priority: 2.1,
			Labels:   []string{checkout.ConstLabelDiscount},
			PerItem: map[string]float64{
				"2": -150,
			},
		},

		// in this case it would be 6% tax for full cart
		checkout.StructPriceAdjustment{
			Code:      "Country-State",
			Name:      "Tax",
			Amount:    6,
			IsPercent: true,
			Priority:  2.5,
			Labels:    []string{checkout.ConstLabelTax},
			PerItem:   map[string]float64{},
		},

		// tax applied in different way for different types of products
		checkout.StructPriceAdjustment{
			Code:     "Product-Addings",
			Name:     "Tax",
			Priority: 2.51,
			Labels:   []string{checkout.ConstLabelTax},
			PerItem: map[string]float64{
				"0": 5,  // 5$ on full cart
				"1": 10, // 10$ on product1
				"2": 7,  // 7$ on product2
			},
		},

		// gift card for a full cart
		checkout.StructPriceAdjustment{
			Code:     "gift-card1",
			Name:     "gift-card",
			Amount:   -999.9999,
			Priority: 3.1,
			Labels:   []string{checkout.ConstLabelGiftCard},
		},
	}

	const DEBUG = false // allows to print values later on in this test

	currentCheckout := new(DefaultCheckout)

	// prevent from executing of calculate function
	currentCheckout.calculationDetailTotals = make(map[int]map[string]float64)
	currentCheckout.calculateFlag = true

	for index, priceAdjustment := range priceAdjustments {
		currentCheckout.applyPriceAdjustment(priceAdjustment)

		// after PA applied it added to checkout internal array with updated amount
		appliedPriceAdjustment := currentCheckout.priceAdjustments[index]
		if DEBUG {
			fmt.Println(currentCheckout.calculateAmount, appliedPriceAdjustment.Amount)
			fmt.Println(currentCheckout.calculationDetailTotals)
		}
		// this value would be the total amount that was applied to grand total
		if appliedPriceAdjustment.Amount == 0 {
			t.Error("Amount is equal to 0")
		}
	}
	if DEBUG {
		fmt.Println(currentCheckout.GetPriceAdjustments(""))
		fmt.Println("Subtotal: ", currentCheckout.GetSubtotal())
		fmt.Println("Shipping: ", currentCheckout.GetShippingAmount())
		fmt.Println("Discount: ", currentCheckout.GetDiscountAmount())
		fmt.Println("Tax: ", currentCheckout.GetTaxAmount())

		fmt.Println("Grandtotal: ", currentCheckout.GetGrandTotal())

	}

	total := currentCheckout.GetSubtotal() + currentCheckout.GetShippingAmount() + currentCheckout.GetDiscountAmount() + currentCheckout.GetTaxAmount()
	if utils.RoundPrice(total) != utils.RoundPrice(currentCheckout.GetGrandTotal()) {
		t.Error("Total obteined from adding part elements is not equal to grandtotal")
	}

	if currentCheckout.calculateAmount < 0 {
		t.Error("Amount is lesser then 0")
	}
}

/*
This output is generated with
const DEBUG = true
defined earlier in this test

400 400
map[2:map[GT:300 ST:300] 1:map[GT:100 ST:100] 0:map[ST:400 GT:400]]
397 -3
map[1:map[GT:99 ST:100 SPA:-1] 0:map[SPA:-3 GT:397 ST:400] 2:map[GT:298 ST:300 SPA:-2]]
497 100
map[1:map[GT:99 ST:100 SPA:-1] 0:map[GT:497 ST:400 SPA:-3 SP:100] 2:map[GT:298 ST:300 SPA:-2]]
347 -150
map[1:map[GT:99 ST:100 SPA:-1] 0:map[GT:347 ST:400 SPA:-3 SP:100 D:-150] 2:map[D:-150 GT:148 ST:300 SPA:-2]]
367.82 20.82
map[0:map[T:20.82 GT:367.82 ST:400 SPA:-3 SP:100 D:-150] 2:map[GT:148 ST:300 SPA:-2 D:-150] 1:map[SPA:-1 GT:99 ST:100]]
389.82 22
map[1:map[GT:109 ST:100 SPA:-1 T:10] 0:map[SP:100 D:-150 T:42.82 GT:389.82 ST:400 SPA:-3] 2:map[ST:300 SPA:-2 D:-150 T:7 GT:155]]
0 -389.82
map[2:map[T:7 GT:155 ST:300 SPA:-2 D:-150] 1:map[GT:109 ST:100 SPA:-1 T:10] 0:map[GC:-389.82 GT:0 ST:400 SPA:-3 SP:100 D:-150 T:42.82]]

[
{ST ST 1 400 false [ST] map[1:100 2:300]}
{sale_price Sale Price 1.1 -3 false [SPA] map[1:-1 2:-2]}
{default Flat Rate 2 100 false [SP] map[]}
{Free item one free item 2.1 -150 false [D] map[2:-150]}
{Country-State Tax 2.5 20.82 true [T] map[]}
{Product-Addings Tax 2.51 22 false [T] map[0:5 1:10 2:7]}
{gift-card1 gift-card 3.1 -389.82 false [GC] map[]}
]

Subtotal:  400
Shipping:  100
Discount:  -542.8199999999999
Tax:  42.82
Grandtotal:  0
*/
