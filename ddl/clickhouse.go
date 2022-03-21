package ddl

import (
	"bytes"
	"fmt"
	"github.com/ixre/gof/db/orm"
	"github.com/ixre/tto"
	"strings"
)

// GenerateClickHouseMergeTreeTableDDL 使用MergeTree引擎生成Clickhouse表的DDL语法
func GenerateClickHouseMergeTreeTableDDL(table *tto.Table) string {
	/*
		创建本地表

		建表语句基本语法如下：


		CREATE TABLE [IF NOT EXISTS] [db.]table_name ON CLUSTER cluster
		(
		    name1 [type1] [DEFAULT|MATERIALIZED|ALIAS expr1],
		    name2 [type2] [DEFAULT|MATERIALIZED|ALIAS expr2],
		    ...
		    INDEX index_name1 expr1 TYPE type1(...) GRANULARITY value1,
		    INDEX index_name2 expr2 TYPE type2(...) GRANULARITY value2
		) ENGINE = engine_name()
		[PARTITION BY expr]
		[ORDER BY expr]
		[PRIMARY KEY expr]
		[SAMPLE BY expr]
		[SETTINGS name=value, ...];

	*/
	buf := bytes.NewBuffer(nil)
	buf.WriteString("CREATE TABLE IF NOT EXISTS ")
	buf.WriteString(table.Name)
	buf.WriteString(" COMMENT '")
	buf.WriteString(table.Comment)
	buf.WriteString("' \n(\n")
	for i, v := range table.Columns {
		buf.WriteString(strings.Repeat(" ", 4))
		buf.WriteString("`" + v.Name + "` ")
		buf.WriteString(getClickhouseColumnType(v.Type))
		buf.WriteString(" COMMENT ")
		buf.WriteString("'" + v.Comment + "'")
		if i < len(table.Columns)-1 {
			buf.WriteString(",")
		}
		buf.WriteString("\n")
	}
	buf.WriteString(") ENGINE = MergeTree\n")
	//buf.WriteString("PARTITION BY toYYYYMM(FlightDate)")
	buf.WriteString("ORDER BY ")
	buf.WriteString(table.Pk)
	buf.WriteString("\n")
	buf.WriteString("SETTINGS index_granularity= 8192 ;")
	return buf.String()
}

func getClickhouseColumnType(typeId int) string {
	switch typeId {
	case orm.TypeBoolean:
		return "Int8"
	case orm.TypeInt64:
		return "Int64"
	case orm.TypeFloat32:
		return "Float32"
	case orm.TypeFloat64:
		return "Float64"
	case orm.TypeInt16:
		return "Int16"
	case orm.TypeInt32:
		return "Int32"
	case orm.TypeString:
		return "String"
	case orm.TypeDecimal:
		return "Decimal64"
	case orm.TypeDateTime:
		return "Date"
	case orm.TypeBytes:
		return "String"
	}
	return fmt.Sprintf("Unknown type id:%d", typeId)
}
