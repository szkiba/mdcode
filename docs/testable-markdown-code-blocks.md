## Testable Markdown Code Blocks

*There is no worse developer experience than bad example code in the documentation*.

It is easy to make mistakes in code slightly more complicated than a *Hello World* program. In fact, in *Hello World* as well.

### How can you avoid it?

It is advisable to test the example codes in the documentation, like all code. Since the example codes also change over time, this testing should be done automatically, built into the documentation release process.

Basically, we can choose from two options if we want to test the example codes embedded in the documentation.

1. Create the example code in a separate file and embed this file, or part of it, in the documentation
2. The example code included in the documentation is extracted into a file, if necessary, supplemented to make it a valid source file

In both cases, finally, testing the example code is a normal testing task. So the example code can be tested using the standard test framework in the given programming environment. This testing step can be integrated into the build process, guaranteeing that the example codes work correctly.

### Required tooling

For the above, we need a tool that can embed source files or a part of them in the documentation or can create valid source files from the example codes in the documentation (supplementing them if necessary).

Today, the most popular format for developer documentation is the Markdown format. Example code pieces can be placed in a Markdown document using so-called fenced code blocks.

The [mdcode](https://github.com/szkiba/mdcode) command-line tool was created for the development and testing of Markdown fenced code blocks. It supports both example code testing/development processes: it is able to keep Markdown code blocks in sync with external source files, or part of them.

### Show me the code!

Using [mdcode](https://github.com/szkiba/mdcode) is super easy. The first line of the fenced code block, the so-called *info-string*, must be supplemented with some metadata after the programming language. Such metadata is the `file`, which specifies the name of the file belonging to the code block. The `mdcode extract` command writes the contents of the code block into this file, and the `mdcode update` command updates the contents of the code block from here.

Sticking to the *Hello World* example:

```js file=hello.js
console.log("Hello, Testable World!")
```

Here's how to add the file metadata:

~~~
```js file=hello.js
console.log("Hello, Testable World!")
```
~~~

The content of this code block can be saved to file (`hello.js`) using the `mdcode extract` command. Then the *Hello World* example can be tested with the following test code:

```js
const assert = require("node:assert")
const test = require("node:test")

test("hello", (t) => {
    console.log = function (message) {
        assert.equal(message, "Hello, Testable World!")
    }
    require("./hello.js")
})
```

This test code can be embedded in the markdown document invisibly:

~~~
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
~~~

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

It is advisable to include the test code in the document as well, because in this way the example codes can be tested by anyone with the following commands:

```sh
mdcode extract testable-markdown-code-blocks.md
node --test
```

Where `testable-markdown-code-blocks.md` is the name of the document (this post, can also be tested). Of course, any test framework can be used for testing.

### A more realistic example

A slightly more realistic case when the example code is only a fragment, for example a function:

<!--<script type="text/markdown">
```js file=factorial.test.js outline=true
const assert = require("node:assert")
const test = require("node:test")

// #region function
// #endregion

const testvect = [1, 1, 2, 6, 24, 120, 720, 5040, 40320, 362880, 3628800]

test("factorial with test vector", (t) => {
    for (var i = 0; i < testvect.length; i++) {
        assert.equal(factorial(i), testvect[i])
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

Many editors (such as Visual Studio Code) support the use of so-called `#region` comments for code folding. Regions marked in this way are usually named for readability. To avoid having to learn new, useless syntax, [mdcode](https://github.com/szkiba/mdcode) uses `#region` comments to mark embeddable parts of the source code. Simply add a `region` metadata to the fenced code block which contains the name of the region you want to embed.

This is how the code block above can be embedded from the region named `function` in the `factorial.test.js` file:

~~~
```js file=factorial.test.js region=function
function factorial(n) {
    if (n > 1) {
        return n * factorial(n - 1)
    }

    return 1
}
```
~~~

To do this, in the `factorial.test.js` file, you simply mark the part of the code you want to embed with a `#region` comment:

```js
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
```

And to make the document self-contained, you can embed the test code in the following way:

~~~
<!--<script type="text/markdown">
```js file=factorial.test.js outline=true
const assert = require("node:assert")
const test = require("node:test")

// #region function
// #endregion

const testvect = [1, 1, 2, 6, 24, 120, 720, 5040, 40320, 362880, 3628800]

test("factorial with test vector", (t) => {
    for (var i = 0; i < testvect.length; i++) {
        assert.equal(factorial(i), testvect[i])
    }
})
```
</script>-->
~~~

It is important to place this invisible code block before the visible code block in the document, because `mdcode extract` processes the code blocks in order.

It is worth noting that this code block only contains the outline of the file, so the  `#region` part is empty. `mdcode extract` will automatically fill in the part between the `#region` comment, and `mdcode update` will automatically delete it during embedding. Such outline code blocks are marked by the `outline` metadata value `true`.

### How do I edit the examples?

The example codes can also be edited within the markdown document. Most IDEs also provide syntax highlighting support within markdown fenced code blocks. However, we don't get actual programming language support (code formatting, IntelliSense, etc.)

It is more convenient to develop and test the example codes in actual source files. The parts you want to embed simply have to be marked with `#region` comments and can be updated in the markdown document using the `mdcode update` command.

The examples in this document and the corresponding tests can be extracted into files with the following command:

```sh
mdcode extract testable-markdown-code-blocks.md
```

As a result, the following files are created:

```
factorial.test.js
hello.js
hello.test.js 
```

These can be modified and then tested with the following command:

```sh
node --test
```


After that, the document can be updated using the command below:

```sh
mdcode update testable-markdown-code-blocks.md
```
