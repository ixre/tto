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
	"log"
	"os"
	"os/user"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"sync"
	"text/template"
	"time"

	"github.com/ixre/gof/db/db"
)

const (
	// 基础URL
	BASE_URL = "base_url"
	// 基础路径
	BASE_PATH = "base_path"
	// 组织名
	ORGANIZATION = "organization"
	// 包名
	PKG = "pkg"
	// 实体后缀,默认：Entity
	ENTITY_SUFFIX = "entity_suffix"
	// 当前时间
	TIME = "time"
	// 版本
	VERSION = "version"
	// 当前用户
	USER = "user"
)

type (
	// 代码处理
	GenerateHandler func(result string, tpl *Template) (ret string)

	// Session 生成器会话
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
		GenerateCode(table *Table, tpl *Template, g GenerateHandler) string
		// 生成所有表的代码, 可引用的对象为global 和 tables
		GenerateCodeByTables(tables []*Table, tpl *Template, g GenerateHandler) string
		// 遍历模板文件夹, 并生成代码, 如果为源代码目标,文件存在,则自动生成添加 .gen后缀
		WalkGenerateCodes(tables []*Table, g GenerateHandler) error
		// 转换表格,如果meta为true,则读取元数据,如果没有则自动生成元数据
		Parses(tables []*db.Table, meta bool) (arr []*Table, err error)
		// 清理生成目录
		Clean() error
	}
)

type Options struct {
	TplDir          string
	AttachCopyright bool
	// 输出目录
	OutputDir string
	// 默认语言,当设置为"java",则会启用小写命名
	MajorLang string
	// 排除模板文件模式,如：*/tmp,tmp/*或者tmp
	ExcludePatterns []string
}

func (g *Options) prepare() {
	if len(g.ExcludePatterns) == 1 && g.ExcludePatterns[0] == "" {
		g.ExcludePatterns = nil
	}
	g.MajorLang = strings.ToLower(strings.TrimSpace(g.MajorLang))
	if len(g.MajorLang) == 0 {
		g.MajorLang = "java"
	}
	if g.TplDir == "" {
		g.TplDir = "./templates"
	}
	if g.OutputDir == "" {
		g.OutputDir = "./output"
	}
}

var _ Session = new(sessionImpl)

type sessionImpl struct {
	// 生成代码变量
	codeVars map[string]interface{}
	// 模板函数
	funcMap map[string]interface{}
	// 使用大写ID
	useUpperId bool
	// 数据库驱动
	driver string
	opt    *Options
}

func (s *sessionImpl) UseUpperId() {
	s.useUpperId = true
}

// 数据库代码生成器
func DBCodeGenerator(driver string, opt *Options) Session {
	if opt == nil {
		opt = &Options{}
	}
	opt.prepare()
	if sort.SearchStrings([]string{"pgsql", "mssql", "mysql"}, driver) == -1 {
		panic("not support db :" + driver)
	}
	return (&sessionImpl{
		driver:     driver,
		opt:        opt,
		codeVars:   make(map[string]interface{}),
		funcMap:    (&internalFunc{}).funcMap(),
		useUpperId: false,
	}).init()
}

func (s *sessionImpl) init() Session {
	// predefine default vars
	s.Var(BASE_URL, "")
	s.Var(BASE_PATH, "")
	s.Var(ENTITY_SUFFIX, "Entity")
	s.Var(ORGANIZATION, "FZE.NET")
	// load global registry
	rd := GetRegistry()
	for _, k := range rd.Keys {
		s.Var(k, rd.Data[k])
	}
	// put system vars
	s.Package("com/your/pkg")
	s.Var(TIME, time.Now().Format("2006/01/02 15:04:05"))
	s.Var(VERSION, BuildVersion)
	// setting user variables
	if usr, err := user.Current(); err == nil {
		s.Var(USER, usr.Username)
	} else {
		s.Var(USER, "")
	}
	s.Var("db", s.driver)
	s.Var("year", time.Now().Format("2006"))
	s.Var(ORGANIZATION, "FZE.NET")
	return s
}

