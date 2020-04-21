package otter

// Config defines the options for the otter beat.
// nolint lll
type Config struct {
	DSN             string `pflag:"dsn, eg. root:8BE4@127.0.0.1:9633/test. shorthand=d"`
	PipelineListURL string `pflag:"pipeline list page url, eg. http://127.0.0.1:2901/pipeline_list.htm?channelId=1. shorthand=p"`
	InfluxWriteURL  string `pflag:"influx writing url. eg. http://beta.isignet.cn:10014/write?db=metrics. shorthand=i"`
	PrintConfig     bool   `pflag:"print config before running. shorthand=P"`
}
