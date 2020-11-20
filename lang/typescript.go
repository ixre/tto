package lang

import (
	"github.com/ixre/gof/db/orm"
)

/**
 * Copyright (C) 2007-2020 56X.NET,All rights reserved.
 *
 * name : typescript
 * author : jarrysix (jarrysix#gmail.com)
 * date : 2020-11-20 11:10
 * description :
 * history :
 */

type Typescript struct {
}

func (t Typescript) ParsePkg(pkg string) string {
	return PkgStyleLikeJava(pkg)
}

func (t Typescript) ParseType(typeId int) string {
	return TsTypes(typeId)
}

func (t Typescript) DefaultValue(typeId int) string {
	return GoValues(typeId)
}

var _ Lang = new(Typescript)

func TsTypes(typeId int) string {
	switch typeId {
	case orm.TypeBoolean:
		return "boolean"
	case orm.TypeInt64:
		return "number"
	case orm.TypeFloat32:
		return "number"
	case orm.TypeFloat64:
		return "number"
	case orm.TypeInt16, orm.TypeInt32, orm.TypeDecimal:
		return "number"
	case orm.TypeString:
		return "string"
	}
	return "any"
}
