package main

import (
	"os"
	"strings"
	
	"github.com/therecipe/qt/core"
	"github.com/therecipe/qt/gui"
	"github.com/therecipe/qt/widgets"
)

var (
	font *gui.QFont
)

func init() {
	font = gui.NewQFont2("Courier", 16, 8, false)
}

func main() {

	widgets.NewQApplication(len(os.Args), os.Args)

	window := widgets.NewQMainWindow(nil, 0)
	window.SetWindowTitle("Text Editor")

	TextEditor := widgets.NewQTextEdit(window)
	doc := gui.NewQTextDocument2(go_code, window)
	doc.SetDefaultFont(font)
	TextEditor.SetDocument(doc)
	_ = New_GolangHighlighter(TextEditor.Document())

	window.SetCentralWidget(TextEditor)

	widgets.QApplication_SetStyle2("fusion")
	window.ShowMaximized()
	widgets.QApplication_Exec()
}

var go_code = `
package main

import (
	"os"

	"github.com/therecipe/qt/core"
	"github.com/therecipe/qt/gui"
	"github.com/therecipe/qt/widgets"
)

type Application struct {
	*widgets.QApplication
	Window    *widgets.QMainWindow
	Statusbar *widgets.QStatusBar
	TextEditor    *widgets.QTextDocument
}

var (
	ap   *Application
	font *gui.QFont
)


func init() {
	font = gui.NewQFont2("verdana", 13, 1, false)
}

func main() {
	ap = &Application{}
	ap.QApplication = widgets.NewQApplication(len(os.Args), os.Args)

	window := widgets.NewQMainWindow(nil, 0)
	ap.Window = window
	ap.Window.SetWindowTitle("Text Editor")

	ap.TextEditor = widgets.NewQTextEdit(window)
	doc := gui.NewQTextDocument()

	panel := widgets.NewQDockWidget("", ap.Window, core.Qt__Widget)
	window.AddDockWidget(core.Qt__LeftDockWidgetArea, panel)
	w := widgets.NewQWidget(ap.Window, core.Qt__Widget)
	w.SetLayout(edittreelayout)
	panel.SetWidget(w)

	// Run App
	widgets.QApplication_SetStyle2("fusion")
	window.ShowMaximized()
	widgets.QApplication_Exec()
}
`

/// https://github.com/visualfc/liteide/blob/master/liteidex/src/plugins/syntaxeditor/golanghighlighter.cpp

type GolangHighlighter struct {
	*gui.QSyntaxHighlighter
	allWords                []string
	highlightingRules       []HighlightingRule
	regexpQuotesAndComment  *core.QRegExp
	keywordFormat           *gui.QTextCharFormat
	identFormat             *gui.QTextCharFormat
	functionFormat          *gui.QTextCharFormat
	numberFormat            *gui.QTextCharFormat
	quotesFormat            *gui.QTextCharFormat
	singleLineCommentFormat *gui.QTextCharFormat
	multiLineCommentFormat  *gui.QTextCharFormat
}

type HighlightingRule struct {
	pattern *core.QRegExp
	format  *gui.QTextCharFormat
}

const (
	STATE_BACKQUOTES         = 0x04
	STATE_SINGLELINE_COMMENT = 0x08
	STATE_MULTILINE_COMMENT  = 0x10
)

