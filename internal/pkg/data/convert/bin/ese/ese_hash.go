// Package ese based on: https://www.exploit-db.com/docs/english/18244-active-domain-offline-hash-dump-&-forensic-analysis.pdf
package ese

import (
	"bytes"
	"crypto/des"
	"crypto/md5"
	"crypto/rc4"
	"encoding/binary"
	"encoding/hex"
	"errors"
	"fmt"
	"log"

	"github.com/Velocidex/ordereddict"
	"www.velocidex.com/golang/go-ese/parser"
)

const (
	attSid = "ATTr589970"
	attPek = "ATTk590689"
	attLm  = "ATTk589879"
	attNt  = "ATTk589914"
)

type adPek struct {
	key  []byte
	data []byte
}

type adHash struct {
	key  []byte
	data []byte
}

func newPek(b []byte) *adPek {
	if len(b) != 76 {
		log.Fatalln("invalid pek data")
	}

	b = b[8:] // skip header

	ad := &adPek{
		key:  make([]byte, 16),
		data: make([]byte, 52),
	}

	copy(ad.key, b[:16])
	copy(ad.data, b[16:])

	return ad
}

func newHash(b []byte) *adHash {
	if len(b) != 40 {
		log.Fatalln("invalid hash data")
	}

	b = b[8:] // skip header

	ad := &adHash{
		key:  make([]byte, 16),
		data: make([]byte, 16),
	}

	copy(ad.key, b[:16])
	copy(ad.data, b[16:])

	return ad
}

func (ad *adPek) String() string {
	return fmt.Sprintf("key: %x data: %x", ad.key, ad.data)
}

func (ad *adHash) String() string {
	return fmt.Sprintf("key: %x data: %x", ad.key, ad.data)
}

func Extract(b, bk []byte) ([]byte, error) {
	buf := bytes.NewBuffer(nil)

	ctx, err := parser.NewESEContext(bytes.NewReader(b))

	if err != nil {
		return nil, err
	}

	ctl, err := parser.ReadCatalog(ctx)

	if err != nil {
		return nil, err
	}

	/* Enum Users
	_ = ctl.DumpTable("datatable", func(row *ordereddict.Dict) error {
		if v, ok := row.Get("ATTm590045"); ok && v != nil {
			println(fmt.Sprintf("%s", v))
			return nil
		}
		return nil
	})
	*/

	sid := getAttribute(ctl, attSid)
	rid := extractRID(sid)
	k1, k2 := deriveKey(rid)

	pek := newPek(getAttribute(ctl, attPek))
	lm := newHash(getAttribute(ctl, attLm))
	nt := newHash(getAttribute(ctl, attNt))

	pek.key = deriveMD5(pek.key, bk, 1000)
	pek.data = decryptRC4(pek.data, pek.key)

	lm.key = deriveMD5(lm.key, pek.data, 1)
	lm.data = decryptRC4(lm.data, lm.key)
	lm.data = decryptDES(lm.data, k1, k2)

	nt.key = deriveMD5(nt.key, pek.data, 1)
	nt.data = decryptRC4(nt.data, nt.key)
	nt.data = decryptDES(nt.data, k1, k2)

	/*
		println("PEK", pek.String())
		println("LM ", lm.String())
		println("NT ", nt.String())
		println("--")
	*/

	buf.WriteString(fmt.Sprintf("LM: %x\n", lm.data))
	buf.WriteString(fmt.Sprintf("NT: %x\n", nt.data))

	return buf.Bytes(), nil
}

func getAttribute(ctl *parser.Catalog, att string) []byte {
	var s string

	_ = ctl.DumpTable("datatable", func(row *ordereddict.Dict) error {
		if v, ok := row.Get(att); ok && v != nil {
			s = fmt.Sprintf("%s", v)
			return errors.New("stop")
		}
		return nil
	})

	b, err := hex.DecodeString(s)

	if err != nil {
		log.Fatalln(err)
	}

	return b
}

func deriveMD5(b, k []byte, n int) []byte {
	r := make([]byte, 16)

	h := md5.New()
	h.Sum(k)

	for i := 0; i < n; i++ {
		h.Sum(b)
	}

	s := h.Sum(nil)

	copy(r, s)

	return r
}

func decryptRC4(b, k []byte) []byte {
	r := make([]byte, len(b))

	c, err := rc4.NewCipher(k)

	if err != nil {
		log.Fatalln(err)
	}

	c.XORKeyStream(r, b)

	return r
}

func decryptDES(b, k1, k2 []byte) []byte {
	var r []byte

	b1 := make([]byte, 8)
	b2 := make([]byte, 8)

	d1, err := des.NewCipher(k1)

	if err != nil {
		log.Fatalln(err)
	}

	d1.Decrypt(b1, b[:8])

	r = append(r, b1...)

	d2, err := des.NewCipher(k2)

	if err != nil {
		log.Fatalln(err)
	}

	d2.Decrypt(b2, b[8:])

	r = append(r, b2...)

	return r
}

func extractRID(sid []byte) uint32 {
	l, s := sid[1], sid[8:]

	return binary.BigEndian.Uint32(s[(l-1)*4 : (l-1)*4+4])
}

func deriveKey(rid uint32) (k1, k2 []byte) {
	k := make([]byte, 4)

	binary.LittleEndian.PutUint32(k, rid)

	b1 := []byte{k[0], k[1], k[2], k[3], k[0], k[1], k[2]}
	b2 := []byte{k[3], k[0], k[1], k[2], k[3], k[0], k[1]}

	return transformKey(b1), transformKey(b2)
}

func transformKey(k []byte) []byte {
	var b []byte

	b = append(b, k[0]>>0x01)
	b = append(b, ((k[0]&0x01)<<6)|k[1]>>2)
	b = append(b, ((k[1]&0x03)<<5)|k[2]>>3)
	b = append(b, ((k[2]&0x07)<<4)|k[3]>>4)
	b = append(b, ((k[3]&0x0f)<<3)|k[4]>>5)
	b = append(b, ((k[4]&0x1f)<<2)|k[5]>>6)
	b = append(b, ((k[5]&0x3f)<<1)|k[6]>>7)
	b = append(b, k[6]&0x7f)

	for i := 0; i < 8; i++ {
		b[i] = (b[i] << 1) & 0xfe
	}

	return b
}
