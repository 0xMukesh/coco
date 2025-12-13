package utils

func IsLetter(ch byte) bool {
	return 'a' <= ch && ch <= 'z' || 'A' <= ch && ch <= 'Z' || ch == '_'
}

func IsNumber(ch byte) bool {
	return '0' <= ch && ch <= '9'
}
