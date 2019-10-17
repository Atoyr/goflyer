package util

import "fmt" 

type Attribute int

const (
	Clear Attribute = iota
	Bold
	Light
	Italic
	Underline
	Blink
	FastBlink
	Reverse
	Hide
	Undo
)

const (
	FgBlack Attribute = iota + 30
	FgRed
	FgGreen
	FgYellow
	FgBlue
	FgMagenta
	FgCyan
	FgWhite
)

const (
	FgHiBlack Attribute = iota + 90
	FgHiRed
	FgHiGreen
	FgHiYellow
	FgHiBlue
	FgHiMagenta
	FgHiCyan
	FgHiWhite
)

const (
	BgBlack Attribute = iota + 40
	BgRed
	BgGreen
	BgYellow
	BgBlue
	BgMagenta
	BgCyan
	BgWhite
)

const (
	BgHiBlack Attribute = iota + 100
	BgHiRed
	BgHiGreen
	BgHiYellow
	BgHiBlue
	BgHiMagenta
	BgHiCyan
	BgHiWhite
)

const FgClear Attribute = 39
const BgClear Attribute = 49
const esc = "\x1b"
const foreground = "38;5;"
const background = "48;5;"

func (a Attribute) String() string {
	attribute := 0
	ret := ""
	if a < -255 {
		// Background
		attribute = (int(a) + 255) * -1
		ret = fmt.Sprintf("%v%v", background, attribute)
	} else if a < 0 {
		// Forground
		attribute = int(a) * -1
		ret = fmt.Sprintf("%v%v", foreground, attribute)
	} else if a < 256 {
		ret = fmt.Sprintf("%v", attribute)
	}

	return ret
}

func GetMultiColorAttribute(colorCode int, isBackground bool) Attribute {
	if colorCode < 0 || 255 < colorCode {
		return 0
	} else {
		val := colorCode * -1
		if isBackground {
			val += -255
		}
		return Attribute(val)
	}
}

func ApplyAttribute(str string, attributes ...Attribute) string {
	ret := str
	codes := attributes[0].String()
	for i := range attributes[1:] {
		codes = fmt.Sprintf("%v;%v", codes, attributes[i+1].String())
	}
	ret = fmt.Sprintf("%v[%vm%v", esc, codes, ret)
	// color clear
	ret = fmt.Sprintf("%v%v[%vm", ret, esc, Clear)
	return ret
}
