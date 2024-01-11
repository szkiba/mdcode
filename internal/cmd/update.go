package cmd

import (
	_ "embed"
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/szkiba/mdcode/internal/mdcode"
	"github.com/szkiba/mdcode/internal/region"
)

//go:embed help/update.md
var updateHelp string

func updateCmd(opts *options) *cobra.Command {
	cmd := &cobra.Command{ //nolint:exhaustruct
		Use:     "update [flags] [filename]",
		Aliases: []string{"u"},
		Short:   "Update markdown code blocks from the file system",
		Long:    updateHelp,
		Args:    checkargs,
		PreRun: func(cmd *cobra.Command, _ []string) {
			opts.createStatus(cmd.ErrOrStderr())
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			return updateRun(source(args), opts)
		},

		DisableAutoGenTag: true,
	}

	dirFlag(cmd, opts)
	quietFlag(cmd, opts)

	return cmd
}

func updateRun(filename string, opts *options) error {
	opts.status("Updating code blocks in %s\n", filename)

	src, err := os.ReadFile(filename)
	if err != nil {
		return err
	}

	modified, res, e := walk(src, func(block *mdcode.Block) error {
		return load(block, opts.dir, opts.status)
	}, opts.filter)
	if e != nil {
		return e
	}

	if modified {
		return os.WriteFile(filename, res, fileMode)
	}

	return nil
}

func load(block *mdcode.Block, dir string, status statusFunc) error {
	filename := block.Meta.Get(metaFile)
	if len(filename) == 0 {
		return nil
	}

	filename = rel(dir, filepath.FromSlash(filename))

	code, err := os.ReadFile(filename)
	if err != nil {
		return err
	}

	code, err = loadTransform(filename, code, block, status)
	if err != nil {
		return err
	}

	block.Code = code

	return nil
}

func loadTransform(filename string, code []byte, block *mdcode.Block, status statusFunc) ([]byte, error) {
	regionname := block.Meta.Get(metaRegion)
	if len(regionname) != 0 {
		status("%s#%s\n", filename, regionname)

		data, ok, err := region.Read(code, regionname)
		if err != nil {
			return nil, err
		}

		if !ok {
			return nil, fmt.Errorf("%w: %s %s", errMissingRegion, filename, regionname)
		}

		return data, nil
	}

	status("%s\n", filename)

	outline := block.Meta.Get(metaOutline)
	if outline == "true" {
		data, ok, err := region.Outline(code)
		if err != nil {
			return nil, err
		}

		if !ok {
			return nil, fmt.Errorf("%w: %s", errNoRegion, filename)
		}

		return data, nil
	}

	return code, nil
}

var errNoRegion = errors.New("no #region")
