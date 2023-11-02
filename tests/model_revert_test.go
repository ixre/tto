package tests

import (
	"os"
	"testing"

	_ "github.com/go-sql-driver/mysql"
	"github.com/ixre/gof/shell"
	"github.com/ixre/tto"
)

var (
	genDir = "./output/"
	tplDir = "../templates/vue"
)

// 加载路径中的模型并生成代码
func TestGenerateByReadedTables(t *testing.T) {
	//txt, _ := ioutil.ReadFile("./templates/table.tb")
	//tables, _ := ReadTables(string(txt), "user")
	tables, _ := tto.ReadModels(".")
	if len(tables) == 0 {
		t.Log("no such tables")
		t.FailNow()
	}
	// 生成自定义代码
	opt := &tto.Options{
		TplDir:          tplDir,
		AttachCopyright: !true,
		OutputDir:       genDir,
		ExcludePatterns: []string{"grid_list.html"},
	}
	dg := tto.DBCodeGenerator("", opt)

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
	shell.Run("go fmt "+genDir, false)
	t.Log("生成成功, 输出目录", genDir)
}
