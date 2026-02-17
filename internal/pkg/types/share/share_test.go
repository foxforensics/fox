package share

import (
	"fmt"
	"testing"
)

func TestParse(t *testing.T) {
	for i, tt := range []struct {
		unc string
		str string
	}{
		{
			"smb://user:pass@127.0.0.1:445/Share/dir/file.ext:stream",
			"//user:****@127.0.0.1:445/Share/dir/file.ext:stream",
		},
		{
			`\\user@host\Share\dir\file.ext:stream`,
			"//user@host:445/Share/dir/file.ext:stream",
		},
		{
			"//user:@host:445/Share/dir/file.ext:stream",
			"//user@host:445/Share/dir/file.ext:stream",
		},
		{
			"//user@host:445/Share/dir/file.ext:stream",
			"//user@host:445/Share/dir/file.ext:stream",
		},
		{
			"//host:445/Share/dir/file.ext:stream",
			"//host:445/Share/dir/file.ext:stream",
		},
		{
			"//host/Share/dir/file.ext:stream",
			"//host:445/Share/dir/file.ext:stream",
		},
		{
			"//host/Share/dir/file.ext",
			"//host:445/Share/dir/file.ext",
		},
		{
			"//host/Share/file.ext",
			"//host:445/Share/file.ext",
		},
		{
			"//host/Share/dir/",
			"//host:445/Share/dir/",
		},
		{
			"//host/Share/",
			"//host:445/Share/",
		},
		{
			`\\host\Share\`,
			"//host:445/Share/",
		},
	} {
		t.Run(fmt.Sprintf("Path%d", i+1), func(t *testing.T) {
			if parse(tt.unc).String() != tt.str {
				t.Fatal("wrong result")
			}
		})
	}
}
