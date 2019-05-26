/**
 * Copyright 2015 @ at3.net.
 * name : thrift_test.go
 * author : jarryliu
 * date : 2016-11-17 13:37
 * description :
 * history :
 */
package generator

import (
	"testing"
)

type testStruct struct {
	Id      int
	Name    string
	Balance float32
	Enabled bool
}

var (
	v = &testStruct{}
)

// 生成Thrift结构
func TestThriftStruct(t *testing.T) {
	data, err := ThriftStruct(v)
	if err != nil {
		t.Error(err)
	} else {
		t.Logf("生成代码如下:\n\n" + string(data) + "\n\n")
	}
}

// 生成结构赋值代码
func TestStructAssignCode(t *testing.T) {
	data, err := StructAssignCode(v)
	if err != nil {
		t.Error(err)
	} else {
		t.Logf("生成代码如下:\n\n" + string(data) + "\n\n")
	}
}

func TestGenByTemplate(t *testing.T) {
	dg := DBCodeGenerator()
	str := "s$${x}"
	result := dg.GenerateCode(&Table{Name: "Person"},
		CodeTemplate(str), "", true, "")
	t.Log("--", result)
}
