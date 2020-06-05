package tto

import (
	"regexp"
	"strings"
)

var (
	predefineRegexp = regexp.MustCompilePOSIX("#\\!([^\\!-]+):([^#]+?)")
	lineJoinRegexp  = regexp.MustCompile("\\s*\\\\(\\s+)")
)

type CodeTemplate struct {
	path      string
	template  string
	predefine map[string]string
}

func NewTemplate(s string, path string) *CodeTemplate {
	return (&CodeTemplate{path: path}).configure(s)
}

func (g *CodeTemplate) isCodeFile() bool {
	ext := g.path[strings.LastIndex(g.path, ".")+1:]
	return ext == "kt" || ext == "java" || ext == "go" || ext == "ts" || ext == "js" ||
		ext == "cs" || ext == "vb" || ext == "py" || ext == "rb" || ext == "cpp" ||
		ext == "c" || ext == "lua" || ext == "pl"
}

func (g *CodeTemplate) configure(s string) *CodeTemplate {
	if len(s) > 1 && s[0] != '/' {
		if g.isCodeFile() {
			s = `/**
 * this file is auto generated by tto v` + BuildVersion + ` !
 * if you want to modify this code,please read guide doc
 * and modify code template later.
 *
 *  please read user guide on https://github.com/ixre/tto
 *
 */
` + s
		}
	}
	g.predefine = make(map[string]string)
	for _, match := range predefineRegexp.FindAllStringSubmatch(s, -1) {
		g.predefine[match[1]] = match[2]
	}
	g.template = g.format(s)
	return g
}

func (g *CodeTemplate) format(s string) string {
	s = predefineRegexp.ReplaceAllString(s, "")
	s = lineJoinRegexp.ReplaceAllString(s, "")
	return s
}

// 文件路径
func (g *CodeTemplate) FilePath() string {
	return g.path
}

// 获取模板内容
func (g *CodeTemplate) String() string {
	return g.template
}

// 获取预定义的参数
func (g *CodeTemplate) Predefine(key string) (string, bool) {
	n, ok := g.predefine[key]
	return n, ok
}

func (g *CodeTemplate) Replace(s, old string, n int) *CodeTemplate {
	g.template = strings.Replace(g.template, s, old, n)
	return g
}
