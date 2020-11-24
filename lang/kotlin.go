package lang

import (
	"fmt"
	"github.com/ixre/gof/db/orm"
)

/**
 * Copyright (C) 2007-2020 56X.NET,All rights reserved.
 *
 * name : kotlin.go
 * author : jarrysix (jarrysix#gmail.com)
 * date : 2020-11-20 11:09
 * description :
 * history :
 */

type KotlinLang struct {
}

func (k KotlinLang) ParsePkg(pkg string) string {
	return PkgStyleLikeJava(pkg)
}

func (k KotlinLang) ParseType(typeId int) string {
	return KotlinTypes(typeId)
}

func (k KotlinLang) DefaultValue(typeId int) string {
	return JavaValues(typeId)
}

var _ Lang = new(KotlinLang)

func KotlinTypes(typeId int) string {
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
		return "Int"
	case orm.TypeString:
		return "String"
	case orm.TypeDecimal:
		return "BigDecimal"
	case orm.TypeDateTime:
		return "Date"
	case orm.TypeBytes:
		return "ByteArray"
	}
	return fmt.Sprintf("Unknown type id:%d", typeId)
}
