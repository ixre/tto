package main

import (
	"database/sql"
	"flag"
	"fmt"
	"github.com/ixre/gof/db"
	"github.com/ixre/gof/db/orm"
	"github.com/ixre/gof/shell"
	"github.com/ixre/tto"
	"github.com/pelletier/go-toml"
	"log"
	"os"
	"runtime"
	"strings"
	"time"
)

type Registry struct {
	tree *toml.Tree
}

func LoadRegistry(path string) (*Registry, error) {
	tree, err := toml.LoadFile(path)
	if err == nil {
		return &Registry{tree: tree}, err
	}
	return nil, err
}
func (r Registry) Contains(key string) bool {
	return r.tree.Has(key)
}
func (r Registry) GetString(key string) string {
	if r.Contains(key) {
		return r.Get(key).(string)
	}
	return ""
}

func (r Registry) Get(key string) interface{} {
	return r.tree.Get(key)
}
func (r Registry) GetBoolean(key string) bool {
	if r.Contains(key) {
		return r.Get(key).(bool)
	}
	return false
}

func main() {
	var genDir string   //输出目录
	var confPath string //设置目录
	var tplDir string   //模板目录
	var table string
	var arch string //代码架构
	var debug bool
	var printVer bool

	flag.StringVar(&genDir, "out", "./output", "path of output directory")
	flag.StringVar(&tplDir, "tpl", "./templates", "path of code templates directory")
	flag.StringVar(&confPath, "conf", "./tto.conf", "config path")
	flag.StringVar(&table, "table", "", "table name or table prefix")
	flag.StringVar(&arch, "arch", "", "program language")
	flag.BoolVar(&debug, "debug", false, "debug mode")
	flag.BoolVar(&printVer, "v", false, "print version")
	flag.Parse()
	if printVer {
		println("tto Generator v" + tto.BuildVersion)
		return
	}
	re, err := LoadRegistry(confPath)
	if err != nil {
		println("[ Gen][ Fail]:", err.Error())
		return
	}
	log.SetFlags(log.Ltime | log.Lshortfile)
	defer crashRecover(debug)
	// 获取包名
	pkgName := "com/pkg"
	if re.Contains("code.pkg") {
		pkgName = re.GetString("code.pkg")
	}
	// 获取bash启动脚本，默认unix系统包含了bash，windows下需指定
	bashExec := ""
	if runtime.GOOS == "windows" {
		if re.Contains("command.bash_path") {
			bashExec = re.GetString("command.bash_path")
		} else {
			println("[ Gen][ Warning]: guest os need config bash path")
		}
	}
	// 清理之前生成的结果
	if err = os.RemoveAll(genDir); err != nil {
		log.Fatalln("[ Gen][ Fail]:", err.Error())
	}
	// 生成之前执行操作
	if err := runBefore(re, bashExec); err != nil {
		log.Fatalln("[ Gen][ Fail]:", err.Error())
	}
	// 初始化生成器
	driver := re.GetString("database.driver")
	dbName := re.GetString("database.name")
	schema := re.GetString("database.schema")
	ds := orm.DialectSession(getDb(driver, re), getDialect(driver))
	dg := tto.DBCodeGenerator()
	dg.Var(tto.PKG, pkgName)
	dg.Var(tto.TIME, time.Now().Format("2006/01/02 15:04:05"))
	if re.GetBoolean("code.id_upper") {
		dg.IdUpper = true
	}
	// 获取表格并转换
	tables, err := dg.ParseTables(ds.TablesByPrefix(dbName, schema, table))
	if err != nil {
		println("[ Gen][ Fail]:", err.Error())
		return
	}
	// 生成代码
	if err := genByArch(arch, dg, tables, genDir, tplDir); err != nil {
		log.Fatalln("[ Gen][ Fail]:", err.Error())
	}
	// 生成之后执行操作
	if err := runAfter(re, bashExec); err != nil {
		log.Fatalln("[ Gen][ Fail]:", err.Error())
	}
	println(fmt.Sprintf("[ Gen][ Success]: generate successfully! all %d tasks.", len(tables)))
}

func runBefore(re *Registry, bashExec string) error {
	beforeRun := strings.TrimSpace(re.GetString("command.before"))
	return execCommand(beforeRun, bashExec)
}

func runAfter(re *Registry, bashExec string) error {
	afterRun := strings.TrimSpace(re.GetString("command.after"))
	return execCommand(afterRun, bashExec)
}

// 执行命令
func execCommand(command string, bashExec string) error {
	// 生成之后执行操作
	if command != "" {
		if bashExec != "" {
			if strings.Contains(bashExec, " ") {
				bashExec = "\"" + bashExec + "\""
			}
			command = bashExec + " " + command
		}
		_, _, err := shell.StdRun(command)
		return err
	}
	return nil
}

// 根据规则生成代码
func genByArch(arch string, dg *tto.Session, tables []*tto.Table,
	genDir string, tplDir string) (err error) {
	// 按架构生成GO代码
	switch arch {
	case "repo":
		err = genGoRepoCode(dg, tables, genDir)
	}
	if err != nil {
		println(fmt.Sprintf("[ Gen][ Error]: generate go code fail! %s", err.Error()))
	}
	// 生成自定义代码
	return dg.WalkGenerateCode(tables, tplDir, genDir)
}

