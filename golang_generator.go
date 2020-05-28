package tto

import (
	"bytes"
	"fmt"
	"sync"
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
	buf.WriteString(title(table.Name, s.IdUpper))
	buf.WriteString(" struct{\n")

	for _, col := range table.Columns {
		if col.Comment != "" {
			buf.WriteString("    // ")
			buf.WriteString(col.Comment)
			buf.WriteString("\n")
		}
		buf.WriteString("    ")
		buf.WriteString(title(col.Name, s.IdUpper))
		buf.WriteString(" ")
		buf.WriteString(s.fn.langType("go", col.Type))
		buf.WriteString(" `")
		buf.WriteString("db:\"")
		buf.WriteString(col.Name)
		buf.WriteString("\"")
		if col.IsPk {
			buf.WriteString(" pk:\"yes\"")
		}
		if col.IsAuto {
			buf.WriteString(" auto:\"yes\"")
		}
		buf.WriteString("`")
		buf.WriteString("\n")
	}

	buf.WriteString("}")
	return buf.String(), goPath
}

// 生成Go仓储代码
func (s *Session) GenerateGoRepoCodes(tables []*Table, targetDir string) (err error) {
	wg := sync.WaitGroup{}
	for _, table := range tables {
		wg.Add(1)
		go func(wg *sync.WaitGroup, tb *Table) {
			defer wg.Done()
			//生成实体
			str, path := s.tableToGoStruct(tb)
			if err = SaveFile(str, targetDir+"/"+path); err != nil {
				println(fmt.Sprintf("[ Gen][ Error]: save file failed! %s", err.Error()))
			}
			//生成仓储结构
			str, path = s.tableToGoRepo(tb, true, "")
			if err = SaveFile(str, targetDir+"/"+path); err != nil {
				println(fmt.Sprintf("[ Gen][ Error]: save file failed! %s", err.Error()))
			}
			//生成仓储接口
			str, path = s.tableToGoIRepo(tb, true, "")
			if err = SaveFile(str, targetDir+"/"+path); err != nil {
				println(fmt.Sprintf("[ Gen][ Error]: save file failed! %s", err.Error()))
			}
		}(&wg, table)
	}
	wg.Wait()
	// 生成仓储工厂
	code := s.GenerateCodeByTables(tables, GoRepoFactoryTemplate)
	path, _ := s.PredefineTargetPath(GoRepoFactoryTemplate, nil)
	return SaveFile(code, targetDir+"/"+path)
}
