package cmd

import (
	_ "embed"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/szkiba/mdcode/internal/mdcode"
	"github.com/szkiba/mdcode/internal/region"
)

//go:embed help/extract.md
var extractHelp string

func extractCmd(opts *options) *cobra.Command {
	cmd := &cobra.Command{ //nolint:exhaustruct
		Use:     "extract [filename]",
		Aliases: []string{"x"},
		Short:   "Extract markdown code blocks to the file system",
		Long:    extractHelp,
		Args:    checkargs,
		PreRun: func(cmd *cobra.Command, _ []string) {
			opts.createStatus(cmd.ErrOrStderr())
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			return extractRun(source(args), opts)
		},

		DisableAutoGenTag: true,
	}

	dirFlag(cmd, opts)
	quietFlag(cmd, opts)

	return cmd
}

func extractRun(filename string, opts *options) error {
	opts.status("Extracting code blocks from %s\n", filename)

	src, err := os.ReadFile(filename)
	if err != nil {
		return err
	}

	_, _, err = walk(src, func(block *mdcode.Block) error {
		return save(block, opts.dir, opts.status)
	}, opts.filter)

	return err
}

func save(block *mdcode.Block, dir string, status statusFunc) error {
	filename := block.Meta.Get(metaFile)
	if len(filename) == 0 {
		return nil
	}

	filename = rel(dir, filepath.FromSlash(filename))

	code, partial, err := saveTransform(filename, block, os.DirFS("."), status)
	if err != nil {
		return err
	}

	if !partial {
		if err := os.MkdirAll(filepath.Dir(filename), dirMode); err != nil {
			return err
		}
	}

	return os.WriteFile(filename, code, fileMode)
}

func saveTransform(filename string, block *mdcode.Block, fsys fs.FS, status statusFunc) ([]byte, bool, error) {
	regionname := block.Meta.Get(metaRegion)
	if len(regionname) == 0 {
		status("%s\n", filename)

		return block.Code, false, nil
	}

	status("%s#%s\n", filename, regionname)

	orig, err := fs.ReadFile(fsys, filepath.ToSlash(filename))
	if err != nil {
		return nil, false, err
	}

	data, mod, err := region.Replace(orig, regionname, block.Code)
	if err != nil {
		return nil, false, err
	}

	if !mod {
		return nil, false, fmt.Errorf("%w: %s %s", errMissingRegion, filename, regionname)
	}

	return data, true, nil
}

func rel(basedir string, filename string) string {
	if len(basedir) == 0 {
		return filepath.Join(".", filename)
	}

	return filepath.Join(basedir, filename)
}

var errMissingRegion = errors.New("missing region")
