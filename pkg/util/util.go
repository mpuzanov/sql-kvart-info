package util

import "strings"

// GetValidFileName убираем недопустимые символы в имени файла
func GetValidFileName(src string) string {
	const noValidChars string = `?.,;:=+*/\"|<>[]! `
	f := func(char rune) bool {
		// признак символа разделителя
		return strings.ContainsRune(noValidChars, char)
	}
	words := strings.FieldsFunc(src, f)
	return strings.Join(words, " ")
}
