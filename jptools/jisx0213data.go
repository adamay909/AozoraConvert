package jptools

import (
	_ "embed" //for embedding data
	"strings"
)

//This defines the map from JIS X 0213:2004 to Unicode.
//Based on data published by Project X0213 (included in
//data directory).

//go:embed jisx0213data.txt
var data string

func init() {
	Utf8of = make(map[string]string)
	lines := strings.Split(data, "\n")

	for _, l := range lines {

		d := strings.Split(l, ",")
		if len(d) != 2 {
			continue
		}

		Utf8of[d[0]] = d[1]

	}
}
