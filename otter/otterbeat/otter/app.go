package otter

import (
	"database/sql"
	"fmt"
	"github.com/bingoohuang/gg/pkg/ss"
	"github.com/bingoohuang/gonet"
	"github.com/bingoohuang/otterbeat/otter/custom/model"
	reporter2 "github.com/bingoohuang/otterbeat/otter/custom/reporter"
	"os"
	"reflect"
	"sort"
	"strconv"
	"strings"
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
	Version                bool   `pflag:"version. shorthand=v"`
	DSN                    string `pflag:"dsn, eg. root:8BE4@127.0.0.1:9633/test. shorthand=d"`
	PipelineListURLPattern string `pflag:"pipeline list page url pattern, eg. http://127.0.0.1:2901/pipeline_list.htm?channelId=%d. shorthand=p"`
	ChannelIds             string `pflag:"multi channelId, eg. 1,3"`
	CustomMetricUrl        string `pflag:"CustomMetricUrl eg. http://192.168.108.11:30008/custom .see : http://192.168.131.51:8090/pages/viewpage.action?pageId=21955287"`
	CustomMetricCategory   string `pflag:"eg. otter-health-check .see : http://192.168.131.51:8090/pages/viewpage.action?pageId=21955287"`
	InfluxWriteURL         string `pflag:"influx writing url. eg. http://beta.isignet.cn:10014/write?db=metrics. shorthand=i"`
	PrintConfig            bool   `pflag:"print config before running. shorthand=P"`

	Interval time.Duration `pflag:"interval config before running(default 60s). shorthand=I"`
}

// nolint
var pipeLinePosters []*PipeLinePoster
var reporter *reporter2.Reporter
var localIps []string // LocalIps 当前机器Ip列表

func (c *Config) init() {
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

	if c.PipelineListURLPattern != "" && c.ChannelIds != "" {
		pipeLinePosters = make([]*PipeLinePoster, 0)
		for _, channelId := range strings.Split(c.ChannelIds, ",") {
			if channelId == "" {
				continue
			}
			pipeLinePoster := NewPipeLinePoster()
			pipeLinePoster.URL = man.URL(fmt.Sprintf(c.PipelineListURLPattern, channelId))
			pipeLinePosters = append(pipeLinePosters, pipeLinePoster)
		}
	}

	if c.CustomMetricUrl != "" && c.CustomMetricCategory != "" {
		var r reporter2.Reporter
		if err := r.Initialize(c.CustomMetricUrl, c.CustomMetricCategory); err != nil {
			logrus.Errorf("failed to initialize reporter %v", err)
			return
		}
		reporter = &r
	}

	if len(localIps) == 0 {
		ipv4s := gonet.ListIpsv4()
		sort.Strings(ipv4s)
		localIps = ipv4s
	}
}

// Run runs  the otter beat by the config.
func (c *Config) Run() {
	c.init()

	ticker := time.NewTicker(c.Interval)
	defer ticker.Stop()

	if c.DSN != "" {
		db, err := sql.Open("mysql", sqlx.CompatibleMySQLDs(c.DSN))
		if err != nil {
			logrus.Errorf("failed to open mysql %s", c.DSN)
		}

		sqlx.DB = db
	}

	for {
		c.collectDB()
		c.collectPipelineListPage()
		c.collectDBForCustomMetrics()

		<-ticker.C
	}
}

func (c *Config) collectDB() {
	if sqlx.DB != nil {
		c.collectTables()
	}
}

func (c *Config) collectPipelineListPage() {
	if pipeLinePosters == nil {
		return
	}
	for _, pipeLinePoster := range pipeLinePosters {
		c.collectOnePipelineListPage(pipeLinePoster)
	}
}

func (c *Config) collectOnePipelineListPage(pipeLinePoster *PipeLinePoster) {
	list, err := pipeLinePoster.GraspPipeLineList()
	if err != nil {
		logrus.Errorf("failed to GraspPipeLineList %s error %v", pipeLinePoster.URL, err)
		return
	}

	for _, l := range list {
		lineInflux := l.ToPipeLineInflux()
		c.writeInfluxLine(lineInflux, 0)
		c.writeCustomMetricsStateAndDelayTime(lineInflux.PipelineID, lineInflux.State, int(lineInflux.DelayTime))
	}
}

const (
	otterbeatDir = "~/.otterbeat/"
)

