#!type:1
#!target:{{.global.pkg}}/mongo/service/{{.table.Name}}_service.go
{{$title := .table.Title}}
{{$shortTitle := .table.ShortTitle}}
{{$p := substr .table.Name 0 1 }}
{{$pkName := lower_title .table.Pk}}
{{$pkProp := .table.PkProp}}
{{$comment := .table.Comment}}
/** #! 主键类型 */
{{$pkType := type "go" .table.PkType }}
/** #! 服务结构名称 */
{{$ifaceName := join (title .table.Title) "Service"}}
{{$structName := join (lower_title .table.Title) "ServiceImpl"}}
package impl

import (
	"errors"
	"{{.global.pkg}}/repo"
	"{{.global.pkg}}/repo/impl"
	"{{.global.pkg}}/repo/model"
	"{{.global.pkg}}/proto"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

var _ I{{$ifaceName}} = new({{$structName}})

type I{{$ifaceName}} interface{
	// Save{{$shortTitle}} 保存{{$comment}}
	Save{{$shortTitle}}(v *model.{{$shortTitle}}) error	
	// Get{{$shortTitle}} 获取{{$comment}}
	Get{{$shortTitle}}({{$pkName}} {{$pkType}}) *model.{{$shortTitle}}
	// Delete{{$shortTitle}} 删除{{$comment}}
	Delete{{$shortTitle}}({{$pkName}} {{$pkType}}) error
	// Paging{{$shortTitle}} 获取{{$comment}}分页数据
	Paging{{$shortTitle}}(filter bson.M, begin int, size int,
	opts ...*options.FindOptions) (int, []*model.{{$shortTitle}}, error)
}

type {{$structName}} struct {
	r repo.I{{.table.Title}}Repository
}

func New{{$shortTitle}}Service(r repo.I{{.table.Title}}Repository) *{{$structName}} {
	return &{{$structName}}{
		r:   r,
	}
}

// Save{{$shortTitle}} 保存{{$comment}}
func ({{$p}} *{{$structName}}) Save{{$shortTitle}}(v *model.{{$shortTitle}}) error {
	var dst *model.{{$shortTitle}}
	{{if equal_any .table.PkType 3 4 5}}\
    if v.{{$pkProp}} > 0 {
    {{else}}
    if v.{{$pkProp}} != "" {
    {{end}}
		if dst = {{$p}}.r.Get{{$shortTitle}}(v.{{.table.PkProp}}); dst == nil{
            return errors.New("no such data")
        }	
	}else{
		{{$c := try_get .columns "create_time"}} \
        {{if $c}}dst.CreateTime = time.Now().Unix(){{end}}
	}
    /** #! 为对象赋值 */\
    {{range $i,$c := exclude .columns $pkName "create_time" "update_time"}}
    {{ $goType := type "go" $c.Type}}\
    {{if eq $goType "int"}}dst.{{$c.Prop}} = int(v.{{$c.Prop}})\
    {{else if eq $goType "int16"}}dst.{{$c.Prop}} = int16(v.{{$c.Prop}})\
    {{else if eq $goType "int32"}}dst.{{$c.Prop}} = int32(v.{{$c.Prop}})\
    {{else}}dst.{{$c.Prop}} = v.{{$c.Prop}}{{end}}{{end}}
    {{$c := try_get .columns "update_time"}}\
    {{if $c}}dst.UpdateTime = time.Now().Unix(){{end}}\
	_, err := {{$p}}.r.Save{{$shortTitle}}(dst)
    return err
}

// Get{{$shortTitle}} 获取{{$comment}}
func ({{$p}} *{{$structName}}) Get{{$shortTitle}}({{$pkName}} {{$pkType}}) *model.{{$shortTitle}} {
	return {{$p}}.r.Get{{$shortTitle}}({{$pkName}})
}


// Delete{{$shortTitle}} 删除{{$comment}}
func ({{$p}} *{{$structName}}) Delete{{$shortTitle}}({{$pkName}} {{$pkType}}) error {
	return {{$p}}.r.Delete{{$shortTitle}}({{$pkName}})
}

// Paging{{$shortTitle}} 获取{{$comment}}分页数据
func ({{$p}} *{{$structName}}) Paging{{$shortTitle}}(filter bson.M, begin int, size int,
	opts ...*options.FindOptions) (int, []*model.{{$shortTitle}}, error){
	return {{$p}}.r.PagingQuery{{$shortTitle}}(filter,begin,size,opts)
}
