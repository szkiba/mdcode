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
