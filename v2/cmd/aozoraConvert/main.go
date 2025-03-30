package main

import (
	_ "embed" //for embedding
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"unicode/utf8"

	ac "github.com/adamay909/AozoraConvert/v2"
)

var inputFile string

var (
	outputFile = flag.String("o", "/dev/stdout", "name of output file; defaults to Stdout")

	sjisOut = flag.Bool("sjis", false, "generate SJIS encoded output")

	format = flag.String("format", "", "output file format; overrides format inferred from file extension of output file name; supported types are txt, html, json")

	jis0208 = flag.Bool("jis0208", false, "output is JIS0208 compatible.")

	frag = flag.Bool("fragment", false, "set to true if input is a fragment of aozorabunko text.")
)

func init() {

	flag.Parse()

	if len(flag.Args()) < 1 {

		log.Println("must provide name of input file as argument")

		return
	}

	inputFile = filepath.Clean(flag.Args()[0])

	*outputFile = filepath.Clean(*outputFile)

	if *sjisOut {

		*jis0208 = true

	}

	if *format == "" {

		switch filepath.Ext(*outputFile) {

		case ".txt":
			*format = "text"

		case ".html":
			*format = "html"

		case ".json":
			*format = "json"

		default:

			log.Println("cannot determine output file type; defaulting to html")

		}

	}

	ac.SetJIS0208(*jis0208)

	ac.SetFragment(*frag)
}

var renderer func(*ac.Node, *strings.Builder)

func main() {

	data := readFile(inputFile)

	switch *format {

	case "text":
		renderer = ac.RenderAozoraText

	case "json":
		renderer = ac.RenderJSON

	case "rawtokens":
		fmt.Println(ac.ListRawTokens(data))
		return

	case "tokens":
		fmt.Println(ac.ListProcessedTokens(data))
		return

	default:
		renderer = ac.RenderHTML

	}

	ast, err := ac.AST(data)

	if err != nil {

		log.Println("\nExiting.")

		return

	}

	w := new(strings.Builder)

	renderer(ast, w)

	converted := w.String()

	if *sjisOut {

		converted = ac.ToSJIS(converted)

	}

	if inputFile == *outputFile {

		fmt.Println("Output and input have same file name. Renaming input to", inputFile+"~")

		os.Rename(inputFile, inputFile+"~")

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
