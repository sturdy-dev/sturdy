package jsontime

import (
	"strconv"
	"time"
)

type Time time.Time

// MarshalJSON is used to convert the timestamp to JSON
func (t Time) MarshalJSON() ([]byte, error) {
	return []byte(strconv.FormatInt(time.Time(t).Unix(), 10)), nil
}

// UnmarshalJSON is used to convert the timestamp from JSON
func (t *Time) UnmarshalJSON(s []byte) (err error) {
	r := string(s)
	q, err := strconv.ParseInt(r, 10, 64)
	if err != nil {
		return err
	}
	*(*time.Time)(t) = time.Unix(q, 0)
	return nil
}

func Zero() Time {
	return Time(time.Unix(0, 0))
}

func FromTimeZeroIfNil(t *time.Time) Time {
	if t == nil {
		return Zero()
	}
	return Time(*t)
}
