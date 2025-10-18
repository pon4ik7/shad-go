//go:build !solution

package reverse

import "unicode/utf8"

func Reverse(input string) string {
	out := make([]byte, 0, len(input))
	for len(input) > 0 {
		r, n := utf8.DecodeLastRuneInString(input)
		if r == utf8.RuneError {
			out = append(out, []byte(string(utf8.RuneError))...)
		} else {
			out = append(out, input[(len(input)-n):]...)
		}
		input = input[:len(input)-n]
	}
	return string(out)
}
