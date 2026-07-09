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
		{dstart: time.Date(2026, 6, 1, 0, 0, 0, 0, time.UTC), repeat: "-1", want: []int{0}, wantErr: true},
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
