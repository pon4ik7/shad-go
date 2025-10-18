//go:build !solution

package spacecollapse

import (
	"unicode/utf8"
)

func CollapseSpaces(input string) string {
	out := make([]byte, 0, len(input))
	check := false
	for len(input) > 0 {
		err, n := utf8.DecodeRuneInString(input)
		if err == utf8.RuneError {
			out = append(out, []byte(string(utf8.RuneError))...)
			check = false
		} else {
			if input[:n] != " " && input[:n] != "\t" && input[:n] != "\n" && input[:n] != "\r" {
				check = false
			} else {
				if !check {
					out = append(out, []byte(" ")...)
					check = true
				}
			}
			if !check {
				out = append(out, input[:n]...)
			}
		}
		input = input[n:]
	}
	return string(out)
}
