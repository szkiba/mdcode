Update all code blocks that meet the filter criteria from the file system

The code blocks are read from the file named in the `file` metadata. The file name is relative to the current directory or to the directory specified with the `--dir` flag.

The code block may include `region` metadata, which contains the name of the region. In this case, the code block is read from the appropriate part of the file marked with the `#region` comment.

The optional argument of the `mdcode update` command is the name of the markdown file. If it is missing, the `README.md` file in the current directory (if it exists) is processed.
