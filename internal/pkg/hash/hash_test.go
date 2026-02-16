package hash

import (
	"testing"

	"github.com/cuhsat/fox/v4/internal/pkg/test"
	"github.com/cuhsat/fox/v4/internal/pkg/types"
)

const (
	Large = "text/bible.txt"
	Small = "fox.txt"
	PeBin = "convert/fox.exe.zst"
	Image = "misc/fox.jpg"
	Input = "FOX123XOF"
)

func BenchmarkSum(b *testing.B) {
	buf := test.Fixture(Small)

	for b.Loop() {
		_, _ = Sum(types.SHA256, buf)
	}
}

func TestSum(t *testing.T) {
	large := test.Fixture(Large)
	small := test.Fixture(Small)
	pebin := test.Fixture(PeBin)
	image := test.Fixture(Image)
	input := []byte(Input)

	for _, tt := range []struct {
		buf []byte
		imp string
		sum string
	}{
		{small, types.ADLER32, "6b0a131c"},
		{image, types.AHASH, "ffefc7c7c3c1c3ff"},
		{small, types.BLAKE2S256, "35622f6446178515ba503412f31eb768b092d878acbe6bf422b3ee47cf0558e7"},
		{small, types.BLAKE2B256, "adb735516f5c008b29e3313869311096fe671bd2bd5f199b1a49ce579e0f0bd2"},
		{small, types.BLAKE2B384, "879b983993f559354807d61b7cdf00c310c093bbb904eb43db9db2b49a83e8a07512d5f24299e458b5374291b23a9f86"},
		{small, types.BLAKE2B512, "6805b54f9ad456d6a217624fa0c992108a35cf52f35b9d7f617533b50804bf08b2af20653469b3b76acc799bd3c905919cc958084179adf2b1475493d5cd1810"},
		{small, types.BLAKE3256, "68aa491620394f724284e35b51551a21ab715f0c38f85cf8ba837233d34ae4a6"},
		{small, types.BLAKE3512, "68aa491620394f724284e35b51551a21ab715f0c38f85cf8ba837233d34ae4a646914037b62a958cc6f769865ea235ae326f8cc4e7eb7010102dfd0d72652d9a"},
		{small, types.CRC16CCITT, "3ff9"},
		{small, types.CRC32C, "afb3f887"},
		{small, types.CRC32IEEE, "7ab53d60"},
		{small, types.CRC64ECMA, "df2fc66f2c50575f"},
		{small, types.CRC64ISO, "66747f552337d269"},
		{image, types.DHASH, "180c0e1e1f171b19"},
		{small, types.FLETCHER4, "670c11a5040000007293332e26000000f95ebb7dd7000000cdb2a56cc1030000"},
		{small, types.FNV1, "847595167a564758d45f1ac5f7b7fad0"},
		{small, types.FNV1A, "8e1fbe2b2d87d680249d1d1135695632"},
		{small, types.GOST2012256, "59b14e45039898838d3a905382c73f9d7d73d8a376e770d2c78c744a45de840b"},
		{small, types.GOST2012512, "c3893dcdaec9998ac50b551a13f296a5599a1e5443dd508b29fafef6cbe369c063a565a90e0bf7fc0d9367dc59f314684f0dea19d983adf98144f3c107c7b681"},
		{small, types.HAS160, "e58d5cfe11171951799249e751e3bafecbf4d4a8"},
		{pebin, types.IMPFUZZY, "24:ibVjh9wOuuTkkboVaXOr6kwmDgUPMztxdEr6Ul:AwOuUjXOmokx0nl"},
		{pebin, types.IMPHASH, "d42595b695fc008ef2c56aabd8efd68e"},
		{pebin, types.IMPHASH0, "e76ebde493cf72a6f596a2c046a67080"},
		{input, types.LM, "74ac61daa7e79d69482bc9e3e9caf5a9"},
		{small, types.LSH256, "32be26fe3aab949b3682a6bd77028f9ffe9338fe605a7e7fd45c6a1dfa2c8585"},
		{small, types.LSH512, "06d60a9e8abba640bc989578af921ef1e13d333a6d017422fb45fff8611bd8321aaad2b25514877ad4a33ca3b5800f8e349dbf906d97999f9538dbfe0c075dec"},
		{small, types.MD2, "9e49ada9a2ccafdafffff50137351626"},
		{small, types.MD4, "faedf7d245748f2939593258a5e96875"},
		{small, types.MD5, "7fe307fda20e805d110b35bcc1f31167"},
		{small, types.MD6, "599f033e751832ce908f22a3b0b0bf316a77f1553bc4c24146caf9fa6b235854"},
		{small, types.MURMUR3, "785ae97135fcdbc8"},
		{input, types.NT, "0b2bee3bac7bddeae69d63d53c0b68f3"},
		{image, types.PHASH, "b193ce6c9c666c31"},
		{small, types.RIPEMD160, "12d7c8698119913bc60a9e1cfeb60853a9015b9d"},
		{small, types.SHAKE128, "6a000450724089944184129ff3fa56cd"},
		{small, types.SHAKE256, "3988412ad260af82eef7a889f3174147cf652ee07061b3606ea87fa37aabe01c"},
		{small, types.SHA1, "b11b92d927f2eb66f0aa17266f7348c0cdfd1105"},
		{small, types.SHA224, "03294a8a0ed498f32919abd3e7fda8332db3a99e2c882009056df8ea"},
		{small, types.SHA256, "b7e664f9009f84aa056fc78008fe24f33bd45795c407162a78b0fd4c6c2e2d08"},
		{small, types.SHA512, "1d34da51ac535e741e1e555bf80a1f4ca784225e1c443ceb6244e624b5548c9892f4b59a9e3776f8843f65b28cde99e9419eb09506feb3c00da9b11e844b58fe"},
		{small, types.SHA3, "96c5ca5658d7a04cb844539bcab4c2ebe503bc16c41f79ba207ab011"},
		{small, types.SHA3224, "96c5ca5658d7a04cb844539bcab4c2ebe503bc16c41f79ba207ab011"},
		{small, types.SHA3256, "40ec86016388c549a4a4954a068989b2b757f6488dce0f1cd4a558ee550129fe"},
		{small, types.SHA3384, "66f3c9e7d5c888ada9e8fc37994eb468239a31f2694e4aa8450c7eecf2803d9ef385f6c632b0a9b66452032e29ffefbb"},
		{small, types.SHA3512, "794a82f57c8448a5221c8cac462541092f2ef198df3d41edbf5f4ea6f19fdf26f98c37a82eec8be367547822aa5f90e23e2b5f9d26be9f9ee6fb0b654de918e1"},
		{small, types.SIPHASH, "013834802422bbba"},
		{small, types.SKEIN224, "11d11196404956666a66085a442bc697b7118a10e36de30db25ff7b4"},
		{small, types.SKEIN256, "6932c3ecc436c8af84a2ded54657fedcf06e534502e51cee810ba7b4374c5923"},
		{small, types.SKEIN384, "5526fcab9893cb9c732f385da9d7e24bb53d232df050b7452aad00402ac98793789b41c7d3044d1f4b8cf704873d2b9e"},
		{small, types.SKEIN512, "ff21b1fc0e73509e8cf3d1a5a283ef0acd9b4ecd5e4658dce10e1c1496210e9fbf39889e684629e131d6bca8dcbb789a93ba9f58b36cf06bda96e506708ceca6"},
		{small, types.SM3, "c86bb17e03669e27465dbb74c5e8c98035e899414697ab3a31f31345d0ccc2c6"},
		{large, types.SSDEEP, "49152:LkD0m3lNkRAA4Ml/Mo3hdWoPPwXj3NfhrZChJl7v6ih7T87/MvwFLSMyJTszqBPh:t"},
		{large, types.TLSH, "T12526a757e784133b1b620334620ea5d9f31ac43e7676ce30585ee03e2356c7996b9be8"},
		{small, types.WHIRLPOOL, "8aa3c190840f0991205e467f37ba57e8ff0350d10d6a07cd0a54efab6b1a529b97c1f8dd54f6489e55f992855aba3cdb3c8bb60ac072fb89c46dba6cc321dabf"},
		{small, types.XXH3, "50b2cde07882a633"},
		{small, types.XXH32, "ec6606ef"},
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
