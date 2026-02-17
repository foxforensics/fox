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
			"smb://user:pass@127.0.0.1:445/Share/",
			"//user@127.0.0.1:445/Share/",
		},
		{
			`\\user@host\Share\`,
			"//user@host:445/Share/",
		},
		{
			"//user:@host:445/Share/",
			"//user@host:445/Share/",
		},
		{
			"//user@host:445/Share/",
			"//user@host:445/Share/",
		},
		{
			"//host:445/Share/",
			"//host:445/Share/",
		},
		{
			"//host/Share/",
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
