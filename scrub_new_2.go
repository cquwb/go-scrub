package scrub

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strings"
)

func MyScrub2(value interface{}, keys map[string]bool) string {
	if value == nil {
		return "null"
	}
	if keys == nil {
		keys = MyDefaultKeys
	}
	savedValues := make([]string, 0, 10)
	myInternalScrub(value, &savedValues, "", keys, true)
	b, _ := json.Marshal(value)
	fmt.Printf("b is %s \n", b)

	myInternalScrub(value, &savedValues, "", keys, false)
	return string(b)
}

//保证value不是Nil
func myInternalScrub(value interface{}, savedValues *[]string, fieldName string, keys map[string]bool, mask bool) {
	targetValue := reflect.ValueOf(value)
	targetType := reflect.TypeOf(value)
	if targetType.Kind() != reflect.Ptr {
		//进来的元素一定是指针，可寻址的
		return
	}
	if !targetValue.IsValid() {
		//必须是合法的。否则不能调用后面的方法
		return
	}
	targetValue = targetValue.Elem()
	//这个地方其实没有必要IsValid了 因为前面确认了是ptr
	if !targetValue.IsValid() {
		return
	}
	targetType = targetValue.Type()

	//缺少了新对象是一个指针的情况

	if targetType.Kind() == reflect.Struct {
		for i := 0; i < targetValue.NumField(); i++ {
			fieldType := targetType.Field(i)
			fieldValue := targetValue.Field(i)
			if !fieldValue.IsValid() {
				continue
			}
			if !fieldValue.CanAddr() {
				continue
			}
			if !fieldValue.CanInterface() {
				continue
			}
			myInternalScrub(fieldValue.Addr().Interface(), savedValues, fieldType.Name, keys, mask)
		}
		return
	}
	if targetType.Kind() == reflect.Slice || targetType.Kind() == reflect.Array {
		for i := 0; i < targetValue.Len(); i++ {
			fieldValue := targetValue.Index(i)
			if !fieldValue.IsValid() {
				continue
			}
			if !fieldValue.CanAddr() {
				continue
			}
			if !fieldValue.CanInterface() {
				continue
			}
			myInternalScrub(fieldValue.Addr().Interface(), savedValues, fieldName, keys, mask)
		}
	}

	if fieldName == "" {
		//一开始进来传递的不是struct 直接什么都不做
		return
	}

	if targetType.Kind() == reflect.String && keys[strings.ToLower(fieldName)] && !targetValue.IsZero() {
		if !targetValue.CanSet() {
			return
		}
		if mask {
			*savedValues = append(*savedValues, targetValue.String())
			targetValue.SetString("********")
		} else {
			targetValue.SetString((*savedValues)[0])
			*savedValues = (*savedValues)[1:]
		}
	}

}
