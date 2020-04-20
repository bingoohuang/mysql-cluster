package otter

import (
	"database/sql"
	"testing"

	"github.com/bingoohuang/gou/lang"
	"github.com/bingoohuang/sqlx"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"

	_ "github.com/mattn/go-sqlite3"
)

func TestPipelineDelayDao(t *testing.T) {
	logrus.SetLevel(logrus.DebugLevel)

	//ds := sqlx.CompatibleMySQLDs("localhost:3311 root/root db=otter")
	//more := sqlx.NewSQLMore("mysql", ds)
	//
	//sqlx.DB = more.Open()

	db, _ := sql.Open("sqlite3", "testdata/otter.db")
	sqlx.DB = db

	lastTime := lang.ParseTime("2006-01-02 15:04:05", "2006-01-02 15:04:05")
	records := Dao.DelayStats(lastTime)
	//assert.Equal(t, []PipelineDelayRecord{}, records)
	assert.True(t, len(records) > 0)

	logRecords := Dao.LogRecords(0)
	assert.True(t, len(logRecords) > 0)

	historyStats := Dao.HistoryStats(lastTime)
	assert.True(t, len(historyStats) > 0)

	tableStats := Dao.TableStats(lastTime)
	assert.True(t, len(tableStats) > 0)

	throughputStats := Dao.ThroughputStats(lastTime)
	assert.True(t, len(throughputStats) > 0)
}
