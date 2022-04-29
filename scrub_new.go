package scrub

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strings"
)

var MyDefaultKeys = map[string]bool{"password": true}

func MyScrub(value interface{}, keys map[string]bool) string {
	if keys == nil {
		keys = MyDefaultKeys

	}
	if value == nil {
		return "null"
	}
	targetValue := reflect.ValueOf(value)
	if targetValue.Kind() != reflect.Ptr {
		panic("value is not ptr")
	}
	realValue := targetValue.Elem()
	if !realValue.IsValid() {
		return "null"
	}
	if realValue.Kind() != reflect.Struct {
		panic("real value is not struct")
	}
	for i := 0; i < realValue.NumField(); i++ {
		fieldValue := realValue.Field(i)
		fieldType := realValue.Type().Field(i)
		if !fieldValue.IsValid() {
			continue
		}
		if !fieldValue.CanSet() {
			continue
		}
		fmt.Printf("set field %s \n", fieldType.Name)
		MySubScrub(fieldValue.Addr().Interface(), fieldType.Name, keys)
		/*
			if keys[strings.ToLower(fieldType.Name)] {
				if fieldType.Type.Kind() == reflect.Slice || fieldType.Type.Kind() == reflect.Array {
					for j := 0; j < fieldValue.Len(); j++ {
						//这里是只有数组的值才有长度，不是数组的类型,切片的类型是没有长度这个值的，会一直扩展,只有数组这个类型有
						fieldValue.Index(j).SetString("********")
					}
				} else if fieldType.Type.Kind() == reflect.String {
					if !fieldValue.CanSet() {
						continue
					}
					fieldValue.SetString("********")
				} else if fieldType.Type.Kind() == reflect.Struct {

				}
			}
		*/
	}
	r, _ := json.Marshal(realValue.Addr().Interface())
	return string(r)
}

func MySubScrub(value interface{}, fieldName string, keys map[string]bool) {
	targetValue := reflect.ValueOf(value)
	targetType := reflect.TypeOf(value)
	fmt.Printf("MySubScrub :type kind:%s filedName:%s \n", targetType.Kind(), fieldName)
	if targetType.Kind() == reflect.Ptr {
		//指针类型，比如字符串就是的
		targetValue = targetValue.Elem()
		targetType = targetValue.Type()
	}
	fmt.Printf("MySubScrub2 :type kind:%s filedName:%s \n", targetType.Kind(), fieldName)
	if targetType.Kind() == reflect.Struct {
		for i := 0; i < targetValue.NumField(); i++ {
			fieldValue := targetValue.Field(i)
			fieldType := targetType.Field(i)
			//fmt.Printf("MySubScrub struct :fieldValue:%v filedType:%v \n", fieldValue, fieldType)
			MySubScrub(fieldValue.Addr().Interface(), fieldType.Name, keys)
		}
		return
	}
	if targetType.Kind() == reflect.Slice || targetType.Kind() == reflect.Array {
		for i := 0; i < targetValue.Len(); i++ {
			//这里是只有数组的值才有长度，不是数组的类型,切片的类型是没有长度这个值的，会一直扩展,只有数组这个类型有
			MySubScrub(targetValue.Index(i).Addr().Interface(), fieldName, keys)
		}
		return
	}

	if targetType.Kind() == reflect.String && keys[strings.ToLower(fieldName)] && !targetValue.IsZero() {
		targetValue.SetString("********")
	}

}
