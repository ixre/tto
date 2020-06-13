package tto

import (
	"fmt"
	"github.com/ixre/gof/db/orm"
)

func PySqlTypes(typeId int, len int) string {
	switch typeId {
	case orm.TypeString:
		if len > 0 {
			if len > 2048 {
				return "Text"
			}
			return fmt.Sprintf("String(%d)", len)
		}
		return "String"
	case orm.TypeBoolean:
		return "Boolean"
	case orm.TypeInt16:
		return "SmallInteger"
	case orm.TypeInt32:
		return "Integer"
	case orm.TypeInt64:
		return "BigInteger"
	case orm.TypeFloat32:
		return "Float"
	case orm.TypeFloat64:
		return "Float"
	}
	return "String"
}
