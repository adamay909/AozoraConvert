package aozoraConvert

const (
	markupNoteStartStr = "\n" + `-------------------------------------------------------
【テキスト中に現れる記号について】` + "\n"

	markupNoteEndStr = "\n" + `-------------------------------------------------------` + "\n"

	markupNoteDelimiterStr = `-------------------------------------------------------`

	rubyStartStr = `《`

	rubyEndStr = `》`

	rubyParentStartStr = `｜`

	gaijiMarkerStr = `※［＃`

	noteStartStr = `［＃`

	noteEndStr = `］`

	bibInfoStartStr = "\n" + `底本：`

	emptyStr = ""

	lineBreakStr = "\n"

	refStartStr = "「"

	refEndStr = "」"

	kunojiStr = "／＼"

	kunojiEndStr = "＼"

	kunojiStrU = "〳〵"

	kunojiDakuStr = "／″＼"

	kunojiDakuStrU = "〴〵"

	blockStartStr = "ここから"

	blockEndStr = "ここで"

	formatEndStr = "終わり"

	referenceMarkStr = "※"

	accentStartStr = "〔"

	accentEndStr = "〕"

	mainTextEndStr = "本文終わり"

	fwqSpaceStr = "？　"

	fwqRegStr = "？ "

	fwexSpaceStr = "！　"

	fwexRegStr = "！ "

	dashStr = "―"

	dashStrLatex = "---"
)

var kuntenchars = []rune("レ一二三四五六七八九十上中下甲乙丙丁天地人元亨利貞乾坤")

// jisage etc.

var indentationMarker = []string{
	"字下げ",
}

var narrowparMarker = []string{
	"字詰め",
}

var bottomalignMarker = []string{
	"地付き",
	"字上げ",
}

// sectioning
var sectionMarker = []string{
	"大見出し",
	"中見出し",
	"小見出し",
}

var inlineSectionMarker = []string{
	"同行中見出し",
	"同行大見出し",
	"同行小見出し",
}

var windowSectionMarker = []string{
	"窓中見出し",
	"窓大見出し",
	"窓小見出し",
}

//emphasis

var lineDecoMarker = []string{
	"二重傍線",
	"傍線",
	"鎖線",
	"破線",
	"波線",
	"罫囲み",
}

var decoMarker = []string{
	"白ゴマ傍点",
	"白丸傍点",
	"黒三角傍点",
	"白三角傍点",
	"二重丸傍点",
	"蛇の目傍点",
	"丸傍点",
	"ばつ傍点",
	"傍点",
}

var decoString = []string{
	"﹆",
	"◦",
	"▲",
	"△",
	"◎ ",
	"◉",
	"•",
	"✕",
	"﹅",
}

// notes
var inlineNoteMarker = []string{
	"割り注",
}

var rubylikeNoteSimpleMarker = []string{
	"ルビ",
	"注記",
}

var rubylikeNoteMarker = []string{
	"ルビ付き",
	"注記付き",
}

var captionMarker = []string{
	"キャプション",
}

// page level formatting
var centeringMarker = []string{
	"ページの左右中央",
}

var paginationMarker = []string{
	"改丁",
	"改ページ",
	"改見開き",
	"改段",
}

//fontspec

var fontShapeMarker = []string{
	"太字",
	"斜体",
}

var fontsizeMarker = []string{

	"大きな文字",
	"小さな文字",
}

//subscript

var offsetMarker = []string{
	"上付き小文字",
	"下付き小文字",
	"行右小書き",
	"行左小書き",
	"上付き小文字",
	"下付き小文字",
}

//text direction

var directionMarker = []string{
	"縦中横",
	"横組み",
}

var warichuLineBreakMarker = []string{
	"改行",
}

var _sectionLevelMarker = []string{
	"大区分",
	"中区分",
	"小区分",
}

var impliedCloserMarker = []string{}

var formattingMarker = []string{}

var rubylikeMarker = []string{}

var impliedOpenerMarker = []string{}

var pairMarker []string

func init() {

	includes := [][]string{
		indentationMarker,
		narrowparMarker,
		bottomalignMarker,
		captionMarker,
		inlineSectionMarker,
		windowSectionMarker,
		sectionMarker,
		lineDecoMarker,
		decoMarker,
		inlineNoteMarker,
		rubylikeNoteSimpleMarker,
		rubylikeNoteMarker,
		fontShapeMarker,
		fontsizeMarker,
		offsetMarker,
		directionMarker,
		centeringMarker,
	}

	for _, s := range includes {

		pairMarker = append(pairMarker, s...)

	}

	includes = [][]string{
		indentationMarker,
		bottomalignMarker,
	}

	for _, s := range includes {

		impliedCloserMarker = append(impliedCloserMarker, s...)

	}

	impliedOpenerMarker = append(impliedOpenerMarker, pairMarker...)

}
