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

		fmt.Println()

		fmt.Println("生成された青空注記:")

		fmt.Println()

		fmt.Println(ac.RenderAozoraText(n))

		if n.HasDescendantOfType("gaiji char") || n.HasDescendantOfType("accent string") || n.HasDescendantOfType("special char") || n.HasDescendantOfType("gaiji node") {

			ac.SetJIS0208(true)

			fmt.Println()

			fmt.Println("JIS0208互換で生成された注記:")

			fmt.Println(ac.RenderAozoraText(n))

			ac.SetJIS0208(false)

		}

		fmt.Println()

		fmt.Println("生成されたHTML:")

		fmt.Println()

		fmt.Println(ac.RenderHTML(n))

		fmt.Println()

		fmt.Println("抽出されたAST:")

		fmt.Println()

		fmt.Println(ac.RenderJSON(n))

		fmt.Println()

		fmt.Println("---------------------------------------------")

	}
}
