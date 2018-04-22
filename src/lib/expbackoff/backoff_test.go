package expbackoff

import (
	"testing"
	"time"
)

func backofftest(try, min, max int, duration time.Duration) time.Duration {
	end := EndingBounds(try, min, max, duration)
	d := RandomDuration(min, end, duration)
	return d
}

func Test_backofftest(t *testing.T) {
	type args struct {
		min      int
		max      int
		try      int
		duration time.Duration
	}
	tests := []struct {
		name     string
		args     args
		validate func(got time.Duration)
	}{
		{
			"try:0 with 0 to 10 using seconds returns zero",
			args{try: 0, min: 0, max: 10, duration: time.Second},
			func(got time.Duration) {
				if got != time.Duration(0) {
					t.Errorf("try 0 should be 0 seconds, got:%v", got)
				}
			},
		},
		{
			"try:1 with 0 to 10 using seconds returns between (0,1) second",
			args{try: 1, min: 0, max: 10, duration: time.Second},
			func(got time.Duration) {
				if got < time.Duration(0) || got > time.Second {
					t.Errorf("try 0 should be between (0,1) seconds, got:%v", got)
				}
			},
		},
		{
			"try:10 with 0 to 10 using seconds returns zero",
			args{try: 10, min: 0, max: 10, duration: time.Second},
			func(got time.Duration) {
				if got < time.Duration(0) || got > time.Second {
					t.Errorf("try 0 should be between (0,1) seconds, got:%v", got)
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := backofftest(tt.args.try, tt.args.min, tt.args.max, tt.args.duration)
			tt.validate(got)
			t.Logf("%v -->> got %v", tt.name, got)
		})
	}
}
