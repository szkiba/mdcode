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
