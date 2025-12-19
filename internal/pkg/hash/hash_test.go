package hash

import (
	"log"
	"os"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/cuhsat/fox/v4/internal/pkg/types"
)

const File = "text/bible.txt"
const Text = "FOX123XOF"

func BenchmarkSum(b *testing.B) {
	buf := Fixture(File)

	b.ResetTimer()

	for b.Loop() {
		_, _ = Sum(types.SHA256, buf)
	}
}

func TestSum(t *testing.T) {
	file := Fixture(File)
	text := []byte(Text)

	for _, tt := range []struct {
		buf []byte
		imp string
		sum string
	}{
		{file, types.ADLER32, "77e8f18a"},
		{file, types.CRC32C, "6e164f51"},
		{file, types.CRC32IEEE, "6de95d61"},
		{file, types.CRC64ECMA, "1c10bfee4250c76f"},
		{file, types.CRC64ISO, "37c0210d10e905f2"},
		{file, types.PE, "aac24600"},
		{text, types.LM, "74ac61daa7e79d69482bc9e3e9caf5a9"},
		{text, types.NT, "0b2bee3bac7bddeae69d63d53c0b68f3"},
		{file, types.MD2, "aceb3e20d985564d17838fc437744843"},
		{file, types.MD4, "fbb9a5a610458386e0ff2bdb4dea1076"},
		{file, types.MD5, "f7ebcc3119549346b871212958dbc203"},
		{file, types.MD6, "7abba14b23ea4438d2009118fe9d1befed73ba7420b5b7952fa8cd3c1a6ce62a"},
		{file, types.SHA1, "7763c40e323d9d57fc151cf9732dc4d5a07eaebf"},
		{file, types.SHA256, "61a54c7611855e09266732d923e64819273baf71b65bbb7c50249083e5b655fd"},
		{file, types.SHA3, "3429b00349d4dae6b707647747c8f3f9c819fb2ed8087fe435a6d126"},
		{file, types.SHA3224, "3429b00349d4dae6b707647747c8f3f9c819fb2ed8087fe435a6d126"},
		{file, types.SHA3256, "96b1fbd188e128c79eb5e4c3436b47785f2894187c2545d9db1dea6580ab5679"},
		{file, types.SHA3384, "273e6a839c571c724ad2d071889c23d2184387e1ed7e891d80e74a5a80b12f15393d26ec1abbe34dc85025402833de92"},
		{file, types.SHA3512, "5ce1d00eaf6da7409d009a3c242597559dc3cd2a0b41d5be9d2c86ae9709b66edeff0465fbdfca0432496ad0f3d839a4d3bf1d039a4161e3910d58d39d52e930"},
		{file, types.SSDEEP, "49152:LkD0m3lNkRAA4Ml/Mo3hdWoPPwXj3NfhrZChJl7v6ih7T87/MvwFLSMyJTszqBPh:t"},
		{file, types.TLSH, "T12526a757e784133b1b620334620ea5d9f31ac43e7676ce30585ee03e2356c7996b9be8"},
		{file, types.XXH3, "996653bb371ee4a1"},
		{file, types.XXH64, "6047f571a76ec9bb"},
	} {
		t.Run(tt.imp, func(t *testing.T) {
			sum, err := Sum(tt.imp, tt.buf)

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
