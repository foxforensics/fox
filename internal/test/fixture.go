package test

import (
	"bytes"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"runtime"
)

const Sample = "fox.txt"

const root = "../../testdata"

func Assert(b []byte) bool {
	return bytes.Equal(b, Fixture(filepath.Join("texts", Sample)))
}

func FixtureFile(name string) string {
	_, c, _, ok := runtime.Caller(0)

	if !ok {
		slog.Error("fixture: runtime error")
		return ""
	}

	return filepath.Join(filepath.Dir(c), root, name)
}

func FixtureDir(names []string) []string {
	var v []string

	for _, name := range names {
		v = append(v, FixtureFile(name))
	}

	return v
}

func Fixture(name string) []byte {
	buf, err := os.ReadFile(FixtureFile(name))

	if err != nil {
		slog.Error(fmt.Sprintf("fixture: %s", err.Error()))
		return []byte(nil)
	}

	return buf
}
