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

// NewLog returns a log interface.
func NewLog(str string, date time.Time) Log {
	tm, typ, content, notes := SplitLog(str)
	if tm.IsZero() {
		return nil
	}
	createdAt := time.Date(date.Year(), date.Month(), date.Day(),
		tm.Hour(), tm.Minute(), 0, 0, piyoLoc)
	return NewLogItem(typ, content, notes, createdAt).Log()
}

const logSeparator = `   `

// SplitLog splits a given str divided by the logSeparator, separating it into
// a time, type, content and notes.
func SplitLog(str string) (time.Time, string, string, string) {
	split := strings.Split(str, logSeparator)
	if len(split) < 3 {
		return time.Time{}, "", "", ""
	}
	tm := piyologutil.ParseTime(split[0])
	fields := strings.Fields(split[1])
	return tm,
		fields[0],
		strings.Join(fields[1:], ` `),
		strings.Join(split[2:], logSeparator)
}

type LogItem struct {
	typ       string
	content   string
	notes     string
	createdAt time.Time
}

// NewLogItem returns a LogItem value.
func NewLogItem(typ, content, notes string, createdAt time.Time) LogItem {
	return LogItem{
		typ:       typ,
		content:   content,
		notes:     notes,
		createdAt: createdAt,
	}
}

func (i LogItem) Log() Log {
	switch i.typ {
	case "母乳", "Nursing":
		return NewNursingLog(i)
	case "ミルク", "Formula":
		return NewFormulaLog(i)
	case "離乳食", "Solid":
		return NewSolidLog(i)
	case "寝る", "Sleep":
		return NewSleepLog(i)
	case "起きる", "Wake-up":
		return NewWakeUpLog(i)
	case "おしっこ", "Pee":
		return NewPeeLog(i)
	case "うんち", "Poop":
		return NewPoopLog(i)
	case "お風呂", "Baths":
		return NewBathsLog(i)
	case "体温", "Body Temp.":
		return NewBodyTemperatureLog(i)
	}
	return i
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

var reAmount = regexp.MustCompile(`^([0-9]+)(.+)$`)

func amountAndUnit(str string) (int, string) {
	sm := reAmount.FindStringSubmatch(str)
	amount, _ := strconv.Atoi(sm[1])
	return amount, sm[2]
}

type NursingLog struct {
	LogItem
	Left   time.Duration // TODO
	Right  time.Duration // TODO
	Amount int
	Unit   string
}

// NewNursingLog returns a NursingLog value.
func NewNursingLog(i LogItem) NursingLog {
	f := strings.Fields(i.content)
	amount, unit := amountAndUnit(strings.Trim(f[len(f)-1], "()"))
	return NursingLog{
		LogItem: i,
		Amount:  amount,
		Unit:    unit,
	}
}

type FormulaLog struct {
	LogItem
	Amount int
	Unit   string
}

// NewFormulaLog returns a FormulaLog value.
func NewFormulaLog(i LogItem) FormulaLog {
	amount, unit := amountAndUnit(i.content)
	return FormulaLog{
		LogItem: i,
		Amount:  amount,
		Unit:    unit,
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
	content := strings.Trim(i.content, "()")
	return WakeUpLog{
		LogItem:  i,
		Duration: piyologutil.ParseDuration(content),
	}
}

type PeeLog struct {
	LogItem
}

// NewPeeLog returns a PeeLog value.
func NewPeeLog(i LogItem) PeeLog {
	return PeeLog{
		LogItem: i,
	}
}

type PoopLog struct {
	LogItem
}

// NewPoopLog returns a PoopLog value.
func NewPoopLog(i LogItem) PoopLog {
	return PoopLog{
		LogItem: i,
	}
}

type BathsLog struct {
	LogItem
}

// NewBathsLog returns a BathsLog value.
func NewBathsLog(i LogItem) BathsLog {
	return BathsLog{
		LogItem: i,
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
