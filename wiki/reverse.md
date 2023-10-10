# 反向生成代码

`tto`可通过定义元数据模型, 逆向转换为数据库元数据. 通过模型反向生成, 在无数据库的情况下也能生成代码.

## 模型定义

按照以下格式进行定义模型类（语法同go语言Struct)，通过`prefix:"user_"`来指定表的前缀

```go
 // [表名] [表备注]
type [结构名] struct{
 // [列备注]
    [列名] [类型] `db:"[字段名]"`
}
```

## 反向生成

命令行生成,传递`model`参数，自动识别目录中的模型文件, 默认后缀名为`.t`.

```shell
tto -model /templates 
```

## 模型逆向

调用`tto.ReadTables`转换为`tto.Tables`, 参照单元测试: [model_revert_test.go](../tests/model_revert_test.go)
