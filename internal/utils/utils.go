package utils

import (
	"fmt"
	"strconv"
	"strings"
)

func IsLetter(ch byte) bool {
	return 'a' <= ch && ch <= 'z' || 'A' <= ch && ch <= 'Z' || ch == '_'
}

func IsDigit(ch byte) bool {
	return '0' <= ch && ch <= '9'
}

func IsEscapeSequenceCode(ch byte) bool {
	return ch == 'n' || ch == 't' || ch == '"' || ch == '\\'
}

func NormalizeQuotedString(str string) string {
	normalizedInput, err := strconv.Unquote(str)
	if err == nil {
		str = "\"" + normalizedInput + "\""
	}

	return str
}

func LlvmStringToGoLiteral(llvmStr string) (string, error) {
	if !strings.HasPrefix(llvmStr, "c") {
		return "", fmt.Errorf("invalid llvm string format: expected c\"...\" prefix")
	}

	if !strings.HasSuffix(llvmStr, "\"") {
		return "", fmt.Errorf("invalid llvm string format: missing closing quote")
	}

	content := llvmStr[2 : len(llvmStr)-1]

	var result strings.Builder
	i := 0

	for i < len(content) {
		if content[i] == '\\' && i+1 < len(content) {
			switch content[i+1] {
			case '0':
				switch content[i+2] {
				case 'A':
					result.WriteString("\n")
				case '0':
					result.WriteString("\x00")
				}

				i += 3
			default:
				result.WriteByte(content[i])
				i++
			}
		} else {
			result.WriteByte(content[i])
			i++
		}
	}

	return result.String(), nil
}
