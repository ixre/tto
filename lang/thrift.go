package lang

import (
	"strconv"

	"github.com/ixre/gof/db/db"
)

type Thrift struct {
}

func (t Thrift) SqlMapType(typeId int, len int) string {
	return t.ParseType(typeId)
}
func (t Thrift) PkgName(pkg string) string {
	return t.PkgPath(pkg)
}
func (t Thrift) PkgPath(pkg string) string {
	return PkgStyleLikeJava(pkg)
}

func (t Thrift) ParseType(typeId int) string {
	return ThriftTypes(typeId)
}

func (t Thrift) DefaultValue(typeId int) string {
	panic("not support")
}

var _ Lang = new(Thrift)

func ThriftTypes(typeId int) string {
	switch typeId {
	case db.TypeString:
		return "string"
	case db.TypeBoolean:
		return "bool"
	case db.TypeInt16:
		return "i16"
	case db.TypeInt32:
		return "i32"
	case db.TypeInt64:
		return "i64"
	case db.TypeFloat32:
		return "f32"
	case db.TypeFloat64:
		return "f64"
	case db.TypeDecimal:
		return "f64"
	case db.TypeDateTime:
		return "string"
	}
	return strconv.Itoa(typeId)
}
