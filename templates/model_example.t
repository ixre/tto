/*
 这是一个模型示例文件，可通过模型逆向生成数据库元数据，且支持根据模型生成代码。

 按照以下格式进行定义模型类（语法同Struct)，通过`prefix:"user_"`来指定表的前缀

```
 // [表名] [表备注]
 struct [结构名] struct{
	// [列备注]
	[列名] [类型] `db:"[字段名]"`
 }
```

 # 命令行生成,传递model参数，自动识别目录中的模型文件(.t后缀)
 ```
 tto -model /templates 
 ```
 # 代码生成
 ```
 调用tto.ReadTables转换为tto.Tables
 ```
*/

// t_user 用户
type User struct{
	// 编号
	Id int64 `prefix:"user_" db:"id" pk:"yes" auto:"yes"`
} 

// t_user_profile 用户资料
type UserProfile struct{
	// 编号
	Id int64 `db:"id" pk:"yes" auto:"yes"`
	// 用户编号
	UserId int64 `db:"user_id"`
	// 用户名
	UserName string `db:"user_name"`
	// 是否启用
	IsEnabled bool `db:"is_enabled"`
} 