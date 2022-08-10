package lang

import (
	"fmt"

	"github.com/ixre/gof/db/db"
)

type KotlinLang struct {
}

func (k KotlinLang) SqlMapType(typeId int, len int) string {
	return k.ParseType(typeId)
}

func (k KotlinLang) PkgName(pkg string) string {
	return k.PkgPath(pkg)
}

func (k KotlinLang) PkgPath(pkg string) string {
	return PkgStyleLikeJava(pkg)
}

func (k KotlinLang) ParseType(typeId int) string {
	return KotlinTypes(typeId)
}

func (k KotlinLang) DefaultValue(typeId int) string {
	return KotlinValues(typeId)
}

var _ Lang = new(KotlinLang)

func KotlinTypes(typeId int) string {
	switch typeId {
	case db.TypeBoolean:
		return "Boolean"
	case db.TypeInt64:
		return "Long"
	case db.TypeFloat32:
		return "Float"
	case db.TypeFloat64:
		return "Double"
	case db.TypeInt16, db.TypeInt32:
		return "Int"
	case db.TypeString:
		return "String"
	case db.TypeDecimal:
		return "BigDecimal"
	case db.TypeDateTime:
		return "Date"
	case db.TypeBytes:
		return "ByteArray"
	}
	return fmt.Sprintf("Unknown type id:%d", typeId)
}

func KotlinValues(typeId int) string {
	switch typeId {
	case db.TypeBoolean:
		return "false"
	case db.TypeInt64:
		return "0L"
	case db.TypeFloat32, db.TypeFloat64:
		return "0F"
	case db.TypeInt16, db.TypeInt32:
		return "0"
	case db.TypeDecimal:
		return "BigDecimal(0.0)"
	case db.TypeString:
		return "\"\""
	case db.TypeDateTime:
		return "Date()"
	}
	return fmt.Sprintf("Unknown type id:%d", typeId)
}
