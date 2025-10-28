package vo

import (
	"errors"
	"fmt"
	"time"
)

var ErrBirthdayAfterNow = errors.New(`birthday is after today`)

// Now is used for mock only.
var Now = time.Now

type Birthday time.Time

func (b Birthday) String() string {
	return time.Time(b).Format(time.DateOnly)
}

func MustNewBirthday(b time.Time) Birthday {
	return Birthday(b)
}

func MustParseBirthday(b string) Birthday {
	bd, err := ParseBirthday(b)
	if err != nil {
		panic(err)
	}
	return bd
}

func ParseBirthday(s string) (Birthday, error) {
	t, err := time.Parse(time.DateOnly, s)
	if err != nil {
		return Birthday{}, fmt.Errorf(`parse birthday %s: %w`, s, err)
	}
	return Birthday(t), nil
}

func (b Birthday) IsZero() bool {
	return time.Time(b).IsZero()
}

func (b Birthday) Age() (int, error) {
	// We dont consider leap year for simplicity.
	const yearInHour = 24 * 365
	birthAt := time.Time(b)
	now := Now()

	if birthAt.After(now) {
		return 0, ErrBirthdayAfterNow
	}

	return int(time.Since(time.Time(b)).Hours()) / yearInHour, nil
}
