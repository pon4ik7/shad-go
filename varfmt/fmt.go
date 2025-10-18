//go:build !solution

package varfmt

import (
	"fmt"
	"strconv"
	"strings"
	"unicode/utf8"
)

func containsBraces(s string) bool {
	for _, i := range s {
		if i == '{' || i == '}' {
			return true
		}
	}
	return false
}

func Sprintf(format string, args ...interface{}) string {
	if !containsBraces(format) {
		return format
	}

	type ref struct {
		start int
		end   int
		idx   int
	}

	maxIdx := -1

	refs := make([]ref, 0, len(args))
	count := 0
	n := len(format)
	i := 0
	for i < n {
		c := format[i]
		if c != '{' {
			_, addr := utf8.DecodeRuneInString(format[i:])
			i += addr
			continue
		}

		// find '{' and try find '}'
		j := i + 1
		for j < n && format[j] != '}' {
			j++
		}
		if j >= n {
			i++
			continue
		}

		// {} situation
		if j == i+1 {
			idx := count
			count++
			refs = append(refs, ref{i, j, idx})
			maxIdx = max(maxIdx, idx)
			i = j + 1
			continue
		}

		// {number} situation
		idx, err := strconv.Atoi(format[i+1 : j])
		if err != nil {
			i++
			continue
		}
		refs = append(refs, ref{i, j, idx})
		count++
		i = j + 1
		maxIdx = max(maxIdx, idx)
	}

	if maxIdx >= len(args) {
		panic("varfmt: argument out of range")
	}
	cache := make([]string, maxIdx+1)
	has := make([]bool, maxIdx+1)
	for _, r := range refs {
		if !has[r.idx] {
			cache[r.idx] = fmt.Sprint(args[r.idx])
			has[r.idx] = true
		}
	}

	estimated := len(format)
	for _, r := range refs {
		estimated += len(cache[r.idx]) - (r.end - r.start)
	}

	builder := strings.Builder{}
	builder.Grow(estimated)
	last := 0
	for _, r := range refs {
		builder.WriteString(format[last:r.start])
		builder.WriteString(cache[r.idx])
		last = r.end + 1
	}
	builder.WriteString(format[last:])
	return builder.String()
}