// 获取数据库连接
func getDb(driver string, r *Registry) *sql.DB {
	//数据库连接字符串
	//root@tcp(127.0.0.1:3306)/db_name?charset=utf8
	var prefix = "database"
	dbCharset := r.GetString(prefix + ".charset")
	if dbCharset == "" {
		dbCharset = "utf8"
	}
	var connStr string
	switch driver {
	case "mysql", "mariadb":
		connStr = fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=%s&loc=Local",
			r.GetString(prefix+".user"),
			r.GetString(prefix+".pwd"),
			r.GetString(prefix+".server"),
			r.Get(prefix+".port").(int64),
			r.GetString(prefix+".name"),
			dbCharset)
	case "postgres", "postgresql":
		connStr = fmt.Sprintf("postgresql://%s:%s@%s:%d/%s?sslmode=disable",
			r.GetString(prefix+".user"),
			r.GetString(prefix+".pwd"),
			r.GetString(prefix+".server"),
			r.Get(prefix+".port").(int64),
			r.GetString(prefix+".name"))
	default:
		panic("not support driver :" + driver)
	}
	conn := db.NewConnector(driver, connStr, nil, false)
	d := conn.Raw()
	d.SetMaxIdleConns(10)
	d.SetMaxIdleConns(5)
	d.SetConnMaxLifetime(time.Second * 10)
	return d
}

func getDialect(driver string) orm.Dialect {
	switch driver {
	case "mysql":
		return &orm.MySqlDialect{}
	case "postgres", "postgresql":
		return &orm.PostgresqlDialect{}
	}
	return nil
}

// 恢复应用
func crashRecover(debug bool) {
	if !debug {
		r := recover()
		if r != nil {
			if _, f, l, ok := runtime.Caller(3); ok {
				log.Println(fmt.Sprintf("[ Gen][ Crash]:file:%s line:%d %v ", f, l, r))
			} else {
				log.Println(fmt.Sprintf("[ Gen][ Crash]: %v", r))
			}
		}
	}
}

// 生成Go代码
func genGoRepoCode(dg *tto.Session, tables []*tto.Table,
	genDir string) error {
	// 生成GoRepo代码
	err := dg.GenerateGoRepoCodes(tables, genDir)
	//格式化代码
	shell.Run("gofmt -w " + genDir)
	return err
}

// 生成Go代码
func genGoCodeOld(dg *tto.Session, tables []*tto.Table,
	genDir string, tplDir string) error {
	// 读取自定义模板
	//listTP, _ := dg.ParseTemplate(tplDir + "/grid_list.html")
	//editTP, _ := dg.ParseTemplate(tplDir + "/entity_edit.html")
	//ctrTpl, _ := dg.ParseTemplate(tplDir + "/entity.html")

	// 生成GoRepo代码
	err := dg.GenerateGoRepoCodes(tables, genDir)

	// 初始化表单引擎
	//fe := &form.Engine{}
	//for _, tb := range tables {
	//	entityPath := genDir + "model/" + tb.Name + ".go"
	//	iRepPath := genDir + "repo/auto_iface_" + tb.Name + "_repo.go"
	//	repPath := genDir + "repo/auto_" + tb.Name + "_repo.go"
	//	dslPath := genDir + "form/" + tb.Name + ".form"
	//	htmPath := genDir + "html/" + tb.Name + ".html"
	//	//生成实体
	//	str := dg.TableToGoStruct(tb)
	//	tto.SaveFile(str, entityPath)
	//	//生成仓储结构
	//	str = dg.TableToGoRepo(tb, true, "model.")
	//	tto.SaveFile(str, repPath)
	//	//生成仓储接口
	//	str = dg.TableToGoIRepo(tb, true, "")
	//	tto.SaveFile(str, iRepPath)
	//	//生成表单DSL
	//	f := fe.TableToForm(tb.Raw)
	//	err = fe.SaveDSL(f, dslPath)
	//	//生成表单
	//	if err == nil {
	//		_, err = fe.SaveHtmlForm(f, form.TDefaultFormHtml, htmPath)
	//	}
	//	if err != nil {
	//		return err
	//	}
	//	// 生成列表文件
	//	str = dg.GenerateCode(tb, listTP, "", true, "")
	//	tto.SaveFile(str, genDir+"html_list/"+tb.Name+"_list.html")
	//	// 生成表单文件
	//	str = dg.GenerateCode(tb, editTP, "", true, "")
	//	tto.SaveFile(str, genDir+"html_edit/"+tb.Name+"_edit.html")
	//	// 生成控制器
	//	str = dg.GenerateCode(tb, ctrTpl, "", true, "")
	//	tto.SaveFile(str, genDir+"mvc/"+tb.Name+"_c.go")
	//}
	//// 生成仓储工厂
	//code := dg.GenerateCodeByTables(tables, tto.TPL_REPO_FACTORY)
	//tto.SaveFile(code, genDir+"repo/auto_repo_factory.go")
	//格式化代码
	shell.Run("gofmt -w " + genDir)
	return err
}
