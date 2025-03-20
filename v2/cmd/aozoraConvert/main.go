package main

import (
	_ "embed" //for embedding
	"flag"
	"log"
	"os"
	"unicode/utf8"

	ac "github.com/adamay909/AozoraConvert/v2"
)

var inputFile string

var (
	outputFile = flag.String("o", "/dev/stdout", "name of output file; defaults to Stdout")

	sjisOut = flag.Bool("sjis", false, "generate SJIS encoded output")

	format = flag.String("format", "html", "output file format; supported types are: html, text, json; defaults to html")

	jis0208 = flag.Bool("jis0208", false, "output is JIS0208 compatible.")
)

func init() {

	flag.Parse()

	if len(flag.Args()) < 1 {

		log.Println("must provide name of input file as argument")

		return
	}

	inputFile = flag.Args()[0]

	if *sjisOut {

		*jis0208 = true

	}
}

var renderer func(*ac.Node) string

func main() {

	ac.SetJIS0208(*jis0208)

	data := readFile(inputFile)

	switch *format {

	case "text":
		renderer = ac.RenderAozoraText

	case "json":
		renderer = ac.RenderJSON

	default:
		renderer = ac.RenderHTML

	}

	converted := renderer(ac.AST(data))

	if *sjisOut {

		converted = ac.ToSJIS(converted)

	}

	writeFile(*outputFile, converted)

}

func readFile(fn string) string {

	inFile, err := os.Open(fn)

	if err != nil {

		log.Println(err)

		return ""
	}

	defer inFile.Close()

	finfo, err := inFile.Stat()

	if err != nil {

		log.Println(err)

		return ""
	}

	data := make([]byte, finfo.Size())

	inFile.Read(data)

	if !utf8.Valid(data) {
		log.Println("Not UTF8. Assuming ShiftJIS.")
		data = []byte(ac.ToUTF8(data))
	}

	return string(data)

}

func writeFile(fn string, text string) {

	var err error

	outFile := new(os.File)

	if fn == "/dev/stdout" {

		outFile = os.Stdout

	} else {

		outFile, err = os.Create(fn)

	}

	if err != nil {

		log.Println(err)

		return
	}

	defer outFile.Close()

	outFile.WriteString(text)

	return

}
