package tools

import "strings"

// ToLowerCamelCase a_b_c => aBC
func ToLowerCamelCase(s string) string {
	arr := ParseWords(s)
	var b strings.Builder
	for i, v := range arr {
		if i == 0 {
			b.WriteString(strings.ToLower(arr[i]))
		} else {
			b.WriteString(wordToUpperCamelCase(v))
		}
	}
	return b.String()
}

// ToUpperCamelCase a_b_c => ABC
func ToUpperCamelCase(s string) string {
	arr := ParseWords(s)
	var b strings.Builder
	for _, v := range arr {
		b.WriteString(wordToUpperCamelCase(v))
	}
	return b.String()
}

// ToLowerSnakeString aBc => a_b_c
func ToLowerSnakeString(s string) string {
	arr := ParseWords(s)
	return strings.ToLower(strings.Join(arr, "_"))
}

// ToUpperSnakeString aBc => A_B_C
func ToUpperSnakeString(s string) string {
	arr := ParseWords(s)
	return strings.ToUpper(strings.Join(arr, "_"))
}

func ParseWords(s string) []string {
	var arr []string
	start := 0
	bytes := []byte(s)
	for i := 1; i < len(bytes); i++ {
		if bytes[i] == '_' {
			arr = append(arr, s[start:i])
			start = i + 1
		} else if bytes[i] >= 'A' && bytes[i] <= 'Z' {
			arr = append(arr, s[start:i])
			start = i
		}
	}
	arr = append(arr, s[start:])
	return arr
}

// AbcAbc => abcAbc
func wordToUpperCamelCase(s string) string {
	if s == "" {
		return s
	}
	arr := []byte(strings.ToLower(s))
	if arr[0] >= 'a' {
		arr[0] = arr[0] - 'a' + 'A'
	}
	return string(arr)
}
