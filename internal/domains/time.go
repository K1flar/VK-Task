package domains

import (
	"fmt"
	"time"
)

const layout = time.DateOnly

type Time time.Time

func (t *Time) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf(`"%s"`, time.Time(*t).Format(layout))), nil
}

func (t *Time) UnmarshalJSON(data []byte) error {
	s := string(data)
	s = s[1 : len(s)-1]
	parsedTime, err := time.Parse(layout, s)
	if err != nil {
		return err
	}
	*t = Time(parsedTime)
	return nil
}
