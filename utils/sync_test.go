package utils

import (
	"testing"
	"math/rand"
)

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
