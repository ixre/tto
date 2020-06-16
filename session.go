/**
 * Copyright 2015 @ at3.net.
 * name : tool.go
 * author : jarryliu
 * date : 2016-11-11 12:19
 * description :
 * history :
 */
package tto

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/ixre/gof/db/orm"
	"github.com/ixre/tto/config"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"sync"
	"text/template"
	"time"
)

var (
	emptyReg       = regexp.MustCompile("\\s+\"\\s*\"\\s*\\n")
	emptyImportReg = regexp.MustCompile("import\\s*\\(([\\n\\s\"]+)\\)")
)

const (
	// 包名
	PKG = "pkg"
	// 当前时间
	TIME = "time"
	// 版本
	VERSION = "version"
	//模型包名
	ModelPkgName = "ModelPkgName"
	//仓储结构包名
	RepoPkgName = "RepoPkgName"
	//仓储接口包名
	IfcePkgName = "IfacePkgName"
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
		//　主键属性
		PkProp string
		// 主键类型编号
		PkType int
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
		Prop string
		// 是否主键
		IsPk bool
		// 是否自动生成
		IsAuto bool
		// 是否不能为空
		NotNull bool
		// 类型
		DbType string
		// 注释
		Comment string
		// 长度
		Length int
		// Go类型
		Type int
		// 输出选项
		Render *config.PropRenderOptions
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
	fn := &internalFunc{}
	return (&Session{
		codeVars: make(map[string]interface{}),
		funcMap:  fn.funcMap(),
		IdUpper:  false,
	}).init()
}

func (s *Session) init() *Session {
	s.Var(ModelPkgName, "model")
	s.Var(RepoPkgName, "repo")
	s.Var(IfcePkgName, "ifce")

	// predefine default vars
	s.Var("url_prefix","")
	// load global registry
	rd := GetRegistry()
	for _,k := range rd.Keys{
		s.Var(k,rd.Data[k])
	}
	// put system vars
	s.Var(PKG, "com/tto/pkg")
	s.Var(TIME, time.Now().Format("2006/01/02 15:04:05"))
	s.Var(VERSION, BuildVersion)

	return s
}

// 获取所有的表
func (s *Session) ParseTables(tbs []*orm.Table, err error) ([]*Table, error) {
	n := make([]*Table, len(tbs))
	for i, tb := range tbs {
		n[i] = parseTable(i, tb, s.IdUpper, false)
	}
	return n, err
}

// 转换表格,如果meta为true,则读取元数据,如果没有则自动生成元数据
func (s *Session) Parses(tbs []*orm.Table, meta bool) (arr []*Table, err error) {
	n := make([]*Table, len(tbs))
	for i, tb := range tbs {
		n[i] = parseTable(i, tb, s.IdUpper, meta)
	}
	return n, err
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
	} else {
		s.codeVars[key] = v
	}
}

// 返回所有的变量
func (s *Session) AllVars() map[string]interface{} {
	return s.codeVars
}

// 转换成为模板
func (s *Session) parseTemplate(file string, attachCopy bool) (*CodeTemplate, error) {
	data, err := ioutil.ReadFile(file)
	if err == nil {
		return NewTemplate(string(data), file, attachCopy), nil
	}
	return NewTemplate("", file, attachCopy), err
}

// 生成代码
func (s *Session) GenerateCode(table *Table, tpl *CodeTemplate) string {
	if table == nil {
		return ""
	}
	var err error
	t := &template.Template{}
	t.Funcs(s.funcMap)
	t, err = t.Parse(tpl.template)
	if err != nil {
		panic(fmt.Sprintf("file:%s - error:%s", tpl.FilePath(), err.Error()))
	}
	//n := s.title(table.Name)
	mp := map[string]interface{}{
		"global": s.codeVars, // 全局变量
		"table":   table,         // 数据表
		"columns": table.Columns, // 列
	}
	buf := bytes.NewBuffer(nil)
	err = t.Execute(buf, mp)
	if err == nil {
		code := buf.String()
		//去除空引用
		code = emptyImportReg.ReplaceAllString(code, "")
		//如果不包含模型，则可能为引用空的包
		code = emptyReg.ReplaceAllString(code, "")
		return s.formatCode(tpl, code)
	}
	log.Println("execute template error:", err.Error())
	return ""
}

