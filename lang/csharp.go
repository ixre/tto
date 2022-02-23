package lang

/**
 * Copyright (C) 2007-2020 56X.NET,All rights reserved.
 *
 * name : java
 * author : jarrysix (jarrysix#gmail.com)
 * date : 2022-02-23 20:36
 * description :
 * history :
 */



import (
	"fmt"
	"github.com/ixre/gof/db/orm"
)

type CSharpLang struct {
}

func (j CSharpLang) SqlMapType(typeId int, len int) string {
	return j.ParseType(typeId)
}

func (j CSharpLang) PkgName(pkg string) string {
	return j.PkgPath(pkg)
}
func (j CSharpLang) PkgPath(pkg string) string {
	return PkgStyleLikeJava(pkg)
}

func (j CSharpLang) ParseType(typeId int) string {
	switch typeId {
	case orm.TypeBoolean:
		return "bool"
	case orm.TypeInt64:
		return "long"
	case orm.TypeFloat32:
		return "float"
	case orm.TypeFloat64:
		return "double"
	case orm.TypeInt16, orm.TypeInt32:
		return "int"
	case orm.TypeString:
		return "string"
	case orm.TypeDecimal:
		return "Decimal"
	case orm.TypeDateTime:
		return "DateTime"
	case orm.TypeBytes:
		return "byte[]"
	}
	return fmt.Sprintf("Unknown type id:%d", typeId)
}

func (j CSharpLang) DefaultValue(typeId int) string {
	switch typeId {
	case orm.TypeBoolean:
		return "false"
	case orm.TypeInt64:
		return "0L"
	case orm.TypeFloat32, orm.TypeFloat64:
		return "0F"
	case orm.TypeInt16, orm.TypeInt32:
		return "0"
	case orm.TypeDecimal:
		return "new Decimal(0.0)"
	case orm.TypeString:
		return "\"\""
	case orm.TypeDateTime:
		return "DateTime.Now"
	case orm.TypeBytes:
		return "new Byte[0]"
	}
	return fmt.Sprintf("Unknown type id:%d", typeId)
}

func (j CSharpLang) ParsePkType(typeId int) string {
	switch typeId {
	case orm.TypeBoolean:
		return "bool"
	case orm.TypeInt64:
		return "long"
	case orm.TypeFloat32:
		return "float"
	case orm.TypeFloat64:
		return "double"
	case orm.TypeInt16, orm.TypeInt32:
		return "int"
	case orm.TypeString:
		return "string"
	case orm.TypeDecimal:
		return "Decimal"
	case orm.TypeDateTime:
		return "DateTime"
	case orm.TypeBytes:
		return "byte[]"
	}
	return fmt.Sprintf("Unknown type id:%d", typeId)
}

var _ Lang = new(CSharpLang)
