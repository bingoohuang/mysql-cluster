package influx

import (
	"testing"
	"time"

	"github.com/bingoohuang/gou/lang"
	"github.com/stretchr/testify/assert"
)

type delay struct {
	_            T         `measurement:"pipeline_delay"`
	ModifiedTime time.Time `influx:"time"`
	ID           uint64    `influx:"tag" name:"delay_id"`
	Delay        float64   `influx:"field"`
	Something    uint64    `influx:"-"`
	IncrID       uint64
	MyNote       string // default as influx field with snake name conversion.
}

func TestCreateLine(t *testing.T) {
	assert.Equal(t, lang.M2(
		`pipeline_delay,delay_id=1 delay=123456,incr_id=100,my_note="测试" 1587371416000000000`, nil),
		lang.M2(ToLine(
			delay{
				ModifiedTime: lang.ParseTime("2006-01-02 15:04:05", "2020-04-20 16:30:16"),
				Something:    332333,
				ID:           1,
				Delay:        123456,
				IncrID:       100,
				MyNote:       "测试",
			},
		)))
}
