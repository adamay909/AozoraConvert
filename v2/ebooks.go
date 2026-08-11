package aozoraconvert

import (
	"archive/zip"
	"bytes"
	"encoding/base64"
	"io"
	"mime"
	"path/filepath"
	"strings"
	"unicode/utf8"

	"github.com/adamay909/AozoraConvert/v2/mobi/records"
)

// Book represents a book from Aozora Bunko
type Book struct {
	Title, Creator, Publisher string
	Files                     []fileData
	UUID                      string
	Body                      *Node
	TOC                       *Node
	URI                       string
	Images                    []records.ImageRecord
	CSS                       string
	Hash                      string
	DateMod                   string
	TxtFileName               string
}

// NewBook returns a new Book.
func NewBook() *Book {
	b := new(Book)
	//b.UUID = uuid.NewString()
	b.UUID = "2bc8df04-20ec-4fb6-a4c8-356942b083b4"
	return b
}

// NewEbookFromZip returns Book from dz which must be zip
// archive containing the Aozorabunko text and any needed graphics files.
func NewEbookFromZip(dz []byte) (bk *Book) {

	var err error

	arch, err := zip.NewReader(bytes.NewReader(dz), int64(len(dz)))
	if err != nil {
		msglog.Println(err)
		return
	}

	bk = NewBook()

	for _, f := range arch.File {

		if filepath.Ext(f.Name) == ".txt" {

			bk.TxtFileName = strings.TrimSuffix(f.Name, ".txt")

			r, _ := f.Open()

			bk.Body, err = AST(readStringFromFile(r))

			if err != nil {
				msglog.Println(err)
				return
			}

			continue
		}

		var fi fileData

		fi.Name = filepath.Base(f.Name)

		fi.Mtype = mime.TypeByExtension(filepath.Ext(fi.Name))

		fi.ID = "file" + strings.TrimSuffix(fi.Name, filepath.Ext(fi.Name))

		r, err := f.Open()

		defer r.Close()
		if err != nil {
			msglog.Println(err)
			return
		}

		fi.Data, err = io.ReadAll(r)

		if err != nil {
			msglog.Println(err)
			return
		}

		bk.Files = append(bk.Files, fi)

		if fi.Mtype == "image/png" || fi.Mtype == "image/jpeg" {
			bk.Images = append(bk.Images, records.ImageRecord{Data: fi.Data, Ext: filepath.Ext(fi.Name)})
		}
	}

	var fi fileData

	fi.ID = "css"

	fi.Name = "aozora.css"

	fi.Data = []byte(AozoraCSS)

	fi.Mtype = "text/css"

	bk.Files = append(bk.Files, fi)

	bk.UUID = "2bc8df04-20ec-4fb6-a4c8-356942b083b4"
	//	bk.UUID = uuid.NewString()

	bk.SetMetadataFromText()

	return bk
}

func readStringFromFile(f io.Reader) string {

	d, err := io.ReadAll(f)

	if err != nil {
		msglog.Println(err)
		return ""
	}

	if !utf8.Valid(d) {
		return ToUTF8(d)
	}

	return string(d)
}

// SetURI sets the path of book within
// Aozora Bunko's file structure.
func (b *Book) SetURI(l string) {
	b.URI = l
	return
}

// GetURI returns the the path of the
// book within Aozora Bunko's file structure.
func (b *Book) GetURI() string {
	return b.URI

}

// SetTitle sets the title to t.
func (b *Book) SetTitle(t string) {
	b.Title = t
	return
}

// SetCreator sets the creator to c.
func (b *Book) SetCreator(c string) {
	b.Creator = c
	return
}

// SetPublisher sets the publisher to p.
func (b *Book) SetPublisher(p string) {
	b.Publisher = p
	return
}

// SetMetadataFromText sets metadata from
// text.
func (b *Book) SetMetadataFromText() {

	w := new(strings.Builder)

	b.Publisher = "青空文庫"

	for _, e := range linearizeNode(b.Body) {

		if e.Attr["type"] == "metadata" {

			for _, n := range linearizeNode(e) {

				if n.Attr["type"] == "meta title" {

					w.Reset()

					renderInnerTextOnly(n, w)

					b.Title = w.String()

					continue

				}

				if n.Attr["type"] == "meta contributor" {

					w.Reset()

					renderInnerTextOnly(n, w)

					b.Creator = b.Creator + "、" + w.String()

				}
			}

			b.Creator = strings.TrimPrefix(b.Creator, "、")

			return
		}
	}

	b.Title = "不明"

	b.Creator = "不明"

	return
}

// EmbedImages embeds the image data in the relevant nodes in the ast of b.l
func (b *Book) EmbedImages() {

	for _, e := range linearizeNode(b.Body) {

		if e.Attr["type"] != "image" {
			continue
		}

		source := e.Attr["file"]

		//find the corresponding file

		for _, fi := range b.Files {

			if fi.Name == source {

				datastr := "data:" + fi.Mtype + ";base64,"

				data := make([]byte, base64.StdEncoding.EncodedLen(len(fi.Data)))

				base64.StdEncoding.Encode(data, fi.Data)

				datastr = datastr + string(data)

				e.SetAttr("data", datastr)

				break
			}
		}

	}

	return
}
