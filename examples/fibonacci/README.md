# fibonacci

> An example of how to use mdcode

In this example, the examples consist of 4 code blocks each. Of course, any markdown text can be included between individual code blocks. The code block describing the outline of the source file is an invisible code block, just like the test code required for testing.

If you look at the [source of this document](https://github.com/szkiba/mdcode/blob/master/examples/fibonacci/README.md?plain=1), you can see how the embedding is done.


**JavaScript**
<!--<script type="text/markdown">
```js file=fibonacci.js outline=true
// #region signature
    // #endregion
    // #region zero
    // #endregion
    // #region one
    // #endregion

    // #region regular
// #endregion

module.exports = fibonacci
```
</script>-->

```js file=fibonacci.js region=signature
function fibonacci(n) {
```

```js file=fibonacci.js region=zero
    if (n < 1) {
        return 0
    }
```

```js file=fibonacci.js region=one
    if (n == 1) {
        return 1
    }
```

```js file=fibonacci.js region=regular
    return fibonacci(n - 1) + fibonacci(n - 2)
}
```

<!--<script type="text/markdown">
```js file=fibonacci.test.js
const assert = require("node:assert");
const test = require("node:test");
const fibonacci = require("./fibonacci");

const testvect = [0, 1, 1, 2, 3, 5, 8, 13, 21, 34, 55, 89, 144, 233, 377, 610, 987, 1597, 2584, 4181];

test("fibonacci with test vector", (t) => {
    for (var i = 0; i < testvect.length; i++) {
        assert.equal(fibonacci(i), testvect[i]);
    }
});
```
</script>-->

**go**
<!--<script type="text/markdown">
```go file=fibonacci.go outline=true
package main

// #region signature
	// #endregion
	// #region zero
	// #endregion
	// #region one
	// #endregion

	// #region regular
// #endregion
```
</script>-->

```go file=fibonacci.go region=signature
func fibonacci(n uint64) uint64 {
```

```go file=fibonacci.go region=zero
	if n == 0 {
		return 0
	}
```

```go file=fibonacci.go region=one
	if n == 1 {
		return 1
	}
```

```go file=fibonacci.go region=regular
	return fibonacci(n-1) + fibonacci(n-2)
}

```

<!--<script type="text/markdown">
```js file=fibonacci_test.go
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
```
</script>-->

### Tasks

To try mdcode here are some easy tasks

#### extract

You can extract the example code and the test code using the following command:

```
mdcode extract
```

#### test

After extractiong, you can run the test using the following commands:

```
go test
node --test
```

#### update

After modifying the files, you can update the code blocks in the document using the following command:

```
mdcode update
```
