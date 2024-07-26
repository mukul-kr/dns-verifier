package reader

import (
	"strings"
)

type TXTHandler struct{}

func (h TXTHandler) Handle(content string) func() ([]string, error) {
	return func() ([]string, error) {
		// Implement TXT handling logic here
		parts := strings.FieldsFunc(content, func(r rune) bool {
			return r == ',' || r == ' '
		})
		return parts, nil
	}
}