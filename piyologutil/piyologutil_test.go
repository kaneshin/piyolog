package piyologutil

import (
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
)

func Test_ParseTime(t *testing.T) {
	tm := time.Time{}.AddDate(-1, 0, 0)
	tests := []struct {
		in  string
		out time.Time
	}{
		{"00:00", tm},
		{"11:30", tm.Add(time.Duration(11)*time.Hour + time.Duration(30)*time.Minute)},
		{"10:25 PM", tm.Add(time.Duration(22)*time.Hour + time.Duration(25)*time.Minute)},
		{"21:45", tm.Add(time.Duration(21)*time.Hour + time.Duration(45)*time.Minute)},
	}

	for _, tt := range tests {
		t.Run(tt.in, func(t *testing.T) {
			out := ParseTime(tt.in)
			if diff := cmp.Diff(tt.out, out); diff != "" {
				t.Errorf("%s", diff)
			}
		})
	}
}

func Test_ParseDuration(t *testing.T) {
	tests := []struct {
		in  string
		out time.Duration
	}{
		{"0h0m", time.Duration(0)},
		{"2h0m", time.Duration(2) * time.Hour},
		{"20m", time.Duration(20) * time.Minute},
		{"11h30m", time.Duration(11)*time.Hour + time.Duration(30)*time.Minute},
		{"10時間25分", time.Duration(10)*time.Hour + time.Duration(25)*time.Minute},
		{"21時間45分", time.Duration(21)*time.Hour + time.Duration(45)*time.Minute},
	}

	for _, tt := range tests {
		t.Run(tt.in, func(t *testing.T) {
			out := ParseDuration(tt.in)
			if diff := cmp.Diff(tt.out, out); diff != "" {
				t.Errorf("%s", diff)
			}
		})
	}
}
