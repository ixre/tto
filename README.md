# 代码生成器

**TTO是一款使用Go编写的代码生成器,可根据模板定制生成代码.**

特点:
- 支持mysql和postgresql数据库
- 支持Go,JAVA,Kotlin,Html,C#语言
- 支持代码模板, 支持模板函数

资源:
- [下载地址](https://github.com/ixre/goex/releases/)
- [Go模板语法-中](http://www.g-var.com/posts/translation/hugo/hugo-21-go-template-primer/)
- [Go模板语法-English](https://golang.org/pkg/text/template/)

## 快速开始

1. 配置数据源
```
下载安装包,解压修改tto.conf文件进行数据源配置.
```
2. 定制修改模板
```
根据实际需求对模板进行修改, 或创建自己的模板. 模板语法请参考: Go Template
```
3. 运行命令生成代码
```bash
tto -conf tto.conf
```

## 预定义语法

预定义语法用来在代码模板中定义一些数据, 在生成代码时预定义语法不输入任何内容.
预定义语法格式为: !预定义参数名:预定义参数值

目前,支持的预定义语法如下:

- !target : 用来定义代码文件存放的目标路径


## 函数

大/小写函数: lower和upper
```
{{lower .table.Name}}
{{upper .table.Name}}
```
单词首字大写函数:title
```
{{title .table.Name}}
```
首字母小写函数: lower_title
```
{{lower_title .table.Name}}
```
语言类型函数: type
```
{{type "go" .columns[0].TypeId}}
```
包名函数: pkg
```
{{pkg "go" .global.Pkg}}
```
默认值函数: default
```
{{default "go" .columns[0].TypeId}}
```

## 模板

模板主要包含三大对象: 

- global
- table
- columns


### global

**用于读取全局变量, global的属性均以大写开头; global为小写.**

输出生成器的版本号
```
// this file created by generate {{.global.Version}}
```
输出包名,包名通过配置文件配置.格式为: com/pkg
```
package {{.global.Pkg}}
```
如果是Java或其他语言, 包名以"."分割, 可使用pkg函数,如:
```
// java package
package {{pkg "java" .global.Pkg}}
// c# namespace
namespace {{pkg "csharp" .global.Pkg}}
```

### table 数据表对象

数据表对象对来返回表的信息,包含如下属性:

- Name: 表名
- Prefix: 表前缀
- Pk: 主键,默认为:id
- PkTypeId: 主键类型编号
- Title: 表名单词首字大写,通常用来表示类型,
  如:user_info对应的Title为UserInfo
- Comment: 表注释
- Engine: 数据库引擎
- Schema: 架构
- Charset: 数据库编码
- Ordinal: 表的序号

### colums 数据列对象

数据列对象存储表的数据列数组, 并且可遍历. 每个数据列都包含如下属性:

- Name: 列名
- Title: 列名首字大写, 同表Title
- IsPk: 是否主键(bool)
- Auto:  是否自动生成(bool)
- NotNull: 是否不能为空(bool)
- Type: 数据类型
- Comment: 注释
- Length: 长度
- TypeId: 类型编号,使用type函数转换为对应语言的类型
- Ordinal: 列的序号

示例:
```
{{range $i,$c := .columns}}
    列名:$c.Name {{if $c.IsPk}}是主键{{end}}, 类型:{{type "java" $c.TypeId}}
{{end}}
```

## 模板示例

以下代码用于生成Java的Pojo对象, 更多示例点击[这里](templates)

```
!target:{{.global.Pkg}}/pojo/{{.table.Title}}Entity.java
package {{pkg "java" .global.Pkg}}.pojo;

import javax.persistence.Basic;
import javax.persistence.Id;
import javax.persistence.Column;
import javax.persistence.Entity;
import javax.persistence.Table;
import javax.persistence.GenerationType;
import javax.persistence.GeneratedValue;

/** {{.table.Comment}} */
@Entity
@Table(name = "{{.table.Name}}", schema = "{{.table.Schema}}")
public class {{.table.Title}}Entity {
    {{range $i,$c := .columns}}{{$type := type "java" $c.TypeId}}
    private {{$type}} {{$c.Name}}
    public void set{{$c.Title}}({{$type}} {{$c.Name}}){
        this.{{$c.Name}} = {{$c.Name}}
    }

    /** {{$c.Comment}} */{{if $c.IsPk}}
    @Id
    @GeneratedValue(strategy = GenerationType.IDENTITY){{else}}
    @Basic{{end}}
    @Column(name = "{{$c.Name}}"
     {{if not $c.NotNull}}, nullable = true{{end}}
     {{if ne $c.Length 0}},length = {{$c.Length}}{{end}})
    public {{$type}} get{{$c.Title}}() {
        return this.{{$c.Name}};
    }
    {{end}}
}

```



**如果您觉得这个项目不错, 请给个star吧.**


<img src="images/cq-alipay.png" width="320" style="display:inline-block"/><img src="images/cq-wxpay.png" width="354" style="display:inline-block"/>
