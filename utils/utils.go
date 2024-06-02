package utils

import (
	"math/rand"
	"strings"
)

const alphabet = "abcdefghijklmnopqrstuvwxyz"

// RandomString generates a random string of length n
func RandomString(n int) string {
	var sb strings.Builder
	k := len(alphabet)

	for i := 0; i < n; i++ {
		c := alphabet[rand.Intn(k)] //nolint:gosec // only used for testing utils
		sb.WriteByte(c)
	}

	return sb.String()
}

func GetDuplicateColumnName(errString string) string {
	s := strings.Split(errString, ":")
	sqlIdentifier := strings.Split(s[len(s)-1], ".")

	return sqlIdentifier[len(sqlIdentifier)-1]
}
