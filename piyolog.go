package piyolog

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"regexp"
	"strconv"
	"strings"
	"time"

	"golang.org/x/text/language"
)

var piyoLoc, _ = time.LoadLocation("Asia/Tokyo")

// SetLocation sets the location.
func SetLocation(loc *time.Location) {
	piyoLoc = loc
}

type section int

const (
	sectionDate section = iota + 1
	sectionBaby
	sectionLogs
	sectionResults
	sectionJournal
	sectionEnd
)

func (s *section) next() section {
	switch *s {
	case sectionDate:
		*s = sectionBaby
	case sectionBaby:
		*s = sectionLogs
	case sectionLogs:
		*s = sectionResults
	case sectionResults:
		*s = sectionJournal
	case sectionJournal:
		*s = sectionEnd
	}
	return *s
}

func (s *section) end() {
	*s = sectionEnd
}

const (
	piyologJa        = "【ぴよログ】"
	piyologEn        = "[PiyoLog]"
	piyologSeparator = "----------"
)

type Data struct {
	Tag     language.Tag
	Entries []Entry
}

type Entry struct {
	section section
	Date    time.Time
	Baby    *Baby
	Logs    []Log
	Results []string
	Journal string
}

type Baby struct {
	Name        string
	DateOfBirth time.Time
}

func newData(str string) (d Data) {
	switch {
	case strings.Contains(str, piyologJa):
		d.Tag = language.Japanese
	case strings.Contains(str, piyologEn):
		d.Tag = language.English
	}
	return d
}

func (d Data) newEntry(str string) *Entry {
	if str == "" {
		return nil
	}
	var layout string
	switch d.Tag {
	case language.Japanese:
		str, _, _ = strings.Cut(str, "(")
		layout = "2006/1/2"
	case language.English:
		_, str, _ = strings.Cut(str, ", ")
		layout = "Jan 2, 2006"
	}
	date, err := time.ParseInLocation(layout, str, piyoLoc)
	if err != nil {
		return nil
	}
	e := &Entry{
		section: sectionDate,
		Date:    date,
	}
	return e
}

var reBaby = regexp.MustCompile(`^(.*) \(([0-9]+)(歳|y)([0-9]+)(か月|m)([0-9]+)(日|d)\)$`)

// newBaby returns au Baby value retrieving from the given value.
func (e Entry) newBaby(str string) *Baby {
	matches := reBaby.FindStringSubmatch(str)
	y, _ := strconv.Atoi(matches[2])
	m, _ := strconv.Atoi(matches[4])
	d, _ := strconv.Atoi(matches[6])
	return &Baby{
		Name:        matches[1],
		DateOfBirth: e.Date.AddDate(-y, -m, -d),
	}
}

var reLog = regexp.MustCompile(`^([0-9:]{5} ?(AM|PM)?)`)

func (e *Entry) apply(line string) {
	switch e.section {
	case sectionDate:
		e.section.next()
		e.apply(line)
	case sectionBaby:
		if line == "" {
			return
		}
		if reBaby.MatchString(line) {
			e.Baby = e.newBaby(line)
			e.section.next()
			return
		}
		// if text doesn't contain a certain baby infomation, move to the next section.
		e.section = sectionLogs
		e.apply(line)
	case sectionLogs:
		if line == "" && len(e.Logs) > 0 {
			e.section.next()
			return
		}
		if reLog.MatchString(line) {
			e.Logs = append(e.Logs, NewLog(e.Date, line))
			return
		}
	case sectionResults:
		if line == "" && len(e.Results) > 0 {
			e.section = sectionJournal
			return
		}
		e.Results = append(e.Results, line)
	case sectionJournal:
		if e.Journal == "" {
			e.Journal = line
		} else {
			e.Journal = fmt.Sprintf("%s\n%s", e.Journal, line)
		}
	}
}

// Parse returns the Data value represented by the string.
// It accepts only export data from PiyoLog. Any other value may return an error.
func Parse(str string) (*Data, error) {
	// replace escape line breaks with unescaped line breaks to be able to scan line by line.
	str = strings.Replace(str, `\n`, "\n", -1)
	// add one separator with TWO new lines to the tail of the file to handle the string as monthly data.
	exportData := fmt.Sprintf("%s\n\n%s\n", str, piyologSeparator)
	// replace "\n\n" (TWO new liens) with the separator with "\n" (one new line) with the separator
	// in order not to parse the new line before the separator.
	exportData = strings.ReplaceAll(
		exportData, fmt.Sprintf("\n\n%s", piyologSeparator), fmt.Sprintf("\n%s", piyologSeparator))

	scanner := bufio.NewScanner(bytes.NewBufferString(exportData))
	// first, parse the head of the file to detect its language.
	if !scanner.Scan() {
		return nil, io.EOF
	}
	head := strings.TrimSpace(scanner.Text())
	data := newData(head)
	switch data.Tag {
	case language.Japanese:
		head = strings.TrimLeft(head, piyologJa)
	case language.English:
		head = strings.TrimLeft(head, piyologEn)
	}

	// generate an entry with the head text.
	entry := data.newEntry(head)
	for scanner.Scan() {
		// handling the file as if monthly data.
		line := scanner.Text()
		if strings.HasPrefix(line, piyologSeparator) {
			if entry != nil {
				entry.section.end()
				data.Entries = append(data.Entries, *entry)
				entry = nil
			}
			continue
		}
		if entry == nil {
			entry = data.newEntry(line)
		} else {
			entry.apply(line)
		}
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return &data, nil
}
