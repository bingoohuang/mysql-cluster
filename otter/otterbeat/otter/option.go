package otter

import (
	"database/sql"
	"io/ioutil"
	"os"
	"path/filepath"
	"reflect"
	"strconv"
	"time"

	"github.com/bingoohuang/gou/str"

	"github.com/bingoohuang/gou/file"
	"github.com/bingoohuang/gou/lang"
	"github.com/bingoohuang/otterbeat/influx"
	"github.com/bingoohuang/sqlx"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/thoas/go-funk"
)

// Config defines the options for the otter beat.
// nolint lll
type Config struct {
	DSN             string `pflag:"dsn, eg. root:8BE4@127.0.0.1:9633/test. shorthand=d"`
	PipelineListURL string `pflag:"pipeline list page url, eg. http://127.0.0.1:2901/pipeline_list.htm?channelId=1. shorthand=p"`
	InfluxWriteURL  string `pflag:"influx writing url. eg. http://beta.isignet.cn:10014/write?db=metrics. shorthand=i"`
	PrintConfig     bool   `pflag:"print config before running. shorthand=P"`

	Interval time.Duration `pflag:"interval config before running(default 60s). shorthand=I"`
}

// SetDefault sets the default value for config items.
func (c *Config) SetDefault() {
	if c.Interval == 0 {
		c.Interval = time.Minute
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

	list, err := GraspPipeLineList(c.PipelineListURL)
	if err != nil {
		logrus.Errorf("failed to GraspPipeLineList %s error %v", c.PipelineListURL, err)
		return
	}

	for _, l := range list {
		c.writeInfluxLine(l.ToPipeLineInflux(), 0)
	}
}

func (c *Config) collectTables() {
	c.timeRead(".otter.DelayStats", Dao.DelayStats,
		func(v interface{}) time.Time { return v.(PipelineDelay).ModifiedTime })
	c.timeRead(".otter.HistoryStats", Dao.HistoryStats,
		func(v interface{}) time.Time { return v.(TableHistoryStat).ModifiedTime })
	c.timeRead(".otter.TableStats", Dao.TableStats,
		func(v interface{}) time.Time { return v.(TableStat).ModifiedTime })
	c.timeRead(".otter.TableStats", Dao.ThroughputStats,
		func(v interface{}) time.Time { return v.(ThroughputStat).ModifiedTime })
	c.intRead(".otter.LogRecords", Dao.LogRecords,
		func(v interface{}) uint64 { return v.(LogRecord).ID })
}

func (c *Config) intRead(filename string, dao interface{}, f func(interface{}) uint64) {
	lastID, err := ReadFileValue(filename, "0")
	if err != nil {
		logrus.Warnf("failed to load %s %v", filename, err)
		return
	}

	logrus.Infof("read %s with value %v", filename, lastID)

	last := str.ParseUint64(lastID)
	items := reflect.ValueOf(dao).Call([]reflect.Value{reflect.ValueOf(last)})[0].Interface()

	logrus.Infof("read %s got %d items %v", filename, SliceLen(items), SliceLeft(items, 3))

	i := 0

	funk.ForEach(items, func(v interface{}) {
		if x := f(v); x > last {
			last = x
		}

		c.writeInfluxLine(v, i)
		i++
	})

	if newID := strconv.FormatUint(last, 10); newID != lastID {
		logrus.Infof("write %s with value %v", filename, newID)

		if err := WriteFileValue(filename, newID); err != nil {
			logrus.Warnf("failed to write %s error %v", newID, err)
		}
	}
}

func (c *Config) timeRead(filename string, dao interface{}, f func(interface{}) time.Time) {
	lastTime, err := ReadFileTime(filename, StartTime)
	if err != nil {
		logrus.Warnf("failed to load %s %v", filename, err)
		return
	}

	logrus.Infof("read %s with value %v", filename, lastTime.Format(TimeFormat))

	items := reflect.ValueOf(dao).Call([]reflect.Value{reflect.ValueOf(lastTime)})[0].Interface()

	logrus.Infof("read %s got %d items %v", filename, SliceLen(items), SliceLeft(items, 3))

	changed := false
	i := 0

	funk.ForEach(items, func(v interface{}) {
		if l := f(v); l.After(lastTime) {
			lastTime = l
			changed = true
		}

		c.writeInfluxLine(v, i)
		i++
	})

	if !changed {
		return
	}

	logrus.Infof("write %s with value %v", filename, lastTime.Format(TimeFormat))

	if err := WriteFileTime(filename, lastTime); err != nil {
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

	if err := influx.Write(c.InfluxWriteURL, line); err != nil {
		logrus.Warnf("failed to influx  write line %v error %v", v, err)
	}
}

const (
	// StartTime defines the start time of the system.
	StartTime = "2006-01-02 15:04:05"
	// TimeFormat defines the format of time to save to the file.
	TimeFormat = "2006-01-02 15:04:05"
)

// ReadFileTime reads the time.Time from the given file.
func ReadFileTime(filename string, defaultValue string) (time.Time, error) {
	v, err := ReadFileValue(filename, defaultValue)
	if err != nil {
		return time.Time{}, err
	}

	return lang.ParseTime(TimeFormat, v), nil
}

// WriteFileTime writes the time.Time to the given file.
func WriteFileTime(filename string, v time.Time) error {
	return WriteFileValue(filename, v.Format(TimeFormat))
}

func ReadFileValue(filename, defaultValue string) (string, error) {
	stat, err := file.StatE(filename)
	if err != nil {
		return "", errors.Wrapf(err, "file.Stat %s", filename)
	}

	if stat == file.NotExists || stat == file.Unknown {
		if err := WriteFileValue(filename, defaultValue); err != nil {
			return "", err
		}
	}

	content, err := ioutil.ReadFile(filename)
	if err != nil {
		return "", errors.Wrapf(err, "ioutil.ReadFile %s", filename)
	}

	return string(content), nil
}

// WriteFileValue writes a string value to the file.
func WriteFileValue(filename string, value string) error {
	dir := filepath.Dir(filename)
	if err := os.MkdirAll(dir, os.ModePerm); err != nil {
		return errors.Wrapf(err, "MkdirAll %s", dir)
	}

	if err := ioutil.WriteFile(filename, []byte(value), 0644); err != nil {
		return errors.Wrapf(err, "WriteFile %s", filename)
	}

	return nil
}

// SliceLen return the length of the slice.
func SliceLen(s interface{}) int {
	return reflect.ValueOf(s).Len()
}

// SliceLeft return the left at most n items of the slice.
func SliceLeft(s interface{}, n int) interface{} {
	v := reflect.ValueOf(s)
	if n >= v.Len() {
		return s
	}

	return v.Slice(0, n).Interface()
}
