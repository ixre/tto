package lang

import (
	"github.com/ixre/gof/db/orm"
	"strconv"
	"strings"
)

/**
 * Copyright (C) 2007-2020 56X.NET,All rights reserved.
 *
 * name : lang
 * author : jarrysix (jarrysix#gmail.com)
 * date : 2020-11-20 11:14
 * description :
 * history :
 */

type Lang interface {
	// parse to lang type
	ParseType(typeId int) string
	// the type of orm mapping
	SqlMapType(typeId int, len int) string
	// get default value of lang
	DefaultValue(typeId int) string
	// parse package
	ParsePkg(pkg string) string
}

var langMap = map[string]Lang{
	"go":         &GoLang{},
	"java":       &JavaLang{},
	"kotlin":     &KotlinLang{},
	"python":     &PythonLang{},
	"typescript": &Typescript{},
	"protobuf":   &Protobuf{},
	"thrift":     &Thrift{},
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
	}
	if l, b := langMap[n]; b {
		return l
	}
	return &CommonLang{}
}

var _ Lang = new(CommonLang)

type CommonLang struct {
}

func (c CommonLang) SqlMapType(typeId int, len int) string {
	return c.ParseType(typeId)
}

func (c CommonLang) ParsePkg(pkg string) string {
	return PkgStyleLikeGo(pkg)
}

func (c CommonLang) ParseType(typeId int) string {
	return strconv.Itoa(typeId)
}

func (c CommonLang) DefaultValue(typeId int) string {
	return CommonValues(typeId)
}

func CommonValues(typeId int) string {
	switch typeId {
	case orm.TypeBoolean:
		return "false"
	case orm.TypeInt64, orm.TypeInt16, orm.TypeInt32:
		return "0"
	case orm.TypeFloat32, orm.TypeFloat64, orm.TypeDecimal:
		return "0.0"
	case orm.TypeString:
		return "\"\""
	}
	return "null"
}

// case "java", "kotlin", "csharp", "py", "thrift", "protobuf":
func PkgStyleLikeJava(pkg string) string {
	return strings.Replace(pkg, "/", ".", -1)
}

// "go", "rust", "php", "python"
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
