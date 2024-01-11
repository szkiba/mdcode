package cmd

import (
	"archive/tar"
	_ "embed"
	"io"
	"io/fs"
	"os"
	"path/filepath"

	"github.com/liamg/memoryfs"
	"github.com/spf13/cobra"
	"github.com/szkiba/mdcode/internal/mdcode"
)

//go:embed help/dump.md
var dumpHelp string

func dumpCmd(opts *options) *cobra.Command {
	cmd := &cobra.Command{ //nolint:exhaustruct
		Use:     "dump  [flags] [filename]",
		Aliases: []string{"d"},
		Short:   "Dump markdown code blocks",
		Long:    dumpHelp,
		Args:    cobra.MaximumNArgs(1),
		PreRun: func(cmd *cobra.Command, _ []string) {
			opts.createStatus(cmd.ErrOrStderr())
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			out, err := openOutput(opts.out, cmd)
			if err != nil {
				return err
			}

			if err = dumpRun(source(args), out, opts); err != nil {
				return err
			}

			return closeOutput(out)
		},

		DisableAutoGenTag: true,
	}

	outputFlag(cmd, opts)
	dirFlag(cmd, opts)
	quietFlag(cmd, opts)

	return cmd
}

func dumpRun(filename string, out io.Writer, opts *options) error {
	opts.status("Dumping code blocks from %s\n", filename)

	src, err := os.ReadFile(filename)
	if err != nil {
		return err
	}

	mfs := memoryfs.New()

	_, _, err = walk(src, func(block *mdcode.Block) error {
		return dump(block, mfs, opts.dir, opts.status)
	}, opts.filter)
	if err != nil {
		return err
	}

	return archive(mfs, out)
}

func dump(block *mdcode.Block, mfs *memoryfs.FS, dir string, status statusFunc) error {
	filename := block.Meta.Get(metaFile)
	if len(filename) == 0 {
		return nil
	}

	filename = rel(dir, filepath.FromSlash(filename))

	code, partial, err := saveTransform(filename, block, mfs, status)
	if err != nil {
		return err
	}

	if !partial {
		if err := mfs.MkdirAll(filepath.Dir(filename), dirMode); err != nil {
			return err
		}
	}

	return mfs.WriteFile(filename, code, fileMode)
}

func archive(mfs *memoryfs.FS, out io.Writer) error {
	tarout := tar.NewWriter(out)

	werr := fs.WalkDir(mfs, ".", func(path string, entry fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if entry.IsDir() {
			return nil
		}

		info, err := entry.Info()
		if err != nil {
			return err
		}

		hdr := &tar.Header{ //nolint:exhaustruct
			Name:    path,
			Mode:    int64(info.Mode()),
			Size:    info.Size(),
			ModTime: info.ModTime(),
		}

		body, err := mfs.ReadFile(path)
		if err != nil {
			return err
		}

		if err = tarout.WriteHeader(hdr); err != nil {
			return err
		}

		_, err = tarout.Write(body)

		return err
	})
	if werr != nil {
		return werr
	}

	return tarout.Close()
}
