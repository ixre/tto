package tto

import (
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/ixre/gof/db/db"
)

// ReadModels 从目录中识别模型并转换为表
func ReadModels(path string) ([]*Table, error) {
	tables := make([]*Table, 0)
	list, err := os.ReadDir(path)
	if err == nil {
		for _, v := range list {
			if strings.HasSuffix(v.Name(), ".t") {
				bytes, err := os.ReadFile(filepath.Join(path, v.Name()))
				if err == nil {
					tb, _ := ReverseParseTable(string(bytes))
					if len(tb) > 0 {
						tables = append(tables, tb...)
					}
				}
			}
		}
	}
	return tables, err
}

var tableRegex = regexp.MustCompile(`///*\s*([^\n]+)\n*\s*type\s+([^\s]+)\s+struct\s*{([^}]+)}`)
var txtColRegex = regexp.MustCompile(`///*\s*([\S]+)\s*([\S]+)\s+([\S]+)\s+` + "`([^`]+)`")
var txtPropRegex = regexp.MustCompile(`([\S]+?)\s*:\s*"(.+?)"`)

// / 从文本中读取表信息,逆向转换表
func ReverseParseTable(txt string) ([]*Table, error) {
	sub := tableRegex.FindAllStringSubmatch(txt, -1)
	tables := make([]*Table, 0)
	for i, v := range sub {
		// 如: user_info 用户列表
		structName := v[2]
		// 解析表名和文档
		tbTxt := strings.Split(v[1], " ")
		tbName, tbDoc := "", ""
		if len(tbTxt) > 1 {
			// 如果设置了与结构名不同的表名
			if tbTxt[0] != structName {
				tbName = tbTxt[0]
			}
			tbDoc = tbTxt[1]
		} else {
			tbDoc = tbTxt[0]
			tbName = joinLowerCase(structName, '_')
		}
		tb := &Table{
			Ordinal:    i,
			Name:       tbName,
			Prefix:     prefix(tbName),
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
		reverseParseColumns(tb, v[3])
		tables = append(tables, tb)
		//log.Println("[ reverse]: find model ", tb.Name, ":", tb.Comment)
	}
	return tables, nil
}

// 转换列
func reverseParseColumns(table *Table, txt string) {
	sub := txtColRegex.FindAllStringSubmatch(txt, -1)
	for i, v := range sub {
		props := parseColumnProps(v[4])
		ormType := GetOrmTypeFromGoType(v[3])
		colName := joinLowerCase(props["db"], '_')
		col := &Column{
			Ordinal: i,
			Name:    colName,
			Prop:    v[2],
			IsPk:    props["pk"] == "yes",
			IsAuto:  props["auto"] == "yes",
			NotNull: props["null"] != "yes",
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
				table.Name = joinLowerCase(prefix+table.Name, '_')
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
