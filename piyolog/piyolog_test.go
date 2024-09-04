package piyolog

import (
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
)

func TestDaily(t *testing.T) {
	dailytests := []struct {
		in  string
		out *Daily
		err error
	}{
		{
			in: `【ぴよログ】2023/8/1(水)`,
			out: &Daily{
				Date: time.Date(2023, time.August, 1, 0, 0, 0, 0, piyoLoc),
				User: User{Name: ""},
			},
		},
		{
			in: `【ぴよログ】2023/12/31(水),
ごふあ (0歳1か月0日)

08:45 AM   ミルク 140ml   たくさん飲んだ
01:55 PM   寝る   
02:45 PM   起きる (0時間50分)   
03:05 PM   体温 36.4°C   
03:50 PM   ミルク 140ml   
07:35 PM   ミルク 200ml   `,
			out: &Daily{
				Date: time.Date(2023, time.December, 31, 0, 0, 0, 0, piyoLoc),
				User: User{Name: "ごふあ"},
				Logs: []Log{
					FormulaLog{
						LogItem: LogItem{
							typ:       "ミルク",
							content:   "140ml",
							notes:     "たくさん飲んだ",
							createdAt: time.Date(2023, time.December, 31, 8, 45, 0, 0, piyoLoc),
						},
						Amount: "140ml",
					},
					LogItem{
						typ:       "寝る",
						content:   "",
						createdAt: time.Date(2023, time.December, 31, 13, 55, 0, 0, piyoLoc),
					},
					LogItem{
						typ:       "起きる",
						content:   "(0時間50分)",
						createdAt: time.Date(2023, time.December, 31, 14, 45, 0, 0, piyoLoc),
					},
					LogItem{
						typ:       "体温",
						content:   "36.4°C",
						createdAt: time.Date(2023, time.December, 31, 15, 5, 0, 0, piyoLoc),
					},
					FormulaLog{
						LogItem: LogItem{
							typ:       "ミルク",
							content:   "140ml",
							createdAt: time.Date(2023, time.December, 31, 15, 50, 0, 0, piyoLoc),
						},
						Amount: "140ml",
					},
					FormulaLog{
						LogItem: LogItem{
							typ:       "ミルク",
							content:   "200ml",
							createdAt: time.Date(2023, time.December, 31, 19, 35, 0, 0, piyoLoc),
						},
						Amount: "200ml",
					},
				},
			},
		},
		{
			in: `【ぴよログ】2023/12/31(水),
ごふあ (0歳1か月0日)

05:05   ミルク 120ml   
08:10   おしっこ   
08:50   ミルク 120ml   
09:40   うんち   `,
			out: &Daily{
				Date: time.Date(2023, time.December, 31, 0, 0, 0, 0, piyoLoc),
				User: User{Name: "ごふあ"},
				Logs: []Log{
					FormulaLog{
						LogItem: LogItem{
							typ:       "ミルク",
							content:   "120ml",
							createdAt: time.Date(2023, time.December, 31, 5, 5, 0, 0, piyoLoc),
						},
						Amount: "120ml",
					},
					LogItem{
						typ:       "おしっこ",
						content:   "",
						createdAt: time.Date(2023, time.December, 31, 8, 10, 0, 0, piyoLoc),
					},
					FormulaLog{
						LogItem: LogItem{
							typ:       "ミルク",
							content:   "120ml",
							createdAt: time.Date(2023, time.December, 31, 8, 50, 0, 0, piyoLoc),
						},
						Amount: "120ml",
					},
					LogItem{
						typ:       "うんち",
						content:   "",
						createdAt: time.Date(2023, time.December, 31, 9, 40, 0, 0, piyoLoc),
					},
				},
			},
		},
		{
			in:  `ごふあ (0歳1か月0日)`,
			err: errMissingDate,
		},
	}

	for _, tt := range dailytests {
		t.Run(tt.in, func(t *testing.T) {
			daily, err := ParseDaily(tt.in)

			// error cases
			if err != nil {
				if diff := cmp.Diff(tt.err, err, cmpopts.EquateErrors()); diff != "" {
					t.Errorf("unexpected error returned: %s", diff)
				}
				return
			}

			// normal cases
			if tt.err != nil {
				t.Errorf("this error must be returned: %v", tt.err)
			}

			out := tt.out
			if diff := cmp.Diff(out.Date, daily.Date); diff != "" {
				t.Errorf("date parse failure: %s", diff)
			}
			if diff := cmp.Diff(out.User, daily.User); diff != "" {
				t.Errorf("date parse failure: %s", diff)
			}
			if diff := cmp.Diff(out.Logs, daily.Logs, cmpopts.EquateComparable(LogItem{})); diff != "" {
				t.Errorf("date parse failure: %s", diff)
			}
		})
	}
}

