package piyologutil

import (
	"strings"
	"time"
)

// ParseTime returns a time.Time value interpreted by the given string,
// such as "20:15" and "07:35 AM".
func ParseTime(str string) time.Time {
	if strings.ContainsAny(str, "AP") {
		// if str has "A" or "P" that indicates "AM" or "PM".
		t, _ := time.Parse("03:04 PM", str)
		return t
	}
	t, _ := time.Parse("15:04", str)
	return t
}

// ParseDuration returns a time.Duration value interpreted by the given string,
// such as "8時間15分", "7h40m" and "20m".
func ParseDuration(str string) time.Duration {
	str = strings.Replace(str, "時間", "h", 1)
	str = strings.Replace(str, "分", "m", 1)
	d, err := time.ParseDuration(str)
	if err != nil {
		return 0
	}
	return d
}
