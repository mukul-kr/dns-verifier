package reader

import (
	"strings"
)

type TerminalHandler struct{}

func (h TerminalHandler) Handle(content string) func() ([]string, error) {
	return func() ([]string, error) {
		// Implement terminal handling logic here
		parts := strings.FieldsFunc(content, func(r rune) bool {
			return r == ',' || r == ' '
		})
		return parts, nil
	}
}
