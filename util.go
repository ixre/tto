package tto

import (
	"bytes"
	"github.com/ixre/gof/util"
	"strings"
	"text/template"
	"unicode"
)



func  prefix(str string) string {
	if i := strings.Index(str, "_"); i != -1 {
		return str[:i]
	}
	for i, l := 1, len(str); i < l-1; i++ {
		if unicode.IsUpper(rune(str[i])) {
			return strings.ToLower(str[:i])
		}
	}
	return ""
}

func title(str string,shortUpper bool) string {
	// 小于3且ID大写，则返回大写
	if shortUpper && len(str) < 3 {
		return strings.ToUpper(str)
	}
	arr := strings.Split(str, "_")
	for i, v := range arr {
		arr[i] = strings.Title(v)
	}
	return strings.Join(arr, "")
}


// 保存到文件
func SaveFile(s string, path string) error {
	return util.BytesToFile([]byte(s), path)
}

func ResolvePathString(pattern string, global map[string]interface{}, table *Table) string {
	s, err := resolveTableString(pattern, global, table)
	if err != nil {
		panic("路径错误:" + pattern)
	}
	return strings.TrimSpace(s)
}

func resolveTableString(tpl string, global map[string]interface{}, table *Table) (string, error) {
	t := &template.Template{}
	t, err := t.Parse(tpl)
	if err == nil {
		mp := map[string]interface{}{
			"global": global,
			"table":  table,
		}
		buf := bytes.NewBuffer(nil)
		err = t.Execute(buf, mp)
		return buf.String(), err
	}
	return "", err
}

