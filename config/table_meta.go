package config

type TableConfig struct {
	// 结构元数据
	Struct *TableMeta `json:"struct"`
	// 字段元数据
	Fields map[string]*ColumnMeta `json:"field"`
}

type TableMeta struct {
}

type ColumnMeta struct {
	// 标题
	Title string `json:"title"`
	// 显示设置
	Render *PropRenderOptions `json:"render"`
}

type PropRenderOptions struct {
	// 是否可见
	Visible bool `json:"visible"`
	// 显示元素
	Element string `json:"element"`
	// 如果Element为select,radio时可用
	Options map[string]string `json:"options"`
}
