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
	dg := DBCodeGenerator()
	str := `<div class="gra-form-field col-md-6">
        <div class="gra-form-label">&nbsp;</div>
        <div class="gra-form-col">
            <div class="gra-btn gra-btn-inline btn-submit">确认订单</div>
        </div>
    </div>
</div>


s$

		${x}

        
        
        
        
        
        
        
        entity["create_time"] = utils.unix2str(entity["create_time"]);
        
        
        

        
        
        
        
        
        
        
        entity["create_time"] = utils.unix2str(entity["create_time"]);
        
        
        

        
        
        
        
        
        
        
        entity["create_time"] = utils.unix2str(entity["create_time"]);
        
        
        

	//hello

`
	result := dg.GenerateCode(&Table{Name: "Person"},
		NewTemplate(str, ""), "", true, "")
	t.Log("--", result)
}
