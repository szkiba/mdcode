# ｍｄｃｏｄｅ

**Markdown code block authoring tool**

The `mdcode` command-line tool allows code blocks embedded in a markdown document to be developed in the usual way. During the development of the code blocks, the usual tools and methods can be used. This makes the embedded codes testable, which is especially important for example codes. There is no worse developer experience than a faulty sample code.

Here is a simple example code for a factorial calculation:

**go**
<!--<script type="text/markdown">
```go file=README_test.go outline=true
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

```go file=README_test.go region=function
func factorial(n uint64) uint64 {
	if n > 1 {
		return n * factorial(n-1)
	}

	return 1
}

```

**JavaScript**

<!--<script type="text/markdown">
```js file=README.test.js outline=true
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

```js file=README.test.js region=function
function factorial(n) {
    if (n > 1) {
        return n * factorial(n - 1)
    }

    return 1
}

```

At first glance, there is nothing special about this. *However, these code blocks are testable!* 

This document also contains the code required for testing in invisible code blocks. If you look at the [source of this document](https://github.com/szkiba/mdcode/blob/master/README.md?plain=1), you can see how easy it is to embed code blocks (even invisibly) using `mdcode`.

Code blocks embedded in this document can be saved to files using the [`mdcode extract`](#mdcode-extract) command. A `README_test.go` and a `README.test.js` file will be created in the current directory. After modification, the code blocks can be updated from these files to the document using the [`mdcode update`](#mdcode-update) command.

More examples can be found in the [examples](examples/) directory.

### Features

- include source files as code blocks in the markdown document
- update the code blocks in the markdown document
- save markdown code blocks to source files
- supports source file fragments using `#region` comments
- supports invisible (not rendered) code blocks
- allows you to add metadata to code blocks
- supports any programming language
- dump code blocks as tar archive

### Use Cases

**Develop the example/tutorial codes as you would any other code**
  - use your favorite IDE, toolchain
  - use any test framework for testing
  - integrate example code testing into the build process
  - use [`mdcode update`](#mdcode-update) to update example code in markdown documents

**Write example code directly in the markdown documemts**
  - use [`mdcode extract`](#mdcode-extract) to extract code blocks and save them to files
  - use any test framework for testing
  - integrate example code testing into the build process

**Create a self-contained markdown tutorial document**
  - use [`mdcode update`](#mdcode-update) to embed source fragments
  - use [`mdcode update`](#mdcode-update) to embed additional files (package.json, go.mod, etc.) as invisible code blocks
  - use [`mdcode extract`](#mdcode-extract) to extract working examples from the markdown documemt

**Save all examples for later use**
  - use [`mdcode dump`](#mdcode-dump) to create tar archive from code blocks

### Install

Precompiled binaries can be downloaded and installed from the [Releases](https://github.com/szkiba/mdcode/releases) page.

If you have a go development environment, the installation can also be done with the following command:

```
go install github.com/szkiba/mdcode@latest
```

It can even be run without installation using the following command:

```
go run github.com/szkiba/mdcode@latest
```

### Usage

Check [CLI Reference](#cli-reference) section for detailed command line usage.

## Concepts

### Metadata

<!-- #region metadata -->
Metadata can be specified for code blocks. These metadata can be used to modify the operation of the subcommands (for example, they can be used for filtering).

The [CommonMark specification](https://spec.commonmark.org/current/) allows the use of a so-called [info-string](https://spec.commonmark.org/current/#info-string) in the fenced code block. The first word of the *info-string* typically indicates the programming language, the meaning of the remaining part is not defined by the specification.

`mdcode` uses the part after the first word of the *info-string* to specify metadata. The metadata can be entered in JSON format and in a simple, space-separated `name="value"` format list (where the use of quotation marks is only necessary for values containing spaces). The latter form is more readable, but the JSON format is more portable.

Example name="value" list metadata:

    ```js file=sample.js region=factorial

    ```

Example JSON metadata:

    ```js {"file":"sample.js","region":"factorial"}

    ```

Metadata used by `mdcode`:

name      | description
----------|-------------------------------------------------
`file`    | name of the file assigned to the code block
`region`  | name of region within file (if any)
`outline` | true if the code block is an outline of the file

The only mandatory metadata is `file`.
<!-- #endregion metadata -->

### Filtering

<!-- #region filtering -->
By default, `mdcode` work with all code blocks in a markdown document. It is possible to filter code blocks based on programming language or metadata. In this case, `mdcode` ignore code blocks that do not meet the filter criteria.

A language filter pattern can be specified using the `--lang` flag. Then only code blocks with a language matching the pattern will be processed. For example, filtering for code blocks containing JavaScript code:

	mdcode --lang js
		
A file name filtering pattern can be specified using the `--file` flag. Then only code blocks with `file` metadata matching the pattern will be processed. For example, filtering for code blocks containing the file named `examples/foo.js` (or parts of it):

	mdcode --file examples/foo.js

The `--meta` flag can be used to specify an arbitrary metadata filtering pattern. Then only code blocks with metadata matching the pattern are processed. For example, filtering for code blocks that have metadata named `name` and its value is `simple`:

	mdcode --meta name=simple

Specifying several different filter criteria (e.g. language and metadata, or two different metadata) each criteria must be met (and relation).

Standard glob patterns can be used in programming language and metadata filter criteria.

pattern          | match
-----------------|--------------------------------------------------------------
`*`              | matches any sequence of non-separator characters
`**`             | matches any sequence of characters
`?`              | matches any single non-separator character
`[` range `]`    | character in range
`[` `!` range `]`| character not in range 
`{` list `}`     | matches any of comma-separated (without spaces) patterns
c                | matches character c (c != `*`, `**`, `?`, `\`, `[`, `{`, `}`)
`\` c            | matches character c

range     | match
----------|-----------------------------------------
c         | matches character c (c != `\`, `-`, `]`)
`\` c     | matches character c
lo `-` hi | matches character c for lo <= c <= hi

Examples of filter pattern use:

    mdcode extract --meta file='examples/**/*.go'
    mdcode extract --lang '{go,js}'

Filtering with frequently used metadata can also be done using dedicated flags.

flag             | shorthand    | equivalent
-----------------|--------------|----------------------
`--file pattern` | `-f pattern` | `--meta file=pattern`
<!-- #endregion filtering -->

### Regions

<!-- #region regions -->
In addition to embedding entire files, `mdcode` supports the use of file regions. Named regions can be used in the source code of any programming language. The beginning of the region is marked by a comment line with the content `#region name` and the end by a comment line with the content `#endregion`.

For example, in the case of programming languages using C-style line comments (C, C++, Java, JavaScript, go, etc.):

    // #region common

    // #endregion

In the case of programming languages using shell-style line comments (Python, sh, bash, etc.):

    # #region common

    # #endregion

Or if only block comments can be used (CSS):

    /* #region common */

    /* #endregion */

Regions marked in this way are used by IDEs to collapse parts of the source code.

In the case of `mdcode`, regions can be referenced with the `region` metadata. If a region is specified for a code block, the subcommand (update or extract) applies only to the specified region of the file. That is, the update command only embeds the specified region from the file to the markdown document, and the extract command overwrites only the specified region in the file.

`mdcode` can handle regions in any programming language, the only requirement is that the comment indicating the beginning and end of the region is placed in a separate line containing only the given comment.
<!-- #endregion regions -->

### Invisible

<!-- #region invisible -->
It is possible to use invisible code blocks. This is useful, for embedding test code or additional files in the markdown document. The invisible code block is also useful if you want to embed the entire file, but only want to display certain parts of it.

A markdown document can contain HTML elements. Unknown or unsupported HTML elements are usually not rendered by markdown renderers. Taking advantage of this, `mdcode` supports hiding code blocks using the standard `<script>` HTML element:

    <script type="text/markdown">
    ```js file=sample.js region=factorial

    ```
    </script>

Unfortunately, the GitHub markdown renderer renders the content of unsupported HTML elements as text. Therefore, `mdcode` also supports the use of a `<script>` element surrounded by an HTML comment to hide a code block.

    <!--<script type="text/markdown">
    ```js file=sample.js region=factorial

    ```
    </script>-->

*It is important to note that the opening character of the comment and the opening tag of the script element must be placed on the same line. Similarly, the closing tag of the script element and the closing tag of the comment must also be placed on the same line.*
<!-- #endregion invisible -->

**Highlighting invisible code block**

To highlight markdown within an HTML script element, the 
[Markdown Script Tag](https://marketplace.visualstudio.com/items?itemName=sissel.markdown-script-tag) extension can be used in Visual Studio Code.

### Outline

<!-- #region outline -->
When using regions, only parts of the source file are embedded in the markdown document. If we want to create a self-contained markdown document, the `true` value of the `outline` metadata can be used for this purpose.

In this case, only parts of the source file other than the region comments are embedded in the markdown document (and the empty region comments).

The outline flag is typically used in an invisible code block preceding the visible regions. Since the `mdcode extract` command processes the code blocks sequentially, the code block marked with an `outline` first overwrites the file, then the code blocks containing the named regions are inserted in their place.
<!-- #endregion outline -->

Here is the invisible outline code block at the beginning of this document:

    <!--<script type="text/markdown">
    ```go file=README_test.go outline=true
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

And here is the code block containing the corresponding region reference:

    ```go file=README_test.go region=function
    func factorial(n uint64) uint64 {
      if n > 1 {
        return n * factorial(n-1)
	    }

	    return 1
    }
    ```

## Development

### Tasks

This section contains a description of the tasks performed during development. If you have the [xc (Markdown defined task runner)](https://github.com/joerdav/xc) command-line tool, individual tasks can be executed simply by using the `xc task-name` command.

<details><summary>Click to expand</summary>

#### lint

Run the static analyzer.

```
golangci-lint run
```

#### test

Run the tests.

```
go test -count 1 -race -coverprofile=build/coverage.txt ./...
```

#### coverage

View the test coverage report.

```
go tool cover -html=build/coverage.txt
```

#### build

Build the executable binary.

This is the easiest way to create an executable binary (although the release process uses the goreleaser tool to create released versions).

```
go build -ldflags="-w -s" -o build/mdcode .
```

#### snapshot

Creating an executable binary with a snapshot version.

The goreleaser command-line tool is used during the release process. During development, it is advisable to create binaries with the same tool from time to time.

```
goreleaser build --snapshot --clean --single-target -o build/mdcode
```

#### doc

Updating the documentation.

Some parts of the documentation, such as the [CLI Reference](#cli-reference), example codes, are automatically generated.

```
go generate
```

#### doc-test

Testing the generated documentation.

```
go generate
go test ./...
node --test
```

#### clean

Delete the build directory.

```
rm -rf build
```

#### all

Run all tasks.

Requires: lint,test,doc,doc-test,build,snapshot

</details>

## CLI Reference

This chapter contains the reference documentation for the command line interface. In addition to the commands described here, additional help topics are available from the command line. To view them, the topic name must be specified as a subcommand.

<!-- #region cli -->
Additional help topics:
* `mdcode filtering` - [Pattern based filtering](#filtering)
* `mdcode invisible` - [Invisible code blocks](#invisible)
* `mdcode metadata` - [Code block metadata](#metadata)
* `mdcode outline` - [Embedding the file structure](#outline)
* `mdcode regions` - [Handling file regions](#regions)
---

## mdcode

Markdown code block authoring tool

### Synopsis

Lists the code blocks (with file metadata) from the markdown document.

The optional argument of the `mdcode` command is the name of the markdown file. If it is missing, the `README.md` file in the current directory (if it exists) is processed.


```
mdcode [filename] [flags]
```

### Flags

```
  -f, --file strings          file filter (default [?*])
  -h, --help                  help for mdcode
      --json                  generate JSON output
  -l, --lang strings          language filter (default [?*])
  -m, --meta stringToString   metadata filter (default [])
  -o, --output string         output file (default: standard output)
```

### SEE ALSO

* [mdcode dump](#mdcode-dump)	 - Dump markdown code blocks
* [mdcode extract](#mdcode-extract)	 - Extract markdown code blocks to the file system
* [mdcode update](#mdcode-update)	 - Update markdown code blocks from the file system

---
## mdcode dump

Dump markdown code blocks

### Synopsis

Dump markdown code blocks as tar archive

Creating a tar format archive from code blocks that meet the filtering criteria. By default, it writes to standard output, but it can also be directed to file with the `--output` flag.

A base directory can be specified with the `--dir` flag, all files will be created under this directory.

The optional argument of the `mdcode dump` command is the name of the markdown file. If it is missing, the `README.md` file in the current directory (if it exists) is processed.


```
mdcode dump [filename] [flags]
```

### Flags

```
  -d, --dir string      base directory name (default ".")
  -h, --help            help for dump
  -o, --output string   output file (default: standard output)
  -q, --quiet           suppress the status output
```

### Global Flags

```
  -f, --file strings          file filter (default [?*])
  -l, --lang strings          language filter (default [?*])
  -m, --meta stringToString   metadata filter (default [])
```

### SEE ALSO

* [mdcode](#mdcode)	 - Markdown code block authoring tool

---
## mdcode extract

Extract markdown code blocks to the file system

### Synopsis

Writing all code blocks matching the filter criteria to the file system

The code blocks are written to the file named in the `file` metadata. The file name is relative to the current directory or to the directory specified with the `--dir` flag.

The code block may include `region` metadata, which contains the name of the region. In this case, the code block is written to the appropriate part of the file marked with the `#region` comment.

The optional argument of the `mdcode extract` command is the name of the markdown file. If it is missing, the `README.md` file in the current directory (if it exists) is processed.


```
mdcode extract [filename] [flags]
```

### Flags

```
  -d, --dir string   base directory name (default ".")
  -h, --help         help for extract
  -q, --quiet        suppress the status output
```

### Global Flags

```
  -f, --file strings          file filter (default [?*])
  -l, --lang strings          language filter (default [?*])
  -m, --meta stringToString   metadata filter (default [])
```

### SEE ALSO

* [mdcode](#mdcode)	 - Markdown code block authoring tool

---
## mdcode update

Update markdown code blocks from the file system

### Synopsis

Update all code blocks that meet the filter criteria from the file system

The code blocks are read from the file named in the `file` metadata. The file name is relative to the current directory or to the directory specified with the `--dir` flag.

The code block may include `region` metadata, which contains the name of the region. In this case, the code block is read from the appropriate part of the file marked with the `#region` comment.

The optional argument of the `mdcode update` command is the name of the markdown file. If it is missing, the `README.md` file in the current directory (if it exists) is processed.


```
mdcode update [filename] [flags]
```

### Flags

```
  -d, --dir string   base directory name (default ".")
  -h, --help         help for update
  -q, --quiet        suppress the status output
```

### Global Flags

```
  -f, --file strings          file filter (default [?*])
  -l, --lang strings          language filter (default [?*])
  -m, --meta stringToString   metadata filter (default [])
```

### SEE ALSO

* [mdcode](#mdcode)	 - Markdown code block authoring tool

<!-- #endregion cli -->