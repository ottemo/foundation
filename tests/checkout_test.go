package tests

import (
	"math/rand"
	"testing"
)

// benchmarks add to cart
func BenchmarkAddToCart(b *testing.B) {
	err := StartAppInTestingMode()
	if err != nil {
		b.Error(err)
	}

	err = MakeSureProductsCount(100)
	if err != nil {
		b.Error(err)
	}

	currentVisitor, err := GetRandomVisitor()
	if err != nil {
		b.Error(err)
	}

	currentCheckout, err := GetNewCheckout(currentVisitor)
	if err != nil {
		b.Error(err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		err = AddRandomProductsToCart(currentCheckout, 1)
		if err != nil {
			b.Error(err)
		}
	}
}

// benchmarks add to cart
func BenchmarkCheckoutSubmit(b *testing.B) {
	err := StartAppInTestingMode()
	if err != nil {
		b.Error(err)
	}

	err = MakeSureProductsCount(100)
	if err != nil {
		b.Error(err)
	}

	currentVisitor, err := GetRandomVisitor()
	if err != nil {
		b.Error(err)
	}

	currentCheckout, err := GetNewCheckout(currentVisitor)
	if err != nil {
		b.Error(err)
	}

	err = AddRandomProductsToCart(currentCheckout, rand.Intn(8)+1)
	if err != nil {
		b.Error(err)
	}

	err = RandomizeShippingAndBillingAddresses(currentCheckout)
	if err != nil {
		b.Error(err)
	}

	err = UpdateShippingAndPaymentMethods(currentCheckout)
	if err != nil {
		b.Error(err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err = currentCheckout.Submit()
		if err != nil {
			b.Error(err)
		}
	}
}

// benchmarks order creation from zero
func BenchmarkFullCheckout(b *testing.B) {
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
		err := FullCheckout()
		if err != nil {
			b.Error(err)
		}
	}
}
