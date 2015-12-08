package context

import (
	"sync"
	"testing"
	"time"
)

func TestMixedTree(t *testing.T) {
	var A, B, C func(testValue interface{})

	const testKey = "test"

	A = func(testValue interface{}) {
		MakeContext(func() {
			if context := GetContext(); context != nil {
				context[testKey] = testValue
			}

			B(testValue)
		})
		B(testValue)
	}

	B = func(testValue interface{}) {
		if context := GetContext(); context != nil {
			if context[testKey] != testValue {
				t.Fatalf("%v != %v", context[testKey], testValue)
			}
		}

		C(testValue)
	}

	C = func(testValue interface{}) {
		if context := GetContext(); context != nil {
			if context[testKey] != testValue {
				t.Fatalf("%v != %v", context[testKey], testValue)
			}
		}
		time.Sleep(time.Second)
	}

	A(1)

	var wg sync.WaitGroup
	for i := 0; i < 1000; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			A(time.Now().Nanosecond())
		}()
	}
	wg.Wait()
}
