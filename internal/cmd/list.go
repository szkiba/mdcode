package cmd

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"sort"
	"strings"

	"github.com/rodaine/table"
	"github.com/szkiba/mdcode/internal/mdcode"
)

func listRun(filename string, out io.Writer, opts *options) error {
	src, err := os.ReadFile(filename)
	if err != nil {
		return err
	}

	blocks, err := unfence(src, func(lang string, meta mdcode.Meta) bool {
		if isScript(lang, meta) {
			return true
		}

		return opts.filter(lang, meta)
	})
	if err != nil {
		return err
	}

	if opts.json {
		return listJSON(out, blocks)
	}

	listTabular(out, blocks)

	return nil
}

func listTabular(out io.Writer, blocks []*mdcode.Block) {
	keys := metaKeys(blocks)
	ikeys := make([]interface{}, 0, len(keys)+1)
	ikeys = append(ikeys, "lang")

	for _, k := range keys {
		ikeys = append(ikeys, k)
	}

	tbl := table.New(ikeys...).WithWriter(out)

	tbl.WithHeaderFormatter(func(format string, vals ...interface{}) string {
		return strings.ToUpper(fmt.Sprintf(format, vals...))
	})

	for _, block := range blocks {
		vals := make([]interface{}, 0, len(ikeys))
		vals = append(vals, block.Lang)

		for _, key := range keys {
			var value interface{}

			if s, has := block.Meta[key]; has {
				value = s
			} else {
				value = ""
			}

			vals = append(vals, value)
		}

		tbl.AddRow(vals...)
	}

	tbl.Print()
}

func listJSON(out io.Writer, blocks []*mdcode.Block) error {
	enc := json.NewEncoder(out)

	for _, b := range blocks {
		if len(b.Lang) != 0 {
			b.Meta["lang"] = b.Lang
		}

		if err := enc.Encode(b.Meta); err != nil {
			return err
		}
	}

	return nil
}

func metaKeys(blocks mdcode.Blocks) []string {
	keyset := make(map[string]struct{})

	for _, block := range blocks {
		for k := range block.Meta {
			keyset[k] = struct{}{}
		}
	}

	keys := make([]string, 0, len(keyset))
	idx := 0

	special := make(map[string]struct{})

	for _, s := range []string{metaName, metaFile, metaOutline, metaRegion} {
		special[s] = struct{}{}

		if _, has := keyset[s]; has {
			keys = append(keys, s)
			idx++
		}
	}

	for k := range keyset {
		if _, has := special[k]; !has {
			keys = append(keys, k)
		}
	}

	sort.Strings(keys[idx:])

	return keys
}

func unfence(src []byte, filter filterFunc) (mdcode.Blocks, error) {
	var blocks mdcode.Blocks

	_, _, err := walk(src, func(block *mdcode.Block) error {
		blocks = append(blocks, block)

		return nil
	}, filter)
	if err != nil {
		return nil, err
	}

	return blocks, err
}
