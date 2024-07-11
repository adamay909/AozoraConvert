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

	azrconvert "github.com/adamay909/AozoraConvert"
)

var (
	web, epub, kindle, verbose bool

	of string

	logfile *os.File
)

func main() {

	defer logfile.Close()

	location := flag.Arg(0)
	if location == "" {
		printmessage("Please specify URL of Aozora Bunko book you want to convert.")

	}

	path, err := url.Parse(location)

	if err != nil {
		printmessage(err)
		return
	}

	if !web && !epub && !kindle {

		printmessage("Please specify until one format to convert to.")
		return
	}

	printmessage("Converting book at " + location)

	b := getbook(path)

	filename := setOutputName(b, location)

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

func init() {

	flag.BoolVar(&web, "web", false, "Convert to UTF-8 encoded, vertical html page.")

	flag.BoolVar(&epub, "epub", false, "Convert to EPUB3.")

	flag.BoolVar(&kindle, "kindle", false, "Convert to azw3 format for Kindle.")

	flag.BoolVar(&verbose, "v", false, "Enable verbose logging to screen and to  azrconvert.log.")

	flag.StringVar(&of, "o", "", "Name of output. Defaults to title of document plus appropriate extension.")

	flag.Parse()

	var err error

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

		return

	}

	log.SetOutput(io.Discard)
}

func printmessage[Q any](m Q) {

	log.Println(m)

	if !verbose {
		fmt.Print(m)

	}
	return

}

func getbook(path *url.URL) *azrconvert.Book {

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

	location := path.String()

	b := azrconvert.NewBookFrom(data)

	b.SetURI(filepath.Dir(location) + "/")

	b.AddFiles()
	b.SetMetadataFromPreamble()

	return b
}

func setOutputName(b *azrconvert.Book, location string) string {

	filename := strings.TrimSuffix(filepath.Base(location), filepath.Ext(location))

	if b.Title != "" {
		filename = b.Title
	}

	if of != "" {
		filename = of
	}

	return filename

}
