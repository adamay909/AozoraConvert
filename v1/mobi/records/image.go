package records

import (
	"image"
	"image/png"
	"io"
	"strings"

	"github.com/adamay909/AozoraConvert/mobi/jfif"
)

type ImageRecord struct {
	Img  image.Image
	Ext  string
	Data []byte
}

/*
	func NewImageRecord(img image.Image) ImageRecord {
		return ImageRecord{
			Img: img,
		}
	}
*/
func (r ImageRecord) Write(w io.Writer) error {
	if len(r.Data) != 0 {
		_, err := w.Write(r.Data)
		return err
	}
	if strings.ToLower(r.Ext) == ".png" {
		return png.Encode(w, r.Img)
	}
	return jfif.Encode(w, r.Img, nil)
}
