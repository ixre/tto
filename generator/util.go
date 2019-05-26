package generator

import (
	"bytes"
	"github.com/ixre/gof/util"
	"text/template"
)

// 保存到文件
func SaveFile(s string, path string) error {
	return util.BytesToFile([]byte(s), path)
}

func ResolvePathString(pattern string,global map[string]interface{},table *Table)string {
	s,err := resolveTableString(pattern,global,table)
	if err != nil{
		panic("路径错误:"+pattern)
	}
	return s
}

func resolveTableString(tpl string,global map[string]interface{},table *Table)(string,error) {
	t := &template.Template{}
	t, err := t.Parse(tpl)
	if err == nil {
		mp := map[string]interface{}{
			"global":global,
			"table": table,
		}
		buf := bytes.NewBuffer(nil)
		err = t.Execute(buf, mp)
		return buf.String(), err
	}
	return "", err
}