// 转换表格,如果meta为true,则读取元数据,如果没有则自动生成元数据
func (s *sessionImpl) Parses(tables []*db.Table, meta bool) (arr []*Table, err error) {
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
	if key == "pkg" {
		panic("please use Package(pkg string)")
	}
	if v == nil {
		delete(s.codeVars, key)
	} else {
		s.codeVars[key] = v
	}
}

func (s sessionImpl) Package(pkg string) {
	pkg = strings.ReplaceAll(pkg, ".", "/")
	s.codeVars[PKG] = pkg
}

// 返回所有的变量
func (s *sessionImpl) AllVars() map[string]interface{} {
	return s.codeVars
}

// 转换成为模板
func (s *sessionImpl) parseTemplate(file string, attachCopy bool) (*Template, error) {
	data, err := os.ReadFile(file)
	if err == nil {
		return NewTemplate(string(data), file, attachCopy), nil
	}
	return NewTemplate("", file, attachCopy), err
}

// 生成代码
func (s *sessionImpl) GenerateCode(table *Table, tpl *Template, g GenerateHandler) string {
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
	tb := s.adapterTable(table, tpl.path)
	mp := map[string]interface{}{
		"global":  s.codeVars, // 全局变量
		"table":   tb,         // 数据表
		"columns": tb.Columns, // 列
	}
	buf := bytes.NewBuffer(nil)
	err = t.Execute(buf, mp)
	if err == nil {
		ret := s.formatCode(tpl, buf.String())
		if g != nil {
			// 格式化代码
			return g(ret, tpl)
		}
		return ret
	}
	log.Println("execute template error:", err.Error())
	return ""
}

// 生成所有表的代码, 可引用的对象为global 和 tables
func (s *sessionImpl) GenerateCodeByTables(tables []*Table, tpl *Template, g GenerateHandler) string {
	if len(tables) == 0 {
		return ""
	}
	var err error
	t := (&template.Template{}).Funcs(s.funcMap)
	t, err = t.Parse(tpl.template)
	if err != nil {
		log.Println("[ app][ fatal]: " + fmt.Sprintf("file:%s - error:%s", tpl.FilePath(), err.Error()))
		return ""
	}
	tbs := make([]*Table, len(tables))
	for i, v := range tables {
		tbs[i] = s.adapterTable(v, tpl.path)
	}
	groups := s.reduceGroup(tbs)
	mp := map[string]interface{}{
		"groups": groups,
		"tables": tbs,
		"global": s.codeVars,
	}
	buf := bytes.NewBuffer(nil)
	err = t.Execute(buf, mp)
	if err == nil {
		ret := s.formatCode(tpl, buf.String())
		if g != nil {
			return g(ret, tpl)
		}
		return ret
	}
	log.Println("execute template error:", err.Error())
	return ""
}

// 获取生成目标代码文件路径
func (s *sessionImpl) predefineTargetPath(tpl *Template, table *Table) (string, error) {
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
	if table == nil {
		return strings.ReplaceAll(tplFilePath, ".tpl", "")
	}
	i := strings.Index(tplFilePath, ".")
	if i != -1 {
		p1 := tplFilePath[:i]
		// 去掉模板文件名的".tpl"后缀
		p2 := strings.Replace(tplFilePath[i+1:], ".tpl", "", -1)
		return strings.Join([]string{p1, "_", table.Name, ".", p2}, "")
	}
	return strings.TrimSpace(tplFilePath + table.Name)
}

var multiLineRegexp = regexp.MustCompile(`(\{|,|>)[\n\r]+?\s*\n+`)
var multiLineRevertRegexp = regexp.MustCompile(`\n+\s*[\n\r]+?\}`)

