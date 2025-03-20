package aozoratext

import (
	_ "embed" //for embedding
	"fmt"
	"testing"
)

//go:embed testdata/dogura_u.txt
var benchmarkdata string

func BenchmarkTT(b *testing.B) {

	//	setTokenizerOption("jis0208")

	for b.Loop() {

		_ = tokenizeAll(benchmarkdata, 0)

	}
}

func BenchmarkRenderer(b *testing.B) {

	for b.Loop() {

		RenderHTML(AST(benchmarkdata))

	}

}

func TestRenderText(t *testing.T) {

	SetJIS0208(true)

	text := RenderAozoraText(AST(benchmarkdata))

	fmt.Println(text)

	//		fmt.Println(toSJIS(text))

}

func TestRenderHTML(t *testing.T) {

	fmt.Println(RenderHTML(AST(benchmarkdata)))

}

func TestDogura(t *testing.T) {

	//	tk := tokenizeAll(benchmarkdata, 0)

	//	fmt.Println(tk.listAllTokens())

	n := getAozoraAST(benchmarkdata)

	fmt.Println(renderHtml(n))

}
