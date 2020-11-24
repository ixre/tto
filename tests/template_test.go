package tests

import (
	"github.com/ixre/tto"
	"testing"
)

func TestTemplate(t *testing.T) {
	tp := tto.NewTemplate(`/**
 * this file is auto generated by tto v{{.global.version}} !
 * if you want to modify this code,please read guide doc
 * and modify code template later.
 *
 * guide please see https://github.com/ixre/tto
 *
 */
!filename:{{.table.Title}}Entity.java
package {{.global.pkg}}.pojo;

 {{range $i,$c := $validateColumns}}\
  column: {{$c.Name}}
 {{end}}

`, "", false)

	t.Log(tp.String())
}

//func TestSubstrN(t *testing.T) {
//	fn := &tto.internalFunc{}
//	s := fn.substrN("admin_user_list", "_", 2)
//	println(s)
//}

func TestMultiQuota(t *testing.T) {
	str := `--, ""
	haha
`
	dg := tto.DBCodeGenerator("")
	result := dg.GenerateCode(&tto.Table{Name: "admin_user"}, tto.NewTemplate(str, "", true))
	t.Log("--", result)
}