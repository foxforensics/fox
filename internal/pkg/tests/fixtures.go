package tests

import (
	"bytes"
	"crypto/rand"
	"fmt"
	"log/slog"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
)

const Fox = "fox.txt"
const Cycles = 5
const Length = 10

func AssertFox(b []byte) bool {
	return bytes.Equal(b, Fixture(filepath.Join("texts", Fox)))
}

func ExecuteMain(args ...string) ([]byte, error) {
	v := append([]string{"run", FixtureFile("../../main.go")}, args...)

	cmd := exec.Command("go", v...)

	return cmd.CombinedOutput()
}

func FixtureFile(name string) string {
	d, err := os.Getwd()

	if err != nil {
		slog.Error(fmt.Sprintf("fixture: %s", err.Error()))
		return ""
	}

	v, err := filepath.Rel(d, filepath.Join(getRoot(), "../../../testdata/corpus", name))

	if err != nil {
		slog.Error(fmt.Sprintf("fixture: %s", err.Error()))
		return ""
	}

	return v
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

func Sample(name string) []byte {
	buf, err := os.ReadFile(filepath.Join(getRoot(), "../../../testdata/golden", name))

	if err != nil {
		slog.Error(fmt.Sprintf("fixture: %s", err.Error()))
		return []byte(nil)
	}

	return buf
}

func Random() [][]byte {
	v := make([][]byte, 0)

	for i := 0; i < Length; i++ {
		b := make([]byte, i)
		_, _ = rand.Read(b)
		v = append(v, b)
	}

	return v
}

func getRoot() string {
	_, c, _, ok := runtime.Caller(0)

	if !ok {
		slog.Error("fixture: runtime error")
		return ""
	}

	return filepath.Dir(c)
}
