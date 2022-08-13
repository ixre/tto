package tto

import (
	"io/ioutil"
	"log"
	"path/filepath"
	"regexp"
	"strings"
	"unicode"

	"github.com/ixre/gof/db/db"
)

func ReadModels(path string) ([]*Table, error) {
	tables := make([]*Table, 0)
	list, err := ioutil.ReadDir(path)
	if err == nil {
		for _, v := range list {
			if strings.HasSuffix(v.Name(), ".t") {
				bytes, err := ioutil.ReadFile(filepath.Join(path, v.Name()))
				if err == nil {
					tb, _ := ReadTables(string(bytes))
					if len(tb) > 0 {
						tables = append(tables, tb...)
					}
				}
			}
		}
	}
	return tables, err
}

var tableRegex = regexp.MustCompile(`///*\s*([^\n]+)\n*\s*type\s+([^\s]+)\sstruct{([^}]+)}`)
var txtColRegex = regexp.MustCompile(`///*\s*([\S]+)\s*([\S]+)\s+([\S]+)\s+` + "`([^`]+)`")
var txtPropRegex = regexp.MustCompile(`([\S]+?)\s*:\s*"(.+?)"`)

/// 从文本中读取表信息
func ReadTables(txt string) ([]*Table, error) {
	sub := tableRegex.FindAllStringSubmatch(txt, -1)
	tables := make([]*Table, 0)
	for i, v := range sub {
		structName := v[2]
		// 解析表名和文档
		tbTxt := strings.Split(v[1]," ")
		tbName,tbDoc := "",""
		if len(tbTxt) > 1{
			// 如果设置了与结构名不同的表名
			if tbTxt[0] != structName{
				tbName = tbTxt[0]
			}
			tbDoc = tbTxt[1]
		}else{
			tbDoc = tbTxt[0]
			tbName = keyFormat(structName)
		}
		tb := &Table{
			Ordinal:    i,
			Name:       tbName,
			Prefix:     "",
			Title:      structName,
			ShortTitle: structName,
			Comment:    tbDoc,
			Engine:     "",
			Schema:     "",
			Charset:    "utf8",
			Pk:         "",
			PkProp:     "",
			PkType:     0,
			Columns:    []*Column{},
		}
		readColumns(tb, v[3])
		tables = append(tables, tb)
		log.Println(i, "->", tb.Name, tb.Comment)
	}

	return tables, nil
}

func readColumns(table *Table, txt string) {
	sub := txtColRegex.FindAllStringSubmatch(txt, -1)
	for i, v := range sub {
		props := parseColumnProps(v[4])
		ormType := GetOrmTypeFromGoType(v[3])
		col := &Column{
			Ordinal: i,
			Name:    props["db"],
			Prop:    v[2],
			IsPk:    props["pk"] == "yes",
			IsAuto:  props["auto"] == "yes",
			NotNull: props["null"] == "no",
			DbType:  props["db_type"],
			Comment: v[1],
			Length:  0,
			Type:    ormType,
			Render:  &PropRenderOptions{},
		}
		if col.IsPk {
			table.Pk = col.Name
			table.PkProp = col.Prop
			table.PkType = col.Type
			if prefix := props["prefix"]; len(prefix) > 0 {
				if !strings.HasSuffix(prefix, "_") {
					prefix += "_"
				}
				table.Name = keyFormat(prefix + table.Name)
				table.Prefix = prefix[:len(prefix)-1]
			}
		}
		table.Columns = append(table.Columns, col)
	}
}

func parseColumnProps(s string) map[string]string {
	props := make(map[string]string, 0)
	if len(s) > 0 {
		for _, v := range txtPropRegex.FindAllStringSubmatch(s, -1) {
			props[v[1]] = v[2]
		}
	}
	return props
}

func keyFormat(s string) string {
	dst := make([]byte, 0)
	for i, b := range strings.TrimSpace(s) {
		if unicode.IsUpper(b) {
			l := byte(unicode.ToLower(b))
			if i == 0 {
				dst = append(dst, l)
			} else {
				dst = append(dst, byte('_'), l)
			}
		} else {
			dst = append(dst, byte(b))
		}
	}
	return string(dst)
}

func GetOrmTypeFromGoType(typeName string) int {
	switch typeName {
	case "string":
		return db.TypeString
	case "bool":
		return db.TypeBoolean
	case "int16":
		return db.TypeInt16
	case "int":
		return db.TypeInt32
	case "int64":
		return db.TypeInt64
	case "float32":
		return db.TypeFloat32
	case "float64":
		return db.TypeDecimal
	case "time.Time":
		return db.TypeDateTime
	case "[]byte":
		return db.TypeBytes
	}
	return 0
}