// 生成所有表的代码, 可引用的对象为global 和 tables
func (s *Session) GenerateCodeByTables(tables []*Table, tpl *CodeTemplate) string {
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
		return s.formatCode(tpl, code)
	}
	log.Println("execute template error:", err.Error())
	return ""
}

// 获取生成目标代码文件路径
func (s *Session) PredefineTargetPath(tpl *CodeTemplate, table *Table) (string, error) {
	if n, ok := tpl.Predefine("target"); ok {
		return ResolvePathString(n, s.AllVars(), table), nil
	}
	return "", errors.New("template not contain predefine command #!target")
}

// 连接文件路径
func (s *Session) defaultTargetPath(tplFilePath string, table *Table) string {
	i := strings.Index(tplFilePath, ".")
	if i != -1 {
		return strings.Join([]string{tplFilePath[:i], "_",
			table.Name, ".", tplFilePath[i+1:]}, "")
	}
	return strings.TrimSpace(tplFilePath + table.Name)
}

// 格式化代码
func (s *Session) formatCode(tpl *CodeTemplate, code string) string {
	// 不格式化代码
	if k, _ := tpl.Predefine("format"); k == "false" {
		return code
	}
	// 去除多行换行
	//code = regexp.MustCompile("(\Data?\n(\\s*\Data?\n)+)").ReplaceAllString(code, "\n\n")
	return code
}

// 遍历模板文件夹, 并生成代码, 如果为源代码目标,文件存在,则自动生成添加 .gen后缀
func (s *Session) WalkGenerateCode(tables []*Table, opt *GenerateOptions) error {
	opt.prepare()
	tplMap := map[string]*CodeTemplate{}
	sliceSize := len(opt.TplDir)
	if opt.TplDir[sliceSize-1] == '/' {
		opt.TplDir = opt.TplDir + "/"
		sliceSize += 1
	}
	err := filepath.Walk(opt.TplDir, func(path string, info os.FileInfo, err error) error {
		// 如果模板名称以"_"开头，则忽略
		if info != nil && !info.IsDir() && s.testName(info.Name(), opt.ExcludeFiles) {
			tp, err := s.parseTemplate(path, opt.AttachCopyright)
			if err != nil {
				return errors.New("template:" + info.Name() + "-" + err.Error())
			}
			s.Resolve(tp)
			tplMap[path[sliceSize:]] = tp
		}
		return nil
	})
	if err != nil {
		return err
	}
	if len(tplMap) == 0 {
		return errors.New("no any code template")
	}
	wg := sync.WaitGroup{}
	for _, tb := range tables {
		for path, tpl := range tplMap {
			wg.Add(1)
			go func(wg *sync.WaitGroup, tpl *CodeTemplate, tb *Table, path string) {
				defer wg.Done()
				str := s.GenerateCode(tb, tpl)
				dstPath, _ := s.PredefineTargetPath(tpl, tb)
				if dstPath == "" {
					dstPath = s.defaultTargetPath(path, tb)
				}
				if err := SaveFile(str, opt.OutputDir+"/"+dstPath); err != nil {
					println(fmt.Sprintf("[ Gen][ Error]: save file failed! %s ,template:%s",
						err.Error(), tpl.FilePath()))
				}
			}(&wg, tpl, tb, path)
		}
	}
	wg.Wait()
	return err
}

// 验证文件名, 是否可以生成
func (s *Session) testName(name string, files []string) bool {
	if name[0] == '_' {
		return false
	}
	if files != nil {
		for _, v := range files {
			if v == name {
				return false
			}
		}
		/*
			i := sort.Search(len(files), func(i int) bool {
				println("---",len(files),files[i],name)
				return files[i] == name
			})
			println("---", i, name, fmt.Sprintf("%#v", files), sort.SearchStrings(files, name))
			if files != nil && sort.SearchStrings(files, name) != -1 {
				return false
			}
		*/
	}
	return true
}

type GenerateOptions struct {
	TplDir          string
	AttachCopyright bool
	OutputDir       string
	ExcludeFiles    []string
}

func (g *GenerateOptions) prepare() {
	if len(g.ExcludeFiles) == 1 && g.ExcludeFiles[0] == "" {
		g.ExcludeFiles = nil
	}
	if g.TplDir == "" {
		g.TplDir = "./templates"
	}
	if g.OutputDir == "" {
		g.OutputDir = "./output"
	}
}
