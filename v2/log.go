package aozoraconvert

import (
	"io"
	"log"
	"os"
)

var clog *log.Logger

var msglog *log.Logger

func init() {
	clog = log.New(os.Stdout, "", 0)

	msglog = log.New(os.Stdout, "", 0)
}

// SetCorrectionLog sets the output destination  regarding
// automatic correction made while parsing in a non-strict way.
func SetCorrectionLog(out io.Writer) {

	clog.SetOutput(out)

}

// SetMessageLog sets the output destination for general messages
// from parser and renderer.
func SetMessageLog(out io.Writer) {

	msglog.SetOutput(out)

}