func New_GolangHighlighter(document *gui.QTextDocument) *GolangHighlighter {

	gh := &GolangHighlighter{
		QSyntaxHighlighter:      gui.NewQSyntaxHighlighter2(document),
		keywordFormat:           gui.NewQTextCharFormat(),
		identFormat:             gui.NewQTextCharFormat(),
		functionFormat:          gui.NewQTextCharFormat(),
		numberFormat:            gui.NewQTextCharFormat(),
		quotesFormat:            gui.NewQTextCharFormat(),
		singleLineCommentFormat: gui.NewQTextCharFormat(),
		multiLineCommentFormat:  gui.NewQTextCharFormat(),
	}

	gh.keywordFormat.SetForeground(gui.NewQBrush4(core.Qt__darkBlue, core.Qt__SolidPattern))
	gh.keywordFormat.SetFontWeight(75) /// core.QFont__Bold
	gh.identFormat.SetForeground(gui.NewQBrush4(core.Qt__darkBlue, core.Qt__SolidPattern))
	gh.functionFormat.SetForeground(gui.NewQBrush4(core.Qt__blue, core.Qt__SolidPattern))
	gh.numberFormat.SetForeground(gui.NewQBrush4(core.Qt__darkMagenta, core.Qt__SolidPattern))
	gh.quotesFormat.SetForeground(gui.NewQBrush4(core.Qt__darkGreen, core.Qt__SolidPattern))
	gh.singleLineCommentFormat.SetForeground(gui.NewQBrush4(core.Qt__darkCyan, core.Qt__SolidPattern))
	gh.multiLineCommentFormat.SetForeground(gui.NewQBrush4(core.Qt__darkCyan, core.Qt__SolidPattern))

	var words []string
	rule := HighlightingRule{}
	highlightingRules := []HighlightingRule{}

	//number
	rule.pattern = core.NewQRegExp2(`(\b|\.)([0-9]+|0[xX][0-9a-fA-F]+|0[0-7]+)(\.[0-9]+)?([eE][+-]?[0-9]+i?)?\b`, core.Qt__CaseSensitive, core.QRegExp__RegExp)
	rule.format = gh.numberFormat
	highlightingRules = append(highlightingRules, rule)

	//function
	rule.pattern = core.NewQRegExp2(`\b[a-zA-Z_][a-zA-Z0-9_]+\s*(?=\()`, core.Qt__CaseSensitive, core.QRegExp__RegExp)
	rule.format = gh.functionFormat
	highlightingRules = append(highlightingRules, rule)

	//indent
	indent := `bool|byte|complex64|complex128|float32|float64|int8|int16|int32|int64|string|uint8|uint16|uint32|uint64|` +
		`int|uint|uintptr|true|false|iota|nil|append|cap|close|closed|complex|copy|imag|len|make|new|panic|print|println|` +
		`real|recover`
	rule.pattern = core.NewQRegExp2(`\b(`+indent+`)\b`, core.Qt__CaseSensitive, core.QRegExp__RegExp)
	rule.format = gh.identFormat
	highlightingRules = append(highlightingRules, rule)

	words = strings.Split(indent, `|`)
	gh.allWords = append(gh.allWords, words...)

	//keyword
	keyword := `break|default|func|interface|select|case|defer|go|map|struct|chan|else|goto|package|switch|` +
		`const|fallthrough|if|range|type|continue|for|import|return|var`
	rule.pattern = core.NewQRegExp2(`\b(`+keyword+`)\b`, core.Qt__CaseSensitive, core.QRegExp__RegExp)
	rule.format = gh.keywordFormat
	highlightingRules = append(highlightingRules, rule)

	words = strings.Split(keyword, `|`)
	gh.allWords = append(gh.allWords, words...)

	gh.highlightingRules = highlightingRules

	//quotes and comment
	gh.regexpQuotesAndComment = core.NewQRegExp2("//|\\\"|'|`|/\\*", core.Qt__CaseSensitive, core.QRegExp__RegExp)

	gh.ConnectHighlightBlock(gh.highlightBlock)

	return gh
}

