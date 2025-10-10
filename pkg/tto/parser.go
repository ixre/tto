package tto

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"regexp"
	"strings"

	"github.com/ixre/gof/db/db"
	"github.com/ixre/gof/log"
	"github.com/ixre/gof/util"
)

// 获取表结构,userMeta 是否使用用户的元数据
// tablePrefix 表前缀
func parseTable(ordinal int, tb *db.Table, tablePrefix string, shortUpper bool, userMeta bool) *Table {
	tableName := tb.Name
	if len(tablePrefix) > 0 {
		// 去掉指定的前缀
		tableName = regexp.MustCompile(`(?i)^`+tablePrefix).ReplaceAllString(tb.Name, "")
	}
	n := &Table{
		Ordinal:    ordinal,
		Name:       tableName,
		Prefix:     prefix(tableName),
		Title:      title(tableName, shortUpper),
		ShortTitle: shortTitle(tableName),
		Comment:    tb.Comment,
		Engine:     tb.Engine,
		Schema:     tb.Schema,
		Charset:    tb.Charset,
		Raw:        tb,
		Pk:         "",
		PkProp:     "",
		PkType:     db.TypeInt32,
		Columns:    make([]*Column, len(tb.Columns)),
	}
	if len(n.Comment) == 0 {
		n.Comment = n.Title
	}
	for i, v := range tb.Columns {
		if v.IsPk && n.Pk != "" {
			n.Pk = v.Name
			n.PkProp = title(v.Name, shortUpper)
			n.PkType = v.Type
		}
		c := &Column{
			Ordinal: i,
			Name:    v.Name,
			Prop:    title(v.Name, shortUpper),
			IsPk:    v.IsPk,
			IsAuto:  v.IsAuto,
			NotNull: v.NotNull,
			DbType:  v.DbType,
			Comment: v.Comment,
			Length:  v.Length,
			Type:    v.Type,
		}
		if len(c.Comment) == 0 {
			c.Comment = c.Prop
		}
		// 兼容JAVA项目int主键
		if CompactMode && c.DbType == "int(11)" {
			c.Type = db.TypeInt32
		}
		n.Columns[i] = c
	}
	// 如果未设置主键，则默认第一个列
	if n.Pk == "" {
		fc := tb.Columns[0]
		n.Pk = fc.Name
		n.PkProp = title(fc.Name, shortUpper)
		n.PkType = fc.Type
	}
	return loadUserMeta(n, userMeta)
}

// 读取并更新元数据
func loadUserMeta(t *Table, userMeta bool) *Table {
	cfg := &TableConfig{}
	dstPath := fmt.Sprintf("./meta-settings/%s.json", t.Name)
	if userMeta {
		f, err := os.Open(dstPath)
		// 如果不存在
		if os.IsNotExist(err) {
			cfg = defaultMetaInfo(t)
			if userMeta {
				flushCfgFile(cfg, dstPath)
			}
		} else {
			bytes, _ := io.ReadAll(f)
			if err = json.Unmarshal(bytes, cfg); err != nil {
				log.Fatalf("[ tto][ fatal]: read user meta file %s failed, reason: %s",
					dstPath, err.Error())
			}
		}
	} else {
		cfg = defaultMetaInfo(t)
	}
	needMerge := false
	// 赋值Render,如果未包含新的字段,则加入
	for _, c := range t.Columns {
		f, b := cfg.Fields[c.Name]
		if !b {
			f = smartField(c)
			cfg.Fields[c.Name] = f
			needMerge = true
		}
		c.Render = f.Render
	}
	// 去掉旧的字段
	if len(t.Columns) != len(cfg.Fields) {
		needMerge = true
		arr := make([]string, 0)
		for k := range cfg.Fields {
			exist := false
			for _, c := range t.Columns {
				if c.Name == k {
					exist = true
					break
				}
			}
			if !exist {
				arr = append(arr, k)
			}
		}
		for _, v := range arr {
			delete(cfg.Fields, v)
		}
	}
	if needMerge {
		flushCfgFile(cfg, dstPath)
	}
	return t
}

func flushCfgFile(cfg *TableConfig, dstFile string) {
	bytes, err := json.MarshalIndent(cfg, "", " ")
	if err == nil {
		err = util.BytesToFile(bytes, dstFile)
	}
	if err != nil {
		log.Fatalf("[ tto][ fatal]: save user meta file %s failed, reason: %s",
			dstFile, err.Error())
	}
}

func defaultMetaInfo(t *Table) *TableConfig {
	cfg := &TableConfig{
		Struct: &TableMeta{},
		Fields: make(map[string]*ColumnMeta, 0),
	}
	for _, c := range t.Columns {
		f := smartField(c)
		f.Title = c.Comment
		cfg.Fields[c.Name] = f
	}
	return cfg
}

// 返回自动生成的字段
func smartField(c *Column) *ColumnMeta {
	ele, options := smartElement(c)
	fd := &ColumnMeta{
		Render: &PropRenderOptions{
			Visible: true,
			Element: ele,
			Options: options,
		},
	}
	return fd
}

func smartElement(c *Column) (string, map[string]string) {
	name := c.Name
	len := c.Length
	if strings.HasPrefix(name, "is_") {
		return "checkbox", map[string]string{"是": "1", "否": "0"}
	}
	if strings.HasSuffix(name, "_time") {
		return "time", map[string]string{}
	}
	if strings.HasSuffix(name, "state") ||
		strings.HasSuffix(name, "status") ||
		strings.HasSuffix(name, "enabled") {
		return "radio", map[string]string{}
	}
	if strings.HasSuffix(name, "_date") {
		return "date", map[string]string{}
	}
	if strings.HasPrefix(name, "upload_") ||
		strings.HasPrefix(name, "file_") ||
		strings.HasSuffix(name, "_image") ||
		strings.HasSuffix(name, "_img") ||
		strings.HasPrefix(name, "attachment") {
		return "upload", map[string]string{}
	}
	if len >= 64 || len == 0 {
		return "textarea", map[string]string{}
	}
	return "input", map[string]string{}
}
