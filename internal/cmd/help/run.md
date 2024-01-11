Extract code blocks to the file system and run shell commands on them

The code blocks are written to the file named in the `file` metadata.

The code block may include `region` metadata, which contains the name of the region. In this case, the code block is written to the appropriate part of the file marked with the `#region` comment.

The optional argument of the `mdcode run` command is the name of the markdown file. If it is missing, the `README.md` file in the current directory (if it exists) is processed.

This can be followed by a double dash (`--`) and then the shell command line to be executed (even a complex command, such as `for`).

Alternatively, the commands to be executed can be embedded in a code block in the document. In this case, the language must be `sh` and it is necessary to name the code block with the metadata `name`. The name of the code block containing the commands can be specified with the `--name` flag (if not, the first code block containing the `sh` language and `name` metadata will be executed).

Code blocks are extracted to a temporary directory. This directory will be the current directory when running the commands. The temporary directory is deleted after executing the commands (deletion can be prevented by using the `--keep` flag). Instead of a temporary directory, the name of the directory to be used can be specified with the `--dir` flag. In this case, of course, the directory is not deleted after executing the commands.
