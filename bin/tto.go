package main

import (
	"database/sql"
	"flag"
	"fmt"
	"github.com/ixre/gof/db"
	"github.com/ixre/gof/db/orm"
	"github.com/ixre/gof/shell"
	"github.com/ixre/tto"
	"io/ioutil"
	"log"
	"os"
	"runtime"
	"strconv"
	"strings"
	"time"
)

func main() {
	cmd := "generate"
	if len(os.Args) > 1 {
		if a := os.Args[1]; !strings.HasPrefix(a, "-") {
			cmd = a
		}
	}
	switch cmd {
	case "update":
		forceUpdate := len(os.Args) > 2 && os.Args[2] == "-y"
		_, _ = tto.DoUpdate(forceUpdate)
	case "generate":
		generate()
	}
}

func checkEveryDay() bool {
	timeFile := os.TempDir() + "/tto_check_time"
	var unix int64
	lastTime, err := ioutil.ReadFile(timeFile)
	if err == nil {
		unix1, _ := strconv.Atoi(string(lastTime))
		unix = int64(unix1)
	}
	if dt := time.Now().Unix(); dt-unix > 3600*24 {
		_ = ioutil.WriteFile(timeFile, []byte(strconv.Itoa(int(dt))), os.ModePerm)
		b, _ := tto.DoUpdate(false)
		return b
	}
	return false
}

func generate() {
	var genDir string    //输出目录
	var confPath string  //设置目录
	var tplDir string    //模板目录
	var majorLang string //主要语言
	var table string
	var excludedTables string
	var arch string //代码架构
	var debug bool
	var printVer bool
	var cleanLast bool
	var compactMode bool
	//var keepLocal bool

	flag.StringVar(&genDir, "o", "./output", "path of output directory")
	flag.StringVar(&tplDir, "t", "./templates", "path of code templates directory")
	flag.StringVar(&majorLang, "m", "go", "major code lang like java or go")
	flag.StringVar(&confPath, "conf", "./tto.conf", "config path")
	flag.StringVar(&table, "table", "", "table name or table prefix")
	flag.StringVar(&excludedTables, "excludes", "", "exclude tables by prefix")
	flag.StringVar(&arch, "arch", "", "program language")
	flag.BoolVar(&cleanLast, "clean", false, "clean last generate files")
	flag.BoolVar(&debug, "debug", false, "debug mode")
	flag.BoolVar(&compactMode, "compact", false, "compact mode for old project")
	//flag.BoolVar(&keepLocal, "local", false, "don't update any new version")
	flag.BoolVar(&printVer, "v", false, "print version")
	flag.Parse()

	if printVer {
		println("tto Generator v" + tto.BuildVersion)
		return
	}
	//if !keepLocal && checkEveryDay() {
	//	os.Exit(0)
	//}
	re, err := tto.LoadRegistry(confPath)
	if err != nil {
		println("[ tto][ fatal]:", err.Error())
		return
	}
	log.SetFlags(log.Ltime | log.Lshortfile)
	defer crashRecover(debug)
	// 兼容模式
	if compactMode {
		tto.CompactMode = true
	}
	// 获取包名
	pkgName := "com/tto/pkg"
	if re.Contains("code.pkg") {
		pkgName = re.GetString("code.pkg")
	}
	// 获取bash启动脚本，默认unix系统包含了bash，windows下需指定
	bashExec := ""
	if runtime.GOOS == "windows" {
		if re.Contains("command.bash_path") {
			bashExec = re.GetString("command.bash_path")
		} else {
			println("[ tto][ warning]: guest os need config bash path")
		}
	}
	// 清理之前生成的结果
	if cleanLast {
		if err = os.RemoveAll(genDir); err != nil {
			log.Fatalln("[ tto][ fatal]:", err.Error())
		}
	}
	// 生成之前执行操作
	if err := runBefore(re, bashExec); err != nil {
		log.Fatalln("[ tto][ fatal]:", err.Error())
	}
	// 初始化生成器
	driver := re.GetString("database.driver")
	dbName := re.GetString("database.name")
	schema := re.GetString("database.schema")
	dialect, dbDriver := tto.GetDialect(driver)
	ds := orm.DialectSession(getDb(driver, re), dialect)
	list, err := ds.TablesByPrefix(dbName, schema, table)
	if err != nil {
		println("[ app][ info]: ", err.Error())
		return
	}
	list = filterTables(list, excludedTables)
	if len(list) == 0 {
		println("[ app][ info]: no any tables")
		return
	}
	// 获取排除的文件名
	excludePatterns := []string{}
	excludePatternParam := re.GetString("code.exclude_patterns")
	if len(excludePatternParam) == 0 {
		excludePatternParam = re.GetString("code.exclude_files")
	}
	if strings.Contains(excludePatternParam, ",") {
		excludePatterns = strings.Split(excludePatternParam, ",")
	} else {
		excludePatterns = strings.Split(excludePatternParam, ";")
	}
	disableAttachCopy := re.GetBoolean("code.disable_attach")
	// 生成自定义代码
	opt := &tto.Options{
		TplDir:          tplDir,
		AttachCopyright: !disableAttachCopy,
		OutputDir:       genDir,
		ExcludePatterns: excludePatterns,
		MajorLang:       majorLang,
	}
	dg := tto.DBCodeGenerator(dbDriver, opt)
	dg.Package(pkgName)
	if re.GetBoolean("code.id_upper") {
		dg.UseUpperId()
	}
	// 实体后缀
	if suffix := re.GetString("code.entity_suffix"); suffix != "" {
		dg.Var(tto.ENTITY_SUFFIX, suffix)
	}
	// 获取表格并转换
	userMeta := re.GetBoolean("code.meta_settings")
	tables, err := dg.Parses(list, userMeta)
	if err != nil {
		println("[ tto][ error]:", err.Error())
		return
	}

	// 生成代码
	if err := genByArch(arch, dg, tables, opt); err != nil {
		log.Fatalln("[ tto][ fatal]:", err.Error())
	}
	// 生成之后执行操作
	if err := runAfter(re, bashExec); err != nil {
		log.Fatalln("[ tto][ fatal]:", err.Error())
	}
	println(fmt.Sprintf("generate successfully! all %d tasks.",
		len(tables)))
}

