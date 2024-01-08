package main

// #region signature
func fibonacci(n uint64) uint64 {
	// #endregion
	// #region zero
	if n == 0 {
		return 0
	}
	// #endregion
	// #region one
	if n == 1 {
		return 1
	}
	// #endregion

	// #region regular
	return fibonacci(n-1) + fibonacci(n-2)
}

// #endregion
