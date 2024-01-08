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
