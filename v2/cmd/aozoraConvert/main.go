package main

import (
	"archive/zip"
	"bytes"
	_ "embed" //for embedding
	"flag"
	"fmt"
	"io"
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

	jis0213 = flag.Bool("jis0213", false, "output is JIS0213 compatible.")

	frag = flag.Bool("fragment", false, "set to true if input is a fragment of aozorabunko text.")

	full = flag.Bool("full", false, "set to true to render a complete and valid output (otherwise, you only get html and latex fragments")

	supFiles = flag.Bool("supportFiles", false, "if set to true, writes supporting files to disk.")

	rubyForEmph = flag.Bool("rubyForEmph", true, "if set, use css text-emphasis for html/epub/azw3  output")
)

func init() {

	flag.Parse()

	if len(flag.Args()) < 1 {

		log.Println("must provide name of input file as argument")

		return
	}

	inputFile = filepath.Clean(flag.Args()[0])

	*outputFile = filepath.Clean(*outputFile)

	if *outputFile == "/dev/stdout" && *supFiles {

		log.Println("support files can only be output together with the file specied with the -o option")

		return

	}

	if *sjisOut {

		*jis0208 = true

		*jis0213 = false

	}

	if *format == "" {

		switch filepath.Ext(*outputFile) {

		case ".txt":
			*format = "text"

		case ".html":
			*format = "html"

		case ".json":
			*format = "json"

		case ".tex":
			*format = "latex"

		case ".epub":
			*format = "epub"

		case ".azw3":
			*format = "azw3"

		default:
			log.Println("cannot determine output file type; defaulting to html")
			*format = "html"

		}

	}

	if *format == "epub" || *format == "azw3" {

		if filepath.Ext(inputFile) != ".zip" {

			log.Println("for ", *format, "output, you need to supply a zip archive as input.")

			return

		}

	}

	if *jis0208 {
		ac.SetJIS0208()
	}

	if *jis0213 {
		ac.SetJIS0213()
	}

	ac.SetFragment(*frag)

	ac.SetStrict(true)

	ac.SetRubyEmph(*rubyForEmph)

}

func main() {

	var resp []byte

	switch *format {

	case "text", "json", "latex", "html":
		resp = parseAndRenderOnly(inputFile)

	case "tokens", "rawtokens":
		resp = tokens(inputFile)

	case "epub", "azw3":
		resp = getEbook(inputFile)

	}

	if inputFile == *outputFile {

		fmt.Println("Output and input have same file name. Renaming input to", inputFile+"~")

		os.Rename(inputFile, inputFile+"~")

	}

	dir := filepath.Dir(*outputFile)

	_, err := os.Open(dir)

	if os.IsNotExist(err) {

		err = os.MkdirAll(dir, 0755)

		if err != nil {

			fmt.Println(err)

			return

		}
	}

	writeFile(*outputFile, resp)

	if *supFiles {

		switch *format {

		case "latex":
			writeFile(filepath.Join(dir, "azcommands.tex"), []byte(ac.LaTeXdefinitions))

		case "html":
			writeFile(filepath.Join(dir, "aozora.css"), []byte(ac.AozoraCSS))

		default:

		}

		if filepath.Ext(inputFile) == ".zip" {

			extractFilesFromZip(inputFile)

		}

	}

}

func readFile(fn string) string {

	data := readData(fn)

	return stringOf(data)

}

func stringOf(data []byte) string {

	resp := string(data)

	if !utf8.Valid(data) {
		log.Println("Not UTF8. Assuming ShiftJIS.")
		resp = ac.ToUTF8(data)
	}

	return resp

}

func readData(fn string) []byte {

	inFile, err := os.Open(fn)

	if err != nil {

		log.Println(err)

		return []byte{}
	}

	defer inFile.Close()

	finfo, err := inFile.Stat()

	if err != nil {

		log.Println(err)

		return []byte{}
	}

	data := make([]byte, finfo.Size())

	inFile.Read(data)

	return data

}

func writeFile(fn string, data []byte) {

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

	outFile.Write(data)

	return

}

func getAozoraTextFromZip(inputFile string) string {

	data := readData(inputFile)

	zfs, err := zip.NewReader(bytes.NewReader(data), int64(len(data)))

	if err != nil {

		fmt.Println(err)

		return ""

	}

	for _, e := range zfs.File {

		if filepath.Ext(e.Name) == ".txt" {

			f, err := e.Open()

			if err != nil {

				fmt.Println(err)

				return ""

			}

			d, err := io.ReadAll(f)

			return stringOf(d)

			break

		}
	}

	log.Println("no Aozora text found")

	return ""
}

func parseAndRenderOnly(inputFile string) []byte {

	var data string

	if filepath.Ext(inputFile) == ".zip" {

		data = getAozoraTextFromZip(inputFile)

	} else {

		data = readFile(inputFile)

	}

	var renderer func(*ac.Node, *strings.Builder) error

	switch *format {

	case "text":
		renderer = ac.RenderAozoraText

	case "json":
		renderer = ac.RenderJSON

	case "latex":
		renderer = ac.RenderLaTeX
		if *full {
			renderer = ac.RenderLaTeXFull
		}

	case "html":
		renderer = ac.RenderHTML
		if *full {
			renderer = ac.RenderHTMLFull
		}

	}

	ast, err := ac.AST(data)

	if err != nil {

		log.Println(err)

		log.Println("\nExiting.")

		return []byte{}

	}

	w := new(strings.Builder)

	err = renderer(ast, w)

	if err != nil {

		log.Println(err)

		log.Println("\nExiting.")

		return []byte{}

	}

	converted := w.String()

	if *sjisOut {

		converted = ac.ToSJIS(converted)

	}

	return []byte(converted)

}

func tokens(inputFile string) []byte {

	var data string

	if filepath.Ext(inputFile) == ".zip" {

		data = getAozoraTextFromZip(inputFile)

	} else {

		data = readFile(inputFile)

	}

	if *format == "rawtokens" {
		return []byte(ac.ListRawTokens(string(data)))
	}

	return []byte(ac.ListProcessedTokens(data))

}

func getEbook(inputFile string) []byte {

	data := readData(inputFile)

	bk := ac.NewEbookFromZip(data)

	if *format == "epub" {

		return bk.RenderEpub()
	}

	return bk.RenderAZW3()

}

func extractFilesFromZip(infile string) {

	dir := filepath.Dir(*outputFile)

	data := readData(inputFile)

	zfs, err := zip.NewReader(bytes.NewReader(data), int64(len(data)))

	if err != nil {

		fmt.Println(err)

		return

	}

	for _, e := range zfs.File {

		if filepath.Ext(e.Name) == ".txt" {

			continue

		}

		f, err := e.Open()

		if err != nil {

			fmt.Println(err)

			return

		}

		d, err := io.ReadAll(f)

		writeFile(filepath.Join(dir, e.Name), d)

	}
}