// nolint lll
func (c *Config) collectTables() {
	c.timeRead(otterbeatDir+"DelayStat.lastTime", Dao.DelayStat, func(v interface{}) time.Time { return v.(DelayStat).ModifiedTime })
	c.timeRead(otterbeatDir+"TableHistoryStat.lastTime", Dao.TableHistoryStat, func(v interface{}) time.Time { return v.(TableHistoryStat).ModifiedTime })
	c.timeRead(otterbeatDir+"TableStat.lastTime", Dao.TableStat, func(v interface{}) time.Time { return v.(TableStat).ModifiedTime })
	c.timeRead(otterbeatDir+"ThroughputStat.lastTime", Dao.ThroughputStat, func(v interface{}) time.Time { return v.(ThroughputStat).ModifiedTime })
	c.intRead(otterbeatDir+"LogRecord.lastID", Dao.LogRecord, func(v interface{}) uint64 { return v.(LogRecord).ID })
}

// 最近 1 分钟异常个数
// http://localhost:8086/query?db=metrics&q=SELECT count(id) from otter_log_record WHERE time >= now() - 1m and pipeline_id = '-1' group by time(1m)

func (c *Config) intRead(filename string, dao interface{}, f func(interface{}) uint64) {
	lastID, err := file.ReadValue(filename, "0")
	if err != nil {
		logrus.Warnf("failed to load %s %v", filename, err)
		return
	}

	logrus.Infof("read %s with value %v", filename, lastID)

	last := str.ParseUint64(lastID)
	items := reflect.ValueOf(dao).Call([]reflect.Value{reflect.ValueOf(last)})[0].Interface()

	if funk.Len(items) == 0 {
		logrus.Infof("read %s got no new items", filename)
		return
	}

	logrus.Infof("read %s got new items %d: %v", filename, funk.Len(items), funk.Left(items, 3))

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

	if funk.Len(items) == 0 {
		logrus.Infof("read %s got no new items", filename)
		return
	}

	logrus.Infof("read %s got new items %d: %v", filename, funk.Len(items), funk.Left(items, 3))

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
		logrus.Infof("influx %s", line)
	} else if i == 4 {
		logrus.Infof("influx ...")
	}

	if c.InfluxWriteURL == "" {
		return
	}

	if err := influx.Write(line); err != nil {
		logrus.Warnf("failed to influx  write line %v error %v", v, err)
	}
}

// 写 2 个属性 pipeline 的工作状态 + 延时 到自定义埋点监控
func (c *Config) writeCustomMetricsStateAndDelayTime(pipelineId, state string, delayTime int) {
	if reporter == nil {
		return
	}
	var r []model.DetectResult
	r = appendVal(r, []string{model.Otter, "pipeline", pipelineId, "state"}, ss.Ifi(state == "工作中", 1, 0))
	r = appendVal(r, []string{model.Otter, "pipeline", pipelineId, "delayTime"}, delayTime)
	reportAsCustom(r)
}

// 写每分钟 插入 更新 删除 同步记录数，以及自身的存活埋点 到自定义埋点监控
func (c *Config) collectDBForCustomMetrics() {
	// 时间范围 1 min 内
	to := time.Now()
	from := to.Add(-time.Minute)
	exceptions := Dao.Exceptions(from, to)
	var r []model.DetectResult

	// 每分钟打一个点证明 otterbeat 还在运行
	r = appendVal(r, []string{string(model.Self), localIps[0], "alive"}, 1)

	// 是否发生异常
	if len(exceptions) > 0 {
		for _, exception := range exceptions {
			r = appendVal(r, []string{model.Otter, "exception"}, exception.Count)
		}
	}

	// 同步记录数
	stats := Dao.TableHistoryStatByPipeline(from, to)
	if len(stats) > 0 {
		for _, stat := range stats {
			r = appendVal(r, []string{model.Otter, "pipeline", stat.PipelineID, "insertCount"}, int(stat.InsertCount))
			r = appendVal(r, []string{model.Otter, "pipeline", stat.PipelineID, "updateCount"}, int(stat.UpdateCount))
			r = appendVal(r, []string{model.Otter, "pipeline", stat.PipelineID, "deleteCount"}, int(stat.DeleteCount))
		}
	}

	reportAsCustom(r)
}

func appendVal(r []model.DetectResult, keys []string, val int) []model.DetectResult {
	r = append(r, model.DetectResult{
		IpPort: localIps[0],
		Value: model.Value{
			MetricType: model.Gauges,
			Name:       "V",
			Val:        val,
		},
		Keys: keys,
	})
	return r
}

func reportAsCustom(r []model.DetectResult) {
	for _, detectResult := range r {
		if err := reporter.Report(detectResult); err != nil {
			logrus.Errorf("failed to report %v", err)
		}
	}
}

const (
	// StartTime defines the start time of the system.
	StartTime = "2006-01-02 15:04:05.000"
)
