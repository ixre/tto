package repo
#!type:0#!lang:go
#!target:{{.global.pkg}}/mongo/repo/{{.table.Name}}_repo.go
{{$title := .table.Title}}
{{$shortTitle := .table.ShortTitle}}

import(
    "{{pkg "go" .global.pkg}}/mongo/repo/model"
    "go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type I{{$title}}Repository interface{
    // Get{{$shortTitle}} 获取{{.table.Comment}}
    Get{{$shortTitle}}(primary interface{})*model.{{$shortTitle}}
    // Save{{$shortTitle}} 保存{{.table.Comment}}
    Save{{$shortTitle}}(v *model.{{$shortTitle}})error
    // Delete{{$shortTitle}} 删除{{.table.Comment}}
    Delete{{$shortTitle}}(primary interface{}) error
    // Paging{{$shortTitle}} 查询分页{{.table.Comment}}数据
    Paging{{$shortTitle}}(filter bson.M, begin int, size int,
	opts ...*options.FindOptions) (int, []*model.{{$shortTitle}}, error)
}
