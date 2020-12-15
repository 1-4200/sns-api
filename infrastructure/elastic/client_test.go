package elastic

import (
	"reflect"
	"testing"
	"time"
)

func Test_convertTime(t *testing.T) {
	type args struct {
		timeStr string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := convertTime(tt.args.timeStr)
			if (err != nil) != tt.wantErr {
				t.Errorf("convertTime() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("convertTime() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_timeDiff(t *testing.T) {
	type args struct {
		t1 time.Time
		t2 time.Time
	}
	tests := []struct {
		name  string
		args  args
		want  int
		want1 int
		want2 int
		want3 int
		want4 int
		want5 int
	}{
		{
			name: "年を跨ぐケース",
			args: args{
				t1: time.Date(2015, 5, 1, 0, 0, 0, 0, time.UTC),
				t2: time.Date(2016, 6, 2, 1, 1, 1, 1, time.UTC),
			},
			want:  1,
			want1: 1,
			want2: 1,
			want3: 1,
			want4: 1,
			want5: 1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1, got2, got3, got4, got5 := timeDiff(tt.args.t1, tt.args.t2)
			if got != tt.want {
				t.Errorf("timeDiff() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("timeDiff() got1 = %v, want %v", got1, tt.want1)
			}
			if got2 != tt.want2 {
				t.Errorf("timeDiff() got2 = %v, want %v", got2, tt.want2)
			}
			if got3 != tt.want3 {
				t.Errorf("timeDiff() got3 = %v, want %v", got3, tt.want3)
			}
			if got4 != tt.want4 {
				t.Errorf("timeDiff() got4 = %v, want %v", got4, tt.want4)
			}
			if got5 != tt.want5 {
				t.Errorf("timeDiff() got5 = %v, want %v", got5, tt.want5)
			}
		})
	}
}

func Test_monthDiff(t *testing.T) {
	type args struct {
		t1 time.Time
		t2 time.Time
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{
			name: "年を跨がないケース",
			args: args{
				t1: time.Date(2016, 5, 1, 0, 0, 0, 0, time.UTC),
				t2: time.Date(2016, 6, 2, 1, 1, 1, 1, time.UTC),
			},
			want: 1,
		},
		{
			name: "年を跨ぐケース",
			args: args{
				t1: time.Date(2015, 5, 1, 0, 0, 0, 0, time.UTC),
				t2: time.Date(2016, 6, 2, 1, 1, 1, 1, time.UTC),
			},
			want: 13,
		},
		{
			name: "同年同月のケース",
			args: args{
				t1: time.Date(2016, 6, 1, 0, 0, 0, 0, time.UTC),
				t2: time.Date(2016, 6, 2, 1, 1, 1, 1, time.UTC),
			},
			want: 0,
		},
		{
			name: "年を跨ぐケース（t1 < t2）",
			args: args{
				t1: time.Date(2016, 6, 1, 0, 0, 0, 0, time.UTC),
				t2: time.Date(2015, 5, 1, 0, 0, 0, 0, time.UTC),
			},
			want: 13,
		},
		{
			name: "年を跨がないケース2",
			args: args{
				t1: time.Date(2020, 5, 19, 0, 0, 0, 0, time.UTC),
				t2: time.Date(2020, 8, 19, 23, 59, 0, 0, time.UTC),
			},
			want: 3,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := monthDiff(tt.args.t1, tt.args.t2); got != tt.want {
				t.Errorf("monthDiff() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_buildIndexByTimeAdd(t *testing.T) {
	type args struct {
		t   time.Time
		num int
	}
	tests := []struct {
		name string
		args args
		want []string
	}{
		{
			name: "年を跨ぐケース",
			args: args{
				t:   time.Date(2016, 1, 1, 0, 0, 0, 0, time.UTC),
				num: 2,
			},
			want: []string{"sns-2016.01", "sns-2016.02", "sns-2016.03"},
		},
		{
			name: "年を跨ぐケース",
			args: args{
				t:   time.Date(2016, 11, 1, 0, 0, 0, 0, time.UTC),
				num: 3,
			},
			want: []string{"sns-2016.11", "sns-2016.12", "sns-2017.01", "sns-2017.02"},
		},
		{
			name: "同年同月のケース",
			args: args{
				t:   time.Date(2016, 9, 1, 0, 0, 0, 0, time.UTC),
				num: 0,
			},
			want: []string{"sns-2016.09"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := buildIndexByTimeAdd(tweetIndex, tt.args.t, tt.args.num); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("buildIndexByTimeAdd() = %v, want %v", got, tt.want)
			}
		})
	}
}
