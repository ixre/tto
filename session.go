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
	"strings"
	"sync"
	"text/template"
	"time"
)


const (
	// 包名
	PKG = "pkg"
	// 当前时间
	TIME = "time"
	// 版本
	VERSION = "version"
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
	// predefine default vars
	s.Var("url_prefix", "")
	// load global registry
	rd := GetRegistry()
	for _, k := range rd.Keys {
		s.Var(k, rd.Data[k])
	}
	// put system vars
	s.Var(PKG, "com/tto/pkg")
	s.Var(TIME, time.Now().Format("2006/01/02 15:04:05"))
	s.Var(VERSION, BuildVersion)

	return s
}

// 转换表格,如果meta为true,则读取元数据,如果没有则自动生成元数据
func (s *Session) Parses(tbs []*orm.Table, meta bool) (arr []*Table, err error) {
	n := make([]*Table, len(tbs))
	for i, tb := range tbs {
		n[i] = parseTable(i, tb, s.IdUpper, meta)
	}
	return n, err
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
		log.Println("[ app][ fatal]: " + fmt.Sprintf("file:%s - error:%s", tpl.FilePath(), err.Error()))
		return ""
	}
	//n := s.title(table.Name)
	mp := map[string]interface{}{
		"global":  s.codeVars,    // 全局变量
		"table":   table,         // 数据表
		"columns": table.Columns, // 列
	}
	buf := bytes.NewBuffer(nil)
	err = t.Execute(buf, mp)
	if err == nil {
		return s.formatCode(tpl, buf.String())
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
	t := (&template.Template{}).Funcs(s.funcMap)
	t, err = t.Parse(tpl.template)
	if err != nil {
		log.Println("[ app][ fatal]: " + fmt.Sprintf("file:%s - error:%s", tpl.FilePath(), err.Error()))
		return ""
	}
	mp := map[string]interface{}{
		"tables": tables,
		"global": s.codeVars,
	}
	buf := bytes.NewBuffer(nil)
	err = t.Execute(buf, mp)
	if err == nil {
		return s.formatCode(tpl, buf.String())
	}
	log.Println("execute template error:", err.Error())
	return ""
}

// 获取生成目标代码文件路径
func (s *Session) PredefineTargetPath(tpl *CodeTemplate, table *Table) (string, error) {
	n, ok := tpl.Predefine("target")
	if !ok {
		return "", errors.New("template not contain predefine command #!target")
	}
	var t = (&template.Template{}).Funcs(s.funcMap)
	t, err := t.Parse(n)
	if err == nil {
		mp := map[string]interface{}{
			"global": s.AllVars(),
			"table":  table,
			"prefix": "",
		}
		// 添加前缀
		if table != nil {
			mp["prefix"] = table.Prefix
		}
		buf := bytes.NewBuffer(nil)
		err = t.Execute(buf, mp)
		return strings.TrimSpace(buf.String()), err
	}
	return "", err
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
	//emptyReg       = regexp.MustCompile("\\s+\"\\s*\"\\s*\\n")
	//regexp.MustCompile("import\\s*\\(([\\n\\s\"]+)\\)")

	//如果不包含模型，则可能为引用空的包
	//code = emptyReg.ReplaceAllString(code, "")
	// 去除多行换行
	//code = regexp.MustCompile("(\Data?\n(\\s*\Data?\n)+)").ReplaceAllString(code, "\n\n")
	return code
}

// 遍历模板文件夹, 并生成代码, 如果为源代码目标,文件存在,则自动生成添加 .gen后缀
func (s *Session) WalkGenerateCode(tables []*Table, opt *GenerateOptions) error {
	opt.prepare()
	tplMap, err := s.findTemplates(opt)
	if err != nil {
		return err
	}
	rc := NewRunnerCalc()
	ch := make(chan int, 3)
	s.generateAllTablesCode(ch, tables, tplMap, opt, &rc)
	s.generateGroupTablesCode(ch, tables, tplMap, opt, &rc)
	s.generateTables(ch, tables, opt, tplMap)
	<-ch
	<-ch
	<-ch
	return err
}

