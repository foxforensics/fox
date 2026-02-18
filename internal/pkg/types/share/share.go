package share

import (
	"fmt"
	"io/fs"
	"log"
	"net"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"syscall"

	"github.com/hirochachacha/go-smb2"
	"golang.org/x/crypto/ssh/terminal"
)

var re = regexp.MustCompile(`^//((.+?)(:(.*?))?@)?(.+?)?/(.+?)/(.*)$`)

type unc struct {
	user string
	pass string
	host string
	root string
	path string
}

type Share struct {
	unc *unc
	cn  net.Conn
	se  *smb2.Session
	fs  *smb2.Share
}

func New(path string) (*Share, string) {
	shr := &Share{unc: parse(path)}
	return shr, filepath.Join(shr.unc.root, shr.unc.path)
}

func (unc *unc) String() string {
	return fmt.Sprintf("//%s/%s/", unc.host, unc.root)
}

func (shr *Share) String() string {
	return shr.unc.String()
}

func (shr *Share) DirFS(path string) fs.FS {
	return shr.fs.DirFS(path)
}

func (shr *Share) Open(path string) (*smb2.File, error) {
	return shr.fs.OpenFile(path, os.O_RDONLY, 0400)
}

func (shr *Share) Mount() {
	h, p, err := net.SplitHostPort(shr.unc.host)

	if err != nil {
		if strings.Contains(err.Error(), "missing port") {
			h = shr.unc.host
			p = "445" // default SMB2/3 port
		} else {
			log.Fatalln(err)
		}
	}

	shr.cn, err = net.Dial("tcp", fmt.Sprintf("%s:%s", h, p))

	if err != nil {
		log.Fatalln(err)
	}

	if shr.unc.pass == "*" {
		print("Password: ")

		b, err := terminal.ReadPassword(syscall.Stdin)

		println("")

		if err != nil {
			log.Fatalln(err)
		}

		shr.unc.pass = string(b)
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
	// trim common smb prefix
	path = strings.TrimSpace(path)
	path = strings.TrimPrefix(path, "smb:")

	// convert all windows backslashes
	if strings.HasPrefix(path, `\\`) {
		path = strings.ReplaceAll(path, `\`, `/`)
	}

	// append trailing slash to share
	if strings.Count(path, `/`) < 4 {
		path += "/"
	}

	// parse unc parts
	group := re.FindStringSubmatch(path)

	return &unc{
		user: group[2],
		pass: group[4],
		host: group[5],
		root: group[6],
		path: group[7],
	}
}
