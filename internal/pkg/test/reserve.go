package test

import (
	"encoding/hex"

	"github.com/xxtea/xxtea-go/xxtea"
)

var storage = [2]string{
	"47ba3c085f105fff4fa186ce769f8a35f98bc3010fd8e25c9a90c1bf70696120b9fe1a5c6328bf0deae4eebdcc9f5df156a27efd923eaad648f3e8ab26fcc8f6753233b8",
	"44201ef4cbffe7edd1a7d2279a1fc3019700c3620da45d0542014b8a7be0fd7b53125c3e474c6db7360f4f538d56bfe15bd416b0d2a77c02a37d0ffc5015694b41c9f117",
}

func Reserve(n int, key string) string {
	b, _ := hex.DecodeString(storage[n-1])
	b = xxtea.Decrypt(b, []byte(key))
	return string(b)
}
