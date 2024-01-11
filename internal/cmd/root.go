// Package cmd contains mdcode CLI interface.
package cmd

import (
	_ "embed"
	"errors"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
)

//nolint:gochecknoglobals
var (
	version = "dev"
	appname = "mdcode"
)

func Execute(args []string, stdout, stderr io.Writer) {
	root := RootCmd()

	root.SetArgs(args)
	root.SetErr(stderr)
	root.SetOut(stdout)

	cobra.CheckErr(root.Execute())
}

//go:embed help/root.md
var rootHelp string

func RootCmd() *cobra.Command {
	opts := new(options)

	cmd := &cobra.Command{ //nolint:exhaustruct
		Use:     appname + " [flags] [filename]",
		Short:   "Markdown code block authoring tool",
		Long:    rootHelp,
		Version: version,
		Args:    checkargs,
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			err := opts.createFilter()
			if err != nil {
				return err
			}

			if flag := cmd.Flag("dir"); flag != nil && !flag.Changed {
				opts.dir = filepath.Dir(source(args))
			}

			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			out, err := openOutput(opts.out, cmd)
			if err != nil {
				return err
			}

			if err = listRun(source(args), out, opts); err != nil {
				return err
			}

			return closeOutput(out)
		},

		SilenceUsage:      true,
		SilenceErrors:     true,
		DisableAutoGenTag: true,
	}

	cmd.SetVersionTemplate(
		`{{with .Name}}{{printf "%s" .}}{{end}}{{printf " version %s\n" .Version}}`,
	)

	globalFlags(cmd, opts)

	outputFlag(cmd, opts)

	cmd.Flags().BoolVar(&opts.json, "json", false, "generate JSON output")

	cmd.AddCommand(updateCmd(opts))
	cmd.AddCommand(extractCmd(opts))
	cmd.AddCommand(dumpCmd(opts))
	cmd.AddCommand(runCmd(opts))

	cmd.AddCommand(metadataTopic(), filteringTopic(), regionsTopic(), invisibleTopic(), outlineTopic())

	return cmd
}

func globalFlags(cmd *cobra.Command, opts *options) {
	flags := cmd.PersistentFlags()

	flags.StringSliceVarP(&opts.file, "file", "f", []string{"?*"}, "file filter")
	flags.StringSliceVarP(&opts.lang, "lang", "l", []string{"?*"}, "language filter")
	flags.StringToStringVarP(&opts.meta, "meta", "m", nil, "metadata filter")
}

func outputFlag(cmd *cobra.Command, opts *options) {
	cmd.Flags().StringVarP(&opts.out, "output", "o", "", "output file (default: standard output)")

	cobra.CheckErr(cmd.MarkFlagFilename("output"))
}

func dirFlag(cmd *cobra.Command, opts *options) {
	cmd.Flags().StringVarP(&opts.dir, "dir", "d", ".", "base directory name")

	cobra.CheckErr(cmd.MarkFlagDirname("dir"))
}

func quietFlag(cmd *cobra.Command, opts *options) {
	cmd.Flags().BoolVarP(&opts.quiet, "quiet", "q", false, "suppress the status output")
}

func checkargs(cmd *cobra.Command, args []string) error {
	_, args = script(cmd, args)

	if len(args) > 1 {
		return errTooManyArg
	}

	if len(args) == 0 {
		if _, err := os.Stat(defaultArg); errors.Is(err, os.ErrNotExist) {
			return errMissingArg
		}
	}

	return nil
}

var (
	errMissingArg = errors.New("the filename argument is missing and " + defaultArg + " is not found")
	errTooManyArg = errors.New("too many arguments")
)

func openOutput(out string, cmd *cobra.Command) (io.Writer, error) {
	if len(out) == 0 {
		return cmd.OutOrStdout(), nil
	}

	return os.Create(out)
}

func closeOutput(out io.Writer) error {
	if closer, ok := out.(io.Closer); ok {
		return closer.Close()
	}

	return nil
}

func source(args []string) string {
	if len(args) == 0 {
		return defaultArg
	}

	return args[0]
}

func script(cmd *cobra.Command, args []string) (string, []string) {
	if cmd.ArgsLenAtDash() < 0 {
		return "", args
	}

	return strings.Join(args[cmd.ArgsLenAtDash():], " "), args[:cmd.ArgsLenAtDash()]
}

const (
	defaultArg = "README.md"

	dirMode  = 0o750
	fileMode = 0o600
)
