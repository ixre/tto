package tto

import (
	"github.com/ixre/gof/db/orm"
	"github.com/ixre/gof/util"
	"os"
	"strings"
	"unicode"
)

func GetDialect(driver string) (dialect orm.Dialect, driverName string) {
	switch driver {
	case "mysql", "mariadb":
		return &orm.MySqlDialect{}, "mysql"
	case "postgres", "postgresql", "pgsql":
		return &orm.PostgresqlDialect{}, "pgsql"
	}
	return nil, "-"
}

func prefix(str string) string {
	if i := strings.Index(str, "_"); i != -1 {
		return str[:i]
	}
	for i, l := 1, len(str); i < l-1; i++ {
		if unicode.IsUpper(rune(str[i])) {
			return strings.ToLower(str[:i])
		}
	}
	return ""
}

func title(str string, shortUpper bool) string {
	// 小于3且ID大写，则返回大写
	if shortUpper && len(str) < 3 {
		return strings.ToUpper(str)
	}
	arr := strings.Split(str, "_")
	for i, v := range arr {
		arr[i] = strings.Title(v)
	}
	return strings.Join(arr, "")
}

// 不包含前缀的较短的Title
func shortTitle(str string) string {
	arr := strings.Split(str, "_")
	if len(arr) == 1 {
		return strings.Title(str)
	}
	n := make([]string, len(arr)-1)
	for i, v := range arr[1:] {
		n[i] = strings.Title(v)
	}
	return strings.Join(n, "")
}

// 将首字母小写
func lowerTitle(s string) string {
	if rune0 := rune(s[0]); unicode.IsUpper(rune0) {
		return string(unicode.ToLower(rune0)) + s[1:]
	}
	return s
}

// 保存到文件
func SaveFile(s string, path string) error {
	// 将路径转为正确的路径
	//path = filepath.Clean(path)
	// 如果保存到自定义目录, 源文件存在时,自动添加.gen后缀
	if !strings.HasPrefix(path, "output/") {
		fi, _ := os.Stat(path)
		if fi != nil {
			path += ".gen"
		}
	}
	return util.BytesToFile([]byte(strings.TrimSpace(s)), path)
}
