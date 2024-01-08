package main

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra/doc"
	"github.com/szkiba/mdcode/internal/cmd"
	"github.com/szkiba/mdcode/internal/region"
)

func checkerr(err error) {
	if err == nil {
		return
	}

	fmt.Fprintf(os.Stderr, "gendoc: error: %s\n", err)
	os.Exit(1)
}

func linkHandler(name string) string {
	link := strings.ReplaceAll(strings.TrimSuffix(name, ".md"), "_", "-")

	return "#" + link
}

func fprintf(out io.Writer, format string, args ...any) {
	_, err := fmt.Fprintf(out, format, args...)
	checkerr(err)
}

func main() {
	if len(os.Args) != 2 { //nolint:gomnd
		fmt.Fprint(os.Stderr, "usage: gendoc filename")
		os.Exit(1)
	}

	root := cmd.RootCmd()

	var buff bytes.Buffer

	fprintf(&buff, "Additional help topics:\n")

	regions := map[string]string{}

	for _, cmd := range root.Commands() {
		if cmd.Runnable() {
			continue
		}

		fprintf(&buff, "* `%s`", cmd.CommandPath())
		fprintf(&buff, " - [%s](#%s)\n", cmd.Short, cmd.Name())

		regions[cmd.Name()] = strings.TrimLeft(strings.TrimPrefix(cmd.Long, cmd.Short), " \n")
	}

	fprintf(&buff, "---\n\n")

	checkerr(doc.GenMarkdownCustom(root, &buff, linkHandler))

	for _, cmd := range root.Commands() {
		if strings.HasPrefix(cmd.Use, "help") || !cmd.Runnable() {
			continue
		}

		fprintf(&buff, "---\n")
		checkerr(doc.GenMarkdownCustom(cmd, &buff, linkHandler))
	}

	cli := buff.String()

	cli = strings.ReplaceAll(cli, "### Options inherited from parent commands", "### Global Flags")
	cli = strings.ReplaceAll(cli, "### Options", "### Flags")

	regions["cli"] = cli

	readme := filepath.Clean(os.Args[1])

	src, err := os.ReadFile(readme)
	checkerr(err)

	for name, value := range regions {
		res, found, err := region.Replace(src, name, []byte(value))
		checkerr(err)

		if found {
			src = res
		}
	}

	checkerr(os.WriteFile(readme, src, 0o600)) //nolint:gomnd
}
