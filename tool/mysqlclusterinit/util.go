package mysqlclusterinit

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"regexp"
	"strings"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"github.com/tkrajina/go-reflector/reflector"
)

// ReplaceFileContent 使用正则表达式查找模式，并且替换正则1号捕获分组为指定的内容
func ReplaceFileContent(filename, regexStr, repl string) error {
	conf, err := ioutil.ReadFile(filename)
	if err != nil {
		return fmt.Errorf("ReadFile %s error %w", filename, err)
	}

	fixed, err := ReplaceContent(string(conf), regexStr, repl)
	if err != nil {
		return err
	}

	return ioutil.WriteFile(filename, []byte(fixed), 0644)
}

// FileExists 检查文件是否存在，并且不是目录
func FileExists(filename string) error {
	if fi, err := os.Stat(filename); err != nil {
		return err
	} else if fi.IsDir() {
		return fmt.Errorf("file %s is a directory", filename)
	}

	return nil
}

// SearchPatternLines 使用正则表达式boundaryRegexStr查找大块，然后在大块中用captureGroup1Regex中的每行寻找匹配
func SearchPatternLinesInFile(filename, boundaryRegexStr, captureGroup1Regex string) ([]string, error) {
	str, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("ReadFile %s error %w", filename, err)
	}

	return SearchPatternLines(string(str), boundaryRegexStr, captureGroup1Regex)
}

// SearchPatternLines 使用正则表达式boundaryRegexStr查找大块，然后在大块中用captureGroup1Regex中的每行寻找匹配
func SearchPatternLines(str, boundaryRegexStr, captureGroup1Regex string) ([]string, error) {
	re, err := regexp.Compile(boundaryRegexStr)
	if err != nil {
		return nil, err
	}

	ce, err := regexp.Compile(captureGroup1Regex)
	if err != nil {
		return nil, err
	}

	var lines []string

	for _, v := range re.FindAllStringSubmatch(str, -1) {
		for _, l := range ce.FindAllStringSubmatch(v[1], -1) {
			lines = append(lines, l[1])
		}
	}

	return lines, nil
}

// ReplaceContent 使用正则表达式查找模式，并且替换正则1号捕获分组为指定的内容
func ReplaceContent(str, regexStr, repl string) (string, error) {
	re, err := regexp.Compile(regexStr)
	if err != nil {
		return "", err
	}

	fixed := ""
	lastIndex := 0

	for _, v := range re.FindAllStringSubmatchIndex(str, -1) {
		if len(v) != 4 {
			return "", fmt.Errorf("regexp %s should have only one capturing group", regexStr)
		}

		fixed += str[lastIndex:v[2]] + repl
		lastIndex = v[3]
	}

	if lastIndex == 0 {
		return "", fmt.Errorf("regexp %s found non submatches", regexStr)
	}

	return fixed + str[lastIndex:], nil
}

// PrettyJSONSilent prettify the JSON encoding of data silently
func PrettyJSONSilent(data interface{}) string {
	p, _ := PrettyJSON(data)
	return p
}

// PrettyJSON prettify the JSON encoding of data
func PrettyJSON(data interface{}) (string, error) {
	buffer := new(bytes.Buffer)
	encoder := json.NewEncoder(buffer)
	encoder.SetIndent("", "\t")

	err := encoder.Encode(data)
	if err != nil {
		return "", err
	}

	return buffer.String(), nil
}

// DeclarePflagsByStruct declares flags from struct fields'name and type
func DeclarePflagsByStruct(structVar interface{}) {
	for _, f := range reflector.New(structVar).Fields() {
		if !f.IsExported() {
			continue
		}

		switch t, _ := f.Get(); t.(type) {
		case []string:
			pflag.StringP(f.Name(), "", "", f.Name())
		case string:
			pflag.StringP(f.Name(), "", "", f.Name())
		case int:
			pflag.IntP(f.Name(), "", 0, f.Name())
		case bool:
			pflag.BoolP(f.Name(), "", false, f.Name())
		}
	}
}

// ViperToStruct read viper value to struct
func ViperToStruct(structVar interface{}) {
	for _, f := range reflector.New(structVar).Fields() {
		if !f.IsExported() {
			continue
		}

		switch t, _ := f.Get(); t.(type) {
		case []string:
			value := viper.GetString(f.Name())
			valueSlice := make([]string, 0)

			for _, v := range strings.Split(value, ",") {
				v = strings.TrimSpace(v)
				if v != "" {
					valueSlice = append(valueSlice, v)
				}
			}

			if len(valueSlice) > 0 {
				_ = f.Set(valueSlice)
			}
		case string:
			value := viper.GetString(f.Name())
			if value = strings.TrimSpace(value); value != "" {
				_ = f.Set(value)
			}
		case int:
			if value := viper.GetInt(f.Name()); value != 0 {
				_ = f.Set(value)
			}
		case bool:
			if value := viper.GetBool(f.Name()); value {
				_ = f.Set(value)
			}
		}
	}
}