func (gh *GolangHighlighter) highlightBlock(stext string) {

	text := core.NewQStringRef3(stext)
	startPos := 0
	endPos := text.Length()
	gh.SetCurrentBlockState(0)

	startPos, endPos, cont := gh.highlightPreBlock(text, startPos, endPos)

	if cont {
		return
	}

	//keyword and func
	for _, rule := range gh.highlightingRules {
		expression := core.NewQRegExp3(rule.pattern)
		index := expression.IndexIn(text.String(), startPos, core.QRegExp__CaretAtZero)
		for index >= 0 {
			length := expression.MatchedLength()
			gh.SetFormat(index, length, rule.format)
			gh.allWords = append(gh.allWords, text.Mid(index, length).String())
			index = expression.IndexIn(text.String(), startPos+index+length, core.QRegExp__CaretAtZero)
		}
	}

	//quote and comment
	for true {
		startPos = gh.regexpQuotesAndComment.IndexIn(text.String(), startPos, core.QRegExp__CaretAtZero)

		if startPos == -1 {
			break
		}

		cap := gh.regexpQuotesAndComment.Cap(0)
		if (cap == "\"") || (cap == "'") || (cap == "`") {
			endPos = gh.findQuotesEndPos(text, startPos+1, cap)

			if endPos == -1 {
				//multiline
				gh.SetFormat(startPos, text.Length()-startPos, gh.quotesFormat)
				if cap == "`" {
					gh.SetCurrentBlockState(STATE_BACKQUOTES)
				}
				return
			} else {
				endPos += 1
				gh.SetFormat(startPos, endPos-startPos, gh.quotesFormat)
				startPos = endPos
			}
		} else if cap == "//" {
			gh.SetFormat(startPos, text.Length()-startPos, gh.singleLineCommentFormat)
			if text.EndsWith("\\", core.Qt__CaseSensitive) {
				gh.SetCurrentBlockState(STATE_SINGLELINE_COMMENT)
			}
			return
		} else if cap == "/*" {
			endPos = text.IndexOf("*/", startPos+2, core.Qt__CaseSensitive)
			if endPos == -1 {
				//multiline
				gh.SetFormat(startPos, text.Length()-startPos, gh.multiLineCommentFormat)
				gh.SetCurrentBlockState(STATE_MULTILINE_COMMENT)
				return
			} else {
				endPos += 2
				gh.SetFormat(startPos, endPos-startPos, gh.multiLineCommentFormat)
				startPos = endPos
			}
		}
	}
}

func (gh *GolangHighlighter) highlightPreBlock(text *core.QStringRef, startPos int, endPos int) (int, int, bool) {

	state := gh.PreviousBlockState()
	
	if state == -1 {
		state = 0
	}

	if state == STATE_BACKQUOTES {
		endPos = gh.findQuotesEndPos(text, startPos, "`")

		if endPos == -1 {
			gh.SetFormat(0, text.Length(), gh.quotesFormat)
			gh.SetCurrentBlockState(STATE_BACKQUOTES)
			return startPos, endPos, true
		} else {
			endPos += 1
			gh.SetFormat(0, endPos-startPos, gh.quotesFormat)
			startPos = endPos
		}
	} else if state == STATE_MULTILINE_COMMENT {
		endPos = text.IndexOf("*/", startPos, core.Qt__CaseSensitive)
		if endPos == -1 {
			gh.SetFormat(0, text.Length(), gh.multiLineCommentFormat)
			gh.SetCurrentBlockState(gh.PreviousBlockState())
			return startPos, endPos, true
		} else {
			endPos += 2
			gh.SetFormat(0, endPos-startPos, gh.multiLineCommentFormat)
			startPos = endPos
		}
	} else if state == STATE_SINGLELINE_COMMENT {
		gh.SetFormat(0, text.Length(), gh.singleLineCommentFormat)
		if text.EndsWith("\\", core.Qt__CaseSensitive) {
			gh.SetCurrentBlockState(STATE_SINGLELINE_COMMENT)
		}
		return startPos, endPos, true
	}
	return startPos, endPos, false
}

func (gh *GolangHighlighter) findQuotesEndPos(text *core.QStringRef, startPos int, endChar string) int {

	stext := text.String()

	for pos := startPos; pos < len(stext); pos++ {
		if string(stext[pos]) == endChar {
			return pos
		} else if string(stext[pos]) == `\` && endChar != "`" { /// ?
			pos++
		}
	}
	return -1
}
