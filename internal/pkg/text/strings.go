package text

import (
	"regexp"
	"strings"
)

type String struct {
	re   *regexp.Regexp
	name []string
}

type Strings struct {
	db    []String
	level int
}

func GetStrings(level int) *Strings {
	str := &Strings{level: level}

	if level > 0 {
		str.db = append(str.db, getLevel1()...)
	}

	if level > 1 {
		str.db = append(str.db, getLevel2()...)
	}

	if level > 2 {
		str.db = append(str.db, getLevel3()...)
	}

	return str
}

func (str *String) String() string {
	return strings.Join(str.name, ",")
}

func (str *Strings) Search(s string) *String {
	for _, v := range str.db {
		if v.re.MatchString(s) {
			return &v
		}
	}

	return &String{name: []string{"unknown"}}
}

func getLevel1() []String {
	return []String{
		{
			regexp.MustCompile(`(([a-fA-F0-9]{1,4}:){7}[a-fA-F0-9]{1,4}|([a-fA-F0-9]{1,4}:){1,7}:|([a-fA-F0-9]{1,4}:){1,6}:[a-fA-F0-9]{1,4}|([a-fA-F0-9]{1,4}:){1,5}(:[a-fA-F0-9]{1,4}){1,2}|([a-fA-F0-9]{1,4}:){1,4}(:[a-fA-F0-9]{1,4}){1,3}|([a-fA-F0-9]{1,4}:){1,3}(:[a-fA-F0-9]{1,4}){1,4}|([a-fA-F0-9]{1,4}:){1,2}(:[a-fA-F0-9]{1,4}){1,5}|[a-fA-F0-9]{1,4}:((:[a-fA-F0-9]{1,4}){1,6})|:((:[a-fA-F0-9]{1,4}){1,7}|:)|fe80:(:[a-fA-F0-9]{0,4}){0,4}%[0-9a-zA-Z]+|::(ffff(:0{1,4})?:)?((25[0-5]|(2[0-4]|1?[0-9])?[0-9])\\.){3}(25[0-5]|(2[0-4]|1?[0-9])?[0-9])|([a-fA-F0-9]{1,4}:){1,4}:((25[0-5]|(2[0-4]|1?[0-9])?[0-9])\\.){3}(25[0-5]|(2[0-4]|1?[0-9])?[0-9]))`),
			[]string{"IPv6"},
		},
		{
			regexp.MustCompile(`(?:(?:25[0-5]|2[0-4][0-9]|1[0-9][0-9]|[1-9]?[0-9])\.){3}(?:25[0-5]|2[0-4][0-9]|1[0-9][0-9]|[1-9]?[0-9])`),
			[]string{"IPv4"},
		},
		{
			regexp.MustCompile(`([a-fA-F0-9]{2}[:-]){5}([a-fA-F0-9]{2})`),
			[]string{"MAC"},
		},
		{
			regexp.MustCompile("[a-z0-9!#$%&'*+/=?^_`{|}~-]+(?:\\.[a-z0-9!#$%&'*+/=?^_`{|}~-]+)*@(?:[a-z0-9](?:[a-z0-9-]*[a-z0-9])?\\.)+[a-z0-9](?:[a-z0-9-]*[a-z0-9])?"),
			[]string{"Mail"},
		},
		{
			regexp.MustCompile(`\\\\[a-z0-9 %._-]+\\[a-z0-9 $%._-]+`),
			[]string{"UNC"},
		},
		{
			regexp.MustCompile("[-a-zA-Z0-9@:%._+~#=]{1,256}\\.[a-zA-Z0-9()]{1,6}\\b[-a-zA-Z0-9()@:%_+.~#?&/=]*"),
			[]string{"URL"},
		},
		{
			regexp.MustCompile(`[a-z0-9+\-.]+://[a-z0-9\-._~%!$&'()*+,;=]+@`),
			[]string{"URL"},
		},
		{
			regexp.MustCompile("[a-fA-F0-9]{8}(?:-[a-fA-F0-9]{4}){3}-[a-fA-F0-9]{12}"),
			[]string{"UUID"},
		},
		{
			regexp.MustCompile(`S-\d-\d+-(\d+-){1,14}\d+`),
			[]string{"SID"},
		},
		{
			regexp.MustCompile("[0-9]{6}?-[0-9]{6}-[0-9]{6}-[0-9]{6}-[0-9]{6}-[0-9]{6}-[0-9]{6}-[0-9]{6}"),
			[]string{"BitLocker"},
		},
		{
			regexp.MustCompile("HKEY_[A-Z_]+\\\\([A-Za-z0-9]+\\\\+)+[A-Za-z0-9]+"),
			[]string{"Registry"},
		},
		{
			regexp.MustCompile(`(?:""?[a-zA-Z]:|\\+[^/:*?<>|]+\\+[^/:*?<>|]*)\\+(?:[^/:*?<>|]+\\+)*\w([^/:*?<>|])*`),
			[]string{"Windows"},
		},
	}
}

