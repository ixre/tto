package lang

import (
	"fmt"

	"github.com/ixre/gof/db/db"
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
	case db.TypeString:
		if len > 0 {
			if len > 2048 {
				return "Text"
			}
			return fmt.Sprintf("String(%d)", len)
		}
		return "String"
	case db.TypeBoolean:
		return "Boolean"
	case db.TypeInt16:
		return "SmallInteger"
	case db.TypeInt32:
		return "Integer"
	case db.TypeInt64:
		return "BigInteger"
	case db.TypeFloat32:
		return "Float"
	case db.TypeFloat64:
		return "Float"
	case db.TypeDecimal:
		return "Decimal"
	}
	return "String"
}

func (p PythonLang) PkgName(pkg string) string {
	return PkgStyleLikeGo(pkg)
}

func (p PythonLang) PkgPath(pkg string) string {
	return pkg
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
	case db.TypeBoolean:
		return "bool"
	case db.TypeFloat32, db.TypeFloat64, db.TypeDecimal:
		return "float"
	case db.TypeInt16, db.TypeInt32, db.TypeInt64:
		return "int"
	case db.TypeString:
		return "str"
	case db.TypeDateTime:
		return "time"
	case db.TypeBytes:
		return "byte[]"
	}
	return "any"
}

func PythonValues(typeId int) string {
	switch typeId {
	case db.TypeBoolean:
		return "False"
	case db.TypeInt64, db.TypeInt16, db.TypeInt32:
		return "0"
	case db.TypeFloat32, db.TypeFloat64, db.TypeDecimal:
		return "0.0"
	case db.TypeString:
		return "\"\""
	}
	return "None"
}
