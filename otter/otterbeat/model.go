package otterbeat

import (
	"time"

	"github.com/bingoohuang/otterbeat/influx"
)

// OtterPipelineDelay defines the result of TCP/HTTP check.
type OtterPipelineDelay struct {
	_            influx.T  `measurement:"otter_pipeline_delay"`
	ModifiedTime time.Time `influxTime:"true"`
	PipelineID   string    `influxTag:"true"`   // 流水线ID
	DelayTime    float64   `influxField:"true"` // 单位s
	IncrID       uint64    `influxField:"true"` // 对应的数据库表自增ID
}
