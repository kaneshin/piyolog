package piyolog

import (
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"golang.org/x/text/language"
)

func Test_newData(t *testing.T) {
	tests := []struct {
		in  string
		out Data
	}{
		{
			in: `【ぴよログ】2022/6/13(木)`,
			out: Data{
				Tag: language.Japanese,
			},
		},
		{
			in: `【ぴよログ】2024年8月`,
			out: Data{
				Tag: language.Japanese,
			},
		},
		{
			in: `[PiyoLog]Thu, Jun 13, 2022`,
			out: Data{
				Tag: language.English,
			},
		},
		{
			in:  `Thu, Jun 13, 2022`,
			out: Data{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.in, func(t *testing.T) {
			data := newData(tt.in)
			if data.Tag != tt.out.Tag {
				t.Errorf("data parse failure: %s, %s", data.Tag, tt.out.Tag)
			}
		})
	}
}

func Test_newEntry(t *testing.T) {
	tests := []struct {
		data Data
		in   string
		out  *Entry
	}{
		{
			data: Data{Tag: language.Japanese},
			in:   `2022/6/13(木)`,
			out: &Entry{
				Date: time.Date(2022, time.June, 13, 0, 0, 0, 0, piyoLoc),
			},
		},
		{
			data: Data{Tag: language.Japanese},
			in:   `2024年8月`,
			out:  nil,
		},
		{
			data: Data{Tag: language.English},
			in:   `Thu, Jun 13, 2022`,
			out: &Entry{
				Date: time.Date(2022, time.June, 13, 0, 0, 0, 0, piyoLoc),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.in, func(t *testing.T) {
			entry := tt.data.newEntry(tt.in)
			if diff := cmp.Diff(tt.out, entry, cmpopts.IgnoreUnexported(Entry{})); diff != "" {
				t.Errorf("entry parse failure: %s", diff)
			}
		})
	}
}

func Test_newBaby(t *testing.T) {
	tests := []struct {
		in  string
		out Baby
	}{
		{
			in: `ごふあ (0歳0か月22日)`,
			out: Baby{
				Name: "ごふあ",
			},
		},
		{
			in: `ごふあ (0y0m22d)`,
			out: Baby{
				Name: "ごふあ",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.in, func(t *testing.T) {
			user := newBaby(tt.in)
			if diff := cmp.Diff(tt.out, user); diff != "" {
				t.Errorf("user parse failure: %s", diff)
			}
		})
	}
}

func Test_Parse(t *testing.T) {
	tests := []struct {
		in  string
		out Data
		err error
	}{
		{
			in:  ``,
			out: Data{},
		},
		{
			in: `【ぴよログ】2023/12/31(水)
ごふあ (0歳1か月0日)

08:45 AM   ミルク 140ml   たくさん飲んだ
01:55 PM   寝る   
02:45 PM   起きる (0時間50分)   
03:05 PM   体温 36.4°C   
03:50 PM   ミルク 140ml   
07:35 PM   ミルク 200ml   

母乳合計　　   左 7分 / 右 5分
ミルク合計　   7回 1140ml
睡眠合計　　   11時間50分
おしっこ合計   2回
うんち合計　   1回

お食い初めだよ


これは改行です



ここまで`,
			out: Data{
				Tag: language.Japanese,
				Entries: []Entry{
					Entry{
						Date: time.Date(2023, time.December, 31, 0, 0, 0, 0, piyoLoc),
						Baby: Baby{Name: "ごふあ"},
						Logs: []Log{
							FormulaLog{
								LogItem: LogItem{
									typ:       "ミルク",
									content:   "140ml",
									notes:     "たくさん飲んだ",
									createdAt: time.Date(2023, time.December, 31, 8, 45, 0, 0, piyoLoc),
								},
								Amount: 140,
								Unit:   "ml",
							},
							SleepLog{
								LogItem: LogItem{
									typ:       "寝る",
									content:   "",
									createdAt: time.Date(2023, time.December, 31, 13, 55, 0, 0, piyoLoc),
								},
							},
							WakeUpLog{
								LogItem: LogItem{
									typ:       "起きる",
									content:   "(0時間50分)",
									createdAt: time.Date(2023, time.December, 31, 14, 45, 0, 0, piyoLoc),
								},
								Duration: time.Duration(50) * time.Minute,
							},
							BodyTemperatureLog{
								LogItem: LogItem{
									typ:       "体温",
									content:   "36.4°C",
									createdAt: time.Date(2023, time.December, 31, 15, 5, 0, 0, piyoLoc),
								},
								Temperature: 36.4,
								Unit:        "°C",
							},
							FormulaLog{
								LogItem: LogItem{
									typ:       "ミルク",
									content:   "140ml",
									createdAt: time.Date(2023, time.December, 31, 15, 50, 0, 0, piyoLoc),
								},
								Amount: 140,
								Unit:   "ml",
							},
							FormulaLog{
								LogItem: LogItem{
									typ:       "ミルク",
									content:   "200ml",
									createdAt: time.Date(2023, time.December, 31, 19, 35, 0, 0, piyoLoc),
								},
								Amount: 200,
								Unit:   "ml",
							},
						},
					},
				},
			},
		},
		{
			in: `【ぴよログ】2023/12/31(水)\nごふあ (0歳1か月0日)\n\n\n08:45 AM   ミルク 140ml   たくさん飲んだ\n01:55 PM   寝る   \n02:45 PM   起きる (0時間50分)   \n03:05 PM   体温 36.4°C   \n03:50 PM   ミルク 140ml   \n07:35 PM   ミルク 200ml   `,
			out: Data{
				Tag: language.Japanese,
				Entries: []Entry{
					Entry{
						Date: time.Date(2023, time.December, 31, 0, 0, 0, 0, piyoLoc),
						Baby: Baby{Name: "ごふあ"},
						Logs: []Log{
							FormulaLog{
								LogItem: LogItem{
									typ:       "ミルク",
									content:   "140ml",
									notes:     "たくさん飲んだ",
									createdAt: time.Date(2023, time.December, 31, 8, 45, 0, 0, piyoLoc),
								},
								Amount: 140,
								Unit:   "ml",
							},
							SleepLog{
								LogItem: LogItem{
									typ:       "寝る",
									content:   "",
									createdAt: time.Date(2023, time.December, 31, 13, 55, 0, 0, piyoLoc),
								},
							},
							WakeUpLog{
								LogItem: LogItem{
									typ:       "起きる",
									content:   "(0時間50分)",
									createdAt: time.Date(2023, time.December, 31, 14, 45, 0, 0, piyoLoc),
								},
								Duration: time.Duration(50) * time.Minute,
							},
							BodyTemperatureLog{
								LogItem: LogItem{
									typ:       "体温",
									content:   "36.4°C",
									createdAt: time.Date(2023, time.December, 31, 15, 5, 0, 0, piyoLoc),
								},
								Temperature: 36.4,
								Unit:        "°C",
							},
							FormulaLog{
								LogItem: LogItem{
									typ:       "ミルク",
									content:   "140ml",
									createdAt: time.Date(2023, time.December, 31, 15, 50, 0, 0, piyoLoc),
								},
								Amount: 140,
								Unit:   "ml",
							},
							FormulaLog{
								LogItem: LogItem{
									typ:       "ミルク",
									content:   "200ml",
									createdAt: time.Date(2023, time.December, 31, 19, 35, 0, 0, piyoLoc),
								},
								Amount: 200,
								Unit:   "ml",
							},
						},
					},
				},
			},
		},
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

お食い初めだよ

----------`,
			out: Data{
				Tag: language.Japanese,
				Entries: []Entry{
					Entry{
						Date: time.Date(2024, time.August, 1, 0, 0, 0, 0, piyoLoc),
						Baby: Baby{Name: "ごふあ"},
						Logs: []Log{
							WakeUpLog{
								LogItem: LogItem{
									typ:       "起きる",
									content:   "(8時間40分)",
									createdAt: time.Date(2024, time.August, 1, 4, 15, 0, 0, piyoLoc),
								},
								Duration: time.Duration(8)*time.Hour + time.Duration(40)*time.Minute,
							},
							FormulaLog{
								LogItem: LogItem{
									typ:       "ミルク",
									content:   "110ml",
									createdAt: time.Date(2024, time.August, 1, 4, 20, 0, 0, piyoLoc),
								},
								Amount: 110,
								Unit:   "ml",
							},
							SleepLog{
								LogItem: LogItem{
									typ:       "寝る",
									content:   "",
									createdAt: time.Date(2024, time.August, 1, 20, 0, 0, 0, piyoLoc),
								},
							},
						},
					},
					Entry{
						Date: time.Date(2024, time.August, 2, 0, 0, 0, 0, piyoLoc),
						Baby: Baby{Name: "ごふあ"},
						Logs: []Log{
							WakeUpLog{
								LogItem: LogItem{
									typ:       "起きる",
									content:   "(8時間40分)",
									createdAt: time.Date(2024, time.August, 2, 4, 15, 0, 0, piyoLoc),
								},
								Duration: time.Duration(8)*time.Hour + time.Duration(40)*time.Minute,
							},
							FormulaLog{
								LogItem: LogItem{
									typ:       "ミルク",
									content:   "110ml",
									createdAt: time.Date(2024, time.August, 2, 4, 20, 0, 0, piyoLoc),
								},
								Amount: 110,
								Unit:   "ml",
							},
							SleepLog{
								LogItem: LogItem{
									typ:       "寝る",
									content:   "",
									createdAt: time.Date(2024, time.August, 2, 20, 0, 0, 0, piyoLoc),
								},
							},
						},
					},
					Entry{
						Date: time.Date(2024, time.August, 4, 0, 0, 0, 0, piyoLoc),
						Baby: Baby{Name: "ごふあ"},
						Logs: []Log{
							WakeUpLog{
								LogItem: LogItem{
									typ:       "起きる",
									content:   "(8時間40分)",
									createdAt: time.Date(2024, time.August, 4, 4, 15, 0, 0, piyoLoc),
								},
								Duration: time.Duration(8)*time.Hour + time.Duration(40)*time.Minute,
							},
							FormulaLog{
								LogItem: LogItem{
									typ:       "ミルク",
									content:   "110ml",
									createdAt: time.Date(2024, time.August, 4, 4, 20, 0, 0, piyoLoc),
								},
								Amount: 110,
								Unit:   "ml",
							},
							SleepLog{
								LogItem: LogItem{
									typ:       "寝る",
									content:   "",
									createdAt: time.Date(2024, time.August, 4, 20, 0, 0, 0, piyoLoc),
								},
							},
						},
					},
				},
			},
		},
		{
			in:  `ごふあ (0歳1か月0日)`,
			out: Data{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.in, func(t *testing.T) {
			data, err := Parse(tt.in)
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

			if len(tt.out.Entries) != len(data.Entries) {
				t.Errorf("wrong length: want %v, got %v", len(tt.out.Entries), len(data.Entries))
				return
			}
			for idx, entry := range data.Entries {
				if diff := cmp.Diff(tt.out.Entries[idx], entry, cmpopts.EquateComparable(LogItem{}), cmpopts.IgnoreUnexported(Entry{})); diff != "" {
					t.Errorf("%s", diff)
				}
			}
		})
	}
}
