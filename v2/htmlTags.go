package aozoratext

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
func (data *htmlTagSpec) AddStringTo(w *strings.Builder) {

	w.WriteString(data.before)

	w.WriteString("<")

	if data.open {

		w.WriteString(data.element)
		if len(data.class) > 0 {
			addToStringsBuilder(w, ` class="`, strings.Join(data.class, " "), `"`)
		}

		if len(data.extraKeyVal) > 0 {

			keys := []string{}

			for k := range data.extraKeyVal {
				keys = append(keys, k)
			}

			sort.Strings(keys)

			for _, k := range keys {
				addToStringsBuilder(w, " ", k, `="`, strings.Join(data.extraKeyVal[k], " "), `"`)
			}
		}

		if len(data.extra) > 0 {
			addToStringsBuilder(w, " ", strings.Join(data.extra, " "))
		}

		if data.id != "" {
			addToStringsBuilder(w, ` id="`, data.id, `"`)
		}

	} else {
		w.WriteString(`/`)
		w.WriteString(data.element)
	}

	w.WriteString(">")

	w.WriteString(data.after)

}

func (data *htmlTagSpec) String() string {

	output.Reset()

	output.WriteString(data.before)

	output.WriteString("<")

	if data.open {

		output.WriteString(data.element)
		if len(data.class) > 0 {
			addToStringsBuilder(output, ` class="`, strings.Join(data.class, " "), `"`)
		}

		if len(data.extraKeyVal) > 0 {

			keys := []string{}

			for k := range data.extraKeyVal {
				keys = append(keys, k)
			}

			sort.Strings(keys)

			for _, k := range keys {
				addToStringsBuilder(output, " ", k, `="`, strings.Join(data.extraKeyVal[k], " "), `"`)
			}
		}

		if len(data.extra) > 0 {
			addToStringsBuilder(output, " ", strings.Join(data.extra, " "))
		}

		if data.id != "" {
			addToStringsBuilder(output, ` id="`, data.id, `"`)
		}

	} else {
		output.WriteString(`/`)
		output.WriteString(data.element)
	}

	output.WriteString(">")

	output.WriteString(data.after)

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
