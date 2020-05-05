package example

import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/ixre/gof/db"
	"github.com/ixre/gof/db/orm"
	"github.com/ixre/gof/shell"
	"github.com/ixre/tto"
	"os"
	"testing"
)

var (
	driver     = "mysql"
	dbName     = ""
	dbPrefix   = "wal_"
	connString = "root:@tcp(127.0.0.1:3306)/go2o?charset=utf8"
	genDir     = "generated_code/"
)

// 生成数据库所有的代码文件
func TestGenAll(t *testing.T) {
	driver = "postgresql"
	connString = "postgres://postgres:123456@127.0.0.1:5432/go2o?sslmode=disable"

	// 初始化生成器
	conn,_ := db.NewConnector(driver, connString, nil, false)
	dialect := getDialect(driver)
	ds := orm.DialectSession(conn.Raw(), dialect)
	dg := tto.DBCodeGenerator()
	list,err := ds.TablesByPrefix(dbName, "", dbPrefix)
	if err != nil{
		println("[ tto][ error]: not found tables", err.Error())
		return
	}
	// 获取表格并转换
	tables, err := dg.Parses(list,true)
	if err != nil {
		t.Error(err)
		return
	}

	// 设置包名
	dg.Var(tto.PKG, "go2o/core")
	// 清理上次生成的代码
	os.RemoveAll(genDir)
	// 生成GoRepo代码
	dg.GenerateGoRepoCodes(tables, genDir)
	// 生成自定义代码
	dg.WalkGenerateCode(tables, "./templates", genDir)
	//格式化代码
	shell.Run("gofmt -w " + genDir)
	t.Log("生成成功, 输出目录", genDir)
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
