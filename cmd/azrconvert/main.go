package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"

	azrconvert "github.com/adamay909/AozoraConvert/azrconvert"
)

var (
	web, zip, epub, epub3, kindle, azw3, verbose bool

	infile, outfile string

	logfile *os.File
)

func init() {

	flag.BoolVar(&web, "web", false, "Convert to UTF-8 encoded, vertical html page.")

	flag.BoolVar(&zip, "zip", false, "Alias for -web.")

	flag.BoolVar(&epub, "epub", false, "Convert to EPUB3.")

	flag.BoolVar(&epub3, "epub3", false, "Alias for -epub.")

	flag.BoolVar(&kindle, "kindle", false, "Convert to azw3 format for Kindle.")

	flag.BoolVar(&azw3, "azw3", false, "Alias for kindle.")

	flag.BoolVar(&verbose, "v", false, "Enable verbose logging to screen and to  azrconvert.log.")

	flag.StringVar(&outfile, "o", "", "Name output  as `name` + extension. Defaults to title of document plus appropriate extension.")

	flag.StringVar(&infile, "i", "", "Convert local `file`.")

	flag.Parse()

	var err error

	log.SetOutput(io.Discard)

	if verbose {
		logfile, err = os.OpenFile("azrconvert.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)

		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		err = logfile.Truncate(0)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		log.SetOutput(io.MultiWriter(logfile, os.Stdout))
	}

	fmt.Println()

	web = web || zip

	epub = epub || epub3

	kindle = kindle || azw3

	return

}

func main() {

	var b *azrconvert.Book
	var location, filename string

	defer logfile.Close()

	//if flag.Arg(0) != "" {
	//		location = flag.Arg(0)
	//	}

	if !web && !epub && !kindle {

		printmessage("Please specify until one format to convert to.")
		return
	}

	if infile == "" {
		b = getbookFromURL(flag.Arg(0))
	} else {
		b = getbookFromLocal(infile)
	}

	filename = setOutputName(b, location)

	if web {

		err := os.WriteFile(filename+".zip", b.RenderWebpagePackage(), 0644)

		if err != nil {
			printmessage(err)
		}

		printmessage("Output written to " + filename + ".zip.")
	}

	if epub {

		err := os.WriteFile(filename+".epub", b.RenderEpub(), 0644)

		if err != nil {
			printmessage(err)
		}

		printmessage("Output written to " + filename + ".epub.")
	}

	if kindle {

		err := os.WriteFile(filename+".azw3", b.RenderAZW3(), 0644)

		if err != nil {
			printmessage(err)
		}

		printmessage("Output written to " + filename + ".azw3.")
	}

}

func printmessage[Q any](m Q) {

	log.Println(m)

	if !verbose {
		fmt.Println(m)

	}
	return

}

func getbookFromLocal(path string) *azrconvert.Book {

	log.Println("Converting from local files won't download any external graphics.")

	data, err := os.ReadFile(path)

	if err != nil {
		printmessage(err)
		logfile.Close()
		os.Exit(1)
	}

	b := azrconvert.NewBookFrom(data)
	b.AddFiles()
	b.SetMetadataFromPreamble()
	return b

}

func getbookFromURL(location string) *azrconvert.Book {

	log.Println("Converting book at " + location)

	if location == "" {
		printmessage("Please specify URL of Aozora Bunko book you want to convert.")

	}

	path, err := url.Parse(location)
	if err != nil {
		printmessage(err)
		logfile.Close()
		os.Exit(1)
	}
	r, err := http.Get(path.String())

	if err != nil {
		printmessage(err)
		logfile.Close()
		os.Exit(1)
	}

	if r.StatusCode != 200 {
		printmessage(r.Status)
		logfile.Close()
		os.Exit(1)
	}

	data, _ := io.ReadAll(r.Body)

	b := azrconvert.NewBookFrom(data)

	b.SetURI(location)

	b.AddFiles()
	b.SetMetadataFromPreamble()

	return b
}

func setOutputName(b *azrconvert.Book, location string) (filename string) {

	if location == "" {
		filename = "output"
	}

	filename = strings.TrimSuffix(filepath.Base(location), filepath.Ext(location))

	if b.Title != "" {
		filename = b.Title
	}

	if outfile != "" {
		filename = outfile
	}

	return filename

}
