package example

import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/ixre/goex/generator"
	"github.com/ixre/gof/db"
	"github.com/ixre/gof/db/orm"
	"github.com/ixre/gof/shell"
	"github.com/ixre/gof/web/form"
	"testing"
)

var (
	driver     = "mysql"
	dbName     = ""
	dbPrefix   = "wal_"
	connString = "root:@tcp(127.0.0.1:3306)/txmall?charset=utf8"
	genDir     = "generated_code/"
)

// 生成数据库所有的代码文件
func TestGenAll(t *testing.T) {
	driver = "postgresql"
	connString = "postgres://postgres:123456@127.0.0.1:5432/go2o?sslmode=disable"

	// 初始化生成器
	conn := db.NewConnector(driver, connString, nil, false).Raw()
	dialect := getDialect(driver)
	ds := orm.DialectSession(conn, dialect)
	dg := generator.DBCodeGenerator()
	// 获取表格并转换
	tables, err := dg.ParseTables(ds.TablesByPrefix(dbName, "", dbPrefix))
	if err != nil {
		t.Error(err)
		return
	}
	// 设置变量
	modelPkg := "github.com/ixre/goex/generator/example/" + genDir + "model"
	modelPkgName := "model"
	dg.Var(generator.VModelPkgName, modelPkgName)
	dg.Var(generator.VModelPkg, modelPkg)
	dg.Var(generator.VIRepoPkg, modelPkg)
	// 读取自定义模板
	listTP, _ := dg.ParseTemplate("code_templates/grid_list.html")
	editTP, _ := dg.ParseTemplate("code_templates/entity_edit.html")
	ctrTpl, _ := dg.ParseTemplate("code_templates/entity_c._go")
	// 初始化表单引擎
	fe := &form.Engine{}
	for _, tb := range tables {
		entityPath := genDir + modelPkgName + "/" + tb.Name + ".go"
		iRepPath := genDir + "repo/auto_i" + tb.Name + "_repo.go"
		repPath := genDir + "repo/auto_" + tb.Name + "_repo.go"
		dslPath := genDir + "form/" + tb.Name + ".form"
		htmPath := genDir + "html/" + tb.Name + ".html"
		//生成实体
		str := dg.TableToGoStruct(tb)
		generator.SaveFile(str, entityPath)
		//生成仓储结构
		str = dg.TableToGoRepo(tb, true, modelPkgName+".")
		generator.SaveFile(str, repPath)
		//生成仓储接口
		str = dg.TableToGoIRepo(tb, true, modelPkgName+".")
		generator.SaveFile(str, iRepPath)
		//生成表单DSL
		f := fe.TableToForm(tb.Raw)
		err = fe.SaveDSL(f, dslPath)
		//生成表单
		if err == nil {
			_, err = fe.SaveHtmlForm(f, form.TDefaultFormHtml, htmPath)
		}
		// 生成列表文件
		str = dg.GenerateCode(tb, listTP, "", true, "")
		generator.SaveFile(str, genDir+"html_list/"+tb.Name+"_list.html")
		// 生成表单文件
		str = dg.GenerateCode(tb, editTP, "", true, "")
		generator.SaveFile(str, genDir+"html_edit/"+tb.Name+"_edit.html")
		// 生成控制器
		str = dg.GenerateCode(tb, ctrTpl, "", true, "")
		generator.SaveFile(str, genDir+"c/"+tb.Name+"_c.go")
	}
	//格式化代码
	shell.Run("gofmt -w " + genDir)
	t.Log("生成成功")
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
