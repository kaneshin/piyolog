package piyolog

import (
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
)

func TestLogType(t *testing.T) {
	tests := []struct {
		date time.Time
		in   string
		out  Log
	}{
		{
			date: time.Date(2023, time.December, 31, 0, 0, 0, 0, piyoLoc),
			in:   `08:45 AM   ミルク 140ml   たくさん飲んだ`,
			out: FormulaLog{
				LogItem: LogItem{
					typ:       "ミルク",
					content:   "140ml",
					notes:     "たくさん飲んだ",
					createdAt: time.Date(2023, time.December, 31, 8, 45, 0, 0, piyoLoc),
				},
				Amount: 140,
				Unit:   "ml",
			},
		},
		{
			date: time.Date(2023, time.December, 31, 0, 0, 0, 0, piyoLoc),
			in:   `02:55   起きる (3時間35分)   `,
			out: WakeUpLog{
				LogItem: LogItem{
					typ:       "起きる",
					content:   "3時間35分",
					notes:     "",
					createdAt: time.Date(2023, time.December, 31, 2, 55, 0, 0, piyoLoc),
				},
				Duration: time.Duration(3)*time.Hour + time.Duration(35)*time.Minute,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.in, func(t *testing.T) {
			l := NewLog(tt.date, tt.in)
			if diff := cmp.Diff(tt.out, l, cmpopts.EquateComparable(LogItem{})); diff != "" {
				t.Errorf("%s", diff)
			}
		})
	}
}
