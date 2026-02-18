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
			"smb://user:pass@127.0.0.1:445/share/",
			"//127.0.0.1:445/share/",
		},
		{
			`\\user:pass@127.0.0.1:445\share\`,
			"//127.0.0.1:445/share/",
		},
		{
			"//user:@host/share/",
			"//host/share/",
		},
		{
			"//user@host/share/",
			"//host/share/",
		},
		{
			"//host/share/",
			"//host/share/",
		},
		{
			"//host/share",
			"//host/share/",
		},
	} {
		t.Run(fmt.Sprintf("Path%d", i+1), func(t *testing.T) {
			if parse(tt.unc).String() != tt.str {
				t.Fatal("wrong result")
			}
		})
	}
}
