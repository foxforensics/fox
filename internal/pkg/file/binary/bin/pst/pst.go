package pst

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"log"

	set "github.com/emersion/go-message/charset"

	"github.com/rotisserie/eris"
	"go.foxforensics.eu/go-pst/v6/pst"
	"golang.org/x/text/encoding"

	"go.foxforensics.eu/fox/v4/internal/pkg/file"
)

type Folder struct {
	Name     string    `json:"name,omitempty"`
	Messages []Message `json:"messages,omitempty"`
}

type Message struct {
	Content     any          `json:"content,omitempty"`
	Attachments []Attachment `json:"attachments,omitempty"`
}

type Attachment struct {
	Filename string `json:"filename,omitempty"`
	MimeType string `json:"mime_type,omitempty"`
	IsHidden bool   `json:"is_hidden,omitempty"`
	Size     int32  `json:"size,omitempty"`
	Data     string `json:"data,omitempty"`
}

func Detect(b []byte) bool {
	return file.HasMagic(b, 0, []byte{
		'!', 'B', 'D', 'N',
	})
}

func Convert(b []byte) ([]byte, error) {
	pst.ExtendCharsets(func(name string, enc encoding.Encoding) {
		set.RegisterEncoding(name, enc)
	})

	pf, err := pst.New(bytes.NewReader(b))

	if err != nil {
		return b, err
	}

	defer func() {
		pf.Cleanup()
	}()

	v, err := getFolders(pf)

	if err != nil {
		return b, err
	}

	return json.Marshal(v)
}

func getFolders(pf *pst.File) ([]Folder, error) {
	var res = make([]Folder, 0)
	var err error

	err = pf.WalkFolders(func(v *pst.Folder) error {
		f := &Folder{Name: v.Name}

		if f.Messages, err = getMessages(v); err != nil {
			log.Printf("warning: %s!\n", err)
		}

		res = append(res, *f)

		return nil
	})

	return res, err
}

func getMessages(f *pst.Folder) ([]Message, error) {
	var res = make([]Message, 0)
	var err error

	it, err := f.GetMessageIterator()

	if eris.Is(err, pst.ErrMessagesNotFound) {
		return res, nil // no messages
	} else if err != nil {
		return nil, err
	}

	for it.Next() {
		v := it.Value()
		m := &Message{Content: v.Properties}

		if m.Attachments, err = getAttachments(v); err != nil {
			log.Printf("warning: %s!\n", err)
		}

		res = append(res, *m)
	}

	return res, nil
}

func getAttachments(m *pst.Message) ([]Attachment, error) {
	var res = make([]Attachment, 0)
	var err error

	it, err := m.GetAttachmentIterator()

	if eris.Is(err, pst.ErrAttachmentsNotFound) {
		return res, nil // no attachments
	} else if err != nil {
		return nil, err
	}

	for it.Next() {
		v := it.Value()
		a := &Attachment{
			Filename: v.GetAttachLongFilename(),
			MimeType: v.GetAttachMimeTag(),
			IsHidden: v.GetAttachmentHidden(),
			Size:     v.GetAttachSize(),
		}

		buf := bytes.NewBuffer(nil)

		if _, err = v.WriteTo(buf); err == nil {
			a.Data = base64.StdEncoding.EncodeToString(buf.Bytes())
		} else {
			log.Printf("warning: %s!\n", err)
		}

		res = append(res, *a)
	}

	return res, nil
}
