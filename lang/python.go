package lang

import (
	"fmt"
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

// python SQL ORM 类型与python类型不一样,需单独处理
func (p PythonLang) SqlMapType(typeId int, len int) string {
	switch typeId {
	case orm.TypeString:
		if len > 0 {
			if len > 2048 {
				return "Text"
			}
			return fmt.Sprintf("String(%d)", len)
		}
		return "String"
	case orm.TypeBoolean:
		return "Boolean"
	case orm.TypeInt16:
		return "SmallInteger"
	case orm.TypeInt32:
		return "Integer"
	case orm.TypeInt64:
		return "BigInteger"
	case orm.TypeFloat32:
		return "Float"
	case orm.TypeFloat64:
		return "Float"
	case orm.TypeDecimal:
		return "Decimal"
	}
	return "String"
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
	case orm.TypeDateTime:
		return "time"
	case orm.TypeBytes:
		return "byte[]"
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
