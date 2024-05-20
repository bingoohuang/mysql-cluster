package otter

import (
	"database/sql"
	"fmt"
	"github.com/bingoohuang/rotatefile"
	"log"
	"testing"
	"time"

	"github.com/bingoohuang/gou/file"

	"github.com/bingoohuang/otterbeat/influx"

	"github.com/bingoohuang/gou/lang"
	"github.com/bingoohuang/sqlx"
	"github.com/stretchr/testify/assert"

	_ "github.com/mattn/go-sqlite3"
)

// nolint
func TestPipelineDelayDao(t *testing.T) {
	log.SetOutput(rotatefile.New())

	//ds := sqlx.CompatibleMySQLDs("localhost:3311 root/root db=otter")
	//more := sqlx.NewSQLMore("mysql", ds)
	//
	//sqlx.DB = more.Open()

	db, _ := sql.Open("sqlite3", "testdata/otter.db")
	sqlx.DB = db

	lastTime := lang.ParseTime(file.TimeFormat, StartTime)
	delayStats := Dao.DelayStat(lastTime)
	//assert.Equal(t, []PipelineDelayRecord{}, delayStats)
	assert.True(t, len(delayStats) > 0)

	logRecords := Dao.LogRecord(0)
	assert.True(t, len(logRecords) > 0)

	historyStats := Dao.TableHistoryStat(lastTime)
	assert.True(t, len(historyStats) > 0)

	tableStats := Dao.TableStat(lastTime)
	assert.True(t, len(tableStats) > 0)

	throughputStats := Dao.ThroughputStat(lastTime)
	assert.True(t, len(throughputStats) > 0)

	exceptions := Dao.Exceptions(time.Now().Add(-10*365*24*time.Hour), time.Now())
	assert.True(t, len(exceptions) > 0)

	tableHistoryStats := Dao.TableHistoryStatByPipeline(time.Now().Add(-10*365*24*time.Hour), time.Now())
	assert.True(t, len(tableHistoryStats) > 0)

	//const influxDBAddr = `http://beta.isignet.cn:10014/write?db=metrics`
	//assert.Nil(t, influx.Write(influxDBAddr, l))

	for _, r := range delayStats {
		l, err := influx.ToLine(r)
		assert.Nil(t, err)
		fmt.Println(l)
	}

	for _, r := range logRecords {
		l, err := influx.ToLine(r)
		assert.Nil(t, err)
		fmt.Println(l)
	}

	for _, r := range historyStats {
		l, err := influx.ToLine(r)
		assert.Nil(t, err)
		fmt.Println(l)
	}

	for _, r := range tableStats {
		l, err := influx.ToLine(r)
		assert.Nil(t, err)
		fmt.Println(l)

	}

	for _, r := range throughputStats {
		l, err := influx.ToLine(r)
		assert.Nil(t, err)
		fmt.Println(l)
	}
}
