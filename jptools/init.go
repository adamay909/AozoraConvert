package jptools

import (
	_ "embed" //for embedding data
)

//This defines the map from JIS X 0213:2004 to Unicode.
//Based on data published by Project X0213 (included in
//data directory).

type jisuni struct {
	jis string
	uni string
}

var initialized = false

func initmap() {

	Utf8of = make(map[string]string)

	for _, d := range data {

		Utf8of[d.jis] = d.uni

	}
	initialized = true

}
