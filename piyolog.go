package piyolog

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"regexp"
	"strings"
	"time"
)

type (
	User struct {
		Name string
	}
	Log interface {
		Type() string
		Content() string
		Notes() string
		CreatedAt() time.Time
	}
	Daily struct {
		Date time.Time
		User User
		Logs []Log
	}
	Monthly []Daily
)

var (
	piyoLoc, _ = time.LoadLocation("Asia/Tokyo")
	reDate     = regexp.MustCompile(`([0-9]{4})/([0-9]{1,2})/([0-9]{1,2})`)
	reUser     = regexp.MustCompile(`(.*) \([0-9]+歳[0-9]+か月[0-9]+日\)$`)
)

var (
	errMissingDate = errors.New("missing date")
)

// SetLocation sets the location.
func SetLocation(loc *time.Location) {
	piyoLoc = loc
}

// NewUser returns au user value retrieving from the given value.
func NewUser(str string) (u User) {
	matches := reUser.FindStringSubmatch(str)
	u.Name = matches[1]
	return u
}

// NewLog returns a log value.
func NewLog(date time.Time, str string) Log {
	i := NewLogItem(date, str)
	switch i.typ {
	case "ミルク", "Formula":
		return NewFormulaLog(i)
	case "離乳食", "Solid":
		return NewSolidLog(i)
	case "寝る", "Sleep":
		return NewSleepLog(i)
	case "起きる", "Wake-up":
		return NewWakeUpLog(i)
	case "体温", "Body Temp.":
		return NewBodyTemperatureLog(i)
	}
	return i
}

func parseDaily(daily *Daily, line string) error {
	if reDate.MatchString(line) {
		matches := reDate.FindStringSubmatch(line)
		dateStr := fmt.Sprintf("%s/%02s/%02s", matches[1], matches[2], matches[3])
		t, err := time.ParseInLocation("2006/01/02", dateStr, piyoLoc)
		if err != nil {
			return err
		}
		daily.Date = t
	}
	if reUser.MatchString(line) {
		daily.User = NewUser(line)
	}
	if reLog.MatchString(line) {
		daily.Logs = append(daily.Logs, NewLog(daily.Date, line))
	}
	return nil
}

// ParseDaily returns the daily value represented by the string. It accepts only export data from PiyoLog. Any other value returns an error.
func ParseDaily(str string) (*Daily, error) {
	daily := new(Daily)

	buf := bytes.NewBufferString(str)
	scanner := bufio.NewScanner(buf)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if err := parseDaily(daily, line); err != nil {
			return nil, err
		}
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	if daily.Date.IsZero() {
		return nil, errMissingDate
	}
	return daily, nil
}

// ParseMonthly returns the monthly value represented by the string. It accepts only export data from PiyoLog. Any other value returns an error.
func ParseMonthly(str string) (Monthly, error) {
	monthly := Monthly{}

	buf := bytes.NewBufferString(str)
	scanner := bufio.NewScanner(buf)
	var daily Daily
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if strings.HasPrefix(line, "-----") {
			if !daily.Date.IsZero() {
				monthly = append(monthly, daily)
			}
			daily = Daily{}
		}
		if err := parseDaily(&daily, line); err != nil {
			return nil, err
		}
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return monthly, nil
}
