package piyolog

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"regexp"
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
)

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
	Baby    Baby
	Logs    []Log
	Results []string
	Journal string
}

type Baby struct {
	Name string
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

// addEntry append a Entry value if the Date of the Entry sets non zero value.
func (d *Data) addEntry(e Entry) {
	if e.Date.IsZero() {
		return
	}
	d.Entries = append(d.Entries, e)
}

var (
	reBaby = regexp.MustCompile(`^(.*) \([0-9]+(歳|y)[0-9]+(か月|m)[0-9]+(日|d)\)$`)
	reLog  = regexp.MustCompile(`^([0-9:]{5} ?(AM|PM)?)`)
)

func (e *Entry) apply(line string) {
	switch e.section {
	case sectionDate:
		e.section = sectionBaby
		e.Baby = newBaby(line)
		return
	case sectionBaby:
		if line == "" {
			e.section = sectionLogs
			return
		}
		// TODO: case sectionLogs:
	// 	if line == "" {
	// 		e.section = sectionResults
	// 		return
	// 	}
	// 	e.Logs = append(e.Logs, NewLog(e.Date, line))
	// 	return
	case sectionResults:
	case sectionJournal:
	}
	if reLog.MatchString(line) {
		e.Logs = append(e.Logs, NewLog(e.Date, line))
		return
	}
	// Result
	// Journal
	journal := e.Journal
	if journal != "" {
		e.Journal = fmt.Sprintf("%s\n%s", journal, line)
	}
}

// newBaby returns au Baby value retrieving from the given value.
func newBaby(str string) (b Baby) {
	matches := reBaby.FindStringSubmatch(str)
	b.Name = matches[1]
	return b
}

// Parse returns the Data value represented by the string.
// It accepts only export data from PiyoLog. Any other value returns an error.
func Parse(str string) (*Data, error) {
	// replace escape line breaks with unescaped line breaks to be able to scan line by line.
	str = strings.Replace(str, `\n`, "\n", -1)
	// add one separator to the tail of the file to handle the string as monthly data.
	exportData := fmt.Sprintf("%s\n%s\n", str, piyologSeparator)
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
		line := strings.TrimSpace(scanner.Text())
		if strings.HasPrefix(line, piyologSeparator) {
			if entry != nil {
				data.addEntry(*entry)
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
