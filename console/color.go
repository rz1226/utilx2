package console

import (
	"fmt"
	"runtime"
	"strconv"
)

const (
	TextBlack = iota + 30
	TextRed
	TextGreen
	TextYellow
	TextBlue
	TextMagenta
	TextCyan
	TextWhite
)

func ConvertToString(input interface{}) string {
	var output string
	switch val := input.(type) {
	case int:
		output = strconv.Itoa(val)
	case string:
		output = val
	case byte:
		output = string(val)
	case bool:
		if val == true {
			output = "true"
		} else {
			output = "false"
		}
	case error:
		output = val.Error()

	default:
		output = ""
	}

	return output
}

func ConvertToBlackStr(info interface{}) string {
	return textAddColor(TextBlack, ConvertToString(info))
}
func PrintBlack(info ...interface{}) {
	fmt.Print(ConvertToBlackStr(fmt.Sprintln(info...)))
}

func ConvertToRedStr(info interface{}) string {
	return textAddColor(TextRed, ConvertToString(info))
}
func PrintRed(info ...interface{}) {
	fmt.Print(ConvertToRedStr(fmt.Sprintln(info...)))
}

func ConvertToGreenStr(info interface{}) string {
	return textAddColor(TextGreen, ConvertToString(info))
}
func PrintGreen(info ...interface{}) {
	fmt.Print(ConvertToGreenStr(fmt.Sprintln(info...)))
}
func ConvertToYellowStr(info interface{}) string {
	return textAddColor(TextYellow, ConvertToString(info))
}
func PrintYellow(info ...interface{}) {
	fmt.Print(ConvertToYellowStr(fmt.Sprintln(info...)))
}

func ConvertToBlueStr(info interface{}) string {
	return textAddColor(TextBlue, ConvertToString(info))
}
func PrintBlue(info ...interface{}) {
	fmt.Print(ConvertToBlueStr(fmt.Sprintln(info...)))
}
func ConvertToMagentaStr(info interface{}) string {
	return textAddColor(TextMagenta, ConvertToString(info))
}
func PrintMagenta(info ...interface{}) {
	fmt.Print(ConvertToMagentaStr(fmt.Sprintln(info...)))
}
func ConvertToCyanStr(info interface{}) string {
	return textAddColor(TextCyan, ConvertToString(info))
}
func PrintCyan(info ...interface{}) {
	fmt.Print(ConvertToCyanStr(fmt.Sprintln(info...)))
}
func ConvertToWhiteStr(info interface{}) string {
	return textAddColor(TextWhite, ConvertToString(info))
}
func PrintWhite(info ...interface{}) {
	fmt.Print(ConvertToWhiteStr(fmt.Sprintln(info...)))
}
func textAddColor(color int, str string) string {
	if runtime.GOOS == "windows" {
		return str
	}
	switch color {
	case TextBlack:
		return fmt.Sprintf("\x1b[0;%dm%s\x1b[0m", TextBlack, str)
	case TextRed:
		return fmt.Sprintf("\x1b[0;%dm%s\x1b[0m", TextRed, str)
	case TextGreen:
		return fmt.Sprintf("\x1b[0;%dm%s\x1b[0m", TextGreen, str)
	case TextYellow:
		return fmt.Sprintf("\x1b[0;%dm%s\x1b[0m", TextYellow, str)
	case TextBlue:
		return fmt.Sprintf("\x1b[0;%dm%s\x1b[0m", TextBlue, str)
	case TextMagenta:
		return fmt.Sprintf("\x1b[0;%dm%s\x1b[0m", TextMagenta, str)
	case TextCyan:
		return fmt.Sprintf("\x1b[0;%dm%s\x1b[0m", TextCyan, str)
	case TextWhite:
		return fmt.Sprintf("\x1b[0;%dm%s\x1b[0m", TextWhite, str)
	default:
		return str
	}
}