// 格式化代码
func (s *sessionImpl) formatCode(tpl *Template, code string) string {

	// 去除`{`后多余的换行
	code = multiLineRegexp.ReplaceAllString(code, "$1\n")
	// 去除`}`前多余的换行
	code = multiLineRevertRegexp.ReplaceAllString(code, "\n}")
	//code = multiLineEndRegexp.ReplaceAllString(code, "$1\n}")
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
func (s *sessionImpl) WalkGenerateCodes(tables []*Table, g GenerateHandler) error {
	tplMap, err := s.findTemplates(s.opt)
	if err != nil {
		return err
	}
	rc := NewRunnerCalc()
	ch := make(chan int, 3)
	s.generateAllTablesCode(ch, tables, tplMap, s.opt, &rc, g)
	s.generateGroupTablesCode(ch, tables, tplMap, s.opt, &rc, g)
	s.generateTables(ch, tables, s.opt, tplMap, g)
	<-ch
	<-ch
	<-ch
	return err
}

func (s *sessionImpl) generateTables(ch chan int, tables []*Table,
	opt *Options, tplMap map[string]*Template,
	g GenerateHandler) {
	wg := sync.WaitGroup{}
	for path, tpl := range tplMap {
		if tpl.Kind() == KindNormal {
			for _, tb := range tables {
				wg.Add(1)
				go func(wg *sync.WaitGroup, tpl *Template, tb *Table, path string) {
					defer wg.Done()
					out := s.GenerateCode(tb, tpl, g)
					s.flushToFile(tpl, tb, path, out, opt)
				}(&wg, tpl, tb, path)
			}
		}
	}
	wg.Wait()
	ch <- 1
}

// 生成所有表的代码
func (s *sessionImpl) generateAllTablesCode(ch chan int, tables []*Table,
	tplMap map[string]*Template,
	opt *Options, rc *RunnerCalc,
	g GenerateHandler) {
	wg := sync.WaitGroup{}
	for path, tpl := range tplMap {
		if tpl.Kind() != KindTables || rc.State(path) {
			continue // 如果生成类型不符合或已经生成,跳过
		}
		wg.Add(1)
		go func(wg *sync.WaitGroup, tpl *Template, tb []*Table, rc *RunnerCalc, path string) {
			defer wg.Done()
			rc.SignState(tpl.path, true)
			out := s.GenerateCodeByTables(tb, tpl, g)
			s.flushToFile(tpl, nil, path, out, opt)
		}(&wg, tpl, tables, rc, path)
	}
	wg.Wait()
	ch <- 1
}

// 按表前缀分组生成代码
func (s *sessionImpl) generateGroupTablesCode(ch chan int, tables []*Table,
	tplMap map[string]*Template,
	opt *Options, rc *RunnerCalc,
	g GenerateHandler) {
	groups := s.reduceGroup(tables)
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
			go func(wg *sync.WaitGroup, tpl *Template, tables []*Table,
				rc *RunnerCalc, path string, key string) {
				defer wg.Done()
				rc.SignState(tpl.path, true)
				out := s.GenerateCodeByTables(tables, tpl, g)
				s.flushToFile(tpl, tables[0], path, out, opt)
			}(&wg, tpl, tables, rc, path, key)
		}
	}
	wg.Wait()
	ch <- 1
}

// 分组,无前缀的所有表归类到一组
func (s *sessionImpl) reduceGroup(tables []*Table) map[string][]*Table {
	groups := make(map[string][]*Table, 0)
	for _, t := range tables {
		prefix := t.Prefix
		if arr, ok := groups[prefix]; ok {
			groups[prefix] = append(arr, t)
		} else {
			groups[prefix] = []*Table{t}
		}
	}
	return groups
}

// 输出到文件
func (s *sessionImpl) flushToFile(tpl *Template, tb *Table, path string, output string, opt *Options) {
	// 获取文件的路径和名称
	dstPath, _ := s.predefineTargetPath(tpl, tb)
	if dstPath == "" {
		dstPath = s.defaultTargetPath(path, tb)
	}
	fileName := dstPath[strings.LastIndexAny(dstPath, "/\\")+1:]
	output = strings.Replace(output, "$file_name$", fileName, 1) // 替换文件名占位符
	savedPath := filepath.Clean(opt.OutputDir + "/" + dstPath)
	if err := SaveFile(output, savedPath); err != nil {
		println(fmt.Sprintf("[ tto][ error]: save file failed! %s ,template:%s", err.Error(), tpl.FilePath()))
	}
}

