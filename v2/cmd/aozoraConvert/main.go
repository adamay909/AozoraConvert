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
	outputFile = flag.String("o", "/dev/stdout", "出力ファイル名. 指定されなければstdout.")

	sjisOut = flag.Bool("sjis", false, "出力ファイルのエンコーディングをShift_JISにする. デフォルトはutf-8.")

	format = flag.String("format", "", "出力ファイル形式. 出力ファイル名に拡張子より優先される. 可能な形式はtex, txt, html, json, epub, azw3")

	jis0208 = flag.Bool("jis0208", false, "外字の置き換えなし.")

	jis0213 = flag.Bool("jis0213", false, "外字の置き換えをJIS0213範囲内に抑える.")

	frag = flag.Bool("fragment", false, "入力ファイルを青空文庫形式ファイルの断片とみなす.")

	full = flag.Bool("full", false, "ヘッダーやプリアンブルを含む完全な形式のファイルを出力する.")

	supFiles = flag.Bool("supportFiles", false, "サポートファイルも出力する.")

	rubyForEmph = flag.Bool("rubyForEmph", true, "出力がHTML系の場合、傍点類にCSSのtext-emphasisを利用する.")

	parsable = flag.Bool("parsable", false, "出力から再度AST抽出可能なようにする.")

	strict = flag.Bool("strict", false, "自動修正を行わない.")
)

func init() {

	flag.Usage = func() {

		fmt.Println("\n", os.Args[0], "[オプション] 入力ファイル名")

		fmt.Println()

		flag.PrintDefaults()
	}

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
			*format = "txt"

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

	ac.SetStrict(*strict)

	ac.SetRubyEmph(*rubyForEmph)

	ac.SetParsable(*parsable)

	//ac.SetVerbose()

}

func main() {

	var resp []byte

	switch *format {

	case "txt", "json", "latex", "html":
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
			if *sjisOut {
				writeFile(filepath.Join(dir, "azcommands.tex"), []byte(ac.ToSJIS(ac.LaTeXdefinitions)))
			} else {
				writeFile(filepath.Join(dir, "azcommands.tex"), []byte(ac.LaTeXdefinitions))
			}

		case "html":
			if *sjisOut {
				writeFile(filepath.Join(dir, "aozora.css"), []byte(ac.ToSJIS(ac.AozoraCSS)))
			} else {
				writeFile(filepath.Join(dir, "aozora.css"), []byte(ac.AozoraCSS))
			}

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

	case "txt":
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

		log.Println("\nErrors while parsing. Exiting.")

		return []byte{}

	}

	w := new(strings.Builder)

	err = renderer(ast, w)

	if err != nil {

		log.Println(err)

		log.Println("\nErrors while rendering. Exiting.")

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
