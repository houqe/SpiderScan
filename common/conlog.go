package common

import (
	"fmt"
	"github.com/fatih/color"
)

//func Black(str string) string {
//	return textColor(textBlack, str)
//}

func Red(str string) {
	red := color.New(color.FgRed, color.Bold)
	red.Printf(str)
}
func Yellow(str string) {
	yellow := color.New(color.FgYellow, color.Bold)
	yellow.Printf(str)
}
func Green(str string) {
	green := color.New(color.FgGreen, color.Bold)
	green.Printf(str)
}

//func Cyan(str string) string {
//	return textColor(textCyan, str)
//}
//func Blue(str string) string {
//	return textColor(textBlue, str)
//}
//func Purple(str string) string {
//	return textColor(textPurple, str)
//}
//func White(str string) string {
//	return textColor(textWhite, str)
//}

func textColor(color int, str string) string {
	return fmt.Sprintf("\x1b[0;%dm%s\x1b[0m", color, str)
}
