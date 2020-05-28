package otter

import (
	"strings"
	"time"

	"github.com/bingoohuang/gonet/man"

	"github.com/bingoohuang/gou/lang"

	"github.com/bingoohuang/gou/str"
	"github.com/bingoohuang/otterbeat/structs"

	"github.com/PuerkitoBio/goquery"
	"github.com/pkg/errors"
)

// PipeLine 获得从PipeLine管理页面上的饿列表数据
//	{
//		"Seq": "3",
//		"Name": "pipeb",
//		"State": "工作中",
//		"Delay": "1.213 s",
//		"LastSyncTime": "2020-04-14 08:27:35",
//		"LastPositionTime": "2020-04-17 07:22:04"
//	}
type PipeLine struct {
	Seq              string `col:"序号"`
	Name             string `col:"Pipeline名字"`
	State            string `col:"mainstem状态"`
	Delay            string `col:"延迟时间"`
	LastSyncTime     string `col:"最后同步时间"`
	LastPositionTime string `col:"最后位点时间"`
}

// PipeLineInflux defines the structure to write to influx DB.
type PipeLineInflux struct {
	Time             time.Time `influx:"time" measurement:"otter_pipeline"`
	PipelineID       string    `influx:"tag"`   // 流水线ID
	State            string    `influx:"field"` // 流水线ID
	DelayTime        float64   `influx:"field"` // 单位ms
	LastSyncTime     time.Time `influx:"field"` // 单位ms
	LastPositionTime time.Time `influx:"field"` // 单位ms
}

// ToPipeLineInflux converts a PipeLine to PipeLineInflux.
func (p PipeLine) ToPipeLineInflux() PipeLineInflux {
	return PipeLineInflux{
		Time:             time.Now(),
		PipelineID:       p.Seq,
		State:            p.State,
		DelayTime:        float64(str.ParseDuration(p.Delay).Milliseconds()),
		LastSyncTime:     lang.ParseTime("2006-01-02 15:04:05", p.LastSyncTime),
		LastPositionTime: lang.ParseTime("2006-01-02 15:04:05", p.LastPositionTime),
	}
}

// PipeLinePoster is a HTTP poster for pipeline list.
type PipeLinePoster struct {
	man.URL

	PipeLineList func() (string, error)
}

// NewPipeLinePoster creates a new PipeLinePoster
func NewPipeLinePoster() *PipeLinePoster { p := new(PipeLinePoster); man.New(p); return p }

// GraspPipeLineList 从PipeLine管理列表页面抓获数据
func (p *PipeLinePoster) GraspPipeLineList() ([]PipeLine, error) {
	if p.URL == "" {
		return nil, nil
	}

	res, err := p.PipeLineList() // nolint gosec
	if err != nil {
		return nil, err
	}

	// Load the HTML document
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(res))
	if err != nil {
		return nil, errors.Wrapf(err, "NewDocumentFromReader")
	}

	pipelines := make([]PipeLine, 0)
	creator := structs.NewCreator(&pipelines)

	doc.Find("table.list tr").Each(func(i int, s *goquery.Selection) {
		columns := make([]string, 0)
		s.Children().Each(func(j int, s *goquery.Selection) {
			columns = append(columns, strings.TrimSpace(s.Text()))
		})

		if i == 0 {
			creator.PrepareColumns(columns)
		} else {
			creator.CreateSliceItem(columns)
		}
	})

	return pipelines, nil
}
