package piyologutil

import (
	"regexp"
	"strconv"
	"strings"
	"time"
)

var reTime = regexp.MustCompile(`^([0-9:]{5}) ?(AM|PM)?$`)

func HourAndMinuteFromTime(str string) (int, int) {
	matches := reTime.FindStringSubmatch(str)
	hm := strings.Split(matches[1], ":")
	h := Atoi(hm[0])
	m := Atoi(hm[1])
	if len(matches) > 2 && matches[2] == "PM" {
		h += 12
	}
	return h, m
}

func DurationFromTime(str string) time.Duration {
	h, m := HourAndMinuteFromTime(str)
	return time.Duration(h)*time.Hour + time.Duration(m)*time.Minute
}

func Atoi(s string) int {
	v, err := strconv.ParseInt(s, 10, 0)
	if err != nil {
		return 0
	}
	return int(v)
}
