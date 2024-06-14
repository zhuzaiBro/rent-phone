package utils

import (
	"reflect"
	"strings"
)

const (
	FORM = "form"
	JSON = "json"
)

func StructToMap(obj interface{}) map[string]interface{} {
	t := reflect.TypeOf(obj)
	v := reflect.ValueOf(obj)
	var data = make(map[string]interface{})
	for i := 0; i < t.NumField(); i++ {
		data[t.Field(i).Tag.Get(FORM)] = v.Field(i).Interface()
	}
	return data
}

func FlattenStruct(obj interface{}) map[string]interface{} {
	result := make(map[string]interface{})
	objValue := reflect.ValueOf(obj)
	objType := objValue.Type()

	for i := 0; i < objValue.NumField(); i++ {
		fieldValue := objValue.Field(i)
		fieldType := objType.Field(i)

		// 获取字段的form标签值
		tag := fieldType.Tag.Get(FORM)
		if tag == "" {
			continue
		}

		// 将标签值按逗号分隔
		tags := strings.Split(tag, ",")

		// 获取字段的值
		value := fieldValue.Interface()

		// 将字段添加到结果map中
		for _, t := range tags {
			result[t] = value
		}
	}

	return result
}
