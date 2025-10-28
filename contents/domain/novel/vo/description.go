package vo

import (
	"errors"

	"github.com/samber/mo"
)

type Description string

func NewDescription(d string) (mo.Option[Description], error) {
	if len(d) > 500 {
		return mo.None[Description](), errors.New(`description too long, at most 500 words.`)
	}

	if len(d) == 0 {
		return mo.Some(Description(`No description yet.`)), nil
	}

	return mo.Some(Description(d)), nil
}
