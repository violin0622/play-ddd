package vo

import (
	"errors"
	"strings"
)

var (
	ErrTitleTooLong     = errors.New(`title exceed max length 256`)
	ErrTitleInvalidChar = errors.New(`title contains invalid charactor`)
)

type Title string

func NewTitle(s string) (Title, error) {
	if len(s) >= 256 {
		return ``, ErrTitleTooLong
	}

	if strings.Contains(s, `@#$%^&*()_+-=\/`) {
		return ``, ErrTitleInvalidChar
	}

	return Title(s), nil
}