func filterTables(tables []*orm.Table, noTable string) []*orm.Table {
	if noTable == "" {
		return tables
	}
	excludes := strings.Split(noTable, ";")
	arr := make([]*orm.Table, 0)
	for _, v := range tables {
		match := false
		for _, k := range excludes {
			if k != "" && strings.Index(v.Name, k) != -1 {
				match = true
				break
			}
		}
		if !match {
			arr = append(arr, v)
		}
	}
	return arr
}

func runBefore(re *tto.Registry, bashExec string) error {
	beforeRun := strings.TrimSpace(re.GetString("command.before"))
	return execCommand(beforeRun, bashExec)
}

func runAfter(re *tto.Registry, bashExec string) error {
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
func genByArch(arch string, dg tto.Session, tables []*tto.Table, opt *tto.Options) (err error) {
	// 按架构生成GO代码
	switch arch {
	case "repo":
		err = genGoRepoCode(dg, tables, opt.OutputDir)
	}
	if err != nil {
		println(fmt.Sprintf("[ tto][ Error]: generate go code fail! %s", err.Error()))
	}
	return dg.WalkGenerateCodes(tables)
}

// 获取数据库连接
func getDb(driver string, r *tto.Registry) *sql.DB {
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
	conn, err := db.NewConnector(driver, connStr, nil, false)
	if err == nil {
		log.Println("[ tto][ init]: connect to database..")
		err = conn.Ping()
	}
	if err != nil {
		_ = conn.Close()
		//如果异常，则显示并退出
		log.Fatalln("[ tto][ init]:" + conn.Driver() + "-" + err.Error())
	}
	d := conn.Raw()
	d.SetMaxIdleConns(10)
	d.SetMaxIdleConns(5)
	d.SetConnMaxLifetime(time.Second * 10)
	return d
}

// 恢复应用
func crashRecover(debug bool) {
	if !debug {
		r := recover()
		if r != nil {
			if _, f, l, ok := runtime.Caller(3); ok {
				log.Println(fmt.Sprintf("[ tto][ crash]:file:%s line:%d %v ", f, l, r))
			} else {
				log.Println(fmt.Sprintf("[ tto][ crash]: %v", r))
			}
		}
	}
}

// 生成Go代码
func genGoRepoCode(dg tto.Session, tables []*tto.Table,
	genDir string) error {
	if ig, b := dg.(tto.GoSession); b {
		// 生成GoRepo代码
		err := ig.GenerateGoRepoCodes(tables, genDir)
		//格式化代码
		_, _, _ = shell.Run("gofmt -w " + genDir)
		return err
	}
	return nil
}
