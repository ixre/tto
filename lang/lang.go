package lang

import (
	"strconv"
	"strings"

	"github.com/ixre/gof/db/db"
)

type Lang interface {
	// ParseType parse to lang type
	ParseType(typeId int) string
	// SqlMapType the type of orm mapping
	SqlMapType(typeId int, len int) string
	// DefaultValue get default value of lang
	DefaultValue(typeId int) string
	// PkgPath parse package path
	PkgPath(pkg string) string
	// PkgName get package name
	PkgName(pkg string) string
}

var langMap = map[string]Lang{
	"go":         &goLang{},
	"java":       &Java{},
	"kotlin":     &kotlin{},
	"python":     &Python{},
	"typescript": &Typescript{},
	"csharp":     &cSharp{},
	"protobuf":   &Protobuf{},
	"thrift":     &Thrift{},
	"dart":       &dart{},
}

func Get(n string) Lang {
	switch n {
	case "ts":
		n = "typescript"
	case "py":
		n = "python"
	case "kt":
		n = "kotlin"
	case "pb", "grpc":
		n = "protobuf"
	case "cs":
		n = "csharp"
	}
	if l, b := langMap[n]; b {
		return l
	}
	return &CommonLang{}
}

var _ Lang = new(CommonLang)

type CommonLang struct {
}

func (c CommonLang) SqlMapType(typeId int, _ int) string {
	return c.ParseType(typeId)
}

func (c CommonLang) PkgName(pkg string) string {
	return PkgStyleLikeGo(pkg)
}

func (c CommonLang) PkgPath(pkg string) string {
	return pkg
}

func (c CommonLang) ParseType(typeId int) string {
	return strconv.Itoa(typeId)
}

func (c CommonLang) DefaultValue(typeId int) string {
	return CommonValues(typeId)
}

func CommonValues(typeId int) string {
	switch typeId {
	case db.TypeBoolean:
		return "false"
	case db.TypeInt64, db.TypeInt16, db.TypeInt32:
		return "0"
	case db.TypeFloat32, db.TypeFloat64, db.TypeDecimal:
		return "0.0"
	case db.TypeString:
		return "\"\""
	}
	return "null"
}

// PkgStyleLikeJava case "java", "kotlin", "csharp", "py", "thrift", "protobuf":
func PkgStyleLikeJava(pkg string) string {
	return strings.Replace(pkg, "/", ".", -1)
}

// PkgStyleLikeGo "go", "rust", "php", "python"
func PkgStyleLikeGo(pkg string) string {
	i := strings.LastIndexAny(pkg, "/.")
	if i != -1 {
		return pkg[i+1:]
	}
	return pkg
}

func FmtPackage(lang string, pkg string) string {
	switch lang {
	case "java", "kotlin", "csharp", "py", "thrift", "protobuf":
		return strings.Replace(pkg, "/", ".", -1)
	case "go", "rust", "php", "python":
		i := strings.LastIndexAny(pkg, "/.")
		if i != -1 {
			return pkg[i+1:]
		}
	}
	return pkg
}
