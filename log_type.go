package piyolog

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/kaneshin/piyolog/piyologutil"
)

type Log interface {
	Type() string
	Content() string
	Notes() string
	CreatedAt() time.Time
}

type LogItem struct {
	typ       string
	content   string
	notes     string
	createdAt time.Time
}

// NewLog returns a log interface.
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

var reLog = regexp.MustCompile(`^([0-9:]{5} (AM|PM)?) +([^ ]+)(.*)`)

// NewLogItem returns a LogItem value.
func NewLogItem(date time.Time, str string) LogItem {
	matches := reLog.FindStringSubmatch(str)
	// set createdAt
	h, m := piyologutil.HourAndMinuteFromTimeString(matches[1])
	const notesSeparator = `   `
	list := strings.Split(matches[4], notesSeparator)
	return LogItem{
		typ:       matches[3],
		content:   strings.Trim(strings.TrimSpace(list[0]), "()"),
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

type FormulaLog struct {
	LogItem
	Amount int
	Unit   string
}

var reAmount = regexp.MustCompile(`^([0-9]+)(.+)$`)

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

type SolidLog struct {
	LogItem
}

// NewSolidLog returns a SolidLog value.
func NewSolidLog(i LogItem) SolidLog {
	return SolidLog{
		LogItem: i,
	}
}

type SleepLog struct {
	LogItem
}

// NewSleepLog returns a SleepLog value.
func NewSleepLog(i LogItem) SleepLog {
	return SleepLog{
		LogItem: i,
	}
}

type WakeUpLog struct {
	LogItem
	Duration time.Duration
}

// NewWakeUpLog returns a WakeUpLog value.
func NewWakeUpLog(i LogItem) WakeUpLog {
	return WakeUpLog{
		LogItem:  i,
		Duration: piyologutil.DurationFromDurationString(i.content),
	}
}

type BodyTemperatureLog struct {
	LogItem
	Temperature float64
	Unit        string
}

var reBodyTemperature = regexp.MustCompile(`([0-9\.]+)(.+)`)

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
