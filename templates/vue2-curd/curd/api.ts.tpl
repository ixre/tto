#!lang:ts＃!name:API和定义文件
#!target:vue/{{name_path .table.Name}}/api.ts
import {request} from '@/utils'
{{$columns := .columns}}
{{$path := join .global.base_path (name_path .table.Name) "/"}}\

// {{.table.Comment}}对象
export class {{.table.Title}} {
    {{range $i,$c := .columns}}// {{$c.Comment}}
    {{$c.Prop}}: {{type "ts" $c.Type}} = {{if eq $c.Render.Element "radio"}} \
        {{default "ts" $c.Type}} + 1 \
    {{else}} \
      {{default "ts" $c.Type}} \
    {{end}};
    {{end}}
}

// {{.table.Comment}}数据映射类
export interface Paging{{.table.Title}} {
  {{range $i,$c := .columns}} \
  // {{$c.Comment}}
  {{$c.Name}}:{{type "ts" $c.Type}};
  {{end}}
}

// 获取{{.table.Comment}}
export const get{{.table.Title}} = (id: {{type "ts" .table.PkType}}, params: any = {}):Promise<{data:{{.table.Title}}}> => request({
  url: `{{$path}}/${id}`,
  method: 'GET',
  params
})

// 查询{{.table.Comment}}列表
export const query{{.table.Title}}List = (params: any = {}) => request({
  url: '{{$path}}',
  method: 'GET',
  params
})

// 创建{{.table.Comment}}
export const create{{.table.Title}} = (data: {{.table.Title}}) => request({
  url: '{{$path}}',
  method: 'POST',
  data
})

// 保存{{.table.Comment}}
export const update{{.table.Title}} = (id: {{type "ts" .table.PkType}}, data: {{.table.Title}}) => request({
  url: `{{$path}}/${id}`,
  method: 'PUT',
  data
})

// 删除{{.table.Comment}}
export const delete{{.table.Title}} = (id: {{type "ts" .table.PkType}}) => request({
  url: `{{$path}}/${id}`,
  method: 'DELETE'
})

// 批量删除{{.table.Comment}}
export const batchDelete{{.table.Title}} = (arr: Array<{{type "ts" .table.PkType}}>) => request({
  url: '{{$path}}',
  method: 'DELETE',
  data: arr
})

// 查询{{.table.Comment}}分页数据
export const queryPaging{{.table.Title}} = (page:number, size:number, params: any):Promise<{
  data:{total:number,rows:Array<Paging{{.table.Title}}>}}> => request({
  url: '{{$path}}/paging',
  method: 'GET',
  params: { page, size,...params }
})