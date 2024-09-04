package piyolog

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"
)

type (
	User struct {
		Name string
	}
	Log interface {
		Type() string
		Notes() string
		CreatedAt() time.Time
	}
	LogItem struct {
		typ       string
		content   string
		notes     string
		createdAt time.Time
	}
	FormulaLog struct {
		LogItem
		Amount string
	}
	SolidLog struct {
		LogItem
	}
)

type (
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
	reLog      = regexp.MustCompile(`^([0-9:]{5}) (AM|PM)? {1,}([^ ]+)(.*)`)
	reNotes    = regexp.MustCompile(``)
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
	matches := reLog.FindStringSubmatch(str)
	// set createdAt
	hm := strings.Split(matches[1], ":")
	h, _ := strconv.Atoi(hm[0])
	m, _ := strconv.Atoi(hm[1])
	if matches[2] == "PM" {
		h += 12
	}
	const notesSeparator = `   `
	list := strings.Split(matches[4], notesSeparator)
	i := LogItem{
		typ:       matches[3],
		content:   strings.TrimSpace(list[0]),
		notes:     strings.TrimSpace(strings.Join(list[1:], notesSeparator)),
		createdAt: time.Date(date.Year(), date.Month(), date.Day(), h, m, 0, 0, piyoLoc),
	}
	switch i.typ {
	case "ミルク", "Formula":
		return FormulaLog{
			LogItem: i,
			Amount:  i.content,
		}
	case "離乳食", "Solid":
		return SolidLog{
			LogItem: i,
		}
	}
	return i
}

func (i LogItem) Type() string {
	return i.typ
}

func (i LogItem) Notes() string {
	return i.notes
}

func (i LogItem) CreatedAt() time.Time {
	return i.createdAt
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
