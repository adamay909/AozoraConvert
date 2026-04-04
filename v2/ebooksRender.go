package aozoraconvert

import (
	"archive/zip"
	"bytes"
	_ "embed" //embed
	"hash/crc32"
	"io"
	"math/rand"
	"mime"
	"path/filepath"
	"strings"
	"time"

	"github.com/adamay909/AozoraConvert/v2/mobi"
	"github.com/adamay909/AozoraConvert/v2/mobi/records"
	"golang.org/x/text/language"
)

type fileData struct {
	Location string
	ID       string
	Name     string
	Data     []byte
	Mtype    string
	CSS      string
}

// RenderEpub returns b as a zipped Epub
// file.
func (b *Book) RenderEpub() []byte {

	oCompatible = true

	defer func() {

		oCompatible = false

	}()

	buf := new(bytes.Buffer)

	w := zip.NewWriter(buf)

	//set mod time
	b.DateMod = time.Now().Format(time.DateOnly) + "T00:00:00Z"

	//write mimetype file
	fh := new(zip.FileHeader)
	fh.Name = "mimetype"
	fh.Method = 0
	mt := []byte("application/epub+zip")
	fh.UncompressedSize64 = uint64(len(mt))
	fh.CompressedSize64 = uint64(len(mt))
	fh.CRC32 = crc32.ChecksumIEEE(mt)
	iw, err := w.CreateRaw(fh)
	if err != nil {
		msglog.Println(err)
	}
	iw.Write(mt)

	//write META-INF
	f, err := w.Create("META-INF/container.xml")
	_, err = f.Write(metainfxml)
	if err != nil {
		msglog.Println(err)
	}
	//write title page

	f, err = w.Create("OEBPF/title.html")
	_, err = f.Write(oebtitle(b))
	if err != nil {
		msglog.Println(err)
	}

	//write main file
	f, err = w.Create("OEBPF/1.html")
	_, err = f.Write(oebmain(b))
	if err != nil {
		msglog.Println(err)
	}

	//write opf
	f, err = w.Create("OEBPF/content.opf")
	_, err = f.Write(contentopf(b))
	if err != nil {
		msglog.Println(err)
	}

	//write Epub3 toc
	f, err = w.Create("OEBPF/toc.html")
	_, err = f.Write(tocep3(b))
	if err != nil {
		msglog.Println(err)
	}

	//write support files
	for _, file := range b.Files {
		f, err = w.Create("OEBPF/" + file.Name)
		if err != nil {
			msglog.Println(err)
		}
		_, err = f.Write(file.Data)
	}

	err = w.Close()
	if err != nil {
		msglog.Println(err)
	}
	return buf.Bytes()
}

// RenderAZW3 returrns b as an AZW3 file
func (b *Book) RenderAZW3() []byte {

	mb := mobi.Book{
		Title:       b.Title,
		Authors:     []string{b.Creator},
		Publisher:   b.Publisher,
		DocType:     "EBOK",
		Language:    language.Japanese,
		FixedLayout: false,
		Vertical:    true,
		RightToLeft: true,
		UniqueID:    rand.Uint32(),
		CSSFlows:    []string{AozoraCSS},
		//	CoverImage:  b.CoverImage,
		//	ThumbImage:  b.CoverImage,
		Images: b.Images,
	}

	w := new(strings.Builder)

	for _, sec := range linearizeNode(b.Body) {

		w.Reset()

		if sec.Attr["type"] == "metadata" {

			renderHTMLForAzw3(sec, w)

			mb.Chapters = append(mb.Chapters, mobi.Chapter{
				Title:  sec.getSectionTitle(),
				Chunks: mobi.Chunks(w.String()),
			})

			continue
		}

		if sec.Attr["type"] != "section" {
			continue
		}

		if sec.sectionLevel() != 1 {
			continue
		}

		renderHTMLForAzw3(sec, w)

		mb.Chapters = append(mb.Chapters, mobi.Chapter{
			Title:  sec.getSectionTitle(),
			Chunks: mobi.Chunks(w.String()),
		})

	}

	if len(mb.Chapters) == 1 {

		w.Reset()

		for _, e := range b.Body.Children() {

			if e.Attr["type"] != "main text" {
				continue
			}

			renderHTMLForAzw3(e, w)

			mb.Chapters = append(mb.Chapters, mobi.Chapter{
				Title:  b.Body.getTitle(),
				Chunks: mobi.Chunks(w.String()),
			})

			break
		}
	}

	w.Reset()

	for _, e := range b.Body.Children() {

		if e.Attr["type"] != "bibliographical info" {
			continue
		}

		renderHTMLForAzw3(e, w)

		mb.Chapters = append(mb.Chapters, mobi.Chapter{
			Title:  "この青空文庫テキストについて",
			Chunks: mobi.Chunks(w.String()),
		})
		break

	}

	buf := new(bytes.Buffer)
	err := mb.Realize().Write(buf)
	if err != nil {
		msglog.Println(err)
	}
	return buf.Bytes()

}

