package main

import (
	_ "embed"
	"fmt"
	"strings"

	//for embedding

	ac "github.com/adamay909/AozoraConvert/v2"
)

//go:embed data.txt
var data string

func main() {

	latexSample()

	return

	lines := strings.Split(data, "::::\n")

	ac.SetFragment(true)

	for _, l := range lines {

		if len(l) == 0 {
			continue
		}

		if strings.HasPrefix(l, "####") {

			fmt.Println("***************************************************")
			fmt.Println(l)

			fmt.Println()

			continue

		}
		n, err := ac.AST(l)

		fmt.Println("もとの青空注記:")

		fmt.Println()

		fmt.Println(l)

		if err != nil {

			fmt.Println("\ncould not handle this case!!")

			fmt.Println(err.Error())

			fmt.Println("---------------------------------------------")

			continue
		}

		w := new(strings.Builder)
		fmt.Println()

		fmt.Println("生成された青空注記:")

		fmt.Println()

		ac.RenderAozoraText(n, w)

		fmt.Println(w)

		if n.HasDescendantOfType("gaiji char") || n.HasDescendantOfType("accent string") || n.HasDescendantOfType("special char") || n.HasDescendantOfType("gaiji node") {

			ac.SetJIS0208(true)

			fmt.Println()

			fmt.Println("JIS0208互換で生成された注記:")

			w.Reset()

			ac.RenderAozoraText(n, w)

			fmt.Println(w)

			ac.SetJIS0208(false)

		}

		fmt.Println()

		fmt.Println("生成されたHTML:")

		fmt.Println()

		w.Reset()

		ac.RenderHTML(n, w)

		fmt.Println(w)

		fmt.Println()

		fmt.Println("生成されたLaTeX:")

		w.Reset()

		fmt.Println()

		ac.RenderLaTeX(n, w)

		fmt.Println(w)

		fmt.Println()

		/*
			fmt.Println("抽出されたAST:")

			fmt.Println()

			w.Reset()

			ac.RenderJSON(n, w)

			fmt.Println(w)

			fmt.Println()
		*/
		fmt.Println("---------------------------------------------")

	}
}

func latexSample() {

	w := new(strings.Builder)

	sections := strings.Split(data, "####")

	ac.SetFragment(true)

	fmt.Println()

	fmt.Println(`\documentclass[a4paper,tate]{jlreq}`)

	fmt.Println(`\input{azcommands}`)

	fmt.Println(`\setlength{\parindent}{0em}`)

	fmt.Println(`\geometry{landscape}`)

	fmt.Println(`\renewcommand{\baselinestretch}{1.5}`)
	fmt.Println(`\begin{document}`)

	for _, s := range sections {

		if len(s) == 0 {
			continue
		}

		idx := strings.Index(s, "\n")

		fmt.Println(`\subsubsection*{` + s[:idx] + `}`)

		lines := strings.Split(s[idx:], "::::")

		for _, l := range lines {

			if len(l) == 0 {
				continue
			}

			if strings.HasPrefix(l, "%%%") {
				continue
			}

			n, _ := ac.AST(l)

			fmt.Println(`\footnotesize`)

			fmt.Println(`\linespread{1.0}`)

			fmt.Println(`\begin{verbatim}`)

			fmt.Println(l)

			fmt.Println(`\end{verbatim}`)

			fmt.Println(`\linespread{1.4}`)
			fmt.Println(`\normalsize`)

			w.Reset()

			fmt.Println(`LATEX OUTPUT:`)

			fmt.Println(`\vskip 2\zw`)

			fmt.Println()

			ac.RenderLaTeX(n, w)

			fmt.Println(w)

			fmt.Println()

			fmt.Println(`\newpage`)
		}
	}

	fmt.Println(`\end{document}`)

}
