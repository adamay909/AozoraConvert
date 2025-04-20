package aozoraconvert

import (
	"sort"
	"strings"
)

type htmlTagSpec struct {
	before, after string
	element       string
	class         []string
	id            string
	extra         []string
	extraKeyVal   map[string][]string
	open          bool
	selfclose     bool
}

func newHtag(elem string) *htmlTagSpec {

	h := new(htmlTagSpec)

	h.setElement(elem)

	h.setOpen()

	h.extraKeyVal = make(map[string][]string)

	return h

}

func newCloseHtag(elem string) *htmlTagSpec {

	h := new(htmlTagSpec)

	h.setElement(elem)

	h.setClose()

	return h

}

func (h *htmlTagSpec) AddStringTo(w *strings.Builder) {

	w.WriteString(h.before)

	w.WriteString("<")

	if h.open {

		w.WriteString(h.element)
		if len(h.class) > 0 {
			addToStringsBuilder(w, ` class="`, strings.Join(h.class, " "), `"`)
		}

		if len(h.extraKeyVal) > 0 {

			keys := []string{}

			for k := range h.extraKeyVal {
				keys = append(keys, k)
			}

			sort.Strings(keys)

			for _, k := range keys {
				addToStringsBuilder(w, " ", k, `="`, strings.Join(h.extraKeyVal[k], " "), `"`)
			}
		}

		if len(h.extra) > 0 {
			addToStringsBuilder(w, " ", strings.Join(h.extra, " "))
		}

		if h.id != "" {
			addToStringsBuilder(w, ` id="`, h.id, `"`)
		}

		if h.selfclose {
			w.WriteString(`/`)
		}

	} else {
		w.WriteString(`/`)
		w.WriteString(h.element)
	}

	w.WriteString(">")

	w.WriteString(h.after)

}

func (h *htmlTagSpec) String() string {

	output.Reset()

	output.WriteString(h.before)

	output.WriteString("<")

	if h.open {

		output.WriteString(h.element)
		if len(h.class) > 0 {
			addToStringsBuilder(output, ` class="`, strings.Join(h.class, " "), `"`)
		}

		if len(h.extraKeyVal) > 0 {

			keys := []string{}

			for k := range h.extraKeyVal {
				keys = append(keys, k)
			}

			sort.Strings(keys)

			for _, k := range keys {
				addToStringsBuilder(output, " ", k, `="`, strings.Join(h.extraKeyVal[k], " "), `"`)
			}
		}

		if len(h.extra) > 0 {
			addToStringsBuilder(output, " ", strings.Join(h.extra, " "))
		}

		if h.id != "" {
			addToStringsBuilder(output, ` id="`, h.id, `"`)
		}

		if h.selfclose {

			output.WriteString(`/`)
		}

	} else {
		output.WriteString(`/`)
		output.WriteString(h.element)
	}

	output.WriteString(">")

	output.WriteString(h.after)

	return output.String()

}

func (h *htmlTagSpec) setElement(c string) {

	h.element = c

}

func (h *htmlTagSpec) addClass(c string) {

	h.class = append(h.class, c)

	return
}

func (h *htmlTagSpec) setBefore(c string) {

	h.before = c

	return

}

func (h *htmlTagSpec) setAfter(c string) {

	h.after = c

	return
}

func (h *htmlTagSpec) setOpen() {

	h.open = true

	return
}

func (h *htmlTagSpec) setClose() {

	h.open = false

	return

}

func (h *htmlTagSpec) setSelfClose() {

	h.selfclose = true

	return

}

func (h *htmlTagSpec) setID(c string) {

	h.id = c

	return
}

func (h *htmlTagSpec) addExtra(c string) {

	h.extra = append(h.extra, c)

	return
}

func (h *htmlTagSpec) addExtraKeyVal(k string, v string) {

	h.extraKeyVal[k] = append(h.extraKeyVal[k], v)

}
