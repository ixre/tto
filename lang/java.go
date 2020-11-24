package lang

import (
	"fmt"
	"github.com/ixre/gof/db/orm"
)

type JavaLang struct {
}

func (j JavaLang) ParsePkg(pkg string) string {
	return PkgStyleLikeJava(pkg)
}

func (j JavaLang) ParseType(typeId int) string {
	return JavaTypes(typeId)
}

func (j JavaLang) DefaultValue(typeId int) string {
	return JavaValues(typeId)
}

/**
 * Copyright (C) 2007-2020 56X.NET,All rights reserved.
 *
 * name : java
 * author : jarrysix (jarrysix#gmail.com)
 * date : 2020-11-20 11:11
 * description :
 * history :
 */
var _ Lang = new(JavaLang)

func JavaTypes(typeId int) string {
	switch typeId {
	case orm.TypeBoolean:
		return "Boolean"
	case orm.TypeInt64:
		return "Long"
	case orm.TypeFloat32:
		return "Float"
	case orm.TypeFloat64:
		return "Double"
	case orm.TypeInt16, orm.TypeInt32:
		return "int"
	case orm.TypeString:
		return "String"
	case orm.TypeDecimal:
		return "BigDecimal"
	case orm.TypeDateTime:
		return "Date"
	case orm.TypeBytes:
		return "Byte[]"
	}
	return fmt.Sprintf("Unknown type id:%d", typeId)
}

func JavaValues(typeId int) string {
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
		return "0.0"
	case orm.TypeString:
		return "\"\""
	}
	return fmt.Sprintf("Unknown type id:%d", typeId)
}
