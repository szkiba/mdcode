const assert = require("node:assert");
const test = require("node:test");

test("hello", (t) => {
    console.log = function (message) {
        assert.equal(message, "Hello, Testable World!");
    };
    require("./hello.js");
});
