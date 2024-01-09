# hello

> An example of how to use mdcode

In this example, the code blocks contain entire files. The test code required for testing is in invisible code blocks.

If you look at the [source of this document](https://github.com/szkiba/mdcode/blob/master/examples/hello/README.md?plain=1), you can see how the embedding is done.

**JavaScript**
```js file=hello.js
console.log("Hello, Testable World!")
```

<!--<script type="text/markdown">
```js file=hello.test.js
const assert = require("node:assert")
const test = require("node:test")

test("hello", (t) => {
    console.log = function (message) {
        assert.equal(message, "Hello, Testable World!")
    }
    require("./hello.js")
})
```
</script>-->

**go**
```go file=hello.go
package main

import "fmt"

func main() {
	fmt.Println("Hello, Testable World!")
}
```

<!--<script type="text/markdown">
```go file=hello_test.go
package main

import (
	"io"
	"os"
	"testing"
)

func Test_main(t *testing.T) {
	orig := os.Stdout

	reader, writer, err := os.Pipe()
	if err != nil {
		t.Error(err)
	}

	os.Stdout = writer

	main()

	if err = writer.Close(); err != nil {
		t.Error(err)
	}

	out, err := io.ReadAll(reader)
	if err != nil {
		t.Error(err)
	}

	os.Stdout = orig

	const expected = "Hello, Testable World!\n"

	if string(out) != expected {
		t.Errorf("\nexpected: %s\nactual:   %s\n", expected, string(out))
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
