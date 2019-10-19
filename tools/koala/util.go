package main

import (
	"strings"
	"text/template"
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
