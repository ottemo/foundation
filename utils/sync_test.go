package utils

import (
	"testing"
	"math/rand"
	"fmt"
)

func BenchmarkPtrMapAccess(b *testing.B) {
	var i uintptr
	x := make(map[uintptr]int)
	for i=0; i<999999; i++ {
		for j:=1; j<rand.Intn(10); j++ {
			i++
		}
		x[i] = 1
	}

	b.ResetTimer()
	for i=0; i<999999; i++ {
		if val, ok := x[i]; ok {
			x[i-1]=val
		}
	}
}

func BenchmarkInterfaceMapAccess(b *testing.B) {
	var i int
	x := make(map[interface{}]int)
	for i=0; i<999999; i++ {
		switch i%3 {
		case 0:
			x[i] = 1
		case 1:
			x[string(i)] = 1
		case 2:
			x[float64(i)] = 1
		}
	}

	b.ResetTimer()
	for i=0; i<999999; i++ {
		if val, ok := x[i]; ok {
			x[i-1]=val
		}
	}
}

// TestLock makes massive attack to the same map from different go-routines which should generate
// "fatal error: concurrent map read and map write", without synchronization
func TestLock(t *testing.T) {
	const scatter = 10;
	x := make(map[int]map[int]float64)

	// m := GetMutex("x")
	// var m sync.Mutex

	for i:=0; i<scatter; i++ {
		x[i] = make(map[int]float64)
		for j:=0; j<scatter; j++ {
			x[i][j] = 0.0;
		}
	}

	finished := make(chan int)
	routines := 9999
	for i:=0; i<routines; i++ {
		go func(i int) {
			acts := rand.Intn(999)
			for j:=0; j<acts; j++ {
				key1 := rand.Intn(scatter)
				key2 := rand.Intn(scatter)

				m, err := SyncMutex(x) // synchronization
				if err != nil {
					t.Error(err)
				}
				m.Lock()

				oldValue := x[key1][key2]
				x[key1][key2] = oldValue + rand.Float64()

				m.Unlock() // synchronization

			}
			finished <- i
		}(i)
	}

	for routines > 0 {
		<- finished
		routines--
	}
}


func TestSyncSet(t *testing.T) {

	const concurrent = 9999
	finished := make(chan int)

	// Test 1: slice access
	A := make([][]int, 0, 10)

	routines := concurrent
	for i:=0; i<routines; i++ {
		go func(i int) {
			err := SyncSet(&A, 1, -1, -1)
			if err != nil {
				t.Error(err)
			}
			finished <- i
		}(i)
	}

	for routines > 0 {
		<- finished
		routines--
	}

	fmt.Println(A)
	return;

	if len(A) != concurrent || A[concurrent-1][0] != 1 {
		t.Error("unexpected result:",
			"len(A) = ", len(A),
			", A[concurrent-1][0] = ", A[concurrent-1][0])
	}

	// Test 2: map access
	B := make(map[string]map[int]map[bool]int)

	routines = concurrent
	for i:=0; i<routines; i++ {
		setter := func(old int) int {
			return old + 1
		}

		go func(i int) {
			err := SyncSet(B, setter, "a", i, true)
			if err != nil {
				t.Error(err)
			}

			err = SyncSet(B, setter, "b", 0, false)
			if err != nil {
				t.Error(err)
			}

			finished <- i
		}(i)
	}

	for routines > 0 {
		<- finished
		routines--
	}

	if len(B["a"]) != concurrent ||
		B["a"][concurrent-1][true] != 1 ||
		B["b"][0][false] != concurrent {

		t.Error("unexpected result: concurrent =", concurrent,
			", len(B[\"a\"]) =", len(B["a"]),
			", B[\"a\"][concurrent-1][true] =", B["a"][0][true],
			", B[\"b\"][0][false] = ", B["b"][0][false])
	}
}