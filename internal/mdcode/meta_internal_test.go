package mdcode

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_parseMeta(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		arg     string
		want    Meta
		wantErr bool
	}{
		{name: "empty", wantErr: false, want: Meta{}, arg: ``},
		{name: "invalid JSON", wantErr: true, want: nil, arg: `{"foo`},
		{name: "invalid shlex", wantErr: true, want: nil, arg: `foo="`},
		{name: "shlex", wantErr: false, want: Meta{"foo": "bar", "answer": "42"}, arg: `foo="bar" answer=42`},
		{name: "brackets", wantErr: false, want: Meta{"foo": "bar", "answer": "42"}, arg: `{foo="bar" answer=42}`},
		{name: "JSON", wantErr: false, want: Meta{"foo": "bar", "answer": 42.0}, arg: `{"foo":"bar","answer":42}`},
		{name: "shlex skip no assign", wantErr: false, want: Meta{"foo": "bar"}, arg: `foo="bar" answer`},
		{name: "shlex empty assign", wantErr: false, want: Meta{"foo": "bar", "answer": ""}, arg: `foo="bar" answer=`},
	}
	for _, test := range tests {
		test := test

		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			got, err := parseMeta([]byte(test.arg))

			if test.wantErr {
				require.Error(t, err)

				return
			}

			require.NoError(t, err)
			require.Equal(t, test.want, got)
		})
	}
}

func TestMeta_Get(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		meta Meta
		arg  string
		want string
	}{
		{name: "nill meta", arg: "foo", want: "", meta: nil},
		{name: "regular", arg: "foo", want: "bar", meta: Meta{"foo": "bar"}},
		{name: "missing", arg: "bar", want: "", meta: Meta{"foo": "bar"}},
		{name: "non string", arg: "answer", want: "42", meta: Meta{"answer": 42}},
	}
	for _, test := range tests {
		test := test

		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			require.Equal(t, test.want, test.meta.Get(test.arg))
		})
	}
}
