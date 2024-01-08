package main

import (
	"io"
	"os"
	"testing"
)

func Test_main(t *testing.T) {
	orig := os.Stdout

	reader, writer, err := os.Pipe()
	if err != nil {
		t.Error(err)
	}

	os.Stdout = writer

	main()

	if err = writer.Close(); err != nil {
		t.Error(err)
	}

	out, err := io.ReadAll(reader)
	if err != nil {
		t.Error(err)
	}

	os.Stdout = orig

	const expected = "Hello, Testable World!\n"

	if string(out) != expected {
		t.Errorf("\nexpected: %s\nactual:   %s\n", expected, string(out))
	}
}
