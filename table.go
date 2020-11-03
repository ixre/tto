package tto

/**
 * Copyright (C) 2007-2020 56X.NET,All rights reserved.
 *
 * name : table.go
 * author : jarrysix (jarrysix#gmail.com)
 * date : 2020-11-03 17:24
 * description :
 * history :
 */

import (
	"github.com/ixre/gof/db/orm"
	"github.com/ixre/tto/config"
)

type (
	// 表
	Table struct {
		// 顺序
		Ordinal int
		// 表名
		Name string
		// 表前缀
		Prefix string
		// 表名单词首字大写
		Title string
		// 表注释
		Comment string
		// 数据库引擎
		Engine string
		// 架构
		Schema string
		// 数据库编码
		Charset string
		// 表
		Raw *orm.Table
		// 主键
		Pk string
		//　主键属性
		PkProp string
		// 主键类型编号
		PkType int
		// 列
		Columns []*Column
	}
	// 列
	Column struct {
		// 顺序
		Ordinal int
		// 列名
		Name string
		// 列名首字大写
		Prop string
		// 是否主键
		IsPk bool
		// 是否自动生成
		IsAuto bool
		// 是否不能为空
		NotNull bool
		// 类型
		DbType string
		// 注释
		Comment string
		// 长度
		Length int
		// Go类型
		Type int
		// 输出选项
		Render *config.PropRenderOptions
	}
)