// source: https://github.com/EricZimmerman/bstrings/blob/master/bstrings/Program.cs
func getLevel2() []String {
	return []String{
		{
			regexp.MustCompile(`[13][a-km-zA-HJ-NP-Z1-9]{25,34}`),
			[]string{"Bitcoin"},
		},
		{
			regexp.MustCompile("Wm[st][0-9a-zA-Z]{94}"),
			[]string{"Aeon"},
		},
		{
			regexp.MustCompile("2[0-9AB][0-9a-zA-Z]{93}"),
			[]string{"Bytecoin"},
		},
		{
			regexp.MustCompile("D[0-9a-zA-Z]{94}"),
			[]string{"Dashcoin"},
		},
		{
			regexp.MustCompile("([7X])[a-zA-Z0-9]{33}"),
			[]string{"Dashcoin2"},
		},
		{
			regexp.MustCompile("6[0-9a-zA-Z]{94}"),
			[]string{"Fantomcoin"},
		},
		{
			regexp.MustCompile("4[0-9AB][0-9a-zA-Z]{93}|4[0-9AB][0-9a-zA-Z]{104}"),
			[]string{"Monero"},
		},
		{
			regexp.MustCompile("Sumoo[0-9a-zA-Z]{94}"),
			[]string{"Sumokoin"},
		},
		{
			regexp.MustCompile(`\A\b[0-9]{5}(?:-[0-9]{4})?\b\z`),
			[]string{"Zip"},
		},
		{
			regexp.MustCompile(`\(?\b[2-9][0-9]{2}\)?[-. ]?[2-9][0-9]{2}[-. ]?[0-9]{4}`),
			[]string{"Phone"},
		},
		{
			regexp.MustCompile(`4[0-9]{12}(?:[0-9]{3})?|[25][1-7][0-9]{14}|6(?:011|5[0-9][0-9])[0-9]{12}|3[47][0-9]{13}|3(?:0[0-5]|[68][0-9])[0-9]{11}|(?:2131|1800|35\d{3})\d{11}`),
			[]string{"CC"},
		},
		{
			regexp.MustCompile(`\d{3}-?\d{2}-?\d{4}`),
			[]string{"SSN"},
		},
	}
}

