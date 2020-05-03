package tto

import (
	"github.com/ixre/gof/util"
	ht "html/template"
	"strconv"
	"strings"
	"unicode"
)

type internalFunc struct {
}

// 返回模板函数
func (t *internalFunc) funcMap() ht.FuncMap {
	fm := make(map[string]interface{})
	fm["boolInt"] = t.boolInt
	fm["isEmpty"] = t.isEmpty
	fm["raw"] = t.rawHtml
	fm["add"] = t.plus
	fm["plus"] = t.plus
	fm["multi"] = t.multi
	fm["mathRemain"] = t.mathRemain
	// 单词首字大写
	fm["title"] = t.title
	// 小写
	fm["lower"] = t.lower
	// 大写
	fm["upper"] = t.upper
	// 首字母小写: 如:{{lower_title .table.Name}}
	fm["lower_title"] = t.lowerTitle
	// 类型: 如:{{type "go" .columns[0].DbType}}
	fm["type"] = t.langType
	// 包名: {{pkg "go" "com/tto/pkg"}}
	fm["pkg"] = t.langPkg
	// 默认值, 如:{{default "go" .columns[0].DbType}}
	fm["default"] = t.langDefaultValue
	// 是否相等，如：{{equal "go" "rust"}
	fm["equal"] = t.equal
	// 包含函数, 如:{{contain .table.Pk "id"}}
	fm["contain"] = t.contain
	// 是否以指定字符开始, 如:{{starts_with .table.Pk "id"}}
	fm["starts_with"] = t.startsWith
	// 是否以指定字符结束, 如:{{ends_with .table.Pk "id"}}
	fm["ends_with"] = t.endsWith
	return fm
}

// 小写
func (t *internalFunc) lower(s string) string {
	return strings.ToLower(s)
}

// 大写
func (t *internalFunc) upper(s string) string {
	return strings.ToUpper(s)
}

// 将首字母小写
func (t *internalFunc) lowerTitle(s string) string {
	if rune0 := rune(s[0]); unicode.IsUpper(rune0) {
		return string(unicode.ToLower(rune0)) + s[1:]
	}
	return s
}

// 将字符串单词首字母大写
func (t *internalFunc) title(s string) string {
	return strings.Title(s)
}

func (t *internalFunc) langType(lang string, typeId int) string {
	switch lang {
	case "go":
		return GoTypes(typeId)
	case "java":
		return JavaTypes(typeId)
	case "kotlin":
		return KotlinTypes(typeId)
	case "thrift":
		return ThriftTypes(typeId)
	case "ts":
		return TsTypes(typeId)
	}
	return strconv.Itoa(typeId)
}

// 将包名替换为.分割, 通常C#,JAVA语言使用"."分割包名
func (t *internalFunc) langPkg(lang string, pkg string) string {
	switch lang {
	case "java", "kotlin", "csharp", "thrift":
		return strings.Replace(pkg, "/", ".", -1)
	case "go", "rust", "php", "python":
		i := strings.LastIndexAny(pkg, "/.")
		if i != -1 {
			return pkg[i+1:]
		}
	}
	return pkg
}

// 返回类型默认值
func (t *internalFunc) langDefaultValue(lang string, typeId int) string {
	switch lang {
	case "go", "thrift","ts":
		return GoValues(typeId)
	case "java":
	case "kotlin":
		return JavaValues(typeId)
	}
	return CommonValues(typeId)
}

// 是否相等，如：{{equal "go" "rust"}
func (t *internalFunc) equal(v1, v2 interface{}) bool {
	return v1 == v2
}

// 是否包含
func (t *internalFunc) contain(v interface{}, s string) bool {
	if v == nil {
		return false
	}
	return strings.Contains(t.str(v), s)
}

// 是否以指定字符开始
func (t *internalFunc) startsWith(v interface{}, s string) bool {
	if v == nil {
		return false
	}
	return strings.HasPrefix(t.str(v), s)
}

// 是否以指定字符结束
func (t *internalFunc) endsWith(v interface{}, s string) bool {
	if v == nil {
		return false
	}
	return strings.HasSuffix(t.str(v), s)
}

// 判断是否为true
func (t *internalFunc) boolInt(i int32) bool {
	return i > 0
}

// 加法
func (t *internalFunc) plus(x, y int) int {
	return x + y
}

// 乘法
func (t *internalFunc) multi(x, y interface{}) interface{} {
	fx, ok := x.(float64)
	if ok {
		switch y.(type) {
		case float32:
			return fx * float64(y.(float32))
		case float64:
			return fx * y.(float64)
		case int:
			return fx * float64(y.(int))
		case int32:
			return fx * float64(y.(int32))
		case int64:
			return fx * float64(y.(int64))
		}
	}
	panic("not support")
}

// I32转为字符
func (t *internalFunc) str(i interface{}) string {
	return util.Str(i)
}

// 是否为空
func (t *internalFunc) isEmpty(s string) bool {
	if s == "" {
		return true
	}
	return strings.TrimSpace(s) == ""
}

// 转换为HTML
func (t *internalFunc) rawHtml(v interface{}) ht.HTML {
	return ht.HTML(util.Str(v))
}

//求余
func (t *internalFunc) mathRemain(i int, j int) int {
	return i % j
}
