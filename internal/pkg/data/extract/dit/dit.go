// Package dit based on:
// - https://www.exploit-db.com/docs/english/18244-active-domain-offline-hash-dump-&-forensic-analysis.pdf
// - https://github.com/C-Sto/gosecretsdump/blob/master/pkg/ditreader/crypto.go
// - https://github.com/Dionach/NtdsAudit/blob/master/src/NtdsAudit/NTCrypto.cs
package dit

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/des"
	"crypto/md5"
	"crypto/rc4"
	"encoding/binary"
	"encoding/hex"
	"errors"
	"fmt"
	"log"
	"slices"

	"github.com/Velocidex/ordereddict"
	"www.velocidex.com/golang/go-ese/parser"
)

const (
	attAcc = "ATTm590045"
	attTyp = "ATTj590126"
	attSid = "ATTr589970"
	attPek = "ATTk590689"
	attLm  = "ATTk589879"
	attNt  = "ATTk589914"
)

var types = []int64{
	0x30000000, // SAM_NORMAL_USER_ACCOUNT
	0x30000001, // SAM_MACHINE_ACCOUNT
	0x30000002, // SAM_TRUST_ACCOUNT
}

var emptyNt = []byte{0x31, 0xD6, 0xCF, 0xE0, 0xD1, 0x6A, 0xE9, 0x31, 0xB7, 0x3C, 0x59, 0xD7, 0xE0, 0xC0, 0x89, 0xC0}
var emptyLm = []byte{0xAA, 0xD3, 0xB4, 0x35, 0xB5, 0x14, 0x04, 0xEE, 0xAA, 0xD3, 0xB4, 0x35, 0xB5, 0x14, 0x04, 0xEE}

var errStop = errors.New("stop")

type Pek []byte
type Hash []byte
type Record struct {
	User string
	Rid  uint32
	Nt   string
	Lm   string
}

func newPek(b, bk []byte) Pek {
	var key []byte

	if len(b) != 76 {
		log.Fatalln("invalid pek data")
	}

	buf := b[8:] // skip header

	switch b[0] {
	case 0x03: // 2016
		key = decryptAes(buf[16:], bk, buf[:16])
		key = key[4:20]

	case 0x02: // 2000
		key = deriveMd5(buf[:16], bk, 1000)
		key = decryptRc4(buf[16:], key)
		key = key[36:]

	default:
		// plain text?
	}

	if len(key) != 16 {
		log.Fatalln("invalid pek length")
	}

	return key
}

func newHash(b, d, k, k1, k2 []byte) Hash {
	if len(b) == 0 {
		return d
	}

	if len(b) != 40 {
		log.Fatalln("invalid hash data")
	}

	buf := b[8:] // skip header

	switch b[0] {
	case 0x13:
		// TODO: AES

	default:
		key := deriveMd5(buf[:16], k, 1)
		buf = decryptRc4(buf[16:], key)
		buf = decryptDes(buf, k1, k2)
	}

	return buf
}

func (v Pek) String() string {
	return hex.EncodeToString(v)
}

func (v Hash) String() string {
	return hex.EncodeToString(v)
}

func (r *Record) String() string {
	return fmt.Sprintf("%s:%d:%s:%s:::",
		r.User,
		r.Rid,
		r.Lm,
		r.Nt,
	)
}

func Extract(b, bootkey []byte) ([]Record, error) {
	var r []Record

	ctx, err := parser.NewESEContext(bytes.NewReader(b))

	if err != nil {
		return nil, err
	}

	ctl, err := parser.ReadCatalog(ctx)

	if err != nil {
		return nil, err
	}

	pek := newPek(getBytes(ctl, attPek), bootkey)

	_ = ctl.DumpTable("datatable", func(row *ordereddict.Dict) error {
		if v, ok := row.Get(attAcc); ok && v != nil {
			t, _ := row.GetInt64(attTyp)

			if slices.Contains(types, t) {
				rec, err := newRecord(row, v.(string), pek)

				if err != nil {
					log.Println(err)
					return nil
				}

				r = append(r, *rec)
			}

			return nil
		}
		return nil
	})

	return r, nil
}

