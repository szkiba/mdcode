Dump markdown code blocks as tar archive

Creating a tar format archive from code blocks that meet the filtering criteria. By default, it writes to standard output, but it can also be directed to file with the `--output` flag.

A base directory can be specified with the `--dir` flag, all files will be created under this directory.

The optional argument of the `mdcode dump` command is the name of the markdown file. If it is missing, the `README.md` file in the current directory (if it exists) is processed.
