package tto

import "strings"

// 版本号
const BuildVersion = "0.3.23"
// 代码页
const ReleaseCodeHome = "https://github.com/ixre/tto"
// 兼容模式
var CompactMode = false

const (
	L_Unknown = ""
	L_GO = "go"
	L_JAVA = "java"
	L_CSharp = "csharp"
	L_TypeScript = "typescript"
	L_Kotlin = "kotlin"
	L_Python = "python"
	L_Thrift = "thrift"
	L_Protobuf = "protobuf"
	L_PHP = "php"
	L_Rust = "rust"
	L_Dart = "dart"
	L_Shell = "shell"
)

var codeFileExtensions = map[string]string{
	"go":".go",
	"java":".java",
	"csharp":".cs",
	"typescript":".ts",
	"kotlin":".kt",
	"python":".py",
	"thrift":".thrift",
	"protobuf":".proto",
	"rust":"*.rs",
	"dart":"*.dart",
}

// 根据文件路径判断语言类型
func GetLangByPath(path string)string{
	i := strings.LastIndex(path,".")
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