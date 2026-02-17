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
			"smb://user:****@127.0.0.1:445/Share/dir/file.ext:stream",
		},
		{
			"//user:@host:445/Share/dir/file.ext:stream",
			"smb://user@host:445/Share/dir/file.ext:stream",
		},
		{
			"//user@host:445/Share/dir/file.ext:stream",
			"smb://user@host:445/Share/dir/file.ext:stream",
		},
		{
			"//host:445/Share/dir/file.ext:stream",
			"smb://host:445/Share/dir/file.ext:stream",
		},
		{
			"//host/Share/dir/file.ext:stream",
			"smb://host:445/Share/dir/file.ext:stream",
		},
		{
			"//host/Share/dir/file.ext",
			"smb://host:445/Share/dir/file.ext",
		},
		{
			"//host/Share/file.ext",
			"smb://host:445/Share/file.ext",
		},
		{
			"//host/Share/dir/",
			"smb://host:445/Share/dir/",
		},
		{
			"//host/Share/",
			"smb://host:445/Share/",
		},
	} {
		t.Run(fmt.Sprintf("%d", i+1), func(t *testing.T) {
			if Parse(tt.unc).String() != tt.str {
				t.Fatal("wrong result")
			}
		})
	}
}
