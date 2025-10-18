//go:build !solution

package speller

import "strings"

var ones = []string{
	"zero", "one", "two", "three", "four", "five", "six", "seven", "eight", "nine",
	"ten", "eleven", "twelve", "thirteen", "fourteen", "fifteen",
	"sixteen", "seventeen", "eighteen", "nineteen",
}

var tens = []string{
	"", "", "twenty", "thirty", "forty", "fifty", "sixty", "seventy", "eighty", "ninety",
}

func helper(x int) string {
	var parts []string
	if x >= 100 {
		parts = append(parts, ones[x/100], "hundred")
		x %= 100
	}
	if x >= 20 {
		t := tens[x/10]
		if x%10 != 0 {
			parts = append(parts, t+"-"+ones[x%10])
		} else {
			parts = append(parts, t)
		}
	} else if x > 0 {
		parts = append(parts, ones[x])
	}
	return strings.Join(parts, " ")
}

func Spell(n int64) string {
	builder := make([]string, 0)
	if n < 0 {
		builder = append(builder, "minus")
		n *= -1
	}
	billion := n / 1_000_000_000
	n %= 1_000_000_000
	million := n / 1_000_000
	n %= 1_000_000
	thousand := n / 1_000
	n %= 1_000
	if billion > 0 {
		builder = append(builder, helper(int(billion)), "billion")
	}
	if million > 0 {
		builder = append(builder, helper(int(million)), "million")
	}
	if thousand > 0 {
		builder = append(builder, helper(int(thousand)), "thousand")
	}
	if n > 0 {
		builder = append(builder, helper(int(n)))
	}
	if million == 0 && thousand == 0 && billion == 0 && n == 0 {
		builder = append(builder, "zero")
	}
	return strings.Join(builder, " ")
}
