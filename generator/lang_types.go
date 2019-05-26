package generator

import (
	"fmt"
	"github.com/ixre/gof/db/orm"
)


func  GoTypes(typeId int) string {
	switch typeId {
	case orm.TypeString:
		return "string"
	case orm.TypeBoolean:
		return "bool"
	case orm.TypeInt16:
		return "int16"
	case orm.TypeInt32:
		return "int32"
	case orm.TypeInt64:
		return "int64"
	case orm.TypeFloat32:
		return "float32"
	case orm.TypeFloat64:
		return "float64"
	}
	return "interface{}"
}


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
	case orm.TypeInt16,orm.TypeInt32:
		return "int"
	case orm.TypeString:
		return "String"
	}
	return fmt.Sprintf("Unknown type id:%d", typeId)
}


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
	case orm.TypeInt16,orm.TypeInt32:
		return "Int"
	case orm.TypeString:
		return "String"
	}
	return fmt.Sprintf("Unknown type id:%d", typeId)
}
