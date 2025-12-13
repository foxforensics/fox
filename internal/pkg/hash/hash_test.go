package hash

import (
	"log"
	"os"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/cuhsat/fox/v4/internal/pkg/types"
)

const file = "bible.txt"

func BenchmarkSum(b *testing.B) {
	buf := Fixture(file)

	b.ResetTimer()

	for b.Loop() {
		_, _ = Sum(types.SHA256, buf)
	}
}

func TestSum(t *testing.T) {
	buf := Fixture(file)

	for _, tt := range []struct {
		imp, sum string
	}{
		{types.CRC32C, "6e164f51"},
		{types.CRC32IEEE, "6de95d61"},
		{types.CRC64ECMA, "1c10bfee4250c76f"},
		{types.CRC64ISO, "37c0210d10e905f2"},
		{types.MD2, "aceb3e20d985564d17838fc437744843"},
		{types.MD4, "fbb9a5a610458386e0ff2bdb4dea1076"},
		{types.MD5, "f7ebcc3119549346b871212958dbc203"},
		{types.SHA1, "7763c40e323d9d57fc151cf9732dc4d5a07eaebf"},
		{types.SHA256, "61a54c7611855e09266732d923e64819273baf71b65bbb7c50249083e5b655fd"},
		{types.SHA3, "3429b00349d4dae6b707647747c8f3f9c819fb2ed8087fe435a6d126"},
		{types.SHA3224, "3429b00349d4dae6b707647747c8f3f9c819fb2ed8087fe435a6d126"},
		{types.SHA3256, "96b1fbd188e128c79eb5e4c3436b47785f2894187c2545d9db1dea6580ab5679"},
		{types.SHA3384, "273e6a839c571c724ad2d071889c23d2184387e1ed7e891d80e74a5a80b12f15393d26ec1abbe34dc85025402833de92"},
		{types.SHA3512, "5ce1d00eaf6da7409d009a3c242597559dc3cd2a0b41d5be9d2c86ae9709b66edeff0465fbdfca0432496ad0f3d839a4d3bf1d039a4161e3910d58d39d52e930"},
		{types.SSDEEP, "49152:LkD0m3lNkRAA4Ml/Mo3hdWoPPwXj3NfhrZChJl7v6ih7T87/MvwFLSMyJTszqBPh:t"},
		{types.TLSH, "T12526a757e784133b1b620334620ea5d9f31ac43e7676ce30585ee03e2356c7996b9be8"},
		{types.XXH3, "996653bb371ee4a1"},
		{types.XXH64, "6047f571a76ec9bb"},
	} {
		t.Run(tt.imp, func(t *testing.T) {
			sum, err := Sum(tt.imp, buf)

			if err != nil {
				t.Error(err)
			}

			if sum != tt.sum {
				t.Fatal("hash sum invalid")
			}
		})
	}
}

func Fixture(name string) []byte {
	_, c, _, ok := runtime.Caller(0)

	if !ok {
		log.Fatalln("runtime error")
	}

	buf, err := os.ReadFile(filepath.Join(filepath.Dir(c), "../../../testdata", name))

	if err != nil {
		log.Fatalln(err)
	}

	return buf
}
