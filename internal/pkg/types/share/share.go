package share

import (
	"fmt"
	"io/fs"
	"log"
	"net"
	"strings"

	"github.com/hirochachacha/go-smb2"
)

type Share struct {
	cn net.Conn
	se *smb2.Session
	fs *smb2.Share

	user string
	pass string
	host string
	path string
}

func New(addr string) *Share {
	shr := &Share{}

	// remove protocol
	addr = strings.TrimPrefix(addr, "smb:")
	addr = strings.TrimPrefix(addr, "//")

	host := strings.SplitN(addr, "@", 1)

	// parse credentials
	if len(host) > 1 {
		cred := strings.SplitN(host[0], ":", 1)

		// set username
		shr.user = cred[0]

		// set password
		if len(cred) > 1 {
			shr.pass = cred[1]
		}
	}

	// parse hostname
	path := strings.SplitN(host[len(host)-1], "/", 1)

	// set hostname
	shr.host = path[0]

	// set path
	if len(path) > 1 {
		shr.path = path[1]
	}

	return shr
}

func (shr *Share) String() string {
	pw := shr.pass

	if len(pw) > 0 {
		pw = "********"
	}

	return fmt.Sprintf("smb://%s:%s@%s/%s",
		shr.user,
		pw,
		shr.host,
		shr.path,
	)
}

func (shr *Share) DirFS(path string) fs.FS {
	return shr.fs.DirFS(path)
}

func (shr *Share) Mount() {
	conn, err := net.Dial("tcp", shr.host)

	if err != nil {
		log.Fatalln(err)
	}

	d := &smb2.Dialer{
		Initiator: &smb2.NTLMInitiator{
			User:     shr.user,
			Password: shr.pass,
		},
	}

	shr.se, err = d.Dial(conn)

	if err != nil {
		log.Fatalln(err)
	}

	shr.fs, err = shr.se.Mount(shr.path)

	if err != nil {
		log.Fatalln(err)
	}
}

func (shr *Share) Umount() {
	_ = shr.fs.Umount()
	_ = shr.se.Logoff()
	_ = shr.cn.Close()
}
