package lang

import (
	"fmt"

	"github.com/ixre/gof/db/db"
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
	case db.TypeBoolean:
		return "bool"
	case db.TypeInt64:
		return "long"
	case db.TypeFloat32:
		return "float"
	case db.TypeFloat64:
		return "double"
	case db.TypeInt16, db.TypeInt32:
		return "int"
	case db.TypeString:
		return "string"
	case db.TypeDecimal:
		return "Decimal"
	case db.TypeDateTime:
		return "DateTime"
	case db.TypeBytes:
		return "byte[]"
	}
	return fmt.Sprintf("Unknown type id:%d", typeId)
}

func (j CSharpLang) DefaultValue(typeId int) string {
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
		return "new Decimal(0.0)"
	case db.TypeString:
		return "\"\""
	case db.TypeDateTime:
		return "DateTime.Now"
	case db.TypeBytes:
		return "new Byte[0]"
	}
	return fmt.Sprintf("Unknown type id:%d", typeId)
}

func (j CSharpLang) ParsePkType(typeId int) string {
	switch typeId {
	case db.TypeBoolean:
		return "bool"
	case db.TypeInt64:
		return "long"
	case db.TypeFloat32:
		return "float"
	case db.TypeFloat64:
		return "double"
	case db.TypeInt16, db.TypeInt32:
		return "int"
	case db.TypeString:
		return "string"
	case db.TypeDecimal:
		return "Decimal"
	case db.TypeDateTime:
		return "DateTime"
	case db.TypeBytes:
		return "byte[]"
	}
	return fmt.Sprintf("Unknown type id:%d", typeId)
}

var _ Lang = new(CSharpLang)
