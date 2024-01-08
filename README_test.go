package main

import "testing"

// #region function
func factorial(n uint64) uint64 {
	if n > 1 {
		return n * factorial(n-1)
	}

	return 1
}

// #endregion

func Test_factorial(t *testing.T) {
	t.Parallel()

	testvect := []uint64{1, 1, 2, 6, 24, 120, 720, 5040, 40320, 362880, 3628800}

	for idx, expected := range testvect {
		if actual := factorial(uint64(idx)); actual != expected {
			t.Errorf("factorial(%d) should be %d but got %d", idx, expected, actual)
		}
	}
}
