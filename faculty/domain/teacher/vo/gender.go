package vo

import (
	"errors"
	"slices"
	"strings"
)

type Gender uint8

const (
	_ Gender = iota
	GenderMale
	GenderFamale
)

func NewGender(g uint8) (Gender, error) {
	if !slices.Contains([]Gender{GenderMale, GenderFamale}, Gender(g)) {
		return 0, ErrInvalidGender
	}
	return Gender(g), nil
}

var ErrInvalidGender = errors.New(`gender must be 'male' or 'famale'`)

func ParseGender(g string) (Gender, error) {
	switch strings.ToLower(g) {
	case `m`, `male`:
		return GenderMale, nil
	case `f`, `famale`:
		return GenderFamale, nil
	default:
		return 0, ErrInvalidGender
	}
}

func MustParseGender(g string) Gender {
	gender, err := ParseGender(g)
	if err != nil {
		panic(err)
	}
	return gender
}

func (g Gender) String() string {
	if g == GenderMale {
		return `Male`
	} else {
		return `Famale`
	}
}
