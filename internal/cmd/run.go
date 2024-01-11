package cmd

import (
	"context"
	_ "embed"
	"errors"
	"fmt"
	"os"
	"regexp"
	"strings"

	"github.com/spf13/cobra"
	"github.com/szkiba/mdcode/internal/mdcode"
	"mvdan.cc/sh/v3/interp"
	"mvdan.cc/sh/v3/syntax"
)

//go:embed help/run.md
var runHelp string

func runCmd(opts *options) *cobra.Command {
	cmd := &cobra.Command{ //nolint:exhaustruct
		Use:     "run [flags] [filename] [-- commands]",
		Aliases: []string{"r"},
		Short:   "Run shell commands on markdown code blocks",
		Long:    runHelp,
		Args:    checkargs,
		PreRun: func(cmd *cobra.Command, _ []string) {
			opts.createStatus(cmd.ErrOrStderr())
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			script, args := script(cmd, args)

			if !cmd.Flag("dir").Changed {
				dir, err := os.MkdirTemp(".", "mdcode-tmp-")
				if err != nil {
					return err
				}

				opts.dir = dir

				if !opts.keep {
					defer os.RemoveAll(dir)
				}
			}

			return runRun(source(args), opts, script)
		},
		DisableAutoGenTag: true,
	}

	dirFlag(cmd, opts)
	quietFlag(cmd, opts)

	cmd.Flags().StringVarP(&opts.name, "name", "n", "", "code block name contains commands")
	cmd.Flags().BoolVarP(&opts.keep, "keep", "k", false, "don't remove temporary directory")

	return cmd
}

var reShell = regexp.MustCompile("(ba|z)?sh")

func isScript(lang string, meta mdcode.Meta) bool {
	return reShell.MatchString(lang) && len(meta.Get(metaName)) != 0
}

func findScript(filename string, opts *options) (string, error) {
	src, err := os.ReadFile(filename)
	if err != nil {
		return "", err
	}

	var script string

	_, _, err = mdcode.Walk(src, func(block *mdcode.Block) error {
		if len(script) != 0 {
			return nil
		}

		if !isScript(block.Lang, block.Meta) {
			return nil
		}

		if len(opts.name) == 0 {
			script = string(block.Code)

			return nil
		}

		if block.Meta.Get(metaName) == opts.name {
			script = string(block.Code)
		}

		return nil
	})
	if err != nil {
		return "", err
	}

	if len(script) == 0 {
		if len(opts.name) != 0 {
			return "", fmt.Errorf("%w: %s", errMissingScript, opts.name)
		}

		return "", fmt.Errorf("%w: '-- commands' argument required", errMissingScript)
	}

	return script, nil
}

func runRun(filename string, opts *options, script string) error {
	if len(script) == 0 {
		value, err := findScript(filename, opts)
		if err != nil {
			return err
		}

		script = value
	}

	if err := extractRun(filename, opts); err != nil {
		return err
	}

	opts.status("Executing in %s\n%s\n", opts.dir, script)

	file, err := syntax.NewParser().Parse(strings.NewReader(script), "")
	if err != nil {
		return err
	}

	runner, err := interp.New(interp.Dir(opts.dir), interp.StdIO(os.Stdin, os.Stdout, os.Stderr))
	if err != nil {
		return err
	}

	return runner.Run(context.TODO(), file)
}

var errMissingScript = errors.New("missing script")
