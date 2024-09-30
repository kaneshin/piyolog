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

type (
	Data struct {
		Tag     language.Tag
		Entries []Entry
	}
	Entry struct {
		Date    time.Time
		Baby    Baby
		Logs    []Log
		Journal string
	}
	Baby struct {
		Name string
	}
)

const (
	piyologJa        = "【ぴよログ】"
	piyologEn        = "[PiyoLog]"
	piyologSeparator = "----------"
)

var (
	reDateJa = regexp.MustCompile(`^([0-9]{4}/[0-9]{1,2}/[0-9]{1,2})`)
	reDateEn = regexp.MustCompile(`^[a-zA-Z]{3}, (.*)$`)
	reBaby   = regexp.MustCompile(`(.*) \([0-9]+(歳|y)[0-9]+(か月|m)[0-9]+(日|d)\)$`)
	reLog    = regexp.MustCompile(`^([0-9:]{5} ?(AM|PM)?)   ([^ ]+)(.*)`)
)

var piyoLoc, _ = time.LoadLocation("Asia/Tokyo")

// SetLocation sets the location.
func SetLocation(loc *time.Location) {
	piyoLoc = loc
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
	var matches []string
	var layout string
	switch d.Tag {
	case language.Japanese:
		matches = reDateJa.FindStringSubmatch(str)
		layout = "2006/1/2"
	case language.English:
		matches = reDateEn.FindStringSubmatch(str)
		layout = "Jan 2, 2006"
	}
	if len(matches) <= 1 {
		return nil
	}
	date, err := time.ParseInLocation(layout, matches[1], piyoLoc)
	if err != nil {
		return nil
	}
	e := &Entry{
		Date: date,
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

func (e *Entry) apply(line string) {
	if reBaby.MatchString(line) {
		e.Baby = newBaby(line)
	}
	if reLog.MatchString(line) {
		e.Logs = append(e.Logs, NewLog(e.Date, line))
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
		if line == "" {
			continue
		}
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
