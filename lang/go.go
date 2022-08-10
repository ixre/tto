package lang

import (
	"regexp"

	"github.com/ixre/gof/db/db"
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
	case db.TypeString:
		return "string"
	case db.TypeBoolean:
		return "bool"
	case db.TypeInt16:
		return "int16"
	case db.TypeInt32:
		return "int"
	case db.TypeInt64:
		return "int64"
	case db.TypeFloat32:
		return "float32"
	case db.TypeFloat64, db.TypeDecimal:
		return "float64"
	case db.TypeDateTime:
		return "time.Time"
	case db.TypeBytes:
		return "[]byte"
	}
	return "interface{}"
}

func GoValues(typeId int) string {
	switch typeId {
	case db.TypeString:
		return "\"\""
	case db.TypeBoolean:
		return "false"
	case db.TypeInt16, db.TypeInt32, db.TypeInt64:
		return "0"
	case db.TypeFloat32, db.TypeFloat64:
		return "0.0"
	}
	return "<unknown>"
}
