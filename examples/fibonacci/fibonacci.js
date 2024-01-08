// #region signature
function fibonacci(n) {
    // #endregion
    // #region zero
    if (n < 1) {
        return 0
    }
    // #endregion
    // #region one
    if (n == 1) {
        return 1
    }
    // #endregion

    // #region regular
    return fibonacci(n - 1) + fibonacci(n - 2)
}
// #endregion

module.exports = fibonacci
