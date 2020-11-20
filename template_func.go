package tto

import (
	"github.com/ixre/gof/types"
	lang2 "github.com/ixre/tto/lang"
	"reflect"
	"strings"
	ht "text/template"
	"unicode"
)

type internalFunc struct {
}

// 返回模板函数
func (t *internalFunc) funcMap() ht.FuncMap {
	fm := make(map[string]interface{})
	fm["boolInt"] = t.boolInt
	fm["isEmpty"] = t.isEmpty
	fm["add"] = t.plus
	fm["plus"] = t.plus
	fm["multi"] = t.multi
	fm["mathRemain"] = t.mathRemain
	// 单词首字大写
	fm["title"] = t.title
	// 将名称转为路径,规则： 替换首个"_"为"/"
	fm["name_path"] = t.nameToPath
	// 小写
	fm["lower"] = t.lower
	// 大写
	fm["upper"] = t.upper
	// 首字母小写: 如:{{lower_title .table.Name}}
	fm["lower_title"] = t.lowerTitle
	// 类型: 如:{{type "go" .columns[0].Type}}
	fm["type"] = t.langType
	// 返回SQL/ORM类型, 如：{{sql_type "py" .columns[0].Type}}
	fm["sql_type"] = t.sqlType
	// 包名: {{pkg "go" "com/tto/pkg"}}
	fm["pkg"] = t.langPkg
	// 默认值, 如:{{default "go" .columns[0].Type}}
	fm["default"] = t.langDefaultValue
	// 是否相等，如：{{equal "go" "rust"}
	fm["equal"] = t.equal
	// 是否与任意值相等,　如：{{equal_any 1 2 3 4}}, 1是否与2,3,4相等
	fm["equal_any"] = t.equalAnyValue
	// 替换,如: {{replace "table_name" "_" "-"}}
	fm["replace"] = t.replace
	// 替换N次,如: {{replace_n "table_name" "_" "-" 1}}
	fm["replace_n"] = t.replaceN
	// 截取第N个字符位置后的字符串,如:{{substr_n "sys_user_list" "_" 1}}得到:user_list
	fm["substr_n"] = t.substrN
	// 截取索引为N的元素,如:{{get_n .tables 0}}
	fm["get_n"] = t.getN
	// 字符组合,如：{{str_join "1" "2" "3" ","}}结果为:1,2,3
	fm["join"] = t.strJoin
	fm["str_join"] = t.strJoin
	// 包含函数, 如:{{contain .table.Pk "id"}}
	fm["contain"] = t.contain
	// 是否以指定字符开始, 如:{{starts_with .table.Pk "id"}}
	fm["starts_with"] = t.startsWith
	// 是否以指定字符结束, 如:{{ends_with .table.Pk "id"}}
	fm["ends_with"] = t.endsWith
	// 返回是否为数组中的最后一个元素索引
	fm["is_last"] = t.isLast
	// 排除列元素, 组成新的列数组, 如：{{ $columns := exclude .columns "id","create_time" }}
	fm["exclude"] = t.exclude
	// 尝试获取一个列,返回列及是否存在的Boolean, 如: {{ $c,$exist := try_get .columns "update_time" }}
	fm["try_get"] = t.tryGet
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

// 将名称转为路径,规则： 替换首个"_"为"/"
func (t *internalFunc) nameToPath(s string) string {
	return strings.Replace(s, "_", "/", 1)
}

func (t *internalFunc) langType(lang string, typeId int) string {
	return lang2.Get(lang).ParseType(typeId)
	//switch lang {
	//case "go":
	//	return l.GoTypes(typeId)
	//case "java":
	//	return l.JavaTypes(typeId)
	//case "kotlin":
	//	return l.KotlinTypes(typeId)
	//case "thrift":
	//	return l.ThriftTypes(typeId)
	//case "protobuf":
	//	return l.ProtobufTypes(typeId)
	//case "ts":
	//	return l.TsTypes(typeId)
	//case "py":
	//	return l.PyTypes(typeId)
	//}
	//return strconv.Itoa(typeId)
}

// 返回SQL/ORM类型
func (t *internalFunc) sqlType(lang string, typeId int, len int) string {
	if lang == "py" {
		return PySqlTypes(typeId, len)
	}
	panic("not support language except for py")
}

// 将包名替换为.分割, 通常C#,JAVA语言使用"."分割包名
func (t *internalFunc) langPkg(lang string, pkg string) string {
	return lang2.Get(lang).ParsePkg(pkg)
}

// 返回类型默认值
func (t *internalFunc) langDefaultValue(lang string, typeId int) string {
	return lang2.Get(lang).DefaultValue(typeId)
	//switch lang {
	//case "go", "thrift","protobuf","ts":
	//	return l.GoValues(typeId)
	//case "java","kotlin":
	//	return l.JavaValues(typeId)
	//case "py":
	//	return l.PythonValues(typeId)
	//}
	//return l.CommonValues(typeId)
}

// 是否相等，如：{{equal "go" "rust"}
func (t *internalFunc) equal(v1, v2 interface{}) bool {
	return v1 == v2
}

// 是否与任意值相等,　如：{{equal_any 1 2 3 4}}, 1是否与2,3,4相等
func (t *internalFunc) equalAnyValue(src interface{}, args ...interface{}) bool {
	for _, v := range args {
		if v == src {
			return true
		}
	}
	return false
}

// 替换,如: {{replace "table_name" "_" "-"}}
func (t *internalFunc) replace(s, oldStr, newStr string) string {
	return t.replaceN(s, oldStr, newStr, -1)
}

// 替换N次,如: {{replace_n "table_name" "_" "-" 1}}
func (t *internalFunc) replaceN(s, oldStr, newStr string, n int) string {
	return strings.Replace(s, oldStr, newStr, n)
}

// 截取第N个字符位置后的字符串,如:{{substr_n "sys_user_list" "_" 1}}得到:user_list
func (t *internalFunc) substrN(s, char string, n int) string {
	i, times := 0, 0
	for {
		i = strings.Index(s, char)
		if i == -1 {
			break
		}
		s = s[i+1:]
		if times++; times == n {
			break
		}
	}
	return s
}

// 截取索引为N的元素,如:{{get_n .tables 0}}
func (t *internalFunc) getN(args interface{}, n int) interface{} {
	kind := reflect.TypeOf(args).Kind()
	if kind == reflect.Slice || kind == reflect.Array {
		value := reflect.ValueOf(args)
		if value.Len()-1 < n {
			return nil
		}
		return value.Index(n).Interface()
	}
	return nil

	//if len(args)-1 < n{
	//	return nil
	//}
	//return args[n]
}

// 字符组合,如：{{str_join "1","2","3" ","}}结果为:1,2,3
func (t *internalFunc) strJoin(s string, args ...string) string {
	l := len(args)
	if  l== 0{
		return s
	}
	if l == 1{
		return s +args[0]
	}
	n := append([]string{s},args[:len(args)-1]...)
	return strings.Join(n,args[l-1])
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

// 返回是否为数组中的最后一个元素索引,如：
// {{$columns := .columns}}
// {{range $,$v := .columns}}
//	  {{if is_last $i .columns}}
//		last column
//	  {{end}}
// {{end}}
func (t *internalFunc) isLast(i int, arr interface{}) bool {
	kind := reflect.TypeOf(arr).Kind()
	if kind == reflect.Slice || kind == reflect.Array {
		return i == reflect.ValueOf(arr).Len()-1
	}
	return false
}

// 排除列元素, 组成新的列数组, 如：{{ $columns := exclude .columns "id","create_time" }}
func (t *internalFunc) exclude(columns []*Column, names ...string) []*Column {
	arr := make([]*Column, 0)
	for _, c := range columns {
		b := false
		for _, n := range names {
			if c.Name == n {
				b = true
				break
			}
		}
		if !b {
			arr = append(arr, c)
		}
	}
	return arr
}

// 尝试获取一个列,返回列. 如果不存在,返回空, 如: {{ $c := try_get .columns "update_time" }}
func (t *internalFunc) tryGet(columns []*Column, name string) *Column {
	for _, c := range columns {
		if c.Name == name {
			return c
		}
	}
	return nil
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
	return types.String(i)
}

// 是否为空
func (t *internalFunc) isEmpty(s string) bool {
	if s == "" {
		return true
	}
	return strings.TrimSpace(s) == ""
}

//求余
func (t *internalFunc) mathRemain(i int, j int) int {
	return i % j
}
