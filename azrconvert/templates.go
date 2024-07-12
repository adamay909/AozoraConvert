package azrconvert

import (
	_ "embed" //for embedding template files
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

//go:embed resources/toc.xml
var tocxml string

func tocTemplate() *template.Template {
	return template.Must(template.New("toc").Parse(tocxml))
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
