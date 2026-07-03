package help

import (
	"bytes"
	"testing"

	"go.foxforensics.eu/fox/v4/internal/cmd"
	"go.foxforensics.eu/fox/v4/internal/cmd/ad"
	"go.foxforensics.eu/fox/v4/internal/cmd/hash"
	"go.foxforensics.eu/fox/v4/internal/cmd/hunt"
	"go.foxforensics.eu/fox/v4/internal/cmd/info"
	"go.foxforensics.eu/fox/v4/internal/cmd/str"
	"go.foxforensics.eu/fox/v4/internal/test"
)

func TestHelp(t *testing.T) {
	for _, tt := range []struct {
		name   string
		output string
		args   []string
	}{
		{
			"Main",
			cmd.Usage,
			[]string{
				"help",
			},
		},
		{
			"Ad",
			ad.Usage,
			[]string{
				"help",
				"ad",
			},
		},
		{
			"Hash",
			hash.Usage,
			[]string{
				"help",
				"hash",
			},
		},
		{
			"Hunt",
			hunt.Usage,
			[]string{
				"help",
				"hunt",
			},
		},
		{
			"Info",
			info.Usage,
			[]string{
				"help",
				"info",
			},
		},
		{
			"Str",
			str.Usage,
			[]string{
				"help",
				"str",
			},
		},
		{
			"Error",
			"exit status 1",
			[]string{
				"help",
				"error",
			},
		},
	} {
		for range test.Cycles {
			t.Run(tt.name, func(t *testing.T) {
				b, err := test.FixtureMain(tt.args...)

				if err != nil {
					b = []byte(err.Error())
				}

				if !bytes.Contains(b, []byte(tt.output)) {
					t.Fatal("output mismatch")
				}
			})
		}
	}
}
