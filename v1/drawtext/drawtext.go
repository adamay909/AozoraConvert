package drawtext

import (
	"bytes"
	_ "embed" //for embedding default font
	"image"
	"image/color"
	"image/jpeg"
	"image/png"
	"log"
	"strings"

	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"
	"golang.org/x/image/math/fixed"
)

//go:embed resources/jpfont.otf
var fontdata []byte

// Canvas holds information about the canvas on which to
// write text.
type Canvas struct {
	drawer                                           *font.Drawer
	marginLeft, marginRight, marginTop, marginBottom float64
	maxX, maxY                                       int
	maxXPr, maxYPr                                   int
	xoff, yoff                                       int
}

// NewCanvas returns a new Canvas with margins set to 10% on each side,
// background color set to bg, and dimensions given by width and height.
func NewCanvas(bg color.Color, width, height int) *Canvas {

	canvas := new(Canvas)

	canvas.drawer = new(font.Drawer)

	canvas.drawer.Dst = newBGimage(bg, width, height)

	canvas.marginLeft = 0.1
	canvas.marginRight = 0.1
	canvas.marginTop = 0.1
	canvas.marginBottom = 0.1

	canvas.maxX = width
	canvas.maxY = height
	canvas.maxXPr = int(float64(width) * (1 - canvas.marginRight - canvas.marginLeft))
	canvas.maxYPr = int(float64(height) * (1 - canvas.marginTop - canvas.marginBottom))
	canvas.xoff = int(float64(width) * canvas.marginLeft)
	canvas.yoff = int(float64(height) * canvas.marginTop)

	canvas.drawer.Dot = fixed.P(canvas.xoff, canvas.yoff)

	return canvas
}

// ImageOf returns the image given by canvas.
func ImageOf(canvas *Canvas) image.Image {
	return canvas.drawer.Dst
}

// SetCursor sets the position of the curser to xpos,ypos.
// The value should be between 0 and 1 representing the position
// relative to the printable size (canvas size minus the margins).
// (0,0) is top left, (1,1) is bottom right. And the cursor
// position represents the top left corner of the bounding box
// of next string to be written.
func SetCursor(canvas *Canvas, xpos, ypos float64) {

	newX := int(float64(canvas.maxXPr)*xpos) + canvas.xoff
	newY := int(float64(canvas.maxYPr)*ypos) + canvas.yoff

	canvas.drawer.Dot = fixed.P(newX, newY)
}

// WriteLine writes a new line containing s onto canvas
// using font face f and color fg.
// Align is "right", "left", or "center". Any other value for align
// will place the text at current cursor location.
func WriteLine(s string, align string, f font.Face, fg color.Color, canvas *Canvas) {

	var xoffset, yoffset int

	canvas.drawer.Src = image.NewUniform(fg)

	canvas.drawer.Face = f

	bounds, _ := canvas.drawer.BoundString(s)

	xmin := bounds.Min.X.Round()
	xmax := bounds.Max.X.Round()
	textwidth := xmax - xmin

	ymin := bounds.Min.Y.Round()
	ymax := bounds.Max.Y.Round()
	textheight := ymax - ymin

	yoffset = int(float64(textheight) * 1.35)
	x0 := canvas.drawer.Dot.X.Round()
	canvas.drawer.Dot = fixed.P(0, canvas.drawer.Dot.Y.Round()+yoffset)

	dot := canvas.drawer.Dot

	switch align {
	case "center":
		xoffset = (canvas.maxX - textwidth) / 2
		dot.X = fixed.I(xoffset)

	case "left":
		xoffset = canvas.xoff
		dot.X = fixed.I(xoffset)

	case "right":
		xoffset = canvas.xoff + canvas.maxXPr - textwidth
		dot.X = fixed.I(xoffset)
	default:
		xoffset = x0
		dot.X = fixed.I(xoffset)
	}

	canvas.drawer.Dot = dot

	canvas.drawer.DrawString(s)

	return
}

