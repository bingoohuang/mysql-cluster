package otterbeat

import (
	"fmt"
	"net/http"
	"strings"

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

// GraspPipeLineList 从PipeLine管理列表页面抓获数据
func GraspPipeLineList(pipelineListURL string) ([]PipeLine, error) {
	res, err := http.Get(pipelineListURL) // nolint gosec
	if err != nil {
		return nil, err
	}

	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("read pipeline list %s returned non-200 %d", pipelineListURL, res.StatusCode)
	}

	// Load the HTML document
	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		return nil, errors.Wrapf(err, "NewDocumentFromReader")
	}

	pipelines := make([]PipeLine, 0)
	creator := NewStructCreator(&pipelines)

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
