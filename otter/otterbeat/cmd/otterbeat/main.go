package main

import (
	"embed"
	"github.com/bingoohuang/gg/pkg/ctl"
	"github.com/bingoohuang/gg/pkg/fla9"
	"github.com/bingoohuang/otterbeat/otter"
	"github.com/bingoohuang/rotatefile"
	_ "github.com/go-sql-driver/mysql"
	"log"
)

// InitAssets is the initial assets.
//
//go:embed assets
var InitAssets embed.FS

func main() {
	pInit := fla9.Bool("init", false, "Create initial ctl and exit")
	pVersion := fla9.Bool("version,v", false, "Create initial ctl and exit")
	confFile := fla9.String("conf,c", "./cnf.toml", "config file")
	fla9.Parse()
	ctl.Config{Initing: *pInit, PrintVersion: *pVersion, InitFiles: &InitAssets}.ProcessInit()
	log.SetOutput(rotatefile.New())
	config, err := otter.ParseConfFile(*confFile)
	if err != nil {
		log.Fatalf("parse configuration, failed: %v", err)
	}

	config.Run()
}
