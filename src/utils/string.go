package utils

import "strings"

func StringJoin(sep string, e ...string) string {
	return strings.Join(e, sep)
}
