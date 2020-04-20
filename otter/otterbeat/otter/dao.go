package otter

import (
	"time"

	"github.com/bingoohuang/sqlx"

	"github.com/bingoohuang/otterbeat/influx"
)

// PipelineDelay maps to table DELAY_STAT record.
type PipelineDelay struct {
	_ influx.T `measurement:"otter_delay_stat"`

	ModifiedTime time.Time `influx:"time" name:"GMT_MODIFIED"`
	PipelineID   string    `influx:"tag"`             // 流水线ID
	DelayTime    float64   `influx:"field"`           // 单位ms
	ID           uint64    `influx:"field" name:"ID"` // 对应的数据库表自增ID
}

// LogRecord maps to table LOG_RECORD record.
type LogRecord struct {
	_ influx.T `measurement:"otter_log_record"`

	ModifiedTime time.Time `influx:"time" name:"GMT_MODIFIED"`
	PipelineID   uint64    `influx:"tag"` // 流水线ID
	InsertCount  uint64    `influx:"field"`
	UpdateCount  uint64    `influx:"field"`
	DeleteCount  uint64    `influx:"field"`
	StartTime    time.Time `influx:"field"`
	EndTime      time.Time `influx:"field"`
	ID           uint64    `influx:"field" name:"ID"` // 对应的数据库表自增ID
}

// TableHistoryStat maps to table TABLE_HISTORY_STAT record.
type TableHistoryStat struct {
	_ influx.T `measurement:"otter_history_stat"`

	ModifiedTime time.Time `influx:"time" name:"GMT_MODIFIED"`
	PipelineID   uint64    `influx:"tag"` // 流水线ID
	InsertCount  uint64    `influx:"field"`
	UpdateCount  uint64    `influx:"field"`
	DeleteCount  uint64    `influx:"field"`
	StartTime    time.Time `influx:"field"`
	EndTime      time.Time `influx:"field"`
	ID           uint64    `influx:"field" name:"ID"` // 对应的数据库表自增ID
}

// TableStat maps to table TABLE_STAT record.
type TableStat struct {
	_ influx.T `measurement:"otter_table_stat"`

	ModifiedTime time.Time `influx:"time" name:"GMT_MODIFIED"`
	PipelineID   uint64    `influx:"tag"` // 流水线ID
	InsertCount  uint64    `influx:"field"`
	UpdateCount  uint64    `influx:"field"`
	DeleteCount  uint64    `influx:"field"`
	ID           uint64    `influx:"field" name:"ID"` // 对应的数据库表自增ID
}

// ThroughputStat maps to table THROUGHPUT_STAT record.
type ThroughputStat struct {
	_ influx.T `measurement:"otter_throughput_stat"`

	ModifiedTime time.Time `influx:"time" name:"GMT_MODIFIED"`
	PipelineID   uint64    `influx:"tag"` // 流水线ID
	TYPE         string    `influx:"field"`
	Number       uint64    `influx:"field"`
	Size         uint64    `influx:"field"`
	StartTime    time.Time `influx:"field"`
	EndTime      time.Time `influx:"field"`
	ID           uint64    `influx:"field" name:"ID"` // 对应的数据库表自增ID
}

// DaoFn defines the dao functions to operate on table DELAY_STAT.
// nolint lll
type DaoFn struct {
	DelayStats      func(lastTime time.Time) []PipelineDelay    `sql:"select ID, DELAY_TIME, PIPELINE_ID, GMT_CREATE, GMT_MODIFIED from DELAY_STAT where GMT_MODIFIED > :1"`
	LogRecords      func(lastID uint64) []LogRecord             `sql:"select ID, NID, CHANNEL_ID, PIPELINE_ID, TITLE, MESSAGE, GMT_MODIFIED from LOG_RECORD where ID > :1"`
	HistoryStats    func(lastTime time.Time) []TableHistoryStat `sql:"select ID, INSERT_COUNT, UPDATE_COUNT, DELETE_COUNT, PIPELINE_ID, START_TIME, END_TIME, GMT_CREATE, GMT_MODIFIED from TABLE_HISTORY_STAT where GMT_MODIFIED > :1"`
	TableStats      func(lastTime time.Time) []TableStat        `sql:"select ID, INSERT_COUNT, UPDATE_COUNT, DELETE_COUNT, PIPELINE_ID, GMT_CREATE, GMT_MODIFIED from TABLE_STAT where GMT_MODIFIED > :1"`
	ThroughputStats func(lastTime time.Time) []ThroughputStat   `sql:"select ID, TYPE, NUMBER, SIZE, PIPELINE_ID, START_TIME, END_TIME, GMT_CREATE, GMT_MODIFIED from THROUGHPUT_STAT where GMT_MODIFIED > :1"`
}

// nolint
var Dao = func() *DaoFn {
	dao := &DaoFn{}
	_ = sqlx.CreateDao(dao, sqlx.WithLogger(&sqlx.DaoLogrus{}))

	return dao
}()
