package ddl

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/ixre/gof/db/db"
	"github.com/ixre/tto"
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
	case db.TypeBoolean:
		return "Int8"
	case db.TypeInt64:
		return "Int64"
	case db.TypeFloat32:
		return "Float32"
	case db.TypeFloat64:
		return "Float64"
	case db.TypeInt16:
		return "Int16"
	case db.TypeInt32:
		return "Int32"
	case db.TypeString:
		return "String"
	case db.TypeDecimal:
		return "Decimal64"
	case db.TypeDateTime:
		return "Date"
	case db.TypeBytes:
		return "String"
	}
	return fmt.Sprintf("Unknown type id:%d", typeId)
}
