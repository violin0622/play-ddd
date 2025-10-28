package vo

import (
	"fmt"
	"net/url"
)

var NoCover = Cover{url.URL{
	Scheme: "https",
	Host:   "localhost",
	Path:   `a-path-to-image`,
}}

type Cover struct {
	url.URL
}

func NewCover(u string) (Cover, error) {
	url, err := url.Parse(u)
	if err != nil {
		return Cover{}, fmt.Errorf(`parse url string: %w`, err)
	}

	return Cover{*url}, nil
}

func (c Cover) String() string { return c.URL.String() }
