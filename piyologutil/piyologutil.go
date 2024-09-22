package piyologutil

import (
	"regexp"
	"strconv"
	"strings"
	"time"
)

// Atoi interprets a string s and returns the corresponding integer value.
// If the given string s can't be parsed by the strconv.ParseInt, it returns 0.
func Atoi(s string) int {
	v, err := strconv.ParseInt(s, 10, 0)
	if err != nil {
		return 0
	}
	return int(v)
}

var reTime = regexp.MustCompile(`^([0-9:]{5}) ?(AM|PM)?$`)

// HourAndMinuteFromTimeString returns two integer values split by a ":" character
// from the given string, such as "20:15" and "07:35 AM". First return value
// represents its hours, second return value represents its minutes.
func HourAndMinuteFromTimeString(str string) (int, int) {
	matches := reTime.FindStringSubmatch(str)
	hm := strings.Split(matches[1], ":")
	h := Atoi(hm[0])
	m := Atoi(hm[1])
	if len(matches) > 2 && matches[2] == "PM" {
		h += 12
	}
	return h, m
}

// DurationFromTimeString returns a time.Duration value interpreted by the given string.
func DurationFromTimeString(str string) time.Duration {
	h, m := HourAndMinuteFromTimeString(str)
	return time.Duration(h)*time.Hour + time.Duration(m)*time.Minute
}

var reDuration = regexp.MustCompile(`^(([0-9]+)(時間|h))?([0-9]+)(分|m)$`)

// HourAndMinuteFromTime returns two integer values split by a ":" character
// from the given string, such as "8時間15分", "7h40m" and "20m". First return value
// represents its hours, second return value represents its minutes.
func HourAndMinuteFromDurationString(str string) (int, int) {
	matches := reDuration.FindStringSubmatch(str)
	if len(matches) < 6 {
		return 0, 0
	}
	return Atoi(matches[2]), Atoi(matches[4])
}

// DurationFromDurationString returns a time.Duration value interpreted by the given string.
func DurationFromDurationString(str string) time.Duration {
	h, m := HourAndMinuteFromDurationString(str)
	return time.Duration(h)*time.Hour + time.Duration(m)*time.Minute
}
