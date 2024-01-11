package cmd

import (
	"fmt"
	"io"
	"strings"
)

const (
	metaFile    = "file"
	metaRegion  = "region"
	metaOutline = "outline"
	metaName    = "name"
)

type statusFunc func(format string, args ...any)

type options struct {
	lang []string
	file []string
	name string
	meta map[string]string

	dir string
	out string

	json bool

	quiet bool
	keep  bool

	filter filterFunc
	status statusFunc
}

func (o *options) createFilter() error {
	addMeta := func(key string, values []string) {
		v, ok := o.meta[key]
		if ok {
			v += ","
		}

		o.meta[key] = v + strings.Join(values, ",")
	}

	if o.meta == nil {
		o.meta = make(map[string]string)
	}

	addMeta(metaFile, o.file)

	var err error

	if o.filter, err = filter(o.lang, o.meta); err != nil {
		return err
	}

	return nil
}

func (o *options) createStatus(stderr io.Writer) {
	if o.quiet {
		o.status = func(format string, args ...any) {}
	} else {
		o.status = func(format string, args ...any) {
			fmt.Fprintf(stderr, format, args...)
		}
	}
}
