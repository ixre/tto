/**
 * Copyright 2015 @ at3.net.
 * name : tool.go
 * author : jarryliu
 * date : 2016-11-11 12:19
 * description :
 * history :
 */
package generator

import (
	"bytes"
	"github.com/ixre/gof/db/orm"
	"io/ioutil"
	"log"
	"regexp"
	"strings"
	"text/template"
	"unicode"
)

var (
	emptyReg             = regexp.MustCompile("\\s+\"\\s*\"\\s*\\n")
	emptyImportReg       = regexp.MustCompile("import\\s*\\(([\\n\\s\"]+)\\)")
	revertTemplateRegexp = regexp.MustCompile("([^\\$])\\${([^\\}]+)\\}")
)

const (
	// 包名
	PKG = "Pkg"
	// 版本
	VERSION = "Version"
	//模型包名
	VModelPkgName = "ModelPkgName"
	//仓储结构包名
	VRepoPkgName = "RepoPkgName"
	//仓储接口包名
	VIRepoPkgName = "IRepoPkgName"
	//仓储结构引用模型包路径
	VModelPkg = "ModelPkg"
	//仓储接口引用模型包路径
	VIRepoPkg = "IRepoPkg"
	// 仓储包路径
	VRepoPkg = "RepoPkg"
)

type (
	// 表
	Table struct {
		// 顺序
		Ordinal int
		// 表名
		Name string
		// 表前缀
		Prefix string
		// 表名单词首字大写
		Title string
		// 表注释
		Comment string
		// 数据库引擎
		Engine string
		// 架构
		Schema string
		// 数据库编码
		Charset string
		// 表
		Raw *orm.Table
		// 主键
		Pk string
		// 主键类型编号
		PkTypeId int
		// 列
		Columns []*Column
	}
	// 列
	Column struct {
		// 顺序
		Ordinal int
		// 列名
		Name string
		// 列名首字大写
		Title string
		// 是否主键
		IsPk bool
		// 是否自动生成
		Auto bool
		// 是否不能为空
		NotNull bool
		// 类型
		Type string
		// 注释
		Comment string
		// 长度
		Length int
		// Go类型
		TypeId int
	}
)
type Session struct {
	// 生成代码变量
	codeVars map[string]interface{}
	// 模板函数
	funcMap map[string]interface{}
	IdUpper bool
}

// 数据库代码生成器
func DBCodeGenerator() *Session {
	return (&Session{
		codeVars: make(map[string]interface{}),
		funcMap:  (&internalFunc{}).funcMap(),
		IdUpper:  false,
	}).init()
}

func (s *Session) init() *Session {
	s.Var(PKG, "com/pkg")
	s.Var(VERSION, "1.0")
	s.Var(VModelPkgName, "model")
	s.Var(VRepoPkgName, "repo")
	s.Var(VIRepoPkgName, "repo")
	s.Var(VModelPkg, "")
	s.Var(VIRepoPkg, "")
	s.Var(VRepoPkg, "")
	return s
}

func (s *Session) title(str string) string {
	// 小于3且ID大写，则返回大写
	if s.IdUpper && len(str) < 3 {
		return strings.ToUpper(str)
	}
	arr := strings.Split(str, "_")
	for i, v := range arr {
		arr[i] = strings.Title(v)
	}
	return strings.Join(arr, "")
}

