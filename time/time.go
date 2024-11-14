package time

import (
	"time"
)

func GetCurrentDate() (*time.Time, error) {
	t := time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day(), time.Now().Hour(), time.Now().Minute(), time.Now().Second(), 0, time.Local)
	return &t, nil
}
