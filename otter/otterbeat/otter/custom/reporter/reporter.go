package reporter

import (
	"fmt"
	"github.com/bingoohuang/gg/pkg/jsoni"
	"github.com/bingoohuang/otterbeat/otter/custom/http"
	"github.com/bingoohuang/otterbeat/otter/custom/model"
	"log"
	"time"
)

// Reporter 上报到自定义监控服务
type Reporter struct {
	Url      string
	Category string
}

func (r *Reporter) Initialize(url, category string) error {
	if url == "" {
		return fmt.Errorf("上报自定义监控，配置文件参数不完整，缺少 Urls。config: %v", url)
	}
	if category == "" {
		return fmt.Errorf("上报自定义监控，配置文件参数不完整，缺少 Category。config: %v", category)
	}

	r.Url = url
	r.Category = category
	return nil
}

func (r *Reporter) Report(detect model.DetectResult) error {
	url := r.Url
	metric := model.Metric{
		IpPort:   detect.IpPort,
		Keys:     detect.Keys,
		Time:     time.Now().Format("20060102150405"),
		Category: r.Category,
		Values:   []model.Value{detect.Value},
	}
	rsp := make(map[interface{}]interface{})
	if err := http.Post(url, "", metric, &rsp); err != nil {
		var request interface{}
		if bytes, err := jsoni.Marshal(metric); err == nil {
			request = string(bytes)
		} else {
			request = metric
		}

		log.Printf("E! Report：%v，request：%v，response：%v，err：%v", url, request, rsp, err)
		return err
	}
	return nil
}
