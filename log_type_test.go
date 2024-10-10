package piyolog

import (
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
)

func Test_Log(t *testing.T) {
	date := time.Date(2023, time.December, 31, 0, 0, 0, 0, piyoLoc)
	var createdAt = func(h, m int) time.Time {
		return date.Add(time.Duration(h)*time.Hour + time.Duration(m)*time.Minute)
	}
	tests := []struct {
		in  string
		out Log
	}{
		{
			in: `23:00   母乳 左 7分 / 右 5分 (50ml)   たくさん飲んだ`,
			out: NursingLog{
				LogItem: LogItem{
					typ:       "母乳",
					content:   "左 7分 / 右 5分 (50ml)",
					notes:     "たくさん飲んだ",
					createdAt: createdAt(23, 00),
				},
				Amount: 50,
				Unit:   "ml",
			},
		}, {
			in: `08:45 AM   ミルク 140ml   たくさん    飲んだ`,
			out: FormulaLog{
				LogItem: LogItem{
					typ:       "ミルク",
					content:   "140ml",
					notes:     "たくさん    飲んだ",
					createdAt: createdAt(8, 45),
				},
				Amount: 140,
				Unit:   "ml",
			},
		}, {
			in: `23:10   離乳食   たくさん食べた`,
			out: SolidLog{
				LogItem: LogItem{
					typ:       "離乳食",
					content:   "",
					notes:     "たくさん食べた",
					createdAt: createdAt(23, 10),
				},
			},
		}, {
			in: `02:55   起きる (3時間35分)   `,
			out: WakeUpLog{
				LogItem: LogItem{
					typ:       "起きる",
					content:   "(3時間35分)",
					notes:     "",
					createdAt: createdAt(2, 55),
				},
				Duration: time.Duration(3)*time.Hour + time.Duration(35)*time.Minute,
			},
		}, {
			in: `08:00 PM   寝る   `,
			out: SleepLog{
				LogItem: LogItem{
					typ:       "寝る",
					content:   "",
					notes:     "",
					createdAt: createdAt(20, 0),
				},
			},
		}, {
			in: `06:40   おしっこ   `,
			out: PeeLog{
				LogItem{
					typ:       "おしっこ",
					content:   "",
					notes:     "",
					createdAt: createdAt(6, 40),
				},
			},
		}, {
			in: `23:15   うんち (少なめ/ふつう/緑)   たくさん出た`,
			out: PoopLog{
				LogItem{
					typ:       "うんち",
					content:   "(少なめ/ふつう/緑)",
					notes:     "たくさん出た",
					createdAt: createdAt(23, 15),
				},
			},
		}, {
			in: `19:10   お風呂   `,
			out: BathsLog{
				LogItem{
					typ:       "お風呂",
					content:   "",
					notes:     "",
					createdAt: createdAt(19, 10),
				},
			},
		}, {
			in: `14:30   体温 36.5°C   `,
			out: BodyTemperatureLog{
				LogItem: LogItem{
					typ:       "体温",
					content:   "36.5°C",
					notes:     "",
					createdAt: createdAt(14, 30),
				},
				Temperature: 36.5,
				Unit:        "°C",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.in, func(t *testing.T) {
			lg := NewLog(tt.in, date)
			if diff := cmp.Diff(tt.out, lg, cmpopts.EquateComparable(LogItem{})); diff != "" {
				t.Errorf("log parse failure: %s", diff)
			}
		})
	}
}
