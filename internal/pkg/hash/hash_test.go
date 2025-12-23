package hash

import (
	"log"
	"os"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/cuhsat/fox/v4/internal/pkg/types"
)

const Large = "text/bible.txt"
const Small = "fox.txt"
const Input = "FOX123XOF"

func BenchmarkSum(b *testing.B) {
	buf := Fixture(Small)

	b.ResetTimer()

	for b.Loop() {
		_, _ = Sum(types.SHA256, buf)
	}
}

func TestSum(t *testing.T) {
	large := Fixture(Large)
	small := Fixture(Small)
	input := []byte(Input)

	for _, tt := range []struct {
		buf []byte
		imp string
		sum string
	}{
		{small, types.ADLER32, "6b0a131c"},
		{small, types.BLAKE2S256, "35622f6446178515ba503412f31eb768b092d878acbe6bf422b3ee47cf0558e7"},
		{small, types.BLAKE2B256, "adb735516f5c008b29e3313869311096fe671bd2bd5f199b1a49ce579e0f0bd2"},
		{small, types.BLAKE2B384, "879b983993f559354807d61b7cdf00c310c093bbb904eb43db9db2b49a83e8a07512d5f24299e458b5374291b23a9f86"},
		{small, types.BLAKE2B512, "6805b54f9ad456d6a217624fa0c992108a35cf52f35b9d7f617533b50804bf08b2af20653469b3b76acc799bd3c905919cc958084179adf2b1475493d5cd1810"},
		{small, types.BLAKE3256, "68aa491620394f724284e35b51551a21ab715f0c38f85cf8ba837233d34ae4a6"},
		{small, types.BLAKE3512, "68aa491620394f724284e35b51551a21ab715f0c38f85cf8ba837233d34ae4a646914037b62a958cc6f769865ea235ae326f8cc4e7eb7010102dfd0d72652d9a"},
		{small, types.CRC32C, "afb3f887"},
		{small, types.CRC32IEEE, "7ab53d60"},
		{small, types.CRC64ECMA, "df2fc66f2c50575f"},
		{small, types.CRC64ISO, "66747f552337d269"},
		{small, types.FNV1, "847595167a564758d45f1ac5f7b7fad0"},
		{small, types.FNV1A, "8e1fbe2b2d87d680249d1d1135695632"},
		{input, types.LM, "74ac61daa7e79d69482bc9e3e9caf5a9"},
		{input, types.NT, "0b2bee3bac7bddeae69d63d53c0b68f3"},
		{small, types.MD2, "9e49ada9a2ccafdafffff50137351626"},
		{small, types.MD4, "faedf7d245748f2939593258a5e96875"},
		{small, types.MD5, "7fe307fda20e805d110b35bcc1f31167"},
		{small, types.MD6, "599f033e751832ce908f22a3b0b0bf316a77f1553bc4c24146caf9fa6b235854"},
		{small, types.SHA1, "b11b92d927f2eb66f0aa17266f7348c0cdfd1105"},
		{small, types.SHA256, "b7e664f9009f84aa056fc78008fe24f33bd45795c407162a78b0fd4c6c2e2d08"},
		{small, types.SHA3, "96c5ca5658d7a04cb844539bcab4c2ebe503bc16c41f79ba207ab011"},
		{small, types.SHA3224, "96c5ca5658d7a04cb844539bcab4c2ebe503bc16c41f79ba207ab011"},
		{small, types.SHA3256, "40ec86016388c549a4a4954a068989b2b757f6488dce0f1cd4a558ee550129fe"},
		{small, types.SHA3384, "66f3c9e7d5c888ada9e8fc37994eb468239a31f2694e4aa8450c7eecf2803d9ef385f6c632b0a9b66452032e29ffefbb"},
		{small, types.SHA3512, "794a82f57c8448a5221c8cac462541092f2ef198df3d41edbf5f4ea6f19fdf26f98c37a82eec8be367547822aa5f90e23e2b5f9d26be9f9ee6fb0b654de918e1"},
		{large, types.SSDEEP, "49152:LkD0m3lNkRAA4Ml/Mo3hdWoPPwXj3NfhrZChJl7v6ih7T87/MvwFLSMyJTszqBPh:t"},
		{large, types.TLSH, "T12526a757e784133b1b620334620ea5d9f31ac43e7676ce30585ee03e2356c7996b9be8"},
		{small, types.XXH3, "50b2cde07882a633"},
		{small, types.XXH64, "d2ff231ddefb0bd0"},
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
