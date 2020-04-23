package main

import (
	"github.com/sirupsen/logrus"

	"github.com/bingoohuang/gou/cnf"
	"github.com/bingoohuang/otterbeat/otter"
	_ "github.com/go-sql-driver/mysql"
)

func main() {
	logrus.SetLevel(logrus.DebugLevel)

	cnf.DeclarePflags()
	cnf.DeclarePflagsByStruct(otter.Config{})

	if err := cnf.ParsePflags("OTTERBEAT"); err != nil {
		panic(err)
	}

	var config otter.Config

	cnf.LoadByPflag(&config)

	config.SetUp()
	config.Run()
}
