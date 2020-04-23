package main

import (
	"github.com/bingoohuang/gou/str"
	"github.com/bingoohuang/gou/sy"
	"github.com/bingoohuang/otterbeat/otter"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gobuffalo/packr/v2"
)

func main() {
	config := otter.Config{}

	box := packr.New("myBox", "assets")

	sy.SetupApp(&sy.AppOption{
		EnvPrefix:   "OTTERBEAT",
		LogLevel:    "debug",
		ConfigBeans: []interface{}{&config},
		CnfTpl:      str.PickFirst(box.FindString("cnf.tpl.toml")),
		CtlTpl:      str.PickFirst(box.FindString("ctl.tpl.sh")),
	})

	config.Run()
}