func (s *Session) generateTables(ch chan int, tables []*Table, opt *GenerateOptions, tplMap map[string]*CodeTemplate) {
	wg := sync.WaitGroup{}
	for path, tpl := range tplMap {
		if tpl.Kind() == KindNormal {
			for _, tb := range tables {
				wg.Add(1)
				go func(wg *sync.WaitGroup, tpl *CodeTemplate, tb *Table, path string) {
					defer wg.Done()
					out := s.GenerateCode(tb, tpl)
					s.flushToFile(tpl, tb, path, out, opt)
				}(&wg, tpl, tb, path)
			}
		}
	}
	wg.Wait()
	ch <- 1
}

// 生成所有表的代码
func (s *Session) generateAllTablesCode(ch chan int, tables []*Table, tplMap map[string]*CodeTemplate, opt *GenerateOptions, rc *RunnerCalc) {
	wg := sync.WaitGroup{}
	for path, tpl := range tplMap {
		if tpl.Kind() != KindTables || rc.State(path) {
			continue // 如果生成类型不符合或已经生成,跳过
		}
		wg.Add(1)
		go func(wg *sync.WaitGroup, tpl *CodeTemplate, tb []*Table, rc *RunnerCalc, path string) {
			defer wg.Done()
			rc.SignState(tpl.path, true)
			out := s.GenerateCodeByTables(tables, tpl)
			s.flushToFile(tpl, nil, path, out, opt)
		}(&wg, tpl, tables, rc, path)
	}
	wg.Wait()
	ch <- 1
}

// 按表前缀分组生成代码
func (s *Session) generateGroupTablesCode(ch chan int, tables []*Table, tplMap map[string]*CodeTemplate, opt *GenerateOptions, rc *RunnerCalc) {
	//　分组,无前缀的所有表归类到一组
	groups := make(map[string][]*Table, 0)
	for _, t := range tables {
		prefix := t.Prefix
		if arr, ok := groups[prefix]; ok {
			groups[prefix] = append(arr, t)
		} else {
			groups[prefix] = []*Table{t}
		}
	}
	wg := sync.WaitGroup{}
	for path, tpl := range tplMap {
		if tpl.Kind() != KindTablePrefix {
			continue
		}
		for prefix, tbs := range groups {
			key := path + "$" + prefix
			if rc.State(key) {
				break
			}
			wg.Add(1)
			//time.Sleep(time.Second/5)
			go func(wg *sync.WaitGroup, tpl *CodeTemplate, tbs []*Table,
				rc *RunnerCalc, path string, key string) {
				defer wg.Done()
				rc.SignState(tpl.path, true)
				out := s.GenerateCodeByTables(tbs, tpl)
				s.flushToFile(tpl, tbs[0], path, out, opt)
			}(&wg, tpl, tbs, rc, path, key)
		}
	}
	wg.Wait()
	ch <- 1
}

func (s *Session) flushToFile(tpl *CodeTemplate, tb *Table, path string, output string, opt *GenerateOptions) {
	dstPath, _ := s.PredefineTargetPath(tpl, tb)
	if dstPath == "" {
		dstPath = s.defaultTargetPath(path, tb)
	}
	savedPath := filepath.Clean(opt.OutputDir + "/" + dstPath)
	if err := SaveFile(output, savedPath); err != nil {
		println(fmt.Sprintf("[ tto][ error]: save file failed! %s ,template:%s", err.Error(), tpl.FilePath()))
	}
}

func (s *Session) findTemplates(opt *GenerateOptions) (map[string]*CodeTemplate, error) {
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
			tplMap[path[sliceSize-1:]] = tp
		}
		return nil
	})
	if err == nil && len(tplMap) == 0 {
		err = errors.New("no any code template")
	}
	return tplMap, err
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
