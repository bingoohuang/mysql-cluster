package mci

import (
	"errors"
	"fmt"
	"reflect"

	"github.com/tkrajina/go-reflector/reflector"
	"golang.org/x/sync/errgroup"
)

// FlattenBeans flattens a slice of beans to a flattened bean.
func FlattenBeans(beans interface{}, flatBean interface{}, tag string) error {
	arrValue := reflect.ValueOf(beans)

	if err := CheckSlice(arrValue); err != nil {
		return err
	}

	if arrValue.Len() == 0 {
		return nil
	}

	if reflect.TypeOf(flatBean).Kind() != reflect.Ptr {
		return fmt.Errorf("pointer type required for flatBean")
	}

	beanKeyName, beanValueName, err := parseFieldValueIndex(arrValue, tag)
	if err != nil {
		return err
	}

	fieldsMap := makeTargetFieldsMap(flatBean, tag)

	var g errgroup.Group

	for i := 0; i < arrValue.Len(); i++ {
		beanVal := arrValue.Index(i)
		name := fmt.Sprintf("%v", beanVal.Field(beanKeyName).Interface())

		if f, ok := fieldsMap[name]; ok {
			v := beanVal.Field(beanValueName).Interface()

			g.Go(func() error { return f.Set(v) })
		}
	}

	return g.Wait()
}

// CheckSlice checks if v is slice or array.
func CheckSlice(v reflect.Value) error {
	switch v.Type().Kind() {
	case reflect.Slice, reflect.Array:
		return nil
	default:
		return errors.New("first argument should be slice or array")
	}
}

func parseFieldValueIndex(arrValue reflect.Value, tag string) (beanKey, beanValue int, err error) {
	beanKey = -1
	beanValue = -1

	for i, f := range reflector.New(arrValue.Index(0).Interface()).Fields() {
		tagValue, _ := f.Tag(tag)

		switch tagValue {
		case "field":
			beanKey = i
		case "value":
			beanValue = i
		}
	}

	if beanKey < 0 || beanValue < 0 {
		err = fmt.Errorf("field or value tag required")
	}

	return
}

func makeTargetFieldsMap(flatBean interface{}, tag string) map[string]reflector.ObjField {
	fieldsMap := make(map[string]reflector.ObjField)

	for _, f := range reflector.New(flatBean).Fields() {
		if tag, _ := f.Tag(tag); tag != "" {
			fieldsMap[tag] = f
		}
	}

	return fieldsMap
}
