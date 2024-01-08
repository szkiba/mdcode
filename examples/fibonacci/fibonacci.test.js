const assert = require("node:assert");
const test = require("node:test");
const fibonacci = require("./fibonacci");

const testvect = [0, 1, 1, 2, 3, 5, 8, 13, 21, 34, 55, 89, 144, 233, 377, 610, 987, 1597, 2584, 4181];

test("fibonacci with test vector", (t) => {
    for (var i = 0; i < testvect.length; i++) {
        assert.equal(fibonacci(i), testvect[i]);
    }
});
