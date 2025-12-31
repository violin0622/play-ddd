package xerr

import "fmt"

func Expect(err *error, f string, a ...any) {
	if err == nil || *err == nil {
		return
	}

	*err = fmt.Errorf(`%s: %w`, fmt.Sprintf(f, a...), *err)
}