// source: https://github.com/s0md3v/Bolt/blob/master/db/hashes.json
func getLevel3() []String {
	return []String{
		{
			regexp.MustCompile("^[a-f0-9]{4}"),
			[]string{
				"CRC-16",
				"CRC-16-CCITT",
				"FCS-16",
			},
		},
		{
			regexp.MustCompile("^[a-f0-9]{8}"),
			[]string{
				"Adler-32",
				"CRC-32B",
				"FCS-32",
				"GHash-32-3",
				"GHash-32-5",
				"FNV-132",
				"Fletcher-32",
				"Joaat",
				"ELF-32",
				"XOR-32",
			},
		},
		{
			regexp.MustCompile("^[a-f0-9]{6}"),
			[]string{
				"CRC-24",
			},
		},
		{
			regexp.MustCompile("^(\\$crc32\\$[a-f0-9]{8}.)?[a-f0-9]{8}"),
			[]string{
				"CRC-32",
			},
		},
		{
			regexp.MustCompile("^\\+[a-z0-9/.]{12}"),
			[]string{
				"Eggdrop IRC Bot",
			},
		},
		{
			regexp.MustCompile("^[a-z0-9/.]{13}"),
			[]string{
				"DES(Unix)",
				"Traditional DES",
				"DEScrypt",
			},
		},
		{
			regexp.MustCompile("^[a-f0-9]{16}"),
			[]string{
				"MySQL323",
				"DES(Oracle)",
				"Half MD5",
				"Oracle 7-10g",
				"FNV-164",
				"CRC-64",
			},
		},
		{
			regexp.MustCompile("^[a-z0-9/.]{16}"),
			[]string{
				"Cisco-PIX(MD5)",
			},
		},
		{
			regexp.MustCompile("^\\([a-z0-9/+]{20}\\)"),
			[]string{
				"Lotus Notes/Domino 6",
			},
		},
		{
			regexp.MustCompile("^_[a-z0-9/.]{19}"),
			[]string{
				"BSDi Crypt",
			},
		},
		{
			regexp.MustCompile("^[a-f0-9]{24}"),
			[]string{
				"CRC-96(ZIP)",
			},
		},
		{
			regexp.MustCompile("^[a-z0-9/.]{24}"),
			[]string{
				"Crypt16",
			},
		},
		{
			regexp.MustCompile("^(\\$md2\\$)?[a-f0-9]{32}"),
			[]string{
				"MD2",
			},
		},
		{
			regexp.MustCompile("^[a-f0-9]{32}(:.+)?"),
			[]string{
				"MD5",
				"MD4",
				"Double MD5",
				"LM",
				"RIPEMD-128",
				"Haval-128",
				"Tiger-128",
				"Skein-256(128)",
				"Skein-512(128)",
				"Lotus Notes/Domino 5",
				"Skype",
				"ZipMonster",
				"PrestaShop",
				"md5(md5(md5($pass)))",
				"md5(strtoupper(md5($pass)))",
				"md5(sha1($pass))",
				"md5($pass.$salt)",
				"md5($salt.$pass)",
				"md5(unicode($pass).$salt)",
				"md5($salt.unicode($pass))",
				"HMAC-MD5 (key = $pass)",
				"HMAC-MD5 (key = $salt)",
				"md5(md5($salt).$pass)",
				"md5($salt.md5($pass))",
				"md5($pass.md5($salt))",
				"md5($salt.$pass.$salt)",
				"md5(md5($pass).md5($salt))",
				"md5($salt.md5($salt.$pass))",
				"md5($salt.md5($pass.$salt))",
				"md5($username.0.$pass)",
			},
		},
		{
			regexp.MustCompile("^(\\$snefru\\$)?[a-f0-9]{32}"),
			[]string{
				"Snefru-128",
			},
		},
		{
			regexp.MustCompile("^(\\$NT\\$)?[a-f0-9]{32}"),
			[]string{
				"NTLM",
			},
		},
		{
			regexp.MustCompile("^([^/:*?\"<>|]{1,20}:)?[a-f0-9]{32}(:[^/:*?\"<>|]{1,20})?"),
			[]string{
				"Domain Cached Credentials",
			},
		},
		{
			regexp.MustCompile("^([^/:*?\"<>|]{1,20}:)?(\\$DCC2\\$10240#[^/:*?\"<>|]{1,20}#)?[a-f0-9]{32}"),
			[]string{
				"Domain Cached Credentials 2",
			},
		},
		{
			regexp.MustCompile("^{SHA}[a-z0-9/+]{27}="),
			[]string{
				"SHA-1(Base64)",
				"Netscape LDAP SHA",
			},
		},
		{
			regexp.MustCompile("^\\$1\\$[a-z0-9/.]{0,8}\\$[a-z0-9/.]{22}(:.*)?"),
			[]string{
				"MD5 Crypt",
				"Cisco-IOS(MD5)",
				"FreeBSD MD5",
			},
		},
		{
			regexp.MustCompile("^0x[a-f0-9]{32}"),
			[]string{
				"Lineage II C4",
			},
		},
		{
			regexp.MustCompile("^\\$H\\$[a-z0-9/.]{31}"),
			[]string{
				"phpBB v3.x",
				"Wordpress v2.6.0/2.6.1",
				"PHPass' Portable Hash",
			},
		},
		{
			regexp.MustCompile("^\\$P\\$[a-z0-9/.]{31}"),
			[]string{
				"Wordpress \u2265 v2.6.2",
				"Joomla \u2265 v2.5.18",
				"PHPass' Portable Hash",
			},
		},
		{
			regexp.MustCompile("^[a-f0-9]{32}:[a-z0-9]{2}"),
			[]string{
				"osCommerce",
				"xt:Commerce",
			},
		},
		{
			regexp.MustCompile("^\\$apr1\\$[a-z0-9/.]{0,8}\\$[a-z0-9/.]{22}"),
			[]string{
				"MD5(APR)",
				"Apache MD5",
				"md5apr1",
			},
		},
		{
			regexp.MustCompile("^{smd5}[a-z0-9$/.]{31}"),
			[]string{
				"AIX(smd5)",
			},
		},
		{
			regexp.MustCompile("^[a-f0-9]{32}:[a-f0-9]{32}"),
			[]string{
				"WebEdition CMS",
			},
		},
		{
			regexp.MustCompile("^[a-f0-9]{32}:.{5}"),
			[]string{
				"IP.Board \u2265 v2+",
			},
		},
		{
			regexp.MustCompile("^[a-f0-9]{32}:.{8}"),
			[]string{
				"MyBB \u2265 v1.2+",
			},
		},
		{
			regexp.MustCompile("^[a-z0-9]{34}"),
			[]string{
				"CryptoCurrency(Address)",
			},
		},
		{
			regexp.MustCompile("^[a-f0-9]{40}(:.+)?"),
			[]string{
				"SHA-1",
				"Double SHA-1",
				"RIPEMD-160",
				"Haval-160",
				"Tiger-160",
				"HAS-160",
				"LinkedIn",
				"Skein-256(160)",
				"Skein-512(160)",
				"MangosWeb Enhanced CMS",
				"sha1(sha1(sha1($pass)))",
				"sha1(md5($pass))",
				"sha1($pass.$salt)",
				"sha1($salt.$pass)",
				"sha1(unicode($pass).$salt)",
				"sha1($salt.unicode($pass))",
				"HMAC-SHA1 (key = $pass)",
				"HMAC-SHA1 (key = $salt)",
				"sha1($salt.$pass.$salt)",
			},
		},
		{
			regexp.MustCompile("^\\*[a-f0-9]{40}"),
			[]string{
				"MySQL5.x",
				"MySQL4.1",
			},
		},
		{
			regexp.MustCompile("^[a-z0-9]{43}"),
			[]string{
				"Cisco-IOS(SHA-256)",
			},
		},
		{
			regexp.MustCompile("^{SSHA}[a-z0-9/+]{38}=="),
			[]string{
				"SSHA-1(Base64)",
				"Netscape LDAP SSHA",
				"nsldaps",
			},
		},
		{
			regexp.MustCompile("^[a-z0-9=]{47}"),
			[]string{
				"Fortigate(FortiOS)",
			},
		},
		{
			regexp.MustCompile("^[a-f0-9]{48}"),
			[]string{
				"Haval-192",
				"Tiger-192",
				"SHA-1(Oracle)",
				"OSX v10.4",
				"OSX v10.5",
				"OSX v10.6",
			},
		},
		{
			regexp.MustCompile("^[a-f0-9]{51}"),
			[]string{
				"Palshop CMS",
			},
		},
		{
			regexp.MustCompile("^[a-z0-9]{51}"),
			[]string{
				"CryptoCurrency(PrivateKey)",
			},
		},
		{
			regexp.MustCompile("^{ssha1}[0-9]{2}\\$[a-z0-9$/.]{44}"),
			[]string{
				"AIX(ssha1)",
			},
		},
		{
			regexp.MustCompile("^0x0100[a-f0-9]{48}"),
			[]string{
				"MSSQL(2005)",
				"MSSQL(2008)",
			},
		},
		{
			regexp.MustCompile("^(\\$md5,rounds=[0-9]+\\$|\\$md5\\$rounds=[0-9]+\\$|\\$md5\\$)[a-z0-9/.]{0,16}(\\$|\\$\\$)[a-z0-9/.]{22}"),
			[]string{
				"Sun MD5 Crypt",
			},
		},
		{
			regexp.MustCompile("^[a-f0-9]{56}"),
			[]string{
				"SHA-224",
				"Haval-224",
				"SHA3-224",
				"Skein-256(224)",
				"Skein-512(224)",
			},
		},
		{
			regexp.MustCompile("^(\\$2[axy]|\\$2)\\$[0-9]{2}\\$[a-z0-9/.]{53}"),
			[]string{
				"Blowfish(OpenBSD)",
				"Woltlab Burning Board 4.x",
				"bcrypt",
			},
		},
		{
			regexp.MustCompile("^[a-f0-9]{40}:[a-f0-9]{16}"),
			[]string{
				"Android PIN",
			},
		},
		{
			regexp.MustCompile("^(S:)?[a-f0-9]{40}(:)?[a-f0-9]{20}"),
			[]string{
				"Oracle 11g/12c",
			},
		},
		{
			regexp.MustCompile("^\\$bcrypt-sha256\\$(2[axy]|2),[0-9]+\\$[a-z0-9/.]{22}\\$[a-z0-9/.]{31}"),
			[]string{
				"bcrypt(SHA-256)",
			},
		},
		{
			regexp.MustCompile("^[a-f0-9]{32}:.{3}"),
			[]string{
				"vBulletin < v3.8.5",
			},
		},
		{
			regexp.MustCompile("^[a-f0-9]{32}:.{30}"),
			[]string{
				"vBulletin \u2265 v3.8.5",
			},
		},
		{
			regexp.MustCompile("^(\\$snefru\\$)?[a-f0-9]{64}"),
			[]string{
				"Snefru-256",
			},
		},
		{
			regexp.MustCompile("^[a-f0-9]{64}(:.+)?"),
			[]string{
				"SHA-256",
				"RIPEMD-256",
				"Haval-256",
				"GOST R 34.11-94",
				"GOST CryptoPro S-Box",
				"SHA3-256",
				"Skein-256",
				"Skein-512(256)",
				"Ventrilo",
				"sha256($pass.$salt)",
				"sha256($salt.$pass)",
				"sha256(unicode($pass).$salt)",
				"sha256($salt.unicode($pass))",
				"HMAC-SHA256 (key = $pass)",
				"HMAC-SHA256 (key = $salt)",
			},
		},
		{
			regexp.MustCompile("^[a-f0-9]{32}:[a-z0-9]{32}"),
			[]string{
				"Joomla < v2.5.18",
			},
		},
		{
			regexp.MustCompile("^[a-f-0-9]{32}:[a-f-0-9]{32}"),
			[]string{
				"SAM(LM_Hash:NT_Hash)",
			},
		},
		{
			regexp.MustCompile("^(\\$chap\\$0\\*)?[a-f0-9]{32}[*:][a-f0-9]{32}(:[0-9]{2})?"),
			[]string{
				"MD5(Chap)",
				"iSCSI CHAP Authentication",
			},
		},
		{
			regexp.MustCompile("^\\$episerver\\$\\*0\\*[a-z0-9/=+]+\\*[a-z0-9/=+]{27,28}"),
			[]string{
				"EPiServer 6.x < v4",
			},
		},
		{
			regexp.MustCompile("^{ssha256}[0-9]{2}\\$[a-z0-9$/.]{60}"),
			[]string{
				"AIX(ssha256)",
			},
		},
		{
			regexp.MustCompile("^[a-f0-9]{80}"),
			[]string{
				"RIPEMD-320",
			},
		},
		{
			regexp.MustCompile("^\\$episerver\\$\\*1\\*[a-z0-9/=+]+\\*[a-z0-9/=+]{42,43}"),
			[]string{
				"EPiServer 6.x \u2265 v4",
			},
		},
		{
			regexp.MustCompile("^0x0100[a-f0-9]{88}"),
			[]string{
				"MSSQL(2000)",
			},
		},
		{
			regexp.MustCompile("^[a-f0-9]{96}"),
			[]string{
				"SHA-384",
				"SHA3-384",
				"Skein-512(384)",
				"Skein-1024(384)",
			},
		},
		{
			regexp.MustCompile("^{SSHA512}[a-z0-9/+]{96}"),
			[]string{
				"SSHA-512(Base64)",
				"LDAP(SSHA-512)",
			},
		},
		{
			regexp.MustCompile("^{ssha512}[0-9]{2}\\$[a-z0-9/.]{16,48}\\$[a-z0-9/.]{86}"),
			[]string{
				"AIX(ssha512)",
			},
		},
		{
			regexp.MustCompile("^[a-f0-9]{128}(:.+)?"),
			[]string{
				"SHA-512",
				"Whirlpool",
				"Salsa10",
				"Salsa20",
				"SHA3-512",
				"Skein-512",
				"Skein-1024(512)",
				"sha512($pass.$salt)",
				"sha512($salt.$pass)",
				"sha512(unicode($pass).$salt)",
				"sha512($salt.unicode($pass))",
				"HMAC-SHA512 (key = $pass)",
				"HMAC-SHA512 (key = $salt)",
			},
		},
		{
			regexp.MustCompile("^[a-f0-9]{136}"),
			[]string{
				"OSX v10.7",
			},
		},
		{
			regexp.MustCompile("^0x0200[a-f0-9]{136}"),
			[]string{
				"MSSQL(2012)",
				"MSSQL(2014)",
			},
		},
		{
			regexp.MustCompile("^\\$ml\\$[0-9]+\\$[a-f0-9]{64}\\$[a-f0-9]{128}"),
			[]string{
				"OSX v10.8",
				"OSX v10.9",
			},
		},
		{
			regexp.MustCompile("^[a-f0-9]{256}"),
			[]string{
				"Skein-1024",
			},
		},
		{
			regexp.MustCompile("^grub\\.pbkdf2\\.sha512\\.[0-9]+\\.([a-f0-9]{128,1000}[a-f0-9]{0,1000}[a-f0-9]{0,48}\\.|[0-9]+\\.)?[a-f0-9]{128}"),
			[]string{
				"GRUB 2",
			},
		},
		{
			regexp.MustCompile("^sha1\\$[a-z0-9]+\\$[a-f0-9]{40}"),
			[]string{
				"Django(SHA-1)",
			},
		},
		{
			regexp.MustCompile("^[a-f0-9]{49}"),
			[]string{
				"Citrix Netscaler",
			},
		},
		{
			regexp.MustCompile("^\\$S\\$[a-z0-9/.]{52}"),
			[]string{
				"Drupal > v7.x",
			},
		},
		{
			regexp.MustCompile("^\\$5\\$(rounds=[0-9]+\\$)?[a-z0-9/.]{0,16}\\$[a-z0-9/.]{43}"),
			[]string{
				"SHA-256 Crypt",
			},
		},
		{
			regexp.MustCompile("^0x[a-f0-9]{4}[a-f0-9]{16}[a-f0-9]{64}"),
			[]string{
				"Sybase ASE",
			},
		},
		{
			regexp.MustCompile("^\\$6\\$(rounds=[0-9]+\\$)?[a-z0-9/.]{0,16}\\$[a-z0-9/.]{86}"),
			[]string{
				"SHA-512 Crypt",
			},
		},
		{
			regexp.MustCompile("^\\$sha\\$[a-z0-9]{1,16}\\$([a-f0-9]{32}|[a-f0-9]{40}|[a-f0-9]{64}|[a-f0-9]{128}|[a-f0-9]{140})"),
			[]string{
				"Minecraft(AuthMe Reloaded)",
			},
		},
		{
			regexp.MustCompile("^sha256\\$[a-z0-9]+\\$[a-f0-9]{64}"),
			[]string{
				"Django(SHA-256)",
			},
		},
		{
			regexp.MustCompile("^sha384\\$[a-z0-9]+\\$[a-f0-9]{96}"),
			[]string{
				"Django(SHA-384)",
			},
		},
		{
			regexp.MustCompile("^crypt1:[a-z0-9+=]{12}:[a-z0-9+=]{12}"),
			[]string{
				"Clavister Secure Gateway",
			},
		},
		{
			regexp.MustCompile("^[a-f0-9]{112}"),
			[]string{
				"Cisco VPN Client(PCF-File)",
			},
		},
		{
			regexp.MustCompile("^[a-f0-9]{1000}[a-f0-9]{329}"),
			[]string{
				"Microsoft MSTSC(RDP-File)",
			},
		},
		{
			regexp.MustCompile("^[^\\\\/:*?\"<>|]{1,20}:{2,3}([^\\\\/:*?\"<>|]{1,20})?:[a-f0-9]{48}:[a-f0-9]{48}:[a-f0-9]{16}"),
			[]string{
				"NetNTLMv1-VANILLA / NetNTLMv1+ESS",
			},
		},
		{
			regexp.MustCompile("^([^\\\\/:*?\"<>|]{1,20}\\\\)?[^\\\\/:*?\"<>|]{1,20}:{2,3}([^\\\\/:*?\"<>|]{1,20}:)?[^\\\\/:*?\"<>|]{1,20}:[a-f0-9]{32}:[a-f0-9]+"),
			[]string{
				"NetNTLMv2",
			},
		},
		{
			regexp.MustCompile("^\\$(krb5pa|mskrb5)\\$([0-9]{2})?\\$.+\\$[a-f0-9]+"),
			[]string{
				"Kerberos 5 AS-REQ Pre-Auth",
			},
		},
		{
			regexp.MustCompile("^\\$scram\\$[0-9]+\\$[a-z0-9/.]{16}\\$sha-1=[a-z0-9/.]{27},sha-256=[a-z0-9/.]{43},sha-512=[a-z0-9/.]{86}"),
			[]string{
				"SCRAM Hash",
			},
		},
		{
			regexp.MustCompile("^[a-f0-9]{40}:[a-f0-9]{0,32}"),
			[]string{
				"Redmine Project Management Web App",
			},
		},
		{
			regexp.MustCompile("^(.+)?\\$[a-f0-9]{16}"),
			[]string{
				"SAP CODVN B (BCODE)",
			},
		},
		{
			regexp.MustCompile("^(.+)?\\$[a-f0-9]{40}"),
			[]string{
				"SAP CODVN F/G (PASSCODE)",
			},
		},
		{
			regexp.MustCompile("^(.+\\$)?[a-z0-9/.+]{30}(:.+)?"),
			[]string{
				"Juniper Netscreen/SSG(ScreenOS)",
			},
		},
		{
			regexp.MustCompile("^0x[a-f0-9]{60}\\s0x[a-f0-9]{40}"),
			[]string{
				"EPi",
			},
		},
		{
			regexp.MustCompile("^[a-f0-9]{40}:[^*]{1,25}"),
			[]string{
				"SMF \u2265 v1.1",
			},
		},
		{
			regexp.MustCompile("^(\\$wbb3\\$\\*1\\*)?[a-f0-9]{40}[:*][a-f0-9]{40}"),
			[]string{
				"Woltlab Burning Board 3.x",
			},
		},
		{
			regexp.MustCompile("^[a-f0-9]{130}(:[a-f0-9]{40})?"),
			[]string{
				"IPMI2 RAKP HMAC-SHA1",
			},
		},
		{
			regexp.MustCompile("^[a-f0-9]{32}:[0-9]+:[a-z0-9_.+-]+@[a-z0-9-]+\\.[a-z0-9-.]+"),
			[]string{
				"Lastpass",
			},
		},
		{
			regexp.MustCompile("^[a-z0-9/.]{16}([:$].+)?"),
			[]string{
				"Cisco-ASA(MD5)",
			},
		},
		{
			regexp.MustCompile("^\\$vnc\\$\\*[a-f0-9]{32}\\*[a-f0-9]{32}"),
			[]string{
				"VNC",
			},
		},
		{
			regexp.MustCompile("^[a-z0-9]{32}(:([a-z0-9-]+\\.)?[a-z0-9-.]+\\.[a-z]{2,7}:.+:[0-9]+)?"),
			[]string{
				"DNSSEC(NSEC3)",
			},
		},
		{
			regexp.MustCompile("^(user-.+:)?\\$racf\\$\\*.+\\*[a-f0-9]{16}"),
			[]string{
				"RACF",
			},
		},
		{
			regexp.MustCompile("^\\$3\\$\\$[a-f0-9]{32}"),
			[]string{
				"NTHash(FreeBSD Variant)",
			},
		},
		{
			regexp.MustCompile("^\\$sha1\\$[0-9]+\\$[a-z0-9/.]{0,64}\\$[a-z0-9/.]{28}"),
			[]string{
				"SHA-1 Crypt",
			},
		},
		{
			regexp.MustCompile("^[a-f0-9]{70}"),
			[]string{
				"hMailServer",
			},
		},
		{
			regexp.MustCompile("^[:$][AB][:$]([a-f0-9]{1,8}[:$])?[a-f0-9]{32}"),
			[]string{
				"MediaWiki",
			},
		},
		{
			regexp.MustCompile("^[a-f0-9]{140}"),
			[]string{
				"Minecraft(xAuth)",
			},
		},
		{
			regexp.MustCompile("^\\$pbkdf2(-sha1)?\\$[0-9]+\\$[a-z0-9/.]+\\$[a-z0-9/.]{27}"),
			[]string{
				"PBKDF2-SHA1(Generic)",
			},
		},
		{
			regexp.MustCompile("^\\$pbkdf2-sha256\\$[0-9]+\\$[a-z0-9/.]+\\$[a-z0-9/.]{43}"),
			[]string{
				"PBKDF2-SHA256(Generic)",
			},
		},
		{
			regexp.MustCompile("^\\$pbkdf2-sha512\\$[0-9]+\\$[a-z0-9/.]+\\$[a-z0-9/.]{86}"),
			[]string{
				"PBKDF2-SHA512(Generic)",
			},
		},
		{
			regexp.MustCompile("^\\$p5k2\\$[0-9]+\\$[a-z0-9/+=-]+\\$[a-z0-9/+-]{27}="),
			[]string{
				"PBKDF2(Cryptacular)",
			},
		},
		{
			regexp.MustCompile("^\\$p5k2\\$[0-9]+\\$[a-z0-9/.]+\\$[a-z0-9/.]{32}"),
			[]string{
				"PBKDF2(Dwayne Litzenberger)",
			},
		},
		{
			regexp.MustCompile("^{FSHP[0123]\\|[0-9]+\\|[0-9]+}[a-z0-9/+=]+"),
			[]string{
				"Fairly Secure Hashed Password",
			},
		},
		{
			regexp.MustCompile("^\\$PHPS\\$.+\\$[a-f0-9]{32}"),
			[]string{
				"PHPS",
			},
		},
		{
			regexp.MustCompile("^[0-9]{4}:[a-f0-9]{16}:[a-f0-9]{1000}[a-f0-9]{1000}[a-f0-9]{80}"),
			[]string{
				"1Password(Agile Keychain)",
			},
		},
		{
			regexp.MustCompile("^[a-f0-9]{64}:[a-f0-9]{32}:[0-9]{5}:[a-f0-9]{608}"),
			[]string{
				"1Password(Cloud Keychain)",
			},
		},
		{
			regexp.MustCompile("^[a-f0-9]{256}:[a-f0-9]{256}:[a-f0-9]{16}:[a-f0-9]{16}:[a-f0-9]{320}:[a-f0-9]{16}:[a-f0-9]{40}:[a-f0-9]{40}:[a-f0-9]{32}"),
			[]string{
				"IKE-PSK MD5",
			},
		},
		{
			regexp.MustCompile("^[a-f0-9]{256}:[a-f0-9]{256}:[a-f0-9]{16}:[a-f0-9]{16}:[a-f0-9]{320}:[a-f0-9]{16}:[a-f0-9]{40}:[a-f0-9]{40}:[a-f0-9]{40}"),
			[]string{
				"IKE-PSK SHA1",
			},
		},
		{
			regexp.MustCompile("^[a-z0-9/+]{27}="),
			[]string{
				"PeopleSoft",
			},
		},
		{
			regexp.MustCompile("^crypt\\$[a-f0-9]{5}\\$[a-z0-9/.]{13}"),
			[]string{
				"Django(DES Crypt Wrapper)",
			},
		},
		{
			regexp.MustCompile("^(\\$django\\$\\*1\\*)?pbkdf2_sha256\\$[0-9]+\\$[a-z0-9]+\\$[a-z0-9/+=]{44}"),
			[]string{
				"Django(PBKDF2-HMAC-SHA256)",
			},
		},
		{
			regexp.MustCompile("^pbkdf2_sha1\\$[0-9]+\\$[a-z0-9]+\\$[a-z0-9/+=]{28}"),
			[]string{
				"Django(PBKDF2-HMAC-SHA1)",
			},
		},
		{
			regexp.MustCompile("^bcrypt(\\$2[axy]|\\$2)\\$[0-9]{2}\\$[a-z0-9/.]{53}"),
			[]string{
				"Django(bcrypt)",
			},
		},
		{
			regexp.MustCompile("^md5\\$[a-f0-9]+\\$[a-f0-9]{32}"),
			[]string{
				"Django(MD5)",
			},
		},
		{
			regexp.MustCompile("^\\{PKCS5S2}[a-z0-9/+]{64}"),
			[]string{
				"PBKDF2(Atlassian)",
			},
		},
		{
			regexp.MustCompile("^md5[a-f0-9]{32}"),
			[]string{
				"PostgreSQL MD5",
			},
		},
		{
			regexp.MustCompile("^\\([a-z0-9/+]{49}\\)"),
			[]string{
				"Lotus Notes/Domino 8",
			},
		},
		{
			regexp.MustCompile("^SCRYPT:[0-9]+:[0-9]:[0-9]:[a-z0-9:/+=]+"),
			[]string{
				"scrypt",
			},
		},
		{
			regexp.MustCompile("^\\$8\\$[a-z0-9/.]{14}\\$[a-z0-9/.]{43}"),
			[]string{
				"Cisco Type 8",
			},
		},
		{
			regexp.MustCompile("^\\$9\\$[a-z0-9/.]{14}\\$[a-z0-9/.]{43}"),
			[]string{
				"Cisco Type 9",
			},
		},
		{
			regexp.MustCompile("^\\$office\\$\\*2007\\*[0-9]{2}\\*[0-9]{3}\\*[0-9]{2}\\*[a-z0-9]{32}\\*[a-z0-9]{32}\\*[a-z0-9]{40}"),
			[]string{
				"Microsoft Office 2007",
			},
		},
		{
			regexp.MustCompile("^\\$office\\$\\*2010\\*[0-9]{6}\\*[0-9]{3}\\*[0-9]{2}\\*[a-z0-9]{32}\\*[a-z0-9]{32}\\*[a-z0-9]{64}"),
			[]string{
				"Microsoft Office 2010",
			},
		},
		{
			regexp.MustCompile("^\\$office\\$\\*2013\\*[0-9]{6}\\*[0-9]{3}\\*[0-9]{2}\\*[a-z0-9]{32}\\*[a-z0-9]{32}\\*[a-z0-9]{64}"),
			[]string{
				"Microsoft Office 2013",
			},
		},
		{
			regexp.MustCompile("^\\$fde\\$[0-9]{2}\\$[a-f0-9]{32}\\$[0-9]{2}\\$[a-f0-9]{32}\\$[a-f0-9]{1000}[a-f0-9]{1000}[a-f0-9]{1000}[a-f0-9]{72}"),
			[]string{
				"Android FDE \u2264 4.3",
			},
		},
		{
			regexp.MustCompile("^\\$oldoffice\\$[01]\\*[a-f0-9]{32}\\*[a-f0-9]{32}\\*[a-f0-9]{32}"),
			[]string{
				"Microsoft Office \u2264 2003 (MD5+RC4)",
				"Microsoft Office \u2264 2003 (MD5+RC4) collider-mode #1",
				"Microsoft Office \u2264 2003 (MD5+RC4) collider-mode #2",
			},
		},
		{
			regexp.MustCompile("^\\$oldoffice\\$[34]\\*[a-f0-9]{32}\\*[a-f0-9]{32}\\*[a-f0-9]{40}"),
			[]string{
				"Microsoft Office \u2264 2003 (SHA1+RC4)",
				"Microsoft Office \u2264 2003 (SHA1+RC4) collider-mode #1",
				"Microsoft Office \u2264 2003 (SHA1+RC4) collider-mode #2",
			},
		},
		{
			regexp.MustCompile("^(\\$radmin2\\$)?[a-f0-9]{32}"),
			[]string{
				"RAdmin v2.x",
			},
		},
		{
			regexp.MustCompile("^{x-issha,\\s[0-9]{4}}[a-z0-9/+=]+"),
			[]string{
				"SAP CODVN H (PWDSALTEDHASH) iSSHA-1",
			},
		},
		{
			regexp.MustCompile("^\\$cram_md5\\$[a-z0-9/+=-]+\\$[a-z0-9/+=-]{52}"),
			[]string{
				"CRAM-MD5",
			},
		},
		{
			regexp.MustCompile("^[a-f0-9]{16}:2:4:[a-f0-9]{32}"),
			[]string{
				"SipHash",
			},
		},
		{
			regexp.MustCompile("^[a-f0-9]{4,}"),
			[]string{
				"Cisco Type 7",
			},
		},
		{
			regexp.MustCompile("^[a-z0-9/.]{13,}"),
			[]string{
				"BigCrypt",
			},
		},
		{
			regexp.MustCompile("^(\\$cisco4\\$)?[a-z0-9/.]{43}"),
			[]string{
				"Cisco Type 4",
			},
		},
		{
			regexp.MustCompile("^bcrypt_sha256\\$\\$(2[axy]|2)\\$[0-9]+\\$[a-z0-9/.]{53}"),
			[]string{
				"Django(bcrypt-SHA256)",
			},
		},
		{
			regexp.MustCompile("^\\$postgres\\$.[^*]+[*:][a-f0-9]{1,32}[*:][a-f0-9]{32}"),
			[]string{
				"PostgreSQL Challenge-Response Authentication (MD5)",
			},
		},
		{
			regexp.MustCompile("^\\$siemens-s7\\$[0-9]\\$[a-f0-9]{40}\\$[a-f0-9]{40}"),
			[]string{
				"Siemens-S7",
			},
		},
		{
			regexp.MustCompile("^(\\$pst\\$)?[a-f0-9]{8}"),
			[]string{
				"Microsoft Outlook PST",
			},
		},
		{
			regexp.MustCompile("^sha256[:$][0-9]+[:$][a-z0-9/+]+[:$][a-z0-9/+]{32,128}"),
			[]string{
				"PBKDF2-HMAC-SHA256(PHP)",
			},
		},
		{
			regexp.MustCompile("^(\\$dahua\\$)?[a-z0-9]{8}"),
			[]string{
				"Dahua",
			},
		},
		{
			regexp.MustCompile("^\\$mysqlna\\$[a-f0-9]{40}[:*][a-f0-9]{40}"),
			[]string{
				"MySQL Challenge-Response Authentication (SHA1)",
			},
		},
		{
			regexp.MustCompile("^\\$pdf\\$[24]\\*[34]\\*128\\*[0-9-]{1,5}\\*1\\*(16|32)\\*[a-f0-9]{32,64}\\*32\\*[a-f0-9]{64}\\*(8|16|32)\\*[a-f0-9]{16,64}"),
			[]string{
				"PDF 1.4 - 1.6 (Acrobat 5 - 8)",
			},
		},
	}
}
