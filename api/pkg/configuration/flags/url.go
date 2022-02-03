package flags

import "net/url"

type URL struct {
	url.URL
}

func (u *URL) UnmarshalFlag(s string) error {
	url, err := url.Parse(s)
	if err != nil {
		return err
	}
	u.URL = *url
	return nil
}

func (u URL) MarshalFlag() string {
	return u.String()
}
