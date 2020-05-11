package otter

import (
	"time"

	"github.com/bingoohuang/sqlx"
)

// DelayStat maps to table DELAY_STAT record.
type DelayStat struct {
	ModifiedTime time.Time `influx:"time" name:"GMT_MODIFIED" measurement:"otter_delay_stat"`
	PipelineID   string    `influx:"tag"`   // 流水线ID
	DelayTime    float64   `influx:"field"` // 单位ms
	ID           uint64    `influx:"field"` // 对应的数据库表自增ID
}

// LogRecord maps to table LOG_RECORD record.
type LogRecord struct {
	ModifiedTime time.Time `influx:"time" name:"GMT_MODIFIED" measurement:"otter_log_record"`
	PipelineID   string    `influx:"tag"` // 流水线ID
	InsertCount  uint64    `influx:"field"`
	UpdateCount  uint64    `influx:"field"`
	DeleteCount  uint64    `influx:"field"`
	StartTime    time.Time `influx:"field"`
	EndTime      time.Time `influx:"field"`
	ID           uint64    `influx:"field"` // 对应的数据库表自增ID
}

// TableHistoryStat maps to table TABLE_HISTORY_STAT record.
type TableHistoryStat struct {
	ModifiedTime time.Time `influx:"time" name:"GMT_MODIFIED" measurement:"otter_history_stat"`
	PipelineID   string    `influx:"tag"` // 流水线ID
	InsertCount  uint64    `influx:"field"`
	UpdateCount  uint64    `influx:"field"`
	DeleteCount  uint64    `influx:"field"`
	StartTime    time.Time `influx:"field"`
	EndTime      time.Time `influx:"field"`
	ID           uint64    `influx:"field"` // 对应的数据库表自增ID
}

// TableStat maps to table TABLE_STAT record.
type TableStat struct {
	ModifiedTime time.Time `influx:"time" name:"GMT_MODIFIED" measurement:"otter_table_stat"`
	PipelineID   string    `influx:"tag"` // 流水线ID
	InsertCount  uint64    `influx:"field"`
	UpdateCount  uint64    `influx:"field"`
	DeleteCount  uint64    `influx:"field"`
	ID           uint64    `influx:"field"` // 对应的数据库表自增ID
}

// ThroughputStat maps to table THROUGHPUT_STAT record.
type ThroughputStat struct {
	ModifiedTime time.Time `influx:"time" name:"GMT_MODIFIED" measurement:"otter_throughput_stat"`
	PipelineID   string    `influx:"tag"` // 流水线ID
	TYPE         string    `influx:"field"`
	Number       uint64    `influx:"field"`
	Size         uint64    `influx:"field"`
	StartTime    time.Time `influx:"field"`
	EndTime      time.Time `influx:"field"`
	ID           uint64    `influx:"field"` // 对应的数据库表自增ID
}

// DaoFn defines the dao functions to operate on table DELAY_STAT.
// nolint lll
type DaoFn struct {
	DelayStat        func(lastTime time.Time) []DelayStat        `sql:"select ID, DELAY_TIME, PIPELINE_ID, GMT_CREATE, GMT_MODIFIED from DELAY_STAT where GMT_MODIFIED > :1"`
	LogRecord        func(lastID uint64) []LogRecord             `sql:"select ID, NID, CHANNEL_ID, PIPELINE_ID, TITLE, MESSAGE, GMT_MODIFIED from LOG_RECORD where ID > :1"`
	TableHistoryStat func(lastTime time.Time) []TableHistoryStat `sql:"select ID, INSERT_COUNT, UPDATE_COUNT, DELETE_COUNT, PIPELINE_ID, START_TIME, END_TIME, GMT_CREATE, GMT_MODIFIED from TABLE_HISTORY_STAT where GMT_MODIFIED > :1"`
	TableStat        func(lastTime time.Time) []TableStat        `sql:"select ID, INSERT_COUNT, UPDATE_COUNT, DELETE_COUNT, PIPELINE_ID, GMT_CREATE, GMT_MODIFIED from TABLE_STAT where GMT_MODIFIED > :1"`
	ThroughputStat   func(lastTime time.Time) []ThroughputStat   `sql:"select ID, TYPE, NUMBER, SIZE, PIPELINE_ID, START_TIME, END_TIME, GMT_CREATE, GMT_MODIFIED from THROUGHPUT_STAT where GMT_MODIFIED > :1"`
}

// nolint
var Dao = func() *DaoFn {
	dao := &DaoFn{}
	_ = sqlx.CreateDao(dao, sqlx.WithLogger(&sqlx.DaoLogrus{}))

	return dao
}()
