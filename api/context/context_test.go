package context

import (
	"math/rand"
	"sync"
	"testing"
	"time"
)

func TestStackContextMixedTree(t *testing.T) {
	var A, B, C, D func(testValue interface{})

	const testKey = "test"

	A = func(testValue interface{}) {
		MakeContext(func() {
			if context := GetContext(); context != nil {
				context[testKey] = testValue
				context[testKey+"A"] = testValue
			} else {
				t.Logf("no context in A %v", testValue)
			}

			B(testValue)
		})
	}

	B = func(testValue interface{}) {
		time.Sleep(time.Duration(rand.Intn(100)) * time.Nanosecond)
		if context := GetContext(); context != nil {
			context[testKey+"B"] = testValue
			if context[testKey] != testValue {
				t.Fatalf("%v != %v, A = %v", context[testKey], testValue, context[testKey+"A"])
			}
		} else {
			t.Logf("no context in B %v", testValue)
		}

		C(testValue)
	}

	C = func(testValue interface{}) {

		MakeContext(func() {
			if context := GetContext(); context != nil {
				context[testKey] = testValue
				context[testKey+"C"] = testValue
			} else {
				t.Logf("no context in C %v", testValue)
			}

			D(testValue)
		})

		D(testValue)
	}

	D = func(testValue interface{}) {
		if context := GetContext(); context != nil {
			context[testKey+"D"] = testValue
			if context[testKey] != testValue {
				t.Fatalf("%v != %v, A = %v, B = %v, C = %v", context[testKey], testValue, context[testKey+"A"], context[testKey+"B"], context[testKey+"C"])
			}
		} else {
			t.Logf("no context in D %v", testValue)
		}
	}

	A(1)

	var wg sync.WaitGroup
	for i := 0; i < 9999; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			A(time.Now().Nanosecond())
		}()
	}
	wg.Wait()
}
