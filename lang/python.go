package lang

import (
	"github.com/ixre/gof/db/orm"
)

/**
 * Copyright (C) 2007-2020 56X.NET,All rights reserved.
 *
 * name : python.go
 * author : jarrysix (jarrysix#gmail.com)
 * date : 2020-11-20 11:10
 * description :
 * history :
 */

type PythonLang struct {
}

func (p PythonLang) ParsePkg(pkg string) string {
	return PkgStyleLikeGo(pkg)
}

func (p PythonLang) ParseType(typeId int) string {
	return PythonValues(typeId)
}

func (p PythonLang) DefaultValue(typeId int) string {
	return PythonValues(typeId)
}

var _ Lang = new(PythonLang)

func PyTypes(typeId int) string {
	switch typeId {
	case orm.TypeBoolean:
		return "bool"
	case orm.TypeFloat32, orm.TypeFloat64, orm.TypeDecimal:
		return "float"
	case orm.TypeInt16, orm.TypeInt32, orm.TypeInt64:
		return "int"
	case orm.TypeString:
		return "str"
	}
	return "any"
}

func PythonValues(typeId int) string {
	switch typeId {
	case orm.TypeBoolean:
		return "False"
	case orm.TypeInt64, orm.TypeInt16, orm.TypeInt32:
		return "0"
	case orm.TypeFloat32, orm.TypeFloat64, orm.TypeDecimal:
		return "0.0"
	case orm.TypeString:
		return "\"\""
	}
	return "None"
}
