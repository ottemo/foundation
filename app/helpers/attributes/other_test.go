package attributes

import "testing"


// TestExternalAttributes tests concurrent execution for an ExternalAttributes
func TestExternalAttributes(t *testing.T) {
	const concurrent = 9999
	finished := make(chan int)

	routines := concurrent
	for i:=0; i<routines; i++ {
		go func(i int) {
			ExampleExternalAttributes()

			finished <- i
		}(i)
	}

	for routines > 0 {
		<- finished
		routines--
	}
}