func (s *sessionImpl) findTemplates(opt *Options) (map[string]*Template, error) {
	tplMap := map[string]*Template{}
	sliceSize := len(opt.TplDir)
	if opt.TplDir[sliceSize-1] == '/' {
		opt.TplDir = opt.TplDir + "/"
		sliceSize += 1
	}
	err := filepath.Walk(opt.TplDir, func(path string, info os.FileInfo, err error) error {
		// 如果模板名称以"_"开头，则忽略
		if info != nil && !info.IsDir() && s.testFilePath(path, info.Name(), opt.ExcludePatterns) {
			tp, err := s.parseTemplate(path, opt.AttachCopyright)
			if err != nil {
				return errors.New("template:" + info.Name() + "-" + err.Error())
			}
			tplMap[path[sliceSize-1:]] = tp
		}
		return nil
	})
	if err == nil && len(tplMap) == 0 {
		err = fmt.Errorf("no any code template in: %s , exclude pattern:%s",
			opt.TplDir, opt.ExcludePatterns)
	}
	return tplMap, err
}

// 验证文件名路径, 是否可以可以被作为模板生成代码
func (s *sessionImpl) testFilePath(path string, fileName string, excludePatterns []string) bool {
	if fileName[0] == '_' || fileName[0] == '~' {
		return false
	}
	if strings.HasSuffix(strings.ToUpper(path), "README.MD") {
		return false
	}
	if !strings.HasSuffix(path, ".tpl") {
		return false
	}
	if excludePatterns == nil {
		return true
	}
	for _, v := range excludePatterns {
		if v == path {
			return false
		}
		// 前后匹配
		if strings.Contains(v, "*") {
			c := strings.Replace(v, "*", "", -1)
			if v[0] == '*' && strings.HasSuffix(path, c) {
				return false
			}
			if v[len(v)-1] == '*' && strings.HasPrefix(path, c) {
				return false
			}
			continue
		}
		// 模糊匹配
		if strings.Contains(path, v) {
			return false
		}
	}
	return true
}

// 根据代码文件类型适配table
func (s *sessionImpl) adapterTable(table *Table, path string) *Table {
	l := getLangByPath(path)
	// 部分语言永远使用大写开头的命名
	switch l {
	case L_GO, L_CSharp, L_Thrift, L_Protobuf, L_PHP, L_Shell:
		return table
	case L_JAVA: // 需要生成getter和setter,故大写
		return table
	case L_Kotlin:
		return s.copyTable(table, true)
	}
	if ml := s.opt.MajorLang; ml == L_JAVA || ml == L_Kotlin {
		return s.copyTable(table, true)
	}
	return table
}

func (s *sessionImpl) copyTable(table *Table, lowerProp bool) *Table {
	prop := func(s string) string {
		if lowerProp {
			return lowerTitle(s)
		}
		return s
	}
	dst := &Table{
		Ordinal:    table.Ordinal,
		Name:       table.Name,
		Prefix:     table.Prefix,
		Title:      table.Title,
		ShortTitle: shortTitle(table.Name),
		Comment:    table.Comment,
		Engine:     table.Engine,
		Schema:     table.Schema,
		Charset:    table.Charset,
		Raw:        table.Raw,
		Pk:         table.Pk,
		PkProp:     prop(table.PkProp),
		PkType:     table.PkType,
		Columns:    make([]*Column, len(table.Columns)),
	}
	for i, v := range table.Columns {
		dst.Columns[i] = &Column{
			Ordinal: v.Ordinal,
			Name:    v.Name,
			Prop:    prop(v.Prop),
			IsPk:    v.IsPk,
			IsAuto:  v.IsAuto,
			NotNull: v.NotNull,
			DbType:  v.DbType,
			Comment: v.Comment,
			Length:  v.Length,
			Type:    v.Type,
			Render:  v.Render,
		}
	}
	return dst
}

// Clean implements Session.
func (s *sessionImpl) Clean() error {
	return os.RemoveAll(s.opt.OutputDir)
}
