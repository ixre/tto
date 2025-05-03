# 代码生成器

**TTO 是一款使用 Go 编写的代码生成器,兼容多种数据库,支持多种语言,支持自定义模板生成代码.**

特点:

- 支持 windows,macos,linux 操作系统.
- 支持 mysql/mariadb、postgresql 和 sqlserver 数据库.
- 支持代码模板和模板函数,内置十多种开箱即用代码模板.
- 支持模型逆向生成数据库和代码.
- 支持 Go,JAVA,Python,Js,C#,Kotlin,Html 等多种语言.
- 支持生成前端代码,内置 vue2/vue3 模板
- 内置 Protobuf 和 Thrift 等 RPC 框架模板
- 自带版本升级和自动更新功能

## 安装

### 命令行工具

在 Linux/Mac 下安装，使用以下命令安装

```bash
curl -L https://raw.githubusercontent.com/ixre/tto/master/install | sh
```

Windows 用户进入下载页面([链接](https://github.com/ixre/tto/releases/)),下载最新版本(文件:tto-generator-client.tar.gz)后解压,将目录中的`tto.exe`文件复制到`C:\windows`下完成安装.　请注意 windows10 以下版本需要复制到`C:\windows\System32`目录.

### 图形界面

`tto`提供基于 B/S 的图形界面，可进行模板管理和代码生成等功能，使用`docker`运行图形界面，参考以下命令：

```bash
docker run -d --name gdp -p 8000:8000 \
    -v $(pwd)/conf:/app/conf\
    -v $(pwd)/storage:/app/storage\
    jarry6/gdp
```

或直接下载运行二进制包后，输入以下网址运行:

```text
http://localhost:8000
```

## 升级

`tto`内置了版本更新功能,命令如下:
`tto update`
_注：在 windows 下升级功能如无法正常使用,可以手动重新安装_

## 快速开始

### 　下载程序包

到下载页面([链接](https://github.com/ixre/tto/releases/)),下载最新版本(文件:tto-generator-bin.tar.gz)后解压;

### 　配置数据库

`tto.conf`为程序的默认配置文件,　打开文件进行找到`[database]`节点配置数据库.

### 　使用模板

您可以直接使用安装包里的模板文件,　或按照您的风格对模板进行修改,　甚至单独创建模板.

### 运行命令生成代码

执行以下命令生成代码,代码会生成到`output`目录

`tto -clean`

但实际应用中,推荐使用脚本文件来完成生成,　您可以参考安装包中的示例脚本文件:`./example.sh`;

在 windows 中可以使用`git-bash`来执行该脚本

## 模板

`tto`模板使用`Go Template`,　具体语法参考:

- [Go 模板语法-中](http://www.g-var.com/posts/translation/hugo/hugo-21-go-template-primer/)
- [Go 模板语法-English](https://golang.org/pkg/text/template/)

### 预定义语法

预定义语法用来在代码模板中定义一些数据, 在生成代码时预定义语法不输入任何内容. 预定义语法格式为: !预定义参数名:预定义参数值

目前,支持的预定义语法如下:

- \#!kind: 定义模板生成类型,0:普通,1:生成所有表 2:按表名前缀生成,默认为 0
- \#!target: 用来定义代码文件存放的目标路径
- \#!append: 是否追加到文件,可选值为:true 和 false , 默认为 false
- \#!format: 是否启用格式化代码，可选值为:true 和 false，默认开启
- \#!lang: 指定当前生成代码的语言 如:

```text
#!target:java/{{.global.pkg}}/pojo/{{.table.Title}}{{.global.entity_suffix}}.java
```

多个预定义表达式可以放在一行

```text
#!format:true#!target:Entity.java
```

### 模板注释

模板注释,使用`/** #! 注释 */`的语法,使用`#!`与普通的代码注释区分

```java
/** #! 这是模板注释,不会出现在生成的代码中 */
```

### 模板函数

获取用户环境变量

```text
{{env "PROJECT_MEMBERS"}}
```

大/小写函数: lower 和 upper

```text
{{lower .table.Name}}
{{upper .table.Name}}
```

单词首字大写函数:title

```text
{{title .table.Name}}
```

首字母小写函数: lower_title

```text
{{lower_title .table.Name}}
```

语言类型函数: type

```text
{{type "go" .columns[0].Type}}
```

返回 SQL/ORM 类型: sql_type

```text
{{sql_type "py" .columns[0].Type .columns[0].Length}}
```

返回 ORM 字段类型,通常在 Java 中使用

```text
{{orm_type "java" 3 }}  // 输出: Integer
```

是否为数值类型

```text
{{num_type .table.PkType}}
```

包函数: pkg, 用于获取包的路径

```text
{{pkg "go" .global.pkg}} # github.com/ixre
```

路径函数: path, 用于生成 url 路径

```text
{{path .global.pkg}} # sys/option/list
```

包名函数:

```text
{{pkg_name "go" "github/com/ixre"}} # ixre
```

默认值函数: default

```text
{{default "go" .columns[0].TypeId}}
```

是否相等

```text
{{equal (3%2) 1}}
```

是否与任意值相等, 如表的主键是否为 int 类型

```text
{{equal_any .table.PkType 3 4 5}}
```

替换, 如将`table_name`替换为:`table-name`

```text
{{replace "table_name" "_" "-"}}
```

替换 N 次, 如将`table_name`替换为:`table-name`

```text
{{replace_n "table_name" "_" "-" 1}}
```

截取字符串函数：substr

```text
{{substr "sys_user_list" 0 3 }} # 结果：sys
{{substr "sys_user_list" 4 }} 结果:sys_list
```

截取第 N 个字符位置后的字符串,如以下语句将输出:user_list

```text
{{substr_n "sys_user_list" "_" 1}}
```

截取索引为 N 的元素

```text
{{$first_table := get_n .tables 0}}
```

字符组合,如以下语句将输出:`1,2,3`

```text
{{join "," "1","2","3"}}
{{$api := join "/" .global.base_path (name_path .table.Name)}}
```

包含函数

```text
{{contain .table.Pk "id"}}
```

是否以指定字符开始

```text
{{starts_with .table.Pk "user_"}}
```

是否以指定字符结束

```text
{{ends_with .table.Pk "_time"}}
```

是否为表的列(数组)的最后一列

```text
{{$columns := .columns}}
{{range $,$v := .columns}}
{{if is_last $i .columns}} last column {{end}}
{{if not (is_last $i .columns) }} not last column {{end}}
{{end}}
```

排除列元素, 组成新的列数组, 如：

```text
{{ $columns := exclude .columns "id","create_time" }}
```

尝试获取一个列,返回列, 如:

```text
{{ $c := try_get .columns "update_time" }}
{{if $c}}prop={{$c.Prop}}{{end}}
```

将名称转为路径,规则： 替换首个"\_"为"/"

```text
{{$path := name_path .table.Name}}
```

计算求余数

```go
{{$mod := mod 3 2}}
```

### 代码模板

模板目录默认为`templates`, 我们可以通过结合内置的函数和语法, 生成项目代码.

模板主要包含三大对象:

- global
- table
- columns

按所有表(前缀分组)模板包含对象：

- global
- tables

### 全局变量(global)

输出生成器的版本号

```text
// this file created by generate {{.global.version}}
```

输出包名,包名通过配置文件配置.格式为: com/pkg

```text
package {{.global.pkg}}
```

如果是 Java 或其他语言, 包名以"."分割, 可使用 pkg 函数,如:

```text
// java package
package {{pkg "java" .global.pkg}}
// c# namespace
namespace {{pkg "csharp" .global.pkg}}
```

输出当前时间

```text
generate time {{.global.time}}
```

获取数据库驱动 可选值：pgsql | mysql, 可针对不同数据库生成代码

```text
{{.global.db}}
```

输出自定义变量 用户可以通过在配置文件的节点`[global]`中进行添加变量,如:

```text
[global]
base_path="/api"
```

使用以下语法读取变量

```text
{{.global.base_path}}
```

### 数据表对象(table)

数据表对象对来返回表的信息,包含如下属性:

- Name: 表名
- Prefix: 表前缀
- Pk: 主键,默认为:id
- PkProp: 主键属性, 首字母大写
- PkType: 主键类型编号
- Title: 表名单词首字大写,通常用来表示类型, 如:user_info 对应的 Title 为 UserInfo
- ShortTitle: 同 title, 但不包含前缀
- Comment: 表注释
- Engine: 数据库引擎
- Schema: 架构
- Charset: 数据库编码
- Ordinal: 表的序号

### 数据列对象(columns)

数据列对象存储表的数据列数组, 并且可遍历. 每个数据列都包含如下属性:

- Name: 列名
- Prop: 列名首字大写, 通常用作属性
- IsPk: 是否主键(bool)
- IsAuto: 是否自动生成(bool)
- NotNull: 是否不能为空(bool)
- DbType: 数据库数据类型
- Comment: 注释
- Length: 长度
- Type: 类型编号,使用 type 函数转换为对应语言的类型
- Ordinal: 列的序号

示例:

```text
{{range $i,$c := .columns}}
    列名:$c.Name {{if $c.IsPk}}是主键{{end}}, 类型:{{type "java" $c.Type}}
{{end}}
```

## 模板示例

以下代码用于生成 Java 的 Pojo 对象, 更多示例点击[这里](templates)

```text
#!target:{{.global.pkg}}/pojo/{{.table.Title}}{{.global.entity_suffix}}.java
package {{pkg "java" .global.pkg}}.pojo;

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
public class {{.table.Title}}{{.global.entity_suffix}} {
    {{range $i,$c := .columns}}{{$type := type "java" $c.Type}}
    private {{$type}} {{$c.Name}}
    public void set{{$c.Prop}}({{$type}} {{$c.Name}}){
        this.{{$c.Name}} = {{$c.Name}}
    }

    /** {{$c.Comment}} */{{if $c.IsPk}}
    @Id
    @GeneratedValue(strategy = GenerationType.IDENTITY){{else}}
    @Basic{{end}}
    @Column(name = "{{$c.Name}}"
     {{if not $c.NotNull}}, nullable = true{{end}}
     {{if ne $c.Length 0}},length = {{$c.Length}}{{end}})
    public {{$type}} get{{$c.Prop}}() {
        return this.{{$c.Name}};
    }
    {{end}}
}

```

## 逆向生成代码

参见代码：[generate_test.go](generate_test.go)

**如果您觉得这个项目不错, 请给个 star 吧.**

<img src="images/cq-alipay.png" width="320" style="display:inline-block"/><img src="images/cq-wxpay.png" width="354" style="display:inline-block"/>
