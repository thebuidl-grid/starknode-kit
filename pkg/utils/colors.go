package utils

import "fmt"

// ANSI Color Codes
const (
	ColorReset  = "\033[0m"
	ColorRed    = "\033[31m"
	ColorGreen  = "\033[32m"
	ColorYellow = "\033[33m"
	ColorBlue   = "\033[34m"
	ColorPurple = "\033[35m"
	ColorCyan   = "\033[36m"
	ColorWhite  = "\033[37m"
	ColorBold   = "\033[1m"
)

// Functions to colorize strings
func Red(s string) string {
	return fmt.Sprintf("%s%s%s", ColorRed, s, ColorReset)
}

func Green(s string) string {
	return fmt.Sprintf("%s%s%s", ColorGreen, s, ColorReset)
}

func Yellow(s string) string {
	return fmt.Sprintf("%s%s%s", ColorYellow, s, ColorReset)
}

func Blue(s string) string {
	return fmt.Sprintf("%s%s%s", ColorBlue, s, ColorReset)
}

func Purple(s string) string {
	return fmt.Sprintf("%s%s%s", ColorPurple, s, ColorReset)
}

func Cyan(s string) string {
	return fmt.Sprintf("%s%s%s", ColorCyan, s, ColorReset)
}

func Bold(s string) string {
    return fmt.Sprintf("%s%s%s", ColorBold, s, ColorReset)
}
