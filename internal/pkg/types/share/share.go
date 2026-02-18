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

var re = regexp.MustCompile(`^//((.+?)(:(.*?))?@)?(.+?)?(:(\d+?))?/(.+)/`)

type unc struct {
	user string
	pass string
	host string
	port string
	root string
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
	var cred = unc.user
	var host = unc.host

	if len(cred) > 0 {
		cred += "@"
	}

	if len(unc.port) > 0 {
		host += ":" + unc.port
	}

	return fmt.Sprintf("//%s%s/%s/", cred, host, unc.root)
}

func (shr *Share) String() string {
	return fmt.Sprintf("//%s/%s/", shr.unc.host, shr.unc.root)
}

func (shr *Share) DirFS(path string) fs.FS {
	return shr.fs.DirFS(path)
}

func (shr *Share) Open(path string) (*smb2.File, error) {
	return shr.fs.OpenFile(path, os.O_RDONLY, 0400)
}

func (shr *Share) Mount() {
	var err error

	addr := fmt.Sprintf("%s:%s", shr.unc.host, shr.unc.port)

	shr.cn, err = net.Dial("tcp", addr)

	if err != nil {
		log.Fatalln(err)
	}

	d := &smb2.Dialer{Initiator: &smb2.NTLMInitiator{
		User:     shr.unc.user,
		Password: shr.unc.pass,
	}}

	shr.se, err = d.Dial(shr.cn)

	if err != nil {
		log.Fatalln(err)
	}

	shr.fs, err = shr.se.Mount(shr.unc.root)

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

	// TODO: does not work with sub dirs
	group := re.FindStringSubmatch(path)

	if len(group[7]) == 0 {
		group[7] = "445"
	}

	return &unc{
		user: group[2],
		pass: group[4],
		host: group[5],
		port: group[7],
		root: group[8],
	}
}
