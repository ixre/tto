package lang

import "github.com/ixre/gof/db/db"

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

func (t Typescript) SqlMapType(typeId int, len int) string {
	return t.ParseType(typeId)
}

func (t Typescript) PkgName(pkg string) string {
	return t.PkgPath(pkg)
}
func (t Typescript) PkgPath(pkg string) string {
	return PkgStyleLikeJava(pkg)
}
func (t Typescript) ParseType(typeId int) string {
	return TsTypes(typeId)
}

func (t Typescript) DefaultValue(typeId int) string {
	return tsValues(typeId)
}

func tsValues(typeId int) string {
	switch typeId {
	case db.TypeString:
		return "''"
	case db.TypeBoolean:
		return "false"
	case db.TypeInt16, db.TypeInt32, db.TypeInt64:
		return "0"
	case db.TypeFloat32, db.TypeFloat64:
		return "0.0"
	}
	return "<unknown>"
}

var _ Lang = new(Typescript)

func TsTypes(typeId int) string {
	switch typeId {
	case db.TypeBoolean:
		return "boolean"
	case db.TypeInt64:
		return "number"
	case db.TypeFloat32:
		return "number"
	case db.TypeFloat64:
		return "number"
	case db.TypeInt16, db.TypeInt32, db.TypeDecimal:
		return "number"
	case db.TypeString:
		return "string"
	case db.TypeDateTime:
		return "Date"
	case db.TypeBytes:
		return "any"
	}
	return "any"
}
