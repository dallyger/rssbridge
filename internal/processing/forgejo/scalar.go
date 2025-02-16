package forgejo

import (
	"strings"
	"time"
)

type Time struct {
	time.Time
}

func (self *Time) UnmarshalJSON(b []byte) error {

	s := strings.Trim(string(b), "\"")
	if s == "null" {
		return nil
	}

	t, err := time.Parse(time.RFC3339, s)
	self.Time = t
	return err
}
