package azrconvert

import (
	_ "embed" //for embedding template files
	"strings"
	"text/template"
)

//go:embed resources/webpage.html
var webpagetemplate string

func webpageTemplate() *template.Template {
	return template.Must(template.New("html").Parse(webpagetemplate))
}

//go:embed resources/oebhtmltemplate.html
var oebhtmltemplate string

func oebHTMLTemplate() *template.Template {
	return template.Must(template.New("oeb").Parse(oebhtmltemplate))
}

//go:embed resources/titlepage.html
var titlepage string

func oebTitleTemplate() *template.Template {
	return template.Must(template.New("oeb").Parse(titlepage))
}

//go:embed resources/metainf.xml
var metainfxml string

func oebMetaInf() *template.Template {
	return template.Must(template.New("oeb").Parse(metainfxml))
}

//go:embed resources/contentopf.xml
var contentopfxml string

func contentopfTemplate() *template.Template {
	return template.Must(template.New("opf").Parse(contentopfxml))
}

//go:embed resources/toc.ncx
var tocncxt string

func tocTemplate() *template.Template {
	return template.Must(template.New("toc").Parse(tocncxt))
}

//go:embed resources/toc.html
var tocxml string

func tocep3Template() *template.Template {
	return template.Must(template.New("tocep3").Parse(tocxml))
}

//go:embed resources/vertical.css
var verticalcss []byte

func verticalCSS() []byte {
	return verticalcss
}

//go:embed resources/aozora.css
var aozoracss []byte

func aozoraCSS() []byte {
	return aozoracss
}

func inlineCSSTemplate() *template.Template {

	t := webpagetemplate

	t = strings.ReplaceAll(t, `<link rel="stylesheet" type="text/css" href="vertical.css"/>`, `<style>`+string(verticalCSS())+`</style>`)

	t = strings.ReplaceAll(t, `<link rel="stylesheet" type="text/css" href="aozora.css"/>`, `<style>`+string(aozoraCSS())+`</style>`)

	return template.Must(template.New("monolithichtml").Parse(t))

}
