package main

import (
	"fmt"
	"math/rand"
	"os"
	"time"

	"github.com/go-mysql-org/go-mysql/canal"
	"github.com/go-mysql-org/go-mysql/mysql"
	"github.com/go-mysql-org/go-mysql/replication"
)

type MyEventHandler struct {
}

func (h *MyEventHandler) OnRotate(header *replication.EventHeader, rotateEvent *replication.RotateEvent) error {
	fmt.Printf("OnRotate, header:  %+v, event: %+v\n", header, rotateEvent)
	return nil
}

func (h *MyEventHandler) OnTableChanged(header *replication.EventHeader, schema string, table string) error {
	fmt.Printf("OnTableChanged, header:  %+v schema: %s, table: %s\n", header, schema, table)
	return nil
}

func (h *MyEventHandler) OnDDL(header *replication.EventHeader, nextPos mysql.Position, queryEvent *replication.QueryEvent) error {
	fmt.Printf("OnDDL, header:  %+v, nextPos: %+v, queryEvent: %+v\n", nextPos, queryEvent)
	return nil
}

func (h *MyEventHandler) OnRow(e *canal.RowsEvent) error {
	fmt.Printf("OnRow, event: %+v\n", e)
	return nil
}

func (h *MyEventHandler) OnXID(header *replication.EventHeader, nextPos mysql.Position) error {
	fmt.Printf("OnXID, header:  %+v, nextPos: %+v\n", header, nextPos)
	return nil
}

func (h *MyEventHandler) OnGTID(header *replication.EventHeader, gtid mysql.GTIDSet) error {
	fmt.Printf("OnGTID, header: %+v, gtid: %+v\n", header, gtid)
	return nil
}

func (h *MyEventHandler) OnPosSynced(header *replication.EventHeader, pos mysql.Position, set mysql.GTIDSet, force bool) error {
	fmt.Printf("OnPosSynced, header:  %+v, pos: %+v, gtidSet: %+v\n, force: %+v", header, pos, set, force)
	return nil
}

func (h *MyEventHandler) String() string { return "MyEventHandler" }

func main() {
	cfg := new(canal.Config)
	cfg.Addr = "127.0.0.1:3306"
	cfg.User = "root"
	cfg.Password = "root"
	cfg.ServerID = uint32(rand.New(rand.NewSource(time.Now().Unix())).Intn(1000)) + 1001

	c, err := canal.NewCanal(cfg)
	if err != nil {
		fmt.Fprintf(os.Stdout, "encounter a error during init canal, and the error: %s", err.Error())
		return
	}

	// Register a handler to handle RowsEvent
	c.SetEventHandler(&MyEventHandler{})

	// Start canal
	c.Run()
	select {}
}
