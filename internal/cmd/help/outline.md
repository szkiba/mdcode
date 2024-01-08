When using regions, only parts of the source file are embedded in the markdown document. If we want to create a self-contained markdown document, the `true` value of the `outline` metadata can be used for this purpose.

In this case, only parts of the source file other than the region comments are embedded in the markdown document (and the empty region comments).

The outline flag is typically used in an invisible code block preceding the visible regions. Since the `mdcode extract` command processes the code blocks sequentially, the code block marked with an `outline` first overwrites the file, then the code blocks containing the named regions are inserted in their place.
