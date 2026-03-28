// Package dit based on:
// https://www.exploit-db.com/docs/english/18244-active-domain-offline-hash-dump-&-forensic-analysis.pdf
// https://github.com/fortra/impacket/blob/master/impacket/examples/secretsdump.py
// https://github.com/C-Sto/gosecretsdump/blob/master/pkg/ditreader/crypto.go
// https://github.com/Dionach/NtdsAudit/blob/master/src/NtdsAudit/NTCrypto.cs
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
	"encoding/json"
	"fmt"
	"log"
	"slices"

	"github.com/Velocidex/ordereddict"
	"go.foxforensics.dev/go-ese/parser"

	"github.com/cuhsat/fox/v4/internal/pkg/text"
)

// row attributes
const (
	accType = "ATTj590126"
	userRow = "ATTm590045"
	userSid = "ATTr589970"
	userUac = "ATTj589832"
	ntHash  = "ATTk589914"
	lmHash  = "ATTk589879"
	pekBin  = "ATTk590689"
)

// user types
var userTypes = []int64{
	0x30000000, // SAM_NORMAL_USER_ACCOUNT
	0x30000001, // SAM_MACHINE_ACCOUNT
	0x30000002, // SAM_TRUST_ACCOUNT
}

// default empty LM hash
var defaultLm = []byte{0xAA, 0xD3, 0xB4, 0x35, 0xB5, 0x14, 0x04, 0xEE, 0xAA, 0xD3, 0xB4, 0x35, 0xB5, 0x14, 0x04, 0xEE}

// default empty NT hash
var defaultNt = []byte{0x31, 0xD6, 0xCF, 0xE0, 0xD1, 0x6A, 0xE9, 0x31, 0xB7, 0x3C, 0x59, 0xD7, 0xE0, 0xC0, 0x89, 0xC0}

type Pek []byte

type Hash []byte

type Flags struct {
	Script                       bool `json:"script,omitempty"`
	AccountDisable               bool `json:"account_disable,omitempty"`
	HomeDirRequired              bool `json:"home_dir_required,omitempty"`
	Lockout                      bool `json:"lockout,omitempty"`
	PasswordNotRequired          bool `json:"password_not_required,omitempty"`
	EncryptedTextPasswordAllowed bool `json:"encrypted_text_password_allowed,omitempty"`
	TemporaryDupAccount          bool `json:"temporary_dup_account,omitempty"`
	NormalAccount                bool `json:"normal_account,omitempty"`
	InterDomainTrustAccount      bool `json:"inter_domain_trust_account,omitempty"`
	WorkstationTrustAccount      bool `json:"workstation_trust_account,omitempty"`
	ServerTrustAccount           bool `json:"server_trust_account,omitempty"`
	DontExpirePassword           bool `json:"dont_expire_password,omitempty"`
	MNSLogonAccount              bool `json:"mns_logon_account,omitempty"`
	SmartCardRequired            bool `json:"smart_card_required,omitempty"`
	TrustedForDelegation         bool `json:"trusted_for_delegation,omitempty"`
	NotDelegated                 bool `json:"not_delegated,omitempty"`
	UseDESOnly                   bool `json:"use_des_only,omitempty"`
	DontPreAuth                  bool `json:"dont_pre_auth,omitempty"`
	PasswordExpired              bool `json:"password_expired,omitempty"`
	TrustedToAuthForDelegation   bool `json:"trusted_to_auth_for_delegation,omitempty"`
	PartialSecrets               bool `json:"partial_secrets,omitempty"`
}

type Record struct {
	Username string `json:"username,omitempty"`
	Rid      uint32 `json:"rid,omitempty"`
	Sid      string `json:"sid,omitempty"`
	Flags    *Flags `json:"flags,omitempty"`
	LmHash   string `json:"lm_hash,omitempty"`
	NtHash   string `json:"nt_hash,omitempty"`
}

func (rec *Record) String() string {
	lm := rec.LmHash
	nt := rec.NtHash

	if lm == fmt.Sprintf("%x", defaultLm) {
		lm = text.AsGray(lm)
	}

	if nt == fmt.Sprintf("%x", defaultNt) {
		nt = text.AsGray(nt)
	}

	return fmt.Sprintf("%s:%d:%s:%s:::",
		rec.Username,
		rec.Rid,
		lm,
		nt,
	)
}

func (rec *Record) ToJSON() string {
	b, _ := json.MarshalIndent(rec, "", "  ")
	return string(b)
}

func (rec *Record) ToJSONL() string {
	b, _ := json.Marshal(rec)
	return string(b)
}

