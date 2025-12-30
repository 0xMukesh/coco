package utils

func IsLetter(char byte) bool {
	return ('a' <= char && char <= 'z') || ('A' <= char && char <= 'Z') || char == '_'
}

func IsDigit(char byte) bool {
	return '0' <= char && char <= '9' || char == '.'
}

func IsWhitespace(char byte) bool {
	return char == ' ' || char == '\t' || char == '\r' || char == '\n'
}
