package lang

/**
 * Copyright (C) 2007-2020 56X.NET,All rights reserved.
 *
 * name : java
 * author : jarrysix (jarrysix#gmail.com)
 * date : 2020-11-20 11:11
 * description :
 * history :
 */

/**
find output/java -name "*.java" | xargs sed -i 's/ int / Integer /g' && \
find output/java -name "*.java" | xargs sed -i 's/ long / Long /g' && \
find output/java -name "*.java" | xargs sed -i 's/ float / Float /g'
*/

import (
	"fmt"
	"github.com/ixre/gof/db/orm"
)

type JavaLang struct {
}

func (j JavaLang) SqlMapType(typeId int, len int) string {
	return j.ParseType(typeId)
}

func (j JavaLang) PkgName(pkg string) string {
	return j.PkgPath(pkg)
}
func (j JavaLang) PkgPath(pkg string) string {
	return PkgStyleLikeJava(pkg)
}

func (j JavaLang) ParseType(typeId int) string {
	return javaTypes(typeId)
}

func (j JavaLang) DefaultValue(typeId int) string {
	return javaValues(typeId)
}

func (j JavaLang) ParsePkType(typeId int) string {
	return javaPkTypes(typeId)
}

var _ Lang = new(JavaLang)

func javaTypes(typeId int) string {
	switch typeId {
	case orm.TypeBoolean:
		return "boolean"
	case orm.TypeInt64:
		return "long"
	case orm.TypeFloat32:
		return "float"
	case orm.TypeFloat64:
		return "double"
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

func javaPkTypes(typeId int) string {
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
		return "Integer"
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

func javaValues(typeId int) string {
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
		return "new BigDecimal(0.0)"
	case orm.TypeString:
		return "\"\""
	case orm.TypeDateTime:
		return "new Date()"
	case orm.TypeBytes:
		return "new Byte[0]"
	}
	return fmt.Sprintf("Unknown type id:%d", typeId)
}
