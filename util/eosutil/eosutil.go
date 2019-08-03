package eosutil

import (
	"strings"
)

const AccountChars = ".12345abcdefghijklmnopqrstuvwxyz"

func CheckEOSAccount(s string) bool {
	if len(s) > 12 {
		return false
	}
	chars := []rune(s)
	for _, c := range chars {
		if !strings.ContainsRune(AccountChars, c) {
			return false
		}
	}
	return true
}
