package azrconvert

import (
	"archive/zip"
	"bytes"
	"image"
	"image/jpeg"
	"image/png"
	"io"
	"log"
	"math/rand"
	"mime"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/adamay909/AozoraConvert/mobi"
	"github.com/adamay909/AozoraConvert/mobi/records"
	"golang.org/x/net/html"
	"golang.org/x/text/language"
)

type fileData struct {
	Location string
	ID       string
	Name     string
	Data     []byte
	Mtype    string
}

// RenderWebpage returns b as a single web page.
// It will be  UTF-8 encoded and css will be
// for vertical reading.
func (b *Book) RenderWebpage() []byte {

	builder := new(strings.Builder)

	err := webpageTemplate().Execute(builder, b)
	if err != nil {
		log.Println(err)
	}
	return []byte(builder.String())
}

// RenderWebpagePackage returns a zip archive containing
// all the files necessary (html, css, graphics files) to
// show the page correctly in a web browser.
func (b *Book) RenderWebpagePackage() []byte {

	buf := new(bytes.Buffer)

	w := zip.NewWriter(buf)

	f, err := w.Create("1.html")
	_, err = f.Write(b.RenderWebpage())

	//write support files
	for _, file := range b.Files {
		f, err = w.Create(file.Name)
		if err != nil {
			log.Println(err)
		}
		_, err = f.Write(file.Data)
	}

	err = w.Close()
	if err != nil {
		log.Println(err)
	}
	return buf.Bytes()
}

// RenderEpub returns b as a zipped Epub
// file.
func (b *Book) RenderEpub() []byte {

	buf := new(bytes.Buffer)

	w := zip.NewWriter(buf)

	//write mimetype file
	f, err := w.Create("mimetype")
	_, err = f.Write([]byte("application/epub+zip"))
	if err != nil {
		log.Println(err)
	}
	//write META-INF
	f, err = w.Create("META-INF/container.xml")
	_, err = f.Write(metainf(b))
	if err != nil {
		log.Println(err)
	}
	// write title image
	f, err = w.Create("cover.png")
	b.genTitlePage()
	img := new(bytes.Buffer)
	_ = png.Encode(img, b.CoverImage)
	_, err = f.Write(img.Bytes())
	if err != nil {
		log.Println(err)
	}
	//write title page

	f, err = w.Create("title.html")
	_, err = f.Write(oebtitle(b))
	if err != nil {
		log.Println(err)
	}

	//write main file
	f, err = w.Create("1.html")
	_, err = f.Write(oebmain(b))
	if err != nil {
		log.Println(err)
	}

	//write opf
	f, err = w.Create("content.opf")
	_, err = f.Write(contentopf(b))
	if err != nil {
		log.Println(err)
	}

	//write toc
	f, err = w.Create("toc.ncx")
	_, err = f.Write(toc(b))
	if err != nil {
		log.Println(err)
	}

	//write support files
	for _, file := range b.Files {
		f, err = w.Create(file.Name)
		if err != nil {
			log.Println(err)
		}
		_, err = f.Write(file.Data)
	}

	err = w.Close()
	if err != nil {
		log.Println(err)
	}
	return buf.Bytes()
}

// RenderAZW3 returrns b as an AZW3 file
func (b *Book) RenderAZW3() []byte {
	var text string

	mb := mobi.Book{
		Title:       b.Title,
		Authors:     []string{b.Creator},
		Publisher:   b.Publisher,
		DocType:     "EBOK",
		Language:    language.Japanese,
		FixedLayout: false,
		RightToLeft: true,
		UniqueID:    rand.Uint32(),
		CSSFlows:    []string{b.CSS + string(verticalCSS()), string(aozoraCSS())},
		CoverImage:  b.CoverImage,
		ThumbImage:  b.CoverImage,
		Images:      b.Images,
	}

	//fix image links
	imc := 0
	for _, t := range b.Body {
		if isImg(t) {
			filename := getAttr(t, "src")
			imc++
			p := strings.ToUpper(strconv.FormatInt(int64(imc), 16))
			for len(p) < 4 {
				p = "0" + p
			}
			p = "kindle:embed:" + p + "?mime=" + mime.TypeByExtension(filepath.Ext(filename))

			delAttr(t, "src")
			setAttr(t, "src", p)
			alt := getAttr(t, "alt")
			delAttr(t, "alt")
			setAttr(t, "alt", alt)

		}
	}

	text = ""
	if b.TopSection != nil {

		var ch *section
		st := 0
		for ch = b.TopSection; ch.nextSibling != nil; ch = ch.nextSibling {
			chap := ""
			if ch == b.TopSection {
				text = renderTokens(b.Body[ch.start+1 : ch.nextSibling.start])
			} else {
				chap = renderTokens(b.Body[ch.start:ch.nextSibling.start])
				text = text + chap
			}

			mb.Chapters = append(mb.Chapters, mobi.Chapter{
				Title:  ch.content,
				Start:  st,
				Length: len(chap)})

			st = len(text)
		}

		chap := renderTokens(b.Body[ch.start:])
		mb.Chapters = append(mb.Chapters, mobi.Chapter{
			Title:  ch.content,
			Start:  st,
			Length: len(chap),
		})
		text = text + chap

	} else {
		text = renderTokens(b.Body[1 : len(b.Body)-1])

		mb.Chapters = append(mb.Chapters, mobi.Chapter{
			Title:  b.Title,
			Start:  0,
			Length: len(text)})

	}

	mb.Html = text + "\n</body>\n</html>"

	buf := new(bytes.Buffer)
	err := mb.Realize().Write(buf)
	if err != nil {
		log.Println(err)
	}
	return buf.Bytes()

}

