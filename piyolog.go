package piyolog

import (
	"bufio"
	"bytes"
	"fmt"
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
		Date time.Time
		User User
		Logs []Log
	}
	User struct {
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
	reUser   = regexp.MustCompile(`(.*) \([0-9]+(歳|y)[0-9]+(か月|m)[0-9]+(日|d)\)$`)
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

func (d Data) newEntry(str string) (e Entry) {
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
		return e
	}
	date, err := time.ParseInLocation(layout, matches[1], piyoLoc)
	if err != nil {
		return e
	}
	e.Date = date
	return e
}

func (d *Data) addEntry(e Entry) {
	if e.Date.IsZero() {
		return
	}
	d.Entries = append(d.Entries, e)
}

func (e *Entry) apply(line string) {
	if reUser.MatchString(line) {
		e.User = newUser(line)
	}
	if reLog.MatchString(line) {
		e.Logs = append(e.Logs, NewLog(e.Date, line))
	}
}

// newUser returns au user value retrieving from the given value.
func newUser(str string) (u User) {
	matches := reUser.FindStringSubmatch(str)
	u.Name = matches[1]
	return u
}

// Parse returns the Data value represented by the string.
// It accepts only export data from PiyoLog. Any other value returns an error.
func Parse(str string) (*Data, error) {
	// add one separator to the tail of the file to handle the string as monthly data.
	exportData := fmt.Sprintf("%s\n%s", str, piyologSeparator)
	buf := bytes.NewBufferString(exportData)
	scanner := bufio.NewScanner(buf)

	// first, parse the head of the file to detect its language.
	scanner.Scan()
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
			data.addEntry(entry)
			// prepare next entry.
			scanner.Scan()
			line := strings.TrimSpace(scanner.Text())
			entry = data.newEntry(line)
			continue
		}
		entry.apply(line)
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return &data, nil
}
