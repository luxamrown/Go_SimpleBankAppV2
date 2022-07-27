package util

import "strings"

func TokenToString(token string) string {
	return strings.Replace(token, "Bearer ", "", -1)
}
