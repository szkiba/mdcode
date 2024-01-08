# factorial

> An example of how to use mdcode

In this example, the code blocks only contain one region (named `function`). The test code required for testing is in invisible code blocks.

If you look at the [source of this document](https://github.com/szkiba/mdcode/blob/master/examples/factorial/README.md?plain=1), you can see how the embedding is done.

**JavaScript**
<!--<script type="text/markdown">
```js file=factorial.test.js outline=true
const assert = require("node:assert");
const test = require("node:test");

// #region function
// #endregion

const testvect = [1, 1, 2, 6, 24, 120, 720, 5040, 40320, 362880, 3628800];

test("factorial with test vector", (t) => {
    for (var i = 0; i < testvect.length; i++) {
        assert.equal(factorial(i), testvect[i]);
    }
})
```
</script>-->

```js file=factorial.test.js region=function
function factorial(n) {
    if (n > 1) {
        return n * factorial(n - 1)
    }

    return 1
}

```

**go**
<!--<script type="text/markdown">
```go file=factorial_test.go outline=true
package main

import "testing"

// #region function
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
```
</script>-->

```go file=factorial_test.go region=function
func factorial(n uint64) uint64 {
	if n > 1 {
		return n * factorial(n-1)
	}

	return 1
}

```


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
