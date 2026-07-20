package tests

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/pavel/go-final-project/api"
)

type mask struct {
	dstart  time.Time
	repeat  string
	want    []int
	wantErr bool
}

func TestGetMaskWeek(t *testing.T) {

	m := []mask{
		{dstart: time.Date(2026, 6, 1, 0, 0, 0, 0, time.UTC), repeat: "2,4", want: []int{1, 3}, wantErr: false},
		{dstart: time.Date(2026, 6, 1, 0, 0, 0, 0, time.UTC), repeat: "1", want: []int{0}, wantErr: false},
		{dstart: time.Date(2026, 6, 1, 0, 0, 0, 0, time.UTC), repeat: "-1", want: []int{}, wantErr: true},
		{dstart: time.Date(2026, 6, 1, 0, 0, 0, 0, time.UTC), repeat: "1,4,", want: []int{}, wantErr: true},
		{dstart: time.Date(2026, 6, 1, 0, 0, 0, 0, time.UTC), repeat: "2,4,4,4", want: []int{1, 3}, wantErr: false},
		{dstart: time.Date(2026, 6, 1, 0, 0, 0, 0, time.UTC), repeat: "5,),+", want: []int{}, wantErr: true},
		{dstart: time.Date(2026, 6, 1, 0, 0, 0, 0, time.UTC), repeat: "5,6,7", want: []int{4, 5, 6}, wantErr: false},
		{dstart: time.Date(2026, 6, 1, 0, 0, 0, 0, time.UTC), repeat: "1,2,3,4,5,6,7", want: []int{0, 1, 2, 3, 4, 5, 6}, wantErr: false},
		{dstart: time.Date(2026, 6, 4, 0, 0, 0, 0, time.UTC), repeat: "1,2", want: []int{4, 5}, wantErr: false},
		{dstart: time.Date(2026, 6, 4, 0, 0, 0, 0, time.UTC), repeat: "3,4,6", want: []int{0, 2, 6}, wantErr: false},
	}

	for _, i := range m {
		v, err := api.GetMaskWeek(i.dstart, i.repeat)

		if i.wantErr {
			assert.Error(t, err)
			continue
		}

		assert.NoError(t, err)
		assert.Equal(t, i.want, v)
	}

}

func TestNextDateCut(t *testing.T) {
	tbl := []nextDate{
		{"20230126", "w 4,5", "20240201"},
	}
	check := func() {
		for _, v := range tbl {

			now, err := time.Parse("20060102", "20240126")
			if err != nil {
				continue
			}

			got, err := api.NextDate(now, v.date, v.repeat)
			if err != nil {
				t.Error(err)
				continue
			}
			assert.Equal(t, v.want, got)

		}
	}
	check()

}
