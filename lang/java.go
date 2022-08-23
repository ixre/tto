package lang

/**
find output/java -name "*.java" | xargs sed -i 's/ int / Integer /g' && \
find output/java -name "*.java" | xargs sed -i 's/ long / Long /g' && \
find output/java -name "*.java" | xargs sed -i 's/ float / Float /g'
*/

import (
	"fmt"

	"github.com/ixre/gof/db/db"
)

var _ Lang = new(Java)

type Java struct {
}

func (j Java) SqlMapType(typeId int, len int) string {
	return j.ParseType(typeId)
}

func (j Java) PkgName(pkg string) string {
	return j.PkgPath(pkg)
}
func (j Java) PkgPath(pkg string) string {
	return PkgStyleLikeJava(pkg)
}

func (j Java) ParseType(typeId int) string {
	return javaTypes(typeId)
}

func (j Java) DefaultValue(typeId int) string {
	return javaValues(typeId)
}

func (j Java) ParsePkType(typeId int) string {
	return javaPkTypes(typeId)
}

var _ Lang = new(Java)

func javaTypes(typeId int) string {
	switch typeId {
	case db.TypeBoolean:
		return "boolean"
	case db.TypeInt64:
		return "long"
	case db.TypeFloat32:
		return "float"
	case db.TypeFloat64:
		return "double"
	case db.TypeInt16, db.TypeInt32:
		return "int"
	case db.TypeString:
		return "String"
	case db.TypeDecimal:
		return "BigDecimal"
	case db.TypeDateTime:
		return "Date"
	case db.TypeBytes:
		return "Byte[]"
	}
	return fmt.Sprintf("Unknown type id:%d", typeId)
}

func javaPkTypes(typeId int) string {
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
		return "Integer"
	case db.TypeString:
		return "String"
	case db.TypeDecimal:
		return "BigDecimal"
	case db.TypeDateTime:
		return "Date"
	case db.TypeBytes:
		return "Byte[]"
	}
	return fmt.Sprintf("Unknown type id:%d", typeId)
}

func javaValues(typeId int) string {
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
		return "new BigDecimal(0.0)"
	case db.TypeString:
		return "\"\""
	case db.TypeDateTime:
		return "new Date()"
	case db.TypeBytes:
		return "new Byte[0]"
	}
	return fmt.Sprintf("Unknown type id:%d", typeId)
}
