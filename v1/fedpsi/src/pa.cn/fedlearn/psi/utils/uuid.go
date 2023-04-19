package utils

import (
	"fmt"
	"strings"

	uuid "github.com/satori/go.uuid"
)

func UUIDStr() string {
	uid := uuid.NewV4()
	return uid.String()
}

func GetWho(user, party string) string {
	return fmt.Sprintf("%s@%s", user, party)
}

func SliceToCommaSeperatedString(s []string) string {
	if s == nil || len(s) == 0 {
		return ""
	}
	return strings.Join(s, ",")
}

func CommaSeperatedStringToSlice(s string) []string {
	if s == "" {
		return make([]string, 0)
	}
	return strings.Split(s, ",")
}