func (s *Session) prefix(str string) string {
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

func (s *Session) goType(goType int) string {
	switch goType {
	case orm.TypeString:
		return "string"
	case orm.TypeBoolean:
		return "bool"
	case orm.TypeInt16:
		return "int16"
	case orm.TypeInt32:
		return "int"
	case orm.TypeInt64:
		return "int64"
	case orm.TypeFloat32:
		return "float32"
	case orm.TypeFloat64:
		return "float64"
	}
	return "interface{}"
}

// 获取所有的表
func (s *Session) ParseTables(tbs []*orm.Table, err error) ([]*Table, error) {
	n := make([]*Table, len(tbs))
	for i, tb := range tbs {
		n[i] = s.parseTable(i, tb)
	}
	return n, err
}

// 获取表结构
func (s *Session) parseTable(ordinal int, tb *orm.Table) *Table {
	n := &Table{
		Name:     tb.Name,
		Prefix:   s.prefix(tb.Name),
		Title:    s.title(tb.Name),
		Comment:  tb.Comment,
		Engine:   tb.Engine,
		Schema:   tb.Schema,
		Charset:  tb.Charset,
		Raw:      tb,
		Pk:       "id",
		PkTypeId: orm.TypeInt32,
		Columns:  make([]*Column, len(tb.Columns)),
	}
	for i, v := range tb.Columns {
		if v.IsPk && n.Pk == "" {
			n.Pk = v.Name
			n.PkTypeId = v.TypeId
		}
		n.Columns[i] = &Column{
			Ordinal: i,
			Name:    v.Name,
			Title:   s.title(v.Name),
			IsPk:    v.IsPk,
			Auto:    v.Auto,
			NotNull: v.NotNull,
			Type:    v.Type,
			Comment: v.Comment,
			Length:  v.Length,
			TypeId:  v.TypeId,
		}
	}
	return n
}

// 表生成结构
func (s *Session) TableToGoStruct(tb *Table) string {
	if tb == nil {
		return ""
	}
	pkgName := ""
	if p, ok := s.codeVars[VModelPkgName]; ok {
		pkgName = p.(string)
	} else {
		pkgName = "model"
	}

	//log.Println(fmt.Sprintf("%#v", tb))
	buf := bytes.NewBufferString("")
	buf.WriteString("package ")
	buf.WriteString(pkgName)

	buf.WriteString("\n// ")
	buf.WriteString(tb.Comment)
	buf.WriteString("\ntype ")
	buf.WriteString(s.title(tb.Name))
	buf.WriteString(" struct{\n")

	for _, col := range tb.Columns {
		if col.Comment != "" {
			buf.WriteString("    // ")
			buf.WriteString(col.Comment)
			buf.WriteString("\n")
		}
		buf.WriteString("    ")
		buf.WriteString(s.title(col.Name))
		buf.WriteString(" ")
		buf.WriteString(s.goType(col.TypeId))
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
	return buf.String()
}

// 解析模板
func (s *Session) Resolve(t *CodeTemplate) *CodeTemplate {
	return t
}

// 添加函数
func (s *Session) Func(funcName string, f interface{}) {
	s.funcMap[funcName] = f
}

// 定义变量或修改变量
func (s *Session) Var(key string, v interface{}) {
	if v == nil {
		delete(s.codeVars, key)
		return
	}
	//if strings.HasSuffix(key, "PkgName") {
	//	if s := v.(string); s != "" && s[len(s)-1] != '.' {
	//		v = s + "."
	//	}
	//}
	s.codeVars[key] = v
}

// 返回所有的变量
func (s *Session) AllVars() map[string]interface{} {
	return s.codeVars
}

// 还原模板的标签: ${...} -> {{...}}, $$ -> $
func (s *Session) revertTemplateVariable(str string) string {
	str = revertTemplateRegexp.ReplaceAllString(str, "$1{{$2}}")
	return strings.Replace(str, "$$", "$", -1)
}

// 转换成为模板
func (s *Session) ParseTemplate(file string) (*CodeTemplate, error) {
	data, err := ioutil.ReadFile(file)
	if err == nil {
		return NewTemplate(string(data)), nil
	}
	return NewTemplate(""), err
}

// 生成代码
func (s *Session) GenerateCode(tb *Table, tpl *CodeTemplate,
	structSuffix string, sign bool, ePrefix string) string {
	if tb == nil {
		return ""
	}
	var err error
	t := &template.Template{}
	t.Funcs(s.funcMap)
	t, err = t.Parse(tpl.template)
	if err != nil {
		panic(err)
	}
	columns := tb.Columns
	//n := s.title(tb.Name)
	n := tb.Title
	r2 := ""
	if sign {
		r2 = n
	}
	mp := map[string]interface{}{
		"global": s.codeVars, // 全局变量
		//"version": s.codeVars[VERSION], // 版本
		//"pkg":     s.codeVars[PKG],     //包名
		"table":   tb,      // 数据表
		"columns": columns, // 列
		//"pk":      tb.Pk,                  // 主键列名

		//---------- 旧的字段 ------------------//
		"VAR":  s.codeVars, // 全局变量
		"T":    tb,         // 数据表
		"R":    n + structSuffix,
		"R2":   r2,
		"E":    n,
		"E2":   ePrefix + n,
		"Ptr":  strings.ToLower(tb.Name[:1]),
		"IsPK": s.title(tb.Pk),
	}
	buf := bytes.NewBuffer(nil)
	err = t.Execute(buf, mp)
	if err == nil {
		code := buf.String()
		//去除空引用
		code = emptyImportReg.ReplaceAllString(code, "")
		//如果不包含模型，则可能为引用空的包
		code = emptyReg.ReplaceAllString(code, "")
		return s.revertTemplateVariable(code)
	}
	log.Println("execute template error:", err.Error())
	return ""
}

func (s *Session) GenerateTablesCode(tables []*Table, tpl *CodeTemplate) string {
	if tables == nil || len(tables) == 0 {
		return ""
	}
	var err error
	t := &template.Template{}
	t.Funcs(s.funcMap)
	t, err = t.Parse(tpl.template)
	if err != nil {
		panic(err)
	}
	mp := map[string]interface{}{
		"VAR":    s.codeVars,
		"Tables": tables,
		"tables": tables,
		"global": s.codeVars,
	}
	buf := bytes.NewBuffer(nil)
	err = t.Execute(buf, mp)
	if err == nil {
		code := buf.String()
		//去除空引用
		code = emptyImportReg.ReplaceAllString(code, "")
		//如果不包含模型，则可能为引用空的包
		code = emptyReg.ReplaceAllString(code, "")
		return s.revertTemplateVariable(code)
	}
	log.Println("execute template error:", err.Error())
	return ""
}

// 表生成仓储结构,sign:函数后是否带签名，ePrefix:实体是否带前缀
func (s *Session) TableToGoRepo(tb *Table,
	sign bool, ePrefix string) string {
	return s.GenerateCode(tb, TPL_ENTITY_REP,
		"Repo", sign, ePrefix)
}

// 表生成仓库仓储接口
func (s *Session) TableToGoIRepo(tb *Table,
	sign bool, ePrefix string) string {
	return s.GenerateCode(tb, TPL_ENTITY_REP_INTERFACE,
		"Repo", sign, ePrefix)
}
