package utils

import (
	"testing"
)

func TestRoundPrice(t *testing.T) {
	if x := RoundPrice(32.87000000001); x != 32.87 {
		t.Error("incorect result:", x)
	}
	if x := RoundPrice(-32.87000000001); x != -32.87 {
		t.Error("incorect result:", x)
	}

	if x := RoundPrice(32.0); x != 32.0 {
		t.Error("incorect result:", x)
	}
	if x := RoundPrice(-32.0); x != -32.0 {
		t.Error("incorect result:", x)
	}

	if x := RoundPrice(32.865); x != 32.87 {
		t.Error("incorect result:", x)
	}
	if x := RoundPrice(-32.865); x != -32.87 {
		t.Error("incorect result:", x)
	}

	if x := RoundPrice(32.8699999999995); x != 32.87 {
		t.Error("incorect result:", x)
	}
	if x := RoundPrice(-32.8699999999995); x != -32.87 {
		t.Error("incorect result:", x)
	}

	if x := RoundPrice(0.0045); x != 0 {
		t.Error("incorect result:", x)
	}
	if x := RoundPrice(-0.0045); x != 0 {
		t.Error("incorect result:", x)
	}

	if x := RoundPrice(0.005); x != 0.01 {
		t.Error("incorect result:", x)
	}
	if x := RoundPrice(-0.005); x != -0.01 {
		t.Error("incorect result:", x)
	}

	if x := 1-RoundPrice(0.000000000000009); x != 1 {
		t.Error("incorect result:", x)
	}
	if x := RoundPrice(0.000000000000009)-1; x != -1 {
		t.Error("incorect result:", x)
	}
}
