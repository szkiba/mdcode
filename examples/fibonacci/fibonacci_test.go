package main

import "testing"

func Test_fibonacci(t *testing.T) {
	t.Parallel()

	testvect := []uint64{0, 1, 1, 2, 3, 5, 8, 13, 21, 34, 55, 89, 144, 233, 377, 610, 987, 1597, 2584, 4181}

	for idx, expected := range testvect {
		if actual := fibonacci(uint64(idx)); actual != expected {
			t.Errorf("fibonacci(%d) should be %d but got %d", idx, expected, actual)
		}
	}
}
