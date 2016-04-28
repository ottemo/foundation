package checkout

import (
	"fmt"
	"github.com/ottemo/foundation/app/models/checkout"
	"testing"
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
			Name:     checkout.ConstLabelSubtotal,
			Priority: checkout.ConstCalculateTargetSubtotal,
			Labels:   []string{checkout.ConstLabelSubtotal},
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

	debug := false // allows to print values

	currentCheckout := new(DefaultCheckout)
	// this 'make' prevents from executing of calculate execution
	currentCheckout.calculationDetailTotals = make(map[int]map[string]float64)

	for index, priceAdjustment := range priceAdjustments {
		currentCheckout.applyPriceAdjustment(priceAdjustment)

		// after PA applied it added to checkout internal array with updated amount
		appliedPriceAdjustment := currentCheckout.priceAdjustments[index]
		if debug {
			fmt.Println(currentCheckout.calculateAmount, appliedPriceAdjustment.Amount)
			fmt.Println(currentCheckout.calculationDetailTotals)
		}
		// this value would be the total amount that was applied to grand total
		if appliedPriceAdjustment.Amount == 0 {
			t.Error("Amount is equal to 0")
		}
	}
	if debug {
		fmt.Println(currentCheckout.GetPriceAdjustments(""))
	}

	if currentCheckout.calculateAmount < 0 {
		t.Error("Amount is lesser then 0")
	}
}

/*
400 400
map[1:map[ST:100 GT:100] 0:map[GT:400 ST:400] 2:map[GT:300 ST:300]]
500 100
map[0:map[SP:100 GT:500 ST:400] 2:map[GT:300 ST:300] 1:map[GT:100 ST:100]]
350 -150
map[1:map[GT:100 ST:100] 0:map[GT:350 ST:250 SP:100] 2:map[GT:150 ST:150]]
371 21
map[1:map[GT:100 ST:100] 0:map[GT:371 ST:250 SP:100 T:21] 2:map[GT:150 ST:150]]
393 22
map[1:map[GT:110 ST:100 T:10] 0:map[GT:393 ST:250 SP:100 T:43] 2:map[ST:150 T:7 GT:157]]
0 -393
map[1:map[T:10 GT:110 ST:100] 0:map[GC:-393 GT:0 ST:250 SP:100 T:43] 2:map[GT:157 ST:150 T:7]]
[
{ST ST 1 400 false [ST] map[1:100 2:300]}
{default Flat Rate 2 100 false [SP] map[]}
{Free item ST 1 -150 false [ST] map[2:-150]}
{Country-State Tax 2.5 21 true [T] map[]}
{Product-Addings Tax 2.51 22 false [T] map[0:5 1:10 2:7]}
{gift-card1 gift-card 3.1 -393 false [GC] map[]}
]
*/
