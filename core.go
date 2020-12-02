package tto

import (
	"github.com/ixre/gof/db/orm"
	"github.com/ixre/tto/config"
	"strings"
)

// 版本号
const BuildVersion = "0.3.23"

// 代码页
const ReleaseCodeHome = "https://github.com/ixre/tto"

// 兼容模式
var CompactMode = false

// 表
type Table struct {
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
type Column struct {
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

const (
	L_Unknown    = ""
	L_GO         = "go"
	L_JAVA       = "java"
	L_CSharp     = "csharp"
	L_TypeScript = "typescript"
	L_Kotlin     = "kotlin"
	L_Python     = "python"
	L_Thrift     = "thrift"
	L_Protobuf   = "protobuf"
	L_PHP        = "php"
	L_Rust       = "rust"
	L_Dart       = "dart"
	L_Shell      = "shell"
)

var codeFileExtensions = map[string]string{
	"go":         ".go",
	"java":       ".java",
	"csharp":     ".cs",
	"typescript": ".ts",
	"kotlin":     ".kt",
	"python":     ".py",
	"thrift":     ".thrift",
	"protobuf":   ".proto",
	"rust":       "*.rs",
	"dart":       "*.dart",
}

// 根据文件路径判断语言类型
func GetLangByPath(path string) string {
	i := strings.LastIndex(path, ".")
	if i != -1 {
		switch path[i:] {
		case ".go":
			return L_GO
		case ".cs":
			return L_CSharp
		case ".ts":
			return L_TypeScript
		case ".kt":
			return L_Kotlin
		case ".py":
			return L_Python
		case ".thrift":
			return L_Thrift
		case ".proto":
			return L_Protobuf
		case "*.php":
			return L_PHP
		case ".rs":
			return L_Rust
		case ".dart":
			return L_Dart
		case ".sh":
			return L_Shell
		}
	}
	return L_Unknown
}
