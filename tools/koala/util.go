package main

import (
	"bytes"
	"fmt"
	"strings"
	"text/template"
	"unicode"
)

var templateFuncMap = template.FuncMap{
	"Capitalize": Capitalize,
}

func Capitalize(str string) string {

	var output string
	chars := []rune(str)
	for i := 0; i < len(chars); i++ {
		if i == 0 {
			if chars[i] < 'a' || chars[i] > 'z' {
				return output
			}

			output += strings.ToUpper(string(chars[i]))
			continue
		}
		output += string(chars[i])
	}
	return output
}

func ToUnderScoreString(name string) string {
	var buffer bytes.Buffer
	for i, r := range name {
		if unicode.IsUpper(r) {
			if i != 0 {
				buffer.WriteString("_")
			}
			buffer.WriteString(fmt.Sprintf("%c", unicode.ToLower(r)))
		} else {
			buffer.WriteString(fmt.Sprintf("%c", unicode.ToLower(r)))
		}
	}

	return buffer.String()
}
