package lang

import (
	"github.com/ixre/gof/db/orm"
	"regexp"
)

/**
 * Copyright (C) 2007-2020 56X.NET,All rights reserved.
 *
 * name : go
 * author : jarrysix (jarrysix#gmail.com)
 * date : 2020-11-20 11:11
 * description :
 * history :
 */

var _ Lang = new(GoLang)

var pkgRegex = regexp.MustCompile("/(com|net|io|cn|org|info)/")

type GoLang struct {
}

func (g GoLang) SqlMapType(typeId int, len int) string {
	return g.ParseType(typeId)
}

func (g GoLang) PkgPath(pkg string) string {
	return pkgRegex.ReplaceAllString(pkg, ".$1/")
}

func (g GoLang) PkgName(pkg string) string {
	return PkgStyleLikeGo(pkg)
}

func (g GoLang) ParseType(typeId int) string {
	return GoTypes(typeId)
}

func (g GoLang) DefaultValue(typeId int) string {
	return GoValues(typeId)
}

func GoTypes(typeId int) string {
	switch typeId {
	case orm.TypeString:
		return "string"
	case orm.TypeBoolean:
		return "bool"
	case orm.TypeInt16:
		return "int16"
	case orm.TypeInt32:
		return "int"
	case orm.TypeInt64:
		return "int64"
	case orm.TypeFloat32:
		return "float32"
	case orm.TypeFloat64, orm.TypeDecimal:
		return "float64"
	case orm.TypeDateTime:
		return "time.Time"
	case orm.TypeBytes:
		return "[]byte"
	}
	return "interface{}"
}

func GoValues(typeId int) string {
	switch typeId {
	case orm.TypeString:
		return "\"\""
	case orm.TypeBoolean:
		return "false"
	case orm.TypeInt16, orm.TypeInt32, orm.TypeInt64:
		return "0"
	case orm.TypeFloat32, orm.TypeFloat64:
		return "0.0"
	}
	return "<unknown>"
}
