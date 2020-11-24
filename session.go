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
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"sync"
	"text/template"
	"time"
)

const (
	// 基础URL
	BASE_URL = "base_url"
	// 基础路径
	BASE_PATH = "base_path"
	// 包名
	PKG = "pkg"
	// 当前时间
	TIME = "time"
	// 版本
	VERSION = "version"
)

type (
	Session interface {
		// add or update variable
		Var(key string, v interface{})
		// set code package
		Package(pkg string)
		// 使用大写ID,默认为false
		UseUpperId()
		// 添加函数
		AddFunc(fnTag string, fnBody interface{})
		// 返回所有的变量
		AllVars() map[string]interface{}
		// 生成代码
		GenerateCode(table *Table, tpl *CodeTemplate) string
		// 生成所有表的代码, 可引用的对象为global 和 tables
		GenerateCodeByTables(tables []*Table, tpl *CodeTemplate) string
		// 遍历模板文件夹, 并生成代码, 如果为源代码目标,文件存在,则自动生成添加 .gen后缀
		WalkGenerateCodes(tables []*Table, opt *GenerateOptions) error
		// 转换表格,如果meta为true,则读取元数据,如果没有则自动生成元数据
		Parses(tables []*orm.Table, meta bool) (arr []*Table, err error)
	}
)

var _ Session = new(sessionImpl)

type sessionImpl struct {
	// 生成代码变量
	codeVars map[string]interface{}
	// 模板函数
	funcMap map[string]interface{}
	// 使用大写ID
	useUpperId bool
	// 数据库驱动
	driver     string
}

func (s *sessionImpl) UseUpperId() {
	s.useUpperId = true
}

// 数据库代码生成器
func DBCodeGenerator(driver string) Session {
	if sort.SearchStrings([]string{"pgsql","mssql","mysql"},driver) == -1{
		panic("not support db :"+driver)
	}
	return (&sessionImpl{
		driver: driver,
		codeVars:   make(map[string]interface{}),
		funcMap:    (&internalFunc{}).funcMap(),
		useUpperId: false,
	}).init()
}

func (s *sessionImpl) init() Session {
	// predefine default vars
	s.Var(BASE_URL, "")
	s.Var(BASE_PATH, "")
	// load global registry
	rd := GetRegistry()
	for _, k := range rd.Keys {
		s.Var(k, rd.Data[k])
	}
	// put system vars
	s.Package("com/your/pkg")
	s.Var("db",s.driver)
	s.Var(TIME, time.Now().Format("2006/01/02 15:04:05"))
	s.Var(VERSION, BuildVersion)
	return s
}

// 转换表格,如果meta为true,则读取元数据,如果没有则自动生成元数据
func (s *sessionImpl) Parses(tables []*orm.Table, meta bool) (arr []*Table, err error) {
	n := make([]*Table, len(tables))
	for i, tb := range tables {
		n[i] = parseTable(i, tb, s.useUpperId, meta)
	}
	return n, err
}

// 添加函数
func (s *sessionImpl) AddFunc(funcName string, f interface{}) {
	s.funcMap[funcName] = f
}

// 定义变量或修改变量
func (s *sessionImpl) Var(key string, v interface{}) {
	if key == "pkg"{
		panic("please use Package(pkg string)")
	}
	if v == nil {
		delete(s.codeVars, key)
	} else {
		s.codeVars[key] = v
	}
}

func (s sessionImpl) Package(pkg string) {
	pkg = strings.ReplaceAll(pkg,".","/")
	s.codeVars[PKG] = pkg
}

// 返回所有的变量
func (s *sessionImpl) AllVars() map[string]interface{} {
	return s.codeVars
}

// 转换成为模板
func (s *sessionImpl) parseTemplate(file string, attachCopy bool) (*CodeTemplate, error) {
	data, err := ioutil.ReadFile(file)
	if err == nil {
		return NewTemplate(string(data), file, attachCopy), nil
	}
	return NewTemplate("", file, attachCopy), err
}

// 生成代码
func (s *sessionImpl) GenerateCode(table *Table, tpl *CodeTemplate) string {
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
func (s *sessionImpl) GenerateCodeByTables(tables []*Table, tpl *CodeTemplate) string {
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
func (s *sessionImpl) predefineTargetPath(tpl *CodeTemplate, table *Table) (string, error) {
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
func (s *sessionImpl) defaultTargetPath(tplFilePath string, table *Table) string {
	i := strings.Index(tplFilePath, ".")
	if i != -1 {
		return strings.Join([]string{tplFilePath[:i], "_",
			table.Name, ".", tplFilePath[i+1:]}, "")
	}
	return strings.TrimSpace(tplFilePath + table.Name)
}

var multiLineRegexp = regexp.MustCompile("\\{\\n{2,}(\\s{4}\\s+)")

// 格式化代码
func (s *sessionImpl) formatCode(tpl *CodeTemplate, code string) string {
	// 去除`{`后多余的换行
	code = multiLineRegexp.ReplaceAllString(code,"{\n$1")
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
func (s *sessionImpl) WalkGenerateCodes(tables []*Table, opt *GenerateOptions) error {
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

func (s *sessionImpl) generateTables(ch chan int, tables []*Table, opt *GenerateOptions, tplMap map[string]*CodeTemplate) {
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
func (s *sessionImpl) generateAllTablesCode(ch chan int, tables []*Table, tplMap map[string]*CodeTemplate, opt *GenerateOptions, rc *RunnerCalc) {
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
func (s *sessionImpl) generateGroupTablesCode(ch chan int, tables []*Table, tplMap map[string]*CodeTemplate, opt *GenerateOptions, rc *RunnerCalc) {
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
		// 文件类型不匹配则跳过
		if tpl.Kind() != KindTablePrefix {
			continue
		}
		for prefix, tables := range groups {
			key := path + "$" + prefix
			if rc.State(key) {
				break
			}
			wg.Add(1)
			//time.Sleep(time.Second/5)
			go func(wg *sync.WaitGroup, tpl *CodeTemplate, tables []*Table,
				rc *RunnerCalc, path string, key string) {
				defer wg.Done()
				rc.SignState(tpl.path, true)
				out := s.GenerateCodeByTables(tables, tpl)
				s.flushToFile(tpl, tables[0], path, out, opt)
			}(&wg, tpl, tables, rc, path, key)
		}
	}
	wg.Wait()
	ch <- 1
}

func (s *sessionImpl) flushToFile(tpl *CodeTemplate, tb *Table, path string, output string, opt *GenerateOptions) {
	dstPath, _ := s.predefineTargetPath(tpl, tb)
	if dstPath == "" {
		dstPath = s.defaultTargetPath(path, tb)
	}
	savedPath := filepath.Clean(opt.OutputDir + "/" + dstPath)
	if err := SaveFile(output, savedPath); err != nil {
		println(fmt.Sprintf("[ tto][ error]: save file failed! %s ,template:%s", err.Error(), tpl.FilePath()))
	}
}

func (s *sessionImpl) findTemplates(opt *GenerateOptions) (map[string]*CodeTemplate, error) {
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
func (s *sessionImpl) testName(name string, files []string) bool {
	if name[0] == '_' {
		return false
	}
	if strings.ToUpper(name) =="README.MD"{
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
