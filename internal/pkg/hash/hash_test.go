package hash

import (
	"testing"

	"github.com/cuhsat/fox/v4/internal/pkg/types"
)

var text = []byte("FOX TEST")

func BenchmarkSum(b *testing.B) {

	b.ResetTimer()

	for b.Loop() {
		_, _ = Sum(types.SHA256, text)
	}
}
