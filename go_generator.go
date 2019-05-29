package tto

import (
	"bytes"
	"fmt"
)

// 表生成仓储结构,sign:函数后是否带签名，ePrefix:实体是否带前缀
func (s *Session) tableToGoRepo(table *Table,
	sign bool, ePrefix string) (string, string) {
	tpl := GoEntityRepTemplate
	path, _ := s.PredefineTargetPath(tpl, table)
	return s.GenerateCode(table, GoEntityRepTemplate,
		"Repo", sign, ePrefix), path
}

// 表生成仓库仓储接口
func (s *Session) tableToGoIRepo(table *Table,
	sign bool, ePrefix string) (string, string) {
	tpl := GoEntityRepIfceTemplate
	path, _ := s.PredefineTargetPath(tpl, table)
	return s.GenerateCode(table, tpl,
		"Repo", sign, ePrefix), path
}

// 表生成结构
func (s *Session) tableToGoStruct(table *Table) (string, string) {
	goPath := fmt.Sprintf("%s/model/%s.go", s.codeVars[PKG], table.Name)
	if table == nil {
		return "", goPath
	}
	pkgName := "model"
	buf := bytes.NewBufferString("")
	buf.WriteString("package ")
	buf.WriteString(pkgName)

	buf.WriteString("\n// ")
	buf.WriteString(table.Comment)
	buf.WriteString("\ntype ")
	buf.WriteString(s.title(table.Name))
	buf.WriteString(" struct{\n")

	for _, col := range table.Columns {
		if col.Comment != "" {
			buf.WriteString("    // ")
			buf.WriteString(col.Comment)
			buf.WriteString("\n")
		}
		buf.WriteString("    ")
		buf.WriteString(s.title(col.Name))
		buf.WriteString(" ")
		buf.WriteString(s.fn.langType("go", col.TypeId))
		buf.WriteString(" `")
		buf.WriteString("db:\"")
		buf.WriteString(col.Name)
		buf.WriteString("\"")
		if col.IsPk {
			buf.WriteString(" pk:\"yes\"")
		}
		if col.Auto {
			buf.WriteString(" auto:\"yes\"")
		}
		buf.WriteString("`")
		buf.WriteString("\n")
	}

	buf.WriteString("}")
	return buf.String(), goPath
}
