package tto

import (
	"io/ioutil"
	"os"
	"testing"

	_ "github.com/go-sql-driver/mysql"
	"github.com/ixre/gof/db"
	"github.com/ixre/gof/db/orm"
	"github.com/ixre/gof/shell"
)

var (
	driver     = "mysql"
	dbName     = ""
	dbPrefix   = "t_CSTD_GoodsMst_2003"
	connString = "root:@tcp(127.0.0.1:3306)/baozhang?charset=utf8"
	genDir     = "./generated_code/"
	tplDir     = "./templates/java/spring/kotlin"
)

// 生成数据库所有的代码文件
func TestGenAll(t *testing.T) {
	driver = "postgresql"
	connString = "postgres://postgres:123456@127.0.0.1:5432/go2o?sslmode=disable"

	driver = "sqlserver"
	connString = "sqlserver://sfDBUser:Jbmeon@008@192.168.16.9:1433?database=DCF19_ERP_TEST_B&encrypt=disable"

	// 初始化生成器
	conn, _ := db.NewConnector(driver, connString, nil, false)
	ds, _ := orm.NewDialectSession(driver, conn.Raw())

	// 生成自定义代码
	opt := &Options{
		TplDir:          tplDir,
		AttachCopyright: true,
		OutputDir:       genDir,
		ExcludePatterns: []string{"grid_list.html"},
	}
	dg := DBCodeGenerator(ds.Driver(), opt)
	list, err := ds.TablesByPrefix(dbName, "", dbPrefix)
	if err != nil {
		println("[ tto][ error]: not found tables", err.Error())
		return
	}
	// 获取表格并转换
	tables, err := dg.Parses(list, true)
	if err != nil {
		t.Error(err)
		return
	}

	// 设置包名
	dg.Package("github.com/ixre/go2o/core")
	// 清理上次生成的代码
	os.RemoveAll(genDir)
	// 生成GoRepo代码
	//dg.GenerateGoRepoCodes(tables, genDir)
	dg.WalkGenerateCodes(tables)
	//格式化代码
	shell.Run("gofmt -w "+genDir, false)
	t.Log("生成成功, 输出目录", genDir)
}

// func getDialect(driver string) orm.Dialect {
// 	switch driver {
// 	case "mysql":
// 		return &orm.MySqlDialect{}
// 	case "postgres", "postgresql":
// 		return &orm.PostgresqlDialect{}
// 	}
// 	return nil
// }

func TestReadTables(t *testing.T) {
	txt, _ := ioutil.ReadFile("./templates/table.tb")
	tables, _ := ReadTables(string(txt))
	t.Log(len(tables))
}

func TestGenerateByReadedTables(t *testing.T) {
	//txt, _ := ioutil.ReadFile("./templates/table.tb")
	//tables, _ := ReadTables(string(txt), "user")
	tables, _ := ReadModels("./templates")
	// 生成自定义代码
	opt := &Options{
		TplDir:          tplDir,
		AttachCopyright: true,
		OutputDir:       genDir,
		ExcludePatterns: []string{"grid_list.html"},
	}
	dg := DBCodeGenerator("", opt)

	// 设置包名
	dg.Package("github.com/ixre/go2o/core")
	// 清理上次生成的代码
	os.RemoveAll(genDir)
	// 生成GoRepo代码
	//dg.GenerateGoRepoCodes(tables, genDir)
	err := dg.WalkGenerateCodes(tables)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	//格式化代码
	shell.Run("gofmt -w "+genDir, false)
	t.Log("生成成功, 输出目录", genDir)
}
