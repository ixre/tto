/**
 * Copyright 2015 @ at3.net.
 * name : thrift_test.go
 * author : jarryliu
 * date : 2016-11-17 13:37
 * description :
 * history :
 */
package tto

import (
	"github.com/ixre/tto/utils"
	"regexp"
	"sync"
	"testing"
	"time"
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
	data, err := utils.ThriftStruct(v)
	if err != nil {
		t.Error(err)
	} else {
		t.Logf("生成代码如下:\n\n" + string(data) + "\n\n")
	}
}

// 生成结构赋值代码
func TestStructAssignCode(t *testing.T) {
	data, err := utils.GoStructAssignCode(v)
	if err != nil {
		t.Error(err)
	} else {
		t.Logf("生成代码如下:\n\n" + string(data) + "\n\n")
	}
}

func TestGenByTemplate(t *testing.T) {
	wg := &sync.WaitGroup{}
	wg.Add(10)
	for i := 0; i < 10; i++ {
		go func(wg *sync.WaitGroup) {
			wg.Done()
			//result := dg.GenerateCode(&Table{Name: "admin_user"}, NewTemplate(str, "", true))
			//t.Log("--", result)
		}(wg)
	}
	wg.Wait()
	println("haha")
	time.Sleep(30 * time.Second)
}


func TestCodeTemplate_String(t *testing.T) {
	var r = regexp.MustCompile("\\{\\n+(\\s{5,})")
	var content = `
message SavePermDeptRequest{
    
    /** ID */
    int64 Id = 1;
    /** 名称 */
    string Name = 2;
`
	t.Log(r.MatchString(content))
	t.Log(r.ReplaceAllString(content,"{$1"))
}