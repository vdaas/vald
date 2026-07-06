package util

import "strings"

func JoinStr(sep string, str ...string) string {
	return strings.Join(str, sep)
}
