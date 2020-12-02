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

func (j JavaLang) ParsePkg(pkg string) string {
	return PkgStyleLikeJava(pkg)
}

func (j JavaLang) ParseType(typeId int) string {
	return JavaTypes(typeId)
}

func (j JavaLang) DefaultValue(typeId int) string {
	return JavaValues(typeId)
}


var _ Lang = new(JavaLang)

func JavaTypes(typeId int) string {
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
