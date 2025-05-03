package main

import (
	"bytes"
	"database/sql"
	"flag"
	"fmt"
	"log"
	"os"
	"runtime"
	"strings"
	"time"

	"github.com/fsnotify/fsnotify"
	_ "github.com/go-sql-driver/mysql"

	"github.com/ixre/gof/db"
	"github.com/ixre/gof/db/dialect"
	"github.com/ixre/gof/db/orm"
	"github.com/ixre/gof/shell"
	"github.com/ixre/gof/util"
	"github.com/ixre/tto"
)

// go run bin/tto.go --model ./templates --watch
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

// func checkEveryDay() bool {
// 	timeFile := os.TempDir() + "/tto_check_time"
// 	var unix int64
// 	lastTime, err := os.ReadFile(timeFile)
// 	if err == nil {
// 		unix1, _ := strconv.Atoi(string(lastTime))
// 		unix = int64(unix1)
// 	}
// 	if dt := time.Now().Unix(); dt-unix > 3600*24 {
// 		_ = os.WriteFile(timeFile, []byte(strconv.Itoa(int(dt))), os.ModePerm)
// 		b, _ := tto.DoUpdate(false)
// 		return b
// 	}
// 	return false
// }

func generate() {
	var genDir string    //输出目录
	var confPath string  //设置目录
	var tplDir string    //模板目录
	var majorLang string //主要语言
	var table string
	var excludedTables string
	var arch string //代码架构
	var verbose bool
	var printVer bool
	var cleanLast bool
	var compactMode bool
	var modelPath string // 模型文件的路径
	var pkgName string   // 包名
	var watch bool
	//var keepLocal bool

	flag.StringVar(&genDir, "o", "./output", "path of output directory")
	flag.StringVar(&tplDir, "t", "./templates", "path of code templates directory")
	flag.BoolVar(&watch, "watch", false, "watch template directory to generate")
	flag.StringVar(&majorLang, "lang", "java", "major code lang like java or go")
	flag.StringVar(&modelPath, "model", "", "path to model directory")
	flag.StringVar(&pkgName, "pkg", "", "the package like 'net.fze.web',it will override file config")
	flag.StringVar(&confPath, "conf", "./tto.conf", "config path")
	flag.StringVar(&table, "table", "", "table name or table prefix")
	flag.StringVar(&excludedTables, "excludes", "", "exclude tables by prefix")
	flag.StringVar(&arch, "arch", "", "program language")
	flag.BoolVar(&cleanLast, "clean", false, "clean last generate files")
	flag.BoolVar(&verbose, "verbose", false, "verbose output")
	flag.BoolVar(&compactMode, "compact", false, "compact mode for old project")
	//flag.BoolVar(&keepLocal, "local", false, "don't update any new version")
	flag.BoolVar(&printVer, "v", false, "print version")
	flag.Parse()

	//confPath = "/data/git/axq/axq-project-demo/generator/tto.conf"
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
	if verbose {
		buf := bytes.NewBuffer(nil)
		buf.WriteString(fmt.Sprintf("package: %s \n", pkgName))
		buf.WriteString(fmt.Sprintf("table : %s* \n", table))
		buf.WriteString(fmt.Sprintf("main program language : %s \n", majorLang))
		fmt.Println(buf.String())
	}

	log.SetFlags(log.Ltime | log.Lmsgprefix)
	log.SetPrefix("[ tto][ info]: ")
	defer crashRecover(verbose)
	// 兼容模式
	if compactMode {
		tto.CompactMode = true
	}
	// 获取包名

	if len(strings.TrimSpace(pkgName)) == 0 {
		if re.Contains("code.pkg") {
			pkgName = re.GetString("code.pkg")
		} else {
			pkgName = "com.tto.pkg"
		}
	}
	orgName := "FZE.NET"
	if re.Contains("global.organization") {
		orgName = re.GetString("global.organization")
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

	// 获取排除的文件名
	var excludePatterns []string
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
	// 初始化生成器
	driver := re.GetString("database.driver")
	dbDriver, dialect := dialect.GetDialect(driver)

	dg := tto.DBCodeGenerator(dbDriver, opt)
	dg.Package(pkgName)
	dg.Var(tto.ORGANIZATION, orgName)
	if re.GetBoolean("code.id_upper") {
		dg.UseUpperId()
	}
	// 实体后缀
	if suffix := re.GetString("code.entity_suffix"); suffix != "" {
		dg.Var(tto.ENTITY_SUFFIX, suffix)
	}
	// 获取表格并转换
	var tables []*tto.Table
	if len(modelPath) == 0 {
		dbName := re.GetString("database.name")
		schema := re.GetString("database.schema")
		ds := orm.DialectSession(getDb(driver, re, verbose), dialect)
				list, err1 := ds.TablesByPrefix(dbName, schema, table)
		if err1 != nil {
			println("[ app][ info]: find table failed ", err1.Error())
		}
		userMeta := re.GetBoolean("code.meta_settings")
		tables, err = dg.Parses(list, userMeta)
	} else {
		tables, err = tto.ReadModels(modelPath)
	}
	if err != nil {
		println("[ tto][ error]:", err.Error())
		return
	}
	// 筛选表
	tables = filterTables(tables, excludedTables)

	h := func() {
		startGenerate(dg, tables, re, bashExec, arch, opt)
	}
	h()
	if watch {
		watchGenerate(opt.TplDir, h)
	}
}

// 监听模板文件变化,并进行生成
func watchGenerate(directory string, h func()) {
	log.Printf("watch directory: %s\n", directory)
	util.FsWatch(func(fsnotify.Event) {
		h()
	}, directory)
}

func startGenerate(dg tto.Session, tables []*tto.Table, re *tto.Registry, bashExec string, arch string, opt *tto.Options) {
	// 生成之前执行操作
	if err := runBefore(re, bashExec); err != nil {
		log.Fatalln(err.Error())
	}
	// 生成代码
	if err := genByArch(arch, dg, tables, opt); err != nil {
		log.Fatalln(err.Error())
	}
	// 生成之后执行操作
	if err := runAfter(re, bashExec); err != nil {
		log.Fatalln(err.Error())
	}
	log.Printf("total %d tasks generate successfully! \r", len(tables))
}

func filterTables(tables []*tto.Table, noTable string) []*tto.Table {
	if noTable == "" {
		return tables
	}
	excludes := strings.Split(strings.ToLower(noTable), ",")
	arr := make([]*tto.Table, 0)
	for _, v := range tables {
		match := false
		for _, k := range excludes {
			if k != "" && strings.Contains(strings.ToLower(v.Name), k) {
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
		_, _, err := shell.Run(command, true)
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
	return dg.WalkGenerateCodes(tables, nil)
}

// 获取数据库连接
func getDb(driver string, r *tto.Registry, debug bool) *sql.DB {
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
	case "postgres", "postgresql", "pgsql":
		connStr = fmt.Sprintf("postgresql://%s:%s@%s:%d/%s?sslmode=disable",
			r.GetString(prefix+".user"),
			r.GetString(prefix+".pwd"),
			r.GetString(prefix+".server"),
			r.Get(prefix+".port").(int64),
			r.GetString(prefix+".name"))
	case "mssql", "sqlserver":
		connStr = fmt.Sprintf("sqlserver://%s:%s@%s:%d?database=%s&encrypt=disable",
			r.GetString(prefix+".user"),
			r.GetString(prefix+".pwd"),
			r.GetString(prefix+".server"),
			r.Get(prefix+".port").(int64),
			r.GetString(prefix+".name"))
	default:
		panic("not support driver :" + driver)
	}
	if debug {
		fmt.Println("driver:", driver)
		fmt.Println("connection string:", connStr)
	}
	conn, err := db.NewConnector(driver, connStr, nil, false)
	if err == nil {
		log.Println("[ tto][ init]: connect to database..")
		err = conn.Ping()
	}
	if err != nil {
		//如果异常，则显示并退出
		log.Fatalln("[ tto][ init]:" + driver + "-" + err.Error())
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
			if _, f, l, ok := runtime.Caller(4); ok {
				log.Printf("[ tto][ crash]:file:%s line:%d %v \n", f, l, r)
			} else {
				log.Printf("[ tto][ crash]: %v \n", r)
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
		_, _, _ = shell.Run("go fmt "+genDir, false)
		return err
	}
	return nil
}
