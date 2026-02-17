package share

import (
	"fmt"
	"io/fs"
	"log"
	"net"
	"os"
	"regexp"
	"strings"

	"github.com/hirochachacha/go-smb2"
)

//goland:noinspection ALL
var re = regexp.MustCompile(`^//((?<user>.+?)(:(?<pass>.*?))?@)?(?<host>.+?)?(:(?<port>\d+?))?/(?<root>.+?)/(?<path>.*?)(:(?<part>.+))?$`)

type unc struct {
	User string
	Pass string
	Host string
	Root string
	Port string
	Path string
	Part string
}

type Share struct {
	unc *unc
	cn  net.Conn
	se  *smb2.Session
	fs  *smb2.Share
}

func New(path string) *Share {
	return &Share{unc: parse(path)}
}

func (unc *unc) String() string {
	var cred string
	var host = unc.Host
	var root = unc.Root
	var path = unc.Path

	if len(unc.User) > 0 {
		cred = unc.User
	}

	if len(unc.Pass) > 0 {
		cred += ":" + strings.Repeat("*", len(unc.Pass))
	}

	if len(cred) > 0 {
		cred += "@"
	}

	if len(unc.Port) > 0 {
		host += ":" + unc.Port
	}

	if len(unc.Part) > 0 {
		path += ":" + unc.Part
	}

	return fmt.Sprintf("//%s%s/%s/%s", cred, host, root, path)
}

func (shr *Share) String() string {
	return fmt.Sprintf("//%s/%s", shr.unc.Host, shr.unc.Root)
}

func (shr *Share) DirFS(path string) fs.FS {
	return shr.fs.DirFS(path)
}

func (shr *Share) Open(path string) (*smb2.File, error) {
	return shr.fs.OpenFile(path, os.O_RDONLY, 0400)
}

func (shr *Share) List() ([]string, error) {
	return shr.se.ListSharenames()
}

func (shr *Share) Mount() {
	var err error

	addr := fmt.Sprintf("%s:%s", shr.unc.Host, shr.unc.Port)

	shr.cn, err = net.Dial("tcp", addr)

	if err != nil {
		log.Fatalln(err)
	}

	d := &smb2.Dialer{Initiator: &smb2.NTLMInitiator{
		User:     shr.unc.User,
		Password: shr.unc.Pass,
	}}

	shr.se, err = d.Dial(shr.cn)

	if err != nil {
		log.Fatalln(err)
	}

	shr.fs, err = shr.se.Mount(shr.unc.Root)

	if err != nil {
		log.Fatalln(err)
	}
}

func (shr *Share) Umount() {
	_ = shr.fs.Umount()
	_ = shr.se.Logoff()
	_ = shr.cn.Close()
}

func parse(path string) *unc {
	path = strings.TrimSpace(path)
	path = strings.TrimPrefix(path, "smb:")

	if strings.HasPrefix(path, `\\`) {
		path = strings.ReplaceAll(path, `\`, `/`)
	}

	match := re.FindStringSubmatch(path)
	group := make(map[string]string, 6)

	for i, k := range re.SubexpNames() {
		if i > 0 && len(k) > 0 {
			group[k] = match[i]
		}
	}

	unc := new(unc)
	unc.User, _ = group["user"]
	unc.Pass, _ = group["pass"]
	unc.Host, _ = group["host"]
	unc.Root, _ = group["root"]
	unc.Port, _ = group["port"]
	unc.Path, _ = group["path"]
	unc.Part, _ = group["part"]

	if len(unc.Port) == 0 {
		unc.Port = "445"
	}

	return unc
}
