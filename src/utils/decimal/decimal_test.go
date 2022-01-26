package decimal

import "testing"

func TestDecimal(t *testing.T) {
	t.Log(NewFromFloat(0.00022875580134782607).StringFixed(2))
}
