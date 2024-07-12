package azrconvert

import (
	_ "embed" //for embedding font file.
	"image/color"
	"log"
	"math/rand"
	"strings"
	"time"

	"github.com/adamay909/AozoraConvert/drawtext"
	"golang.org/x/image/font/opentype"
)

//go:embed resources/jpfont.otf
var fontdata []byte

func (b *Book) genTitlePage() {

	bg, fg := colorpair()
	xsize := 1200
	ysize := 1600

	canvas := drawtext.NewCanvas(bg, xsize, ysize)

	normalfont := drawtext.MakeNewFontFace(getFont(), 18)
	smallfont := drawtext.MakeNewFontFace(getFont(), 16)
	largefont := drawtext.MakeNewFontFace(getFont(), 28)

	drawtext.SetCursor(canvas, 0, 0)

	drawtext.WriteLine(b.Creator, "left", normalfont, fg, canvas)

	l, _ := drawtext.BreakLines(b.Title, largefont, canvas)

	drawtext.SetCursor(canvas, 0, 0.3-float64((len(l)-1)*3/40))

	titleLines := breakTitle(b.Title)

	for _, t := range titleLines {
		drawtext.AddCenteredText(t, largefont, fg, canvas)
	}
	drawtext.SetCursor(canvas, 0, 0.9)

	drawtext.WriteLine(b.Publisher, "right", smallfont, fg, canvas)

	b.CoverImage = drawtext.ImageOf(canvas)

	//	return drawtext.RenderCanvas(canvas, "png")
}

func breakTitle(in string) (out []string) {
	out = strings.Split(in, "─")
	if len(out) == 2 {
		out[1] = strings.TrimSuffix(out[1], "─")
	}
	return
}

func colorpair() (bg, fg color.Color) {

	rand.Seed(time.Now().UnixNano() / 100019)

	r := rand.Intn(4)

	switch r {
	case 0:
		bg = color.RGBA{R: uint8(196), G: uint8(216), B: uint8(196), A: uint8(255)}
		fg = color.RGBA{R: uint8(18), G: uint8(66), B: uint8(18), A: uint8(255)}

	case 1:
		bg = color.RGBA{R: uint8(246), G: uint8(234), B: uint8(245), A: uint8(255)}
		fg = color.RGBA{R: uint8(79), G: uint8(26), B: uint8(74), A: uint8(255)}

	case 2:
		bg = color.RGBA{R: uint8(201), G: uint8(208), B: uint8(242), A: uint8(255)}
		fg = color.RGBA{R: uint8(40), G: uint8(49), B: uint8(92), A: uint8(255)}

	default:
		bg = color.RGBA{R: uint8(245), G: uint8(236), B: uint8(220), A: uint8(255)}
		fg = color.RGBA{R: uint8(73), G: uint8(59), B: uint8(34), A: uint8(255)}
	}

	return
}

func getFont() *opentype.Font {
	fontd, err := opentype.Parse(fontdata)
	if err != nil {
		log.Println(err)
	}
	return fontd
}
