package utils

import (
	"bytes"
	"errors"
	"reflect"
)

// 生成结构赋值代码
func GoStructAssignCode(v interface{}) ([]byte, error) {
	vt := reflect.TypeOf(v)
	if vt.Kind() == reflect.Ptr {
		vt = vt.Elem()
	}
	if vt.Kind() != reflect.Struct {
		return nil, errors.New("value is not struct")
	}
	buf := bytes.NewBufferString("v := &")
	buf.WriteString(vt.Name())
	buf.WriteString(" {\n")
	for i, n := 0, vt.NumField(); i < n; i++ {
		f := vt.Field(i)
		buf.WriteString("    ")
		buf.WriteString(f.Name)
		buf.WriteString(" : src.")
		buf.WriteString(f.Name)
		buf.WriteString(",\n")
	}
	buf.WriteString("}")
	return buf.Bytes(), nil
}
