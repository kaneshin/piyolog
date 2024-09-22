package piyologutil

import (
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
)

func Test_Atoi(t *testing.T) {
	atoitests := []struct {
		in  string
		out int
	}{
		{"0", 0}, {"1", 1}, {"-1", -1},
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

func Test_HourAndMinute_Duration_FromTimeString(t *testing.T) {
	atoitests := []struct {
		in     string
		out1   int
		out2   int
		outdur time.Duration
	}{
		{"00:00", 0, 0, time.Duration(0)},
		{"11:30", 11, 30, time.Duration(11)*time.Hour + time.Duration(30)*time.Minute},
		{"10:25 PM", 22, 25, time.Duration(22)*time.Hour + time.Duration(25)*time.Minute},
		{"21:45", 21, 45, time.Duration(21)*time.Hour + time.Duration(45)*time.Minute},
	}

	for _, tt := range atoitests {
		t.Run(tt.in, func(t *testing.T) {
			out1, out2 := HourAndMinuteFromTimeString(tt.in)
			if diff := cmp.Diff(tt.out1, out1); diff != "" {
				t.Errorf("%s", diff)
			}
			if diff := cmp.Diff(tt.out2, out2); diff != "" {
				t.Errorf("%s", diff)
			}
			outdur := DurationFromTimeString(tt.in)
			if diff := cmp.Diff(tt.outdur, outdur); diff != "" {
				t.Errorf("%s", diff)
			}
		})
	}
}

func Test_HourAndMinute_Duration_FromDurationString(t *testing.T) {
	atoitests := []struct {
		in     string
		out1   int
		out2   int
		outdur time.Duration
	}{
		{"0h0m", 0, 0, time.Duration(0)},
		{"2h0m", 2, 0, time.Duration(2) * time.Hour},
		{"20m", 0, 20, time.Duration(20) * time.Minute},
		{"11h30m", 11, 30, time.Duration(11)*time.Hour + time.Duration(30)*time.Minute},
		{"10時間25分", 10, 25, time.Duration(10)*time.Hour + time.Duration(25)*time.Minute},
		{"21時間45分", 21, 45, time.Duration(21)*time.Hour + time.Duration(45)*time.Minute},
	}

	for _, tt := range atoitests {
		t.Run(tt.in, func(t *testing.T) {
			out1, out2 := HourAndMinuteFromDurationString(tt.in)
			if diff := cmp.Diff(tt.out1, out1); diff != "" {
				t.Errorf("%s", diff)
			}
			if diff := cmp.Diff(tt.out2, out2); diff != "" {
				t.Errorf("%s", diff)
			}
			outdur := DurationFromDurationString(tt.in)
			if diff := cmp.Diff(tt.outdur, outdur); diff != "" {
				t.Errorf("%s", diff)
			}
		})
	}
}
