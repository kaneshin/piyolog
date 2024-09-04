package piyolog

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"
)

var (
	reLog             = regexp.MustCompile(`^([0-9:]{5}) (AM|PM)? {1,}([^ ]+)(.*)`)
	reAmount          = regexp.MustCompile(`^([0-9]+)(.+)$`)
	reSleepLength     = regexp.MustCompile(`([0-9]+)[^0-9]*([0-9]+)`)
	reBodyTemperature = regexp.MustCompile(`([0-9\.]+)(.+)`)
)

type (
	LogItem struct {
		typ       string
		content   string
		notes     string
		createdAt time.Time
	}
	SleepLog struct {
		LogItem
	}
	WakeUpLog struct {
		LogItem
		SleepLength time.Duration
	}
	FormulaLog struct {
		LogItem
		Amount int
		Unit   string
	}
	SolidLog struct {
		LogItem
	}
	BodyTemperatureLog struct {
		LogItem
		Temperature float64
		Unit        string
	}
)

// NewLogItem returns a LogItem value.
func NewLogItem(date time.Time, str string) LogItem {
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
	return LogItem{
		typ:       matches[3],
		content:   strings.TrimSpace(list[0]),
		notes:     strings.TrimSpace(strings.Join(list[1:], notesSeparator)),
		createdAt: time.Date(date.Year(), date.Month(), date.Day(), h, m, 0, 0, piyoLoc),
	}
}

func (i LogItem) Type() string {
	return i.typ
}

func (i LogItem) Content() string {
	return i.content
}

func (i LogItem) Notes() string {
	return i.notes
}

func (i LogItem) CreatedAt() time.Time {
	return i.createdAt
}

func (i LogItem) String() string {
	return fmt.Sprintf("%s %s %s", i.createdAt.Format("15:04"), i.typ, i.content)
}

// NewFormulaLog returns a FormulaLog value.
func NewFormulaLog(i LogItem) FormulaLog {
	sm := reAmount.FindStringSubmatch(i.content)
	amount, _ := strconv.Atoi(sm[1])
	return FormulaLog{
		LogItem: i,
		Amount:  amount,
		Unit:    sm[2],
	}
}

// NewSolidLog returns a SolidLog value.
func NewSolidLog(i LogItem) SolidLog {
	return SolidLog{
		LogItem: i,
	}
}

// NewSleepLog returns a SleepLog value.
func NewSleepLog(i LogItem) SleepLog {
	return SleepLog{
		LogItem: i,
	}
}

// NewWakeUpLog returns a WakeUpLog value.
func NewWakeUpLog(i LogItem) WakeUpLog {
	sm := reSleepLength.FindStringSubmatch(i.content)
	h, _ := strconv.Atoi(sm[1])
	m, _ := strconv.Atoi(sm[2])
	return WakeUpLog{
		LogItem:     i,
		SleepLength: time.Duration(h)*time.Hour + time.Duration(m)*time.Minute,
	}
}

// NewBodyTemperatureLog returns a BodyTemperatureLog value.
func NewBodyTemperatureLog(i LogItem) BodyTemperatureLog {
	sm := reBodyTemperature.FindStringSubmatch(i.content)
	temp, _ := strconv.ParseFloat(sm[1], 64)
	return BodyTemperatureLog{
		LogItem:     i,
		Temperature: temp,
		Unit:        sm[2],
	}
}
