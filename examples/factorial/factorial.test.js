const assert = require("node:assert")
const test = require("node:test")

// #region function
function factorial(n) {
    if (n > 1) {
        return n * factorial(n - 1)
    }

    return 1
}
// #endregion

const testvect = [1, 1, 2, 6, 24, 120, 720, 5040, 40320, 362880, 3628800]

test("factorial with test vector", (t) => {
    for (var i = 0; i < testvect.length; i++) {
        assert.equal(factorial(i), testvect[i])
    }
})
