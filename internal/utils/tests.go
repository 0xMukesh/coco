package utils

import "fmt"

func TestMismatchErrorBuilder(testIdx int, v string, expected, got any) string {
	return fmt.Sprintf("[test #%d] %s mismatch. expected=%v, got=%v", testIdx, v, expected, got)
}
