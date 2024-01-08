package cmd

import (
	"fmt"
	"strings"

	"github.com/gobwas/glob"
	"github.com/szkiba/mdcode/internal/mdcode"
)

type filterFunc func(string, mdcode.Meta) bool

func filter(langs []string, metas map[string]string) (filterFunc, error) {
	var (
		langGlob glob.Glob
		metaGlob map[string]glob.Glob
	)

	comp, err := src2glob("", langs...)
	if err != nil {
		return nil, err
	}

	langGlob = comp

	metaGlob = make(map[string]glob.Glob)

	for key, value := range metas {
		if len(value) != 0 {
			comp, err = src2glob(key, value)
			if err != nil {
				return nil, err
			}

			metaGlob[key] = comp
		}
	}

	return func(lang string, meta mdcode.Meta) bool {
		if langGlob != nil && !langGlob.Match(lang) {
			return false
		}

		for k, g := range metaGlob {
			v, has := meta[k]
			if !has || !g.Match(fmt.Sprint(v)) {
				return false
			}
		}

		return true
	}, nil
}

func src2glob(key string, src ...string) (glob.Glob, error) { //nolint:ireturn
	if len(src) == 0 {
		return nil, nil
	}

	var separators []rune

	if key == metaFile {
		separators = append(separators, '/', '\\')
	}

	g, err := glob.Compile(fmt.Sprintf("{%s}", strings.Join(src, ",")), separators...)
	if err != nil {
		return nil, err
	}

	return g, nil
}
