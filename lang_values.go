package tto

import (
	"fmt"
	"github.com/ixre/gof/db/orm"
)

func GoValues(typeId int) string {
	switch typeId {
	case orm.TypeString:
		return "\"\""
	case orm.TypeBoolean:
		return "false"
	case orm.TypeInt16, orm.TypeInt32, orm.TypeInt64:
		return "0"
	case orm.TypeFloat32, orm.TypeFloat64:
		return "0.0"
	}
	return "<unknown>"
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
	case orm.TypeString:
		return "null"
	}
	return fmt.Sprintf("Unknown type id:%d", typeId)
}

func CommonValues(typeId int) string {
	switch typeId {
	case orm.TypeBoolean:
		return "false"
	case orm.TypeInt64, orm.TypeInt16, orm.TypeInt32:
		return "0"
	case orm.TypeFloat32, orm.TypeFloat64:
		return "0.0"
	case orm.TypeString:
		return "\"\""
	}
	return "null"
}
