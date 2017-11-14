package colors

import "fmt"

const (
	WHITE = iota
	BLACK
	NAVY
	GREEN
	RED
	MAROON
	PURPLE
	OLIVE
	YELLOW
	LIGHTGREEN
	TEAL
	CYAN
	ROYALBLUE
	MAGENTA
	GRAY
	LIGHTGRAY
)

// Color formats a string with the given color
func Color(str string, color int) string {
	return fmt.Sprintf("\x03%d%s\x03", color, str)
}

func Gray(str string) string {
	return Color(str, GRAY)
}