func Extract(b, bootkey []byte) ([]Record, []Pek, error) {
	var r []Record

	ctx, err := parser.NewESEContext(bytes.NewReader(b), int64(len(b)))

	if err != nil {
		return nil, nil, err
	}

	ctl, err := parser.ReadCatalog(ctx)

	if err != nil {
		return nil, nil, err
	}

	keys := getKeys(ctl, pekBin, bootkey)

	_ = ctl.DumpTable("datatable", func(row *ordereddict.Dict) error {
		if v, ok := row.Get(userRow); ok && v != nil {
			t, _ := row.GetInt64(accType)

			if slices.Contains(userTypes, t) {
				rec, err := newRecord(row, v.(string), keys)

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

	return r, keys, nil
}

func newPek(b, bootkey []byte) Pek {
	var key []byte

	buf := b[8:] // skip header

	switch b[0] {
	case 0x03: // 2016
		key = decryptAes(buf[16:], bootkey, buf[:16])
		key = key[36:52]

	case 0x02: // 2000
		key = deriveMd5(buf[:16], bootkey, 1000)
		key = decryptRc4(buf[16:], key)
		key = key[36:]

	default:
		// plain text?
	}

	if len(key) != 16 {
		log.Println("invalid pek length")
		return []byte{}
	}

	return key
}

func newHash(b, def, key1, key2 []byte, pek []Pek) Hash {
	if len(b) == 0 {
		return def
	}

	buf := b[8:] // skip header

	switch b[0] {
	case 0x13:
		buf = decryptAes(buf[20:36], pek[b[4]], buf[:16])

	default:
		key := deriveMd5(buf[:16], pek[0], 1)
		buf = decryptRc4(buf[16:], key)
	}

	return decryptDes(buf, key1, key2)
}

func newRecord(row *ordereddict.Dict, usr string, pek []Pek) (*Record, error) {
	sid := getBytesFromRow(row, userSid)
	rid := extractRid(sid)
	uac, _ := row.GetInt64(userUac)
	k1, k2 := deriveKey(rid)

	return &Record{
		Username: usr,
		Rid:      rid,
		Flags:    extractFlags(uac),
		LmHash:   hex.EncodeToString(newHash(getBytesFromRow(row, lmHash), defaultLm, k1, k2, pek)),
		NtHash:   hex.EncodeToString(newHash(getBytesFromRow(row, ntHash), defaultNt, k1, k2, pek)),
	}, nil
}

func getRow(row *ordereddict.Dict, key string) any {
	if v, ok := row.Get(key); ok && v != nil {
		return v
	}

	return nil
}

func getBytesFromRow(row *ordereddict.Dict, key string) []byte {
	if v := getRow(row, key); v != nil {
		b, _ := hex.DecodeString(v.(string))
		return b
	}

	return nil
}

func getKeys(ctl *parser.Catalog, att string, key []byte) []Pek {
	var keys []Pek

	_ = ctl.DumpTable("datatable", func(row *ordereddict.Dict) error {
		if v, ok := row.Get(att); ok && v != nil {
			b, _ := hex.DecodeString(v.(string))
			keys = append(keys, newPek(b, key))
		}
		return nil
	})

	return keys
}

func extractRid(sid []byte) uint32 {
	l, s := sid[1], sid[8:]

	return binary.BigEndian.Uint32(s[(l-1)*4 : (l-1)*4+4])
}

func extractFlags(v int64) *Flags {
	return &Flags{
		Script:                       v|1 == v,
		AccountDisable:               v|2 == v,
		HomeDirRequired:              v|8 == v,
		Lockout:                      v|6 == v,
		PasswordNotRequired:          v|32 == v,
		EncryptedTextPasswordAllowed: v|128 == v,
		TemporaryDupAccount:          v|256 == v,
		NormalAccount:                v|512 == v,
		InterDomainTrustAccount:      v|2048 == v,
		WorkstationTrustAccount:      v|4096 == v,
		ServerTrustAccount:           v|8192 == v,
		DontExpirePassword:           v|65536 == v,
		MNSLogonAccount:              v|131072 == v,
		SmartCardRequired:            v|262144 == v,
		TrustedForDelegation:         v|524288 == v,
		NotDelegated:                 v|1048576 == v,
		UseDESOnly:                   v|2097152 == v,
		DontPreAuth:                  v|4194304 == v,
		PasswordExpired:              v|8388608 == v,
		TrustedToAuthForDelegation:   v|16777216 == v,
		PartialSecrets:               v|67108864 == v,
	}
}

func decryptDes(b, key1, key2 []byte) []byte {
	var p []byte

	p1 := make([]byte, 8)
	p2 := make([]byte, 8)

	c1, err := des.NewCipher(key1)

	if err != nil {
		log.Println(err)
		return p
	}

	c1.Decrypt(p1, b[:8])

	p = append(p, p1...)

	c2, err := des.NewCipher(key2)

	if err != nil {
		log.Println(err)
		return p
	}

	c2.Decrypt(p2, b[8:])

	p = append(p, p2...)

	return p
}

func decryptAes(b, key, iv []byte) []byte {
	p := make([]byte, len(b))

	c, err := aes.NewCipher(key)

	if err != nil {
		log.Println(err)
		return p
	}

	d := cipher.NewCBCDecrypter(c, iv)
	d.CryptBlocks(p, b)

	return p
}

func decryptRc4(b, key []byte) []byte {
	p := make([]byte, len(b))

	c, err := rc4.NewCipher(key)

	if err != nil {
		log.Println(err)
		return p
	}

	c.XORKeyStream(p, b)

	return p
}

func deriveMd5(b, key []byte, rounds int) []byte {
	r := make([]byte, 16)

	h := md5.New()
	h.Write(key)

	for i := 0; i < rounds; i++ {
		h.Write(b)
	}

	s := h.Sum(nil)

	copy(r, s)

	return r
}

func deriveKey(rid uint32) ([]byte, []byte) {
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
