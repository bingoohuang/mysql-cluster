package influx

import (
	"fmt"
	"reflect"
	"strings"
	"time"

	"github.com/bingoohuang/gonet"
	"github.com/bingoohuang/strcase"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

// ToLine creates a new line string suitable for the influxdb line protocol.
func ToLine(v interface{}) (string, error) {
	su, err := ParseLine(v)
	if err != nil {
		return "", err
	}

	return su.ToLine()
}

// nolint
var (
	tType    = reflect.TypeOf((*T)(nil)).Elem()
	timeType = reflect.TypeOf((*time.Time)(nil)).Elem()
)

func ParseLine(v interface{}) (*Line, error) {
	rv := reflect.ValueOf(v)
	rt := rv.Type()

	if rv.Kind() == reflect.Ptr {
		rv = rv.Elem()
		rt = rt.Elem()
	}

	if rt.Kind() != reflect.Struct {
		return nil, fmt.Errorf(
			"only struct or pointer to struct supported, %v is illegal", rt)
	}

	su := &Line{Tags: make([]Tag, 0), Fields: make([]Field, 0)}

	for i := 0; i < rt.NumField(); i++ {
		su.buildSu(rt.Field(i), rv)
	}

	if su.Measurement == "" {
		su.Measurement = strcase.ToSnake(rt.Name())
	}

	return su, nil
}

func (su *Line) buildSu(rtf reflect.StructField, rv reflect.Value) {
	if rtf.PkgPath != "" {
		return
	}

	if rtf.Type == tType {
		su.Measurement = rtf.Tag.Get("measurement")
		return
	}

	influxTag := rtf.Tag.Get("influx")
	if influxTag == "-" {
		return
	}

	if rtf.Type == timeType {
		if influxTag == "time" {
			su.Time = rv.FieldByIndex(rtf.Index).Interface().(time.Time)
			return
		}
	}

	name := rtf.Tag.Get("name")
	if name == "" {
		name = strcase.ToSnake(rtf.Name)
	}

	if influxTag == "tag" {
		su.Tags = append(su.Tags, Tag{K: name,
			V: fmt.Sprintf("%v", rv.FieldByIndex(rtf.Index).Interface()),
		})

		return
	}

	su.Fields = append(su.Fields, Field{K: name,
		V: rv.FieldByIndex(rtf.Index).Interface(),
	})
}

// LineProtocol format inputs to line protocol
// https://docs.influxdata.com/influxdb/v1.7/write_protocols/line_protocol_tutorial/
func (su *Line) ToLine() (string, error) {
	tagExpr := ""

	for i, v := range su.Tags {
		if i > 0 {
			tagExpr += ","
		}

		tagExpr += fmt.Sprintf("%s=%s",
			escapeSpecialChars(v.K), escapeSpecialChars(v.V))
	}

	fieldsExpr := ""

	for i, v := range su.Fields {
		if i > 0 {
			fieldsExpr += ","
		}

		r, err := toInfluxRepresentation(v.V)
		if err != nil {
			return "", errors.Wrapf(err, "toInfluxRepresentation %+v", v.V)
		}

		fieldsExpr += fmt.Sprintf("%s=%s", escapeSpecialChars(v.K), r)
	}

	if su.Time.IsZero() {
		su.Time = time.Now()
	}

	// construct line protocol string
	return fmt.Sprintf("%s,%s %s %d",
		su.Measurement, tagExpr, fieldsExpr, uint64(su.Time.UnixNano())), nil
}

// Line represents a structure to generate a influxdb line protocol.
type Line struct {
	Measurement string
	Time        time.Time
	Tags        []Tag
	Fields      []Field
}

// Tag represents a tag pair of a measurement record in the influxdb.
type Tag struct {
	K string
	V string
}

// Field represents a field of a measurement record in the influxdb.
type Field struct {
	K string
	V interface{}
}

// T is a helper empty struct to declare influx measurement name.
type T struct{}

func escapeSpecialChars(in string) string {
	str := strings.Replace(in, ",", `\,`, -1)
	str = strings.Replace(str, "=", `\=`, -1)
	str = strings.Replace(str, " ", `\ `, -1)

	return str
}

// toInfluxRepresentation 将val转换为Influx表示形式
func toInfluxRepresentation(val interface{}) (string, error) {
	switch v := val.(type) {
	case string:
		return stringToInfluxRepresentation(v)
	case []byte:
		return stringToInfluxRepresentation(string(v))
	case int32, int64, int16, int8, int, uint32, uint64, uint16, uint8, uint:
		return fmt.Sprintf("%d", v), nil
	case float64, float32:
		return fmt.Sprintf("%g", v), nil
	case bool:
		return fmt.Sprintf("%t", v), nil
	case time.Time:
		return fmt.Sprintf("%d", uint64(v.UnixNano())), nil
	case time.Duration:
		return fmt.Sprintf("%d", uint64(v.Nanoseconds())), nil
	default:
	}

	if s, ok := val.(fmt.Stringer); ok {
		return stringToInfluxRepresentation(s.String())
	}

	return "", fmt.Errorf("%+v: unsupported type for Influx Line Protocol", val)
}

func stringToInfluxRepresentation(v string) (string, error) {
	if len(v) > 64000 { // nolint gomnd
		return "", fmt.Errorf("string too long (%d characters, max. 64K)", len(v))
	}

	return fmt.Sprintf("%q", v), nil
}

// Write 写入打点值
// refer https://github.com/DCSO/fluxline/blob/master/encoder.go
func Write(influxDBWriteAddr, line string) error {
	req, err := gonet.Post(influxDBWriteAddr)
	if err != nil {
		return err
	}

	req.Body([]byte(line))

	rsp, err := req.SendOut()
	if err != nil {
		return err
	}

	rspBody, err := req.ReadResponseBody(rsp)
	if err != nil {
		logrus.Warnf("influx write error %v", err)
		return err
	}

	logrus.Infof("influx write %s returned status %s msg %s", line, rsp.Status, string(rspBody))

	return nil
}
