package context

import (
	"sync"
	"testing"
	"time"
	"math/rand"
)

func TestMixedTree(t *testing.T) {
	var A, B, C func(testValue interface{})

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
		// B(testValue)
	}

	B = func(testValue interface{}) {
		time.Sleep(time.Duration(rand.Intn(1000)) * time.Nanosecond)
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
		if context := GetContext(); context != nil {
			context[testKey+"C"] = testValue
			if context[testKey] != testValue {
				t.Fatalf("%v != %v, A = %v, B = %v", context[testKey], testValue, context[testKey+"A"], context[testKey+"B"])
			}
		} else {
			t.Logf("no context in C %v", testValue)
		}
		time.Sleep(time.Second)
	}

	  //A(1)

	var wg sync.WaitGroup
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			A(time.Now().Nanosecond())
		}()
	}
	wg.Wait()
}
