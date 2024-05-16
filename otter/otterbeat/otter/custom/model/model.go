package model

type DetectResult struct {
	IpPort string
	Value  Value
	Keys   []string
}

type MetricType string

const (
	Gauges      MetricType = "Gauges"      // 瞬时值,常用场景，当前活动连接数 mean()
	Counters               = "Counters"    // 计数器，汇报当前时间需要度量指标的数量，常用场景，网络流量，会对值做速率运算 non_negative_derivative(max())
	HealthCheck            = "HealthCheck" // 存活状态，通常用于记录应用存活状态 1:存活 0:未存活 mean()
)

type MetricsK1 string

const (
	Self  MetricsK1 = "self"
	Otter           = "otter"
)

// Value 部分
type Value struct {
	MetricType MetricType `json:"metricType"`
	Name       string     `json:"name"`
	Val        int        `json:"val"`
}

// Metric 自定义上报结构体
type Metric struct {
	IpPort   string   `json:"ipPort"`
	Keys     []string `json:"keys"`
	Time     string   `json:"time"`
	Category string   `json:"category"`
	Values   []Value  `json:"values"`
}