func TestMonthly(t *testing.T) {
	monthlytests := []struct {
		in  string
		out Monthly
		err error
	}{
		{
			in: `【ぴよログ】2024年8月
----------
2024/8/1(木)
ごふあ (0歳2か月10日)

04:15 AM   起きる (8時間40分)   
04:20 AM   ミルク 110ml   
08:00 PM   寝る   

母乳合計　　   左 0分 / 右 0分
ミルク合計　   7回 790ml
睡眠合計　　   12時間35分
おしっこ合計   3回
うんち合計　   1回

----------
2024/8/2(金)
ごふあ (0歳2か月11日)

04:15 AM   起きる (8時間40分)   
04:20 AM   ミルク 110ml   
08:00 PM   寝る   

母乳合計　　   左 0分 / 右 0分
ミルク合計　   8回 750ml
睡眠合計　　   13時間50分
おしっこ合計   4回
うんち合計　   1回

----------
2024/8/4(土)
ごふあ (0歳2か月12日)

04:15 AM   起きる (8時間40分)   
04:20 AM   ミルク 110ml   
08:00 PM   寝る   

母乳合計　　   左 0分 / 右 0分
ミルク合計　   7回 750ml
睡眠合計　　   14時間0分
おしっこ合計   2回
うんち合計　   0回

----------`,
			out: Monthly{
				Daily{
					Date: time.Date(2024, time.August, 1, 0, 0, 0, 0, piyoLoc),
					User: User{Name: "ごふあ"},
					Logs: []Log{
						LogItem{
							typ:       "起きる",
							content:   "(8時間40分)",
							createdAt: time.Date(2024, time.August, 1, 4, 15, 0, 0, piyoLoc),
						},
						FormulaLog{
							LogItem: LogItem{
								typ:       "ミルク",
								content:   "110ml",
								createdAt: time.Date(2024, time.August, 1, 4, 20, 0, 0, piyoLoc),
							},
							Amount: "110ml",
						},
						LogItem{
							typ:       "寝る",
							content:   "",
							createdAt: time.Date(2024, time.August, 1, 20, 0, 0, 0, piyoLoc),
						},
					},
				},
				Daily{
					Date: time.Date(2024, time.August, 2, 0, 0, 0, 0, piyoLoc),
					User: User{Name: "ごふあ"},
					Logs: []Log{
						LogItem{
							typ:       "起きる",
							content:   "(8時間40分)",
							createdAt: time.Date(2024, time.August, 2, 4, 15, 0, 0, piyoLoc),
						},
						FormulaLog{
							LogItem: LogItem{
								typ:       "ミルク",
								content:   "110ml",
								createdAt: time.Date(2024, time.August, 2, 4, 20, 0, 0, piyoLoc),
							},
							Amount: "110ml",
						},
						LogItem{
							typ:       "寝る",
							content:   "",
							createdAt: time.Date(2024, time.August, 2, 20, 0, 0, 0, piyoLoc),
						},
					},
				},
				Daily{
					Date: time.Date(2024, time.August, 4, 0, 0, 0, 0, piyoLoc),
					User: User{Name: "ごふあ"},
					Logs: []Log{
						LogItem{
							typ:       "起きる",
							content:   "(8時間40分)",
							createdAt: time.Date(2024, time.August, 4, 4, 15, 0, 0, piyoLoc),
						},
						FormulaLog{
							LogItem: LogItem{
								typ:       "ミルク",
								content:   "110ml",
								createdAt: time.Date(2024, time.August, 4, 4, 20, 0, 0, piyoLoc),
							},
							Amount: "110ml",
						},
						LogItem{
							typ:       "寝る",
							content:   "",
							createdAt: time.Date(2024, time.August, 4, 20, 0, 0, 0, piyoLoc),
						},
					},
				},
			},
		},
	}

	for _, tt := range monthlytests {
		t.Run(tt.in, func(t *testing.T) {
			monthly, err := ParseMonthly(tt.in)
			// error cases
			if err != nil {
				if diff := cmp.Diff(tt.err, err, cmpopts.EquateErrors()); diff != "" {
					t.Errorf("unexpected error returned: %s", diff)
				}
				return
			}

			// normal cases
			if tt.err != nil {
				t.Errorf("this error must be returned: %v", tt.err)
			}

			if len(tt.out) != len(monthly) {
				t.Errorf("wrong length: want %v, got %v", len(tt.out), len(monthly))
				return
			}
			for idx, daily := range monthly {
				if diff := cmp.Diff(tt.out[idx], daily, cmpopts.EquateComparable(LogItem{})); diff != "" {
					t.Errorf("%s", diff)
				}
			}
		})
	}
}