// AddFiles adds the CSS style files as well as
// all the image files requested by the book.
func (b *Book) AddFiles() {

	dir := `https://` + strings.TrimPrefix(filepath.Dir(b.URI), `https:/`) + `/`

	for _, t := range b.Body {

		if isImg(t) {
			var fi fileData
			var err error

			alt := getAttr(t, "alt")
			path := getAttr(t, "src")
			if strings.HasPrefix(path, `http`) {
				fi.Location = path
			} else {
				fi.Location = dir + path
			}

			t.Type = html.SelfClosingTagToken //need this to make sure tag self-closes

			log.Println("Adding File", fi.Location)

			ext := filepath.Ext(fi.Location)
			fi.Name = strconv.Itoa(len(b.Files) + 1)
			for len(fi.Name) < 5 {
				fi.Name = "0" + fi.Name
			}
			fi.Name = fi.Name + ext
			fi.Mtype = mime.TypeByExtension(filepath.Ext(fi.Name))
			fi.Data, err = downloadFile(fi.Location)
			if err != nil {
				log.Println("Could not add", fi.Location)
				continue
			}

			fi.ID = "image" + strconv.Itoa(len(b.Files))

			b.Files = append(b.Files, fi)

			var im image.Image

			if fi.Mtype == "image/png" {
				r := bytes.NewReader(fi.Data)
				im, err = png.Decode(r)
			}
			if fi.Mtype == "image/jpeg" {
				r := bytes.NewReader(fi.Data)
				im, err = jpeg.Decode(r)
			}
			b.Images = append(b.Images, records.ImageRecord{Data: fi.Data, Ext: ext})

			//fix css
			if err != nil {
				log.Println("Could not determine size of image", fi.Location)
				continue
			}

			width, height := getSize(im)

			b.CSS = b.CSS + "." + fi.ID + " {\nheight: " + strconv.Itoa(height) + "px;\n width: " + strconv.Itoa(width) + "px;\n}\n"

			t.Attr = nil
			setAttr(t, "class", fi.ID)
			setAttr(t, "src", fi.Name)
			setAttr(t, "alt", alt)
		}
	}

	//Make sure we add aozora.css
	var fi fileData

	fi.ID = "css"
	fi.Name = "aozora.css"
	fi.Data = aozoraCSS()
	fi.Mtype = "text/css"
	b.Files = append(b.Files, fi)

	fi.ID = "vertical_css"
	fi.Name = "vertical.css"
	fi.Data = []byte(string(verticalCSS()) + b.CSS)
	fi.Mtype = "text/css"
	b.Files = append(b.Files, fi)

	return
}

func getSize(im image.Image) (w, h int) {

	w = im.Bounds().Max.X - 1 - im.Bounds().Min.X
	h = im.Bounds().Max.Y - 1 - im.Bounds().Min.Y

	return
}

// RenderBody renders the html body of bk.
func (b *Book) RenderBody() string {
	if len(b.Body) == 0 {
		return ""
	}
	return renderTokens(b.Body)

}

func contains(e string, c []fileData) bool {
	if len(c) == 0 {
		return false
	}
	for i := range c {
		if e == c[i].Location {
			return true
		}
	}
	return false

}

func oebmain(b *Book) []byte {

	builder := new(strings.Builder)
	err := oebHTMLTemplate().Execute(builder, b)
	if err != nil {
		log.Println(err)
	}

	return []byte(builder.String())
}

func oebtitle(b *Book) []byte {

	builder := new(strings.Builder)
	err := oebTitleTemplate().Execute(builder, b)
	if err != nil {
		log.Println(err)
	}

	return []byte(builder.String())
}

func metainf(b *Book) []byte {

	builder := new(strings.Builder)
	err := oebMetaInf().Execute(builder, b)
	if err != nil {
		log.Println(err)
	}

	return []byte(builder.String())
}

func contentopf(b *Book) []byte {

	builder := new(strings.Builder)
	err := contentopfTemplate().Execute(builder, b)
	if err != nil {
		log.Println(err)
	}

	return []byte(builder.String())
}

func toc(b *Book) []byte {

	builder := new(strings.Builder)
	err := tocTemplate().Execute(builder, b)
	if err != nil {
		log.Println(err)
	}

	return []byte(builder.String())
}

func downloadFile(location string) (data []byte, err error) {
	path, _ := url.Parse(location)

	/*if path.Host == "" {

		return getLocalFile(path.Path)
	}*/
	r, err := http.Get(path.String())
	if err != nil {
		log.Println(err)
		return
	}
	log.Println("file downloaded")
	data, err = io.ReadAll(r.Body)
	if err != nil {
		log.Println(err)
		return
	}

	r.Body.Close()
	return data, err
}

func getLocalFile(path string) (data []byte) {
	data, err := os.ReadFile(path)
	if err != nil {
		log.Println(err)
	}
	return
}

func renderTokens(in []*html.Token) string {

	w := new(strings.Builder)

	for _, t := range in {
		w.WriteString(t.String())
	}

	return string(prettifyEmptyLines([]byte(w.String())))
}