// RenderCanvas encodes the image of canvas.
// filetype is given as the extension ("png", "jpg").
// Default filetype is "png".
func RenderCanvas(canvas *Canvas, filetype string) []byte {

	filetype = strings.TrimPrefix(filetype, ".")

	w := new(bytes.Buffer)

	switch filetype {

	case "png":
		png.Encode(w, canvas.drawer.Dst)

	case "jpg":
		jpeg.Encode(w, canvas.drawer.Dst, &jpeg.Options{Quality: 94})

	case "jpeg":
		jpeg.Encode(w, canvas.drawer.Dst, &jpeg.Options{Quality: 94})

	default:
		png.Encode(w, canvas.drawer.Dst)
	}

	return w.Bytes()
}

func newBGimage(bg color.Color, xsize, ysize int) *image.RGBA {
	dstimg := image.NewRGBA(image.Rect(0, 0, xsize, ysize))
	for x := 0; x < xsize; x++ {
		for y := 0; y < ysize; y++ {
			dstimg.Set(x, y, bg)
		}
	}
	return dstimg
}

// DrawGrid draws a grid on canvas. n gives the number
// of grid lines.
func DrawGrid(canvas *Canvas, n int) {
	dstimg := canvas.drawer.Dst

	black := color.Black
	for x := 0; x < canvas.maxX; x = x + canvas.maxX/n {
		for y := 0; y < canvas.maxY; y++ {
			dstimg.Set(x, y, black)
		}
	}

	for y := 0; y < canvas.maxY; y = y + canvas.maxY/n {
		for x := 0; x < canvas.maxX; x++ {
			dstimg.Set(x, y, black)
		}
	}
	return
}

// MakeNewFontFace makes a new font face for writing.
func MakeNewFontFace(f *opentype.Font, size int) font.Face {
	face, err := opentype.NewFace(f, &opentype.FaceOptions{Size: float64(size), DPI: 300})
	if err != nil {
		log.Println(err)
	}
	return face
}

// AddCenteredText adds a centered text.
func AddCenteredText(s string, f font.Face, fg color.Color, canvas *Canvas) {

	canvas.drawer.Face = f

	lines, _ := BreakLines(s, f, canvas)

	for i := 0; i < len(lines); i++ {
		WriteLine(lines[i], "center", f, fg, canvas)
	}
	return
}

// BreakLines breaks s into one or more lines to make sure each
// part fits within the margins of canvas. It also returns the
// expected maximal line height.
func BreakLines(s string, f font.Face, canvas *Canvas) (lines []string, lineHeight float64) {

	canvas.drawer.Face = f
	bounds, _ := canvas.drawer.BoundString(s)

	xmin := bounds.Min.X.Round()
	xmax := bounds.Max.X.Round()
	textwidth := xmax - xmin

	ymin := bounds.Min.Y.Round()
	ymax := bounds.Max.Y.Round()
	lineHeight = float64(ymax-ymin) / float64(canvas.maxYPr)

	divisor := textwidth/canvas.maxXPr + 1

	sRunes := []rune(s)
	strlen := len(sRunes) / divisor

	var line []rune

	for i := 0; i < len(sRunes); i++ {
		line = append(line, sRunes[i])
		if len(line) == strlen {
			lines = append(lines, string(line))
			line = nil
		}
	}
	lines = append(lines, string(line))
	//	for ok := kinsoku(lines); !ok; kinsoku(lines) {
	//	}
	return lines, lineHeight
}

// ensure compliance with Japanese line breaking rules.
// returns ok=true if no changes were needed.
func kinsoku(lines []string) (ok bool) {

	ok = true

	for i := range lines {
		r := []rune(lines[i])
		last := len(r) - 1
		if last < 1 {
			continue
		}
		if strings.ContainsAny(string(r[0]), ",)]｝、〕〉》」』】〙〗〟’”｠»ゝゞーァィゥェォッャュョヮヵヶぁぃぅぇぉっゃゅょゎゕゖㇰㇱㇲㇳㇴㇵㇶㇷㇸㇹㇷ゚ㇺㇻㇼㇽㇾㇿ々〻‐゠–〜～?!‼⁇⁈⁉・:;/。.") {
			if i == 0 {
				continue
			}
			lines[i-1] = lines[i-1] + string(r[0])
			lines[i] = string(r[1:])
			ok = false
		}
		if strings.ContainsAny(string(r[last]), "([｛〔〈《「『【〘〖〝‘“｟«") {
			if i == len(lines)-1 {
				continue
			}
			ok = false
			lines[i+1] = string(r[last]) + lines[i+1]
			lines[i] = string(r[:last-1])
		}
	}

	return
}