// RenderPackage returns a zip package for the given format
func (b *Book) RenderPackage(format string) []byte {

	var renderer func(*Node, *strings.Builder) error

	switch format {

	case "tex":
		renderer = RenderLaTeXFull

	case "txt":
		renderer = RenderAozoraText

	case "html":
		renderer = RenderHTMLFull

	case "json":
		renderer = RenderJSON

	default:
		format = "txt"
		renderer = RenderAozoraText

	}

	buf := new(bytes.Buffer)

	zw := zip.NewWriter(buf)

	w := new(strings.Builder)

	renderer(b.Body, w)

	fname := b.TxtFileName

	if fname == "" {
		fname = "1"
	}

	f, _ := zw.Create(fname + "." + format)

	f.Write([]byte(w.String()))

	for _, file := range b.Files {

		if filepath.Ext(file.Name) == ".css" {
			continue
		}

		f, _ = zw.Create(file.Name)

		f.Write(file.Data)

	}

	if format == "tex" {

		f, _ = zw.Create("azcommands.tex")

		f.Write([]byte(LaTeXdefinitions))
	}

	if format == "html" {

		f, _ = zw.Create("aozora.css")

		f.Write([]byte(AozoraCSS))
	}

	zw.Close()

	return buf.Bytes()

}

// RenderMonolithicHTML returns the book with the
// images embedded into the HTML.
func (b *Book) RenderMonolithicHTML() []byte {

	b.EmbedImages()

	w := new(strings.Builder)

	renderHTMLMonolithic(b.Body, w)

	return []byte(w.String())

}

//go:embed resources/oebhtmltemplate.html
var oebhtmltemplate string

func oebmain(b *Book) []byte {

	resp := strings.ReplaceAll(oebhtmltemplate, `{{.Title}}`, b.Title)

	resp = strings.ReplaceAll(resp, `{{.Creator}}`, b.Creator)

	resp = strings.ReplaceAll(resp, `{{.Publisher}}`, b.Publisher)

	w := new(strings.Builder)

	for _, e := range b.Body.Children() {

		if e.Attr["type"] == "main text" {

			renderEpubHTML(e, w)

		}
		if e.Attr["type"] == "bibliographical info" {

			err := RenderHTML(e, w)

			if err != nil {
				return []byte{}
			}
		}
	}

	resp = strings.ReplaceAll(resp, `{{.RenderBody}}`, w.String())

	return []byte(resp)

}

//go:embed resources/titlepage.html
var titlepage string

func oebtitle(b *Book) []byte {

	resp := strings.ReplaceAll(titlepage, `{{.Creator}}`, b.Creator)

	resp = strings.ReplaceAll(resp, `{{.Title}}`, b.Title)

	w := new(strings.Builder)

	for _, e := range b.Body.Children() {

		if e.Attr["type"] == "metadata" {

			renderEpubHTML(e, w)

			break
		}
	}

	resp = strings.ReplaceAll(resp, `{{.Metadata}}`, w.String())

	return []byte(resp)

}

//go:embed resources/metainf.xml
var metainfxml []byte

func metainf(b *Book) []byte {

	return metainfxml
}

//go:embed resources/contentopf.xml
var contentopfstr string

func contentopf(b *Book) []byte {

	w := new(strings.Builder)

	for _, e := range b.Files {

		addToStringsBuilder(w, `<item id="`, e.ID, `" href="`, e.Name, `" media-type="`, e.Mtype, `" />`)

	}

	resp := strings.ReplaceAll(contentopfstr, "{{.FileList}}", w.String())

	resp = strings.ReplaceAll(resp, "{{.Creator}}", b.Creator)

	resp = strings.ReplaceAll(resp, "{{.Title}}", b.Title)

	resp = strings.ReplaceAll(resp, "{{.Publisher}}", b.Publisher)

	resp = strings.ReplaceAll(resp, "{{.UUID}}", b.UUID)

	resp = strings.ReplaceAll(resp, "{{.DateMod}}", b.DateMod)

	return []byte(resp)

}

//go:embed resources/toc.html
var epubtoc string

func tocep3(b *Book) []byte {

	w := new(strings.Builder)

	renderNavHTML(b.Body, w)

	resp := strings.ReplaceAll(epubtoc, `{{.RenderEP3TOC}}`, w.String())

	resp = strings.ReplaceAll(resp, "#", "1.html#")

	resp = strings.ReplaceAll(resp, "{{.Creator}}", b.Creator)

	resp = strings.ReplaceAll(resp, "{{.Title}}", b.Title)

	return []byte(resp)

}

func (b *Book) addFilesFromZip(arch *zip.Reader) {

	for _, f := range arch.File {

		var fi fileData

		fi.Name = filepath.Base(f.Name)

		if fi.Name == "1.html" {
			continue
		}

		fi.Mtype = mime.TypeByExtension(filepath.Ext(fi.Name))
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

		b.Files = append(b.Files, fi)

		if fi.Mtype == "image/png" || fi.Mtype == "image/jpeg" {
			b.Images = append(b.Images, records.ImageRecord{Data: fi.Data, Ext: filepath.Ext(fi.Name)})
		}
	}
}