func newRecord(row *ordereddict.Dict, usr string, pek []byte) (*Record, error) {
	rid := extractRid(getRowBytes(row, attSid))

	k1, k2 := deriveKey(rid)

	return &Record{
		User: usr,
		Rid:  rid,
		Nt:   newHash(getRowBytes(row, attNt), emptyNt, pek, k1, k2).String(),
		Lm:   newHash(getRowBytes(row, attLm), emptyLm, pek, k1, k2).String(),
	}, nil
}

func getBytes(ctl *parser.Catalog, att string) []byte {
	var b []byte

	_ = ctl.DumpTable("datatable", func(row *ordereddict.Dict) error {
		if v, ok := row.Get(att); ok && v != nil {
			b, _ = hex.DecodeString(v.(string))
			return errStop
		}
		return nil
	})

	return b
}

func getRowBytes(row *ordereddict.Dict, key string) []byte {
	if v := getRow(row, key); v != nil {
		b, _ := hex.DecodeString(v.(string))
		return b
	}

	return nil
}

func getRow(row *ordereddict.Dict, key string) any {
	if v, ok := row.Get(key); ok && v != nil {
		return v
	}

	return nil
}

func extractRid(sid []byte) uint32 {
	l, s := sid[1], sid[8:]

	return binary.BigEndian.Uint32(s[(l-1)*4 : (l-1)*4+4])
}

func decryptDes(b, k1, k2 []byte) []byte {
	var p []byte

	p1 := make([]byte, 8)
	p2 := make([]byte, 8)

	c1, err := des.NewCipher(k1)

	if err != nil {
		log.Fatalln(err)
	}

	c1.Decrypt(p1, b[:8])

	p = append(p, p1...)

	c2, err := des.NewCipher(k2)

	if err != nil {
		log.Fatalln(err)
	}

	c2.Decrypt(p2, b[8:])

	p = append(p, p2...)

	return p
}

func decryptAes(b, k, v []byte) []byte {
	p := make([]byte, len(b))

	c, err := aes.NewCipher(k)

	if err != nil {
		log.Fatalln(err)
	}

	d := cipher.NewCBCDecrypter(c, v)
	d.CryptBlocks(p, b)

	return p
}

func decryptRc4(b, k []byte) []byte {
	p := make([]byte, len(b))

	c, err := rc4.NewCipher(k)

	if err != nil {
		log.Fatalln(err)
	}

	c.XORKeyStream(p, b)

	return p
}

func deriveMd5(b, k []byte, n int) []byte {
	r := make([]byte, 16)

	h := md5.New()
	h.Write(k)

	for i := 0; i < n; i++ {
		h.Write(b)
	}

	s := h.Sum(nil)

	copy(r, s)

	return r
}

func deriveKey(rid uint32) (k1, k2 []byte) {
	k := make([]byte, 4)

	binary.LittleEndian.PutUint32(k, rid)

	b1 := []byte{
		k[0], k[1], k[2], k[3],
		k[0], k[1], k[2],
	}

	b2 := []byte{
		k[3], k[0], k[1], k[2],
		k[3], k[0], k[1],
	}

	return transformKey(b1), transformKey(b2)
}

func transformKey(b []byte) []byte {
	var key []byte

	key = append(key, b[0]>>0x01)
	key = append(key, ((b[0]&0x01)<<6)|b[1]>>2)
	key = append(key, ((b[1]&0x03)<<5)|b[2]>>3)
	key = append(key, ((b[2]&0x07)<<4)|b[3]>>4)
	key = append(key, ((b[3]&0x0f)<<3)|b[4]>>5)
	key = append(key, ((b[4]&0x1f)<<2)|b[5]>>6)
	key = append(key, ((b[5]&0x3f)<<1)|b[6]>>7)
	key = append(key, b[6]&0x7f)

	for i := 0; i < 8; i++ {
		key[i] = (key[i] << 1) & 0xfe
	}

	return key
}
