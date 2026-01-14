package utils

func IsLetter(ch byte) bool {
	return 'a' <= ch && ch <= 'z' || 'A' <= ch && ch <= 'Z' || ch == '_'
}

func IsDigit(ch byte) bool {
	return '0' <= ch && ch <= '9'
}

func IsEscapeSequenceCode(ch byte) bool {
	return ch == 'n' || ch == 't' || ch == '"' || ch == '\'' || ch == '\\'
}
