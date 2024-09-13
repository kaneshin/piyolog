package piyologutil

import (
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
)

func Test_HourAndMinuteFromTime(t *testing.T) {
	atoitests := []struct {
		in   string
		out1 int
		out2 int
	}{
		{"00:00", 0, 0},
		{"11:30", 11, 30},
		{"10:25 PM", 22, 25},
		{"21:45", 21, 45},
	}

	for _, tt := range atoitests {
		t.Run(tt.in, func(t *testing.T) {
			out1, out2 := HourAndMinuteFromTime(tt.in)
			if diff := cmp.Diff(tt.out1, out1); diff != "" {
				t.Errorf("%s", diff)
			}
			if diff := cmp.Diff(tt.out2, out2); diff != "" {
				t.Errorf("%s", diff)
			}
		})
	}
}

func Test_DurationFromTime(t *testing.T) {
	atoitests := []struct {
		in  string
		out time.Duration
	}{
		{"00:00", time.Duration(0)},
		{"11:30", time.Duration(11)*time.Hour + time.Duration(30)*time.Minute},
		{"10:25 PM", time.Duration(22)*time.Hour + time.Duration(25)*time.Minute},
		{"21:45", time.Duration(21)*time.Hour + time.Duration(45)*time.Minute},
	}

	for _, tt := range atoitests {
		t.Run(tt.in, func(t *testing.T) {
			out := DurationFromTime(tt.in)
			if diff := cmp.Diff(tt.out, out); diff != "" {
				t.Errorf("%s", diff)
			}
		})
	}
}

func Test_Atoi(t *testing.T) {
	atoitests := []struct {
		in  string
		out int
	}{
		{"0", 0},
		{"1", 1},
		{"-1", -1},
	}

	for _, tt := range atoitests {
		t.Run(tt.in, func(t *testing.T) {
			out := Atoi(tt.in)
			if diff := cmp.Diff(tt.out, out); diff != "" {
				t.Errorf("%s", diff)
			}
		})
	}
}
