package tto

import (
	"fmt"
	"github.com/ixre/gof/db/orm"
	"strconv"
)

func GoTypes(typeId int) string {
	switch typeId {
	case orm.TypeString:
		return "string"
	case orm.TypeBoolean:
		return "bool"
	case orm.TypeInt16:
		return "int16"
	case orm.TypeInt32:
		return "int"
	case orm.TypeInt64:
		return "int64"
	case orm.TypeFloat32:
		return "float32"
	case orm.TypeFloat64:
		return "float64"
	}
	return "interface{}"
}

func ThriftTypes(typeId int) string {
	switch typeId {
	case orm.TypeString:
		return "string"
	case orm.TypeBoolean:
		return "bool"
	case orm.TypeInt16:
		return "i16"
	case orm.TypeInt32:
		return "i32"
	case orm.TypeInt64:
		return "i64"
	case orm.TypeFloat32:
		return "f32"
	case orm.TypeFloat64:
		return "f64"
	}
	return strconv.Itoa(typeId)
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
	case orm.TypeInt16, orm.TypeInt32:
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
	case orm.TypeInt16, orm.TypeInt32:
		return "Int"
	case orm.TypeString:
		return "String"
	}
	return fmt.Sprintf("Unknown type id:%d", typeId)
}

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
	case orm.TypeInt16, orm.TypeInt32:
		return "number"
	case orm.TypeString:
		return "string"
	}
	return "any"
}

func PyTypes(typeId int) string {
	switch typeId {
	case orm.TypeBoolean:
		return "bool"
	case orm.TypeFloat32,orm.TypeFloat64:
		return "float"
	case orm.TypeInt16, orm.TypeInt32, orm.TypeInt64:
		return "int"
	case orm.TypeString:
		return "str"
	}
	return "any"
}
