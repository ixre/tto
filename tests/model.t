
// 数据源
type DataSource struct {
	// 主键
	Id int `db:"id" pk:"yes" auto:"yes"`
	// 名称
	Name string `db:"name"`
	// 数据库类型
	Type string `db:"type"`
	// 数据库服务器地址
	Server string `db:"server"`
	// 数据库端口
	Port int `db:"port"`
	// 数据库名称
	DbName string `db:"dbName"`
	// 数据库用户
	Username string `db:"username"`
	// 数据库密码
	Password string `db:"password"`
	// 数据库编码
	Charset string `db:"charset"`
	// 创建时间
	CreateTime int `db:"createTime"`
}