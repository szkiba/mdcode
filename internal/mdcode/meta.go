package mdcode

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strings"

	"github.com/google/shlex"
)

type Meta map[string]interface{}

func (m Meta) Get(name string) string {
	if m == nil {
		return ""
	}

	value, has := m[name]
	if !has {
		return ""
	}

	if s, ok := value.(string); ok {
		return s
	}

	return fmt.Sprint(value)
}

var (
	reJSON     = regexp.MustCompile(`^\s*{\s*["}]`)
	reBrackets = regexp.MustCompile(`^\s*{(.*)}$`)
)

func parseMeta(input []byte) (Meta, error) {
	if len(input) == 0 {
		return Meta{}, nil
	}

	if reJSON.Match(input) {
		var meta Meta

		err := json.Unmarshal(input, &meta)
		if err != nil {
			return nil, err
		}

		return meta, nil
	}

	if subs := reBrackets.FindSubmatch(input); subs != nil {
		input = subs[1]
	}

	words, err := shlex.Split(string(input))
	if err != nil {
		return nil, err
	}

	dict := make(Meta)

	for _, word := range words {
		idx := strings.IndexRune(word, '=')
		if idx >= 0 && idx < len(word) {
			dict[word[:idx]] = word[idx+1:]
		}
	}

	return dict, nil
}
