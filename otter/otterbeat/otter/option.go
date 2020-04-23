package otter

import (
	"database/sql"
	"fmt"
	"os"
	"reflect"
	"strconv"
	"time"

	"github.com/bingoohuang/gonet/man"

	"github.com/bingoohuang/gou/enc"

	"github.com/bingoohuang/gou/str"

	"github.com/bingoohuang/funk"
	"github.com/bingoohuang/gou/file"
	"github.com/bingoohuang/otterbeat/influx"
	"github.com/bingoohuang/sqlx"
	"github.com/sirupsen/logrus"
)

// Config defines the options for the otter beat.
// nolint lll
type Config struct {
	Version         bool   `pflag:"version. shorthand=v"`
	DSN             string `pflag:"dsn, eg. root:8BE4@127.0.0.1:9633/test. shorthand=d"`
	PipelineListURL string `pflag:"pipeline list page url, eg. http://127.0.0.1:2901/pipeline_list.htm?channelId=1. shorthand=p"`
	InfluxWriteURL  string `pflag:"influx writing url. eg. http://beta.isignet.cn:10014/write?db=metrics. shorthand=i"`
	PrintConfig     bool   `pflag:"print config before running. shorthand=P"`

	Interval time.Duration `pflag:"interval config before running(default 60s). shorthand=I"`
}

// SetUp sets the default value for config items.
func (c *Config) SetUp() {
	if c.PrintConfig {
		fmt.Printf("Config%s\n", enc.JSONPretty(c))
	}

	if c.Version {
		fmt.Println("version: v0.0.1")
		os.Exit(0)
	}

	if c.Interval == 0 {
		c.Interval = time.Minute
	}

	if c.InfluxWriteURL != "" {
		influx.PostMan.URL = man.URL(c.InfluxWriteURL)
	}

	if c.PipelineListURL != "" {
		PostMan.URL = man.URL(c.PipelineListURL)
	}
}

// Run runs  the otter beat by the config.
func (c *Config) Run() {
	ticker := time.NewTicker(c.Interval)
	defer ticker.Stop()

	if c.DSN != "" {
		db, err := sql.Open("mysql", sqlx.CompatibleMySQLDs(c.DSN))
		if err != nil {
			logrus.Errorf("failed to open mysql %s", c.DSN)
		}

		sqlx.DB = db
	}

	for range ticker.C {
		c.collectDB()
		c.collectPipelineListPage()
	}
}

func (c *Config) collectDB() {
	if sqlx.DB != nil {
		c.collectTables()
	}
}

func (c *Config) collectPipelineListPage() {
	if c.PipelineListURL == "" {
		return
	}

	list, err := GraspPipeLineList()
	if err != nil {
		logrus.Errorf("failed to GraspPipeLineList %s error %v", c.PipelineListURL, err)
		return
	}

	for _, l := range list {
		c.writeInfluxLine(l.ToPipeLineInflux(), 0)
	}
}

const (
	otterbeatDir = "~/.otterbeat/"
)

func (c *Config) collectTables() {
	c.timeRead(otterbeatDir+"DelayStats.lastTime", Dao.DelayStats,
		func(v interface{}) time.Time { return v.(PipelineDelay).ModifiedTime })
	c.timeRead(otterbeatDir+"HistoryStats.lastTime", Dao.HistoryStats,
		func(v interface{}) time.Time { return v.(TableHistoryStat).ModifiedTime })
	c.timeRead(otterbeatDir+"TableStats.lastTime", Dao.TableStats,
		func(v interface{}) time.Time { return v.(TableStat).ModifiedTime })
	c.timeRead(otterbeatDir+"TableStats.lastTime", Dao.ThroughputStats,
		func(v interface{}) time.Time { return v.(ThroughputStat).ModifiedTime })
	c.intRead(otterbeatDir+"LogRecords.lastID", Dao.LogRecords,
		func(v interface{}) uint64 { return v.(LogRecord).ID })
}

func (c *Config) intRead(filename string, dao interface{}, f func(interface{}) uint64) {
	lastID, err := file.ReadValue(filename, "0")
	if err != nil {
		logrus.Warnf("failed to load %s %v", filename, err)
		return
	}

	logrus.Infof("read %s with value %v", filename, lastID)

	last := str.ParseUint64(lastID)
	items := reflect.ValueOf(dao).Call([]reflect.Value{reflect.ValueOf(last)})[0].Interface()

	logrus.Infof("read %s got %d items %v", filename, funk.Len(items), funk.Left(items, 3))

	funk.ForEach(items, func(i int, v interface{}) {
		if x := f(v); x > last {
			last = x
		}

		c.writeInfluxLine(v, i)
	})

	if newID := strconv.FormatUint(last, 10); newID != lastID {
		logrus.Infof("write %s with value %v", filename, newID)

		if err := file.WriteValue(filename, newID); err != nil {
			logrus.Warnf("failed to write %s error %v", newID, err)
		}
	}
}

func (c *Config) timeRead(filename string, dao interface{}, f func(interface{}) time.Time) {
	lastTime, err := file.ReadTime(filename, StartTime)
	if err != nil {
		logrus.Warnf("failed to load %s %v", filename, err)
		return
	}

	logrus.Infof("read %s with value %v", filename, lastTime.Format(file.TimeFormat))

	items := reflect.ValueOf(dao).Call([]reflect.Value{reflect.ValueOf(lastTime)})[0].Interface()

	logrus.Infof("read %s got %d items %v", filename, funk.Len(items), funk.Left(items, 3))

	changed := false

	funk.ForEach(items, func(i int, v interface{}) {
		if l := f(v); l.After(lastTime) {
			lastTime = l
			changed = true
		}

		c.writeInfluxLine(v, i)
	})

	if !changed {
		return
	}

	logrus.Infof("write %s with value %v", filename, lastTime.Format(file.TimeFormat))

	if err := file.WriteTime(filename, lastTime); err != nil {
		logrus.Warnf("failed to write %s error %v", filename, err)
	}
}

func (c *Config) writeInfluxLine(v interface{}, i int) {
	line, err := influx.ToLine(v)
	if err != nil {
		logrus.Warnf("failed to influx  line %v error %v", v, err)
		return
	}

	// nolint gomnd
	if i < 3 {
		logrus.Infof("[InfluxDB Line] %s", line)
	} else if i == 4 {
		logrus.Infof("[InfluxDB Line] ...")
	}

	if c.InfluxWriteURL == "" {
		return
	}

	if err := influx.Write(line); err != nil {
		logrus.Warnf("failed to influx  write line %v error %v", v, err)
	}
}

const (
	// StartTime defines the start time of the system.
	StartTime = "2006-01-02 15:04:05"
)
