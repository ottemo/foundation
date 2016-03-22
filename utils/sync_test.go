package utils

import (
	"testing"
	"math/rand"
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

				m := GetMutex(x) // synchronization
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
