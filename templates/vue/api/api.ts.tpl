#!lang:ts＃!name:API和定义文件
#!target:vue/api/{{.table.Title}}Api.ts
import {request} from '../utils'
{{$entityName := .table.Title}}
{{$columns := .columns}}
{{$path := join .global.base_path (name_path .table.Name) "/"}}\
{{$pkName := lower_title .table.Pk}}
{{$pkType := type "ts" .table.PkType}}

// {{.table.Comment}}对象
export class {{$entityName}} {
    {{range $i,$c := .columns}}// {{$c.Comment}}
    {{$c.Prop}}: {{type "ts" $c.Type}} = {{if eq $c.Render.Element "radio"}} \
        {{default "ts" $c.Type}} + 1 \
    {{else}} \
      {{default "ts" $c.Type}} \
    {{end}};
    {{end}}
}

// {{.table.Comment}}数据映射类
export interface Paging{{$entityName}} {
  {{range $i,$c := .columns}} \
  // {{$c.Comment}}
  {{$c.Name}}:{{type "ts" $c.Type}};
  {{end}}
}

// 获取{{.table.Comment}}
export const get{{$entityName}} = (id: {{type "ts" .table.PkType}}, params: any = {}):Promise<{data:{{$entityName}}}> => request({
  url: `{{$path}}/${id}`,
  method: 'GET',
  params
})

// 查询{{.table.Comment}}列表
export const query{{$entityName}}List = (params: any = {}):Promise<{data:Array<{{$entityName}}>}> => request({
  url: '{{$path}}',
  method: 'GET',
  params
})

// 创建{{.table.Comment}}
export const create{{$entityName}} = (data: {{$entityName}}) => request({
  url: '{{$path}}',
  method: 'POST',
  data
})

// 保存{{.table.Comment}}
export const update{{$entityName}} = (id: {{$pkType}}, data: {{$entityName}}) => request({
  url: `{{$path}}/${id}`,
  method: 'PUT',
  data
})

// 删除{{.table.Comment}}
export const delete{{$entityName}} = ({{$pkName}}: Array<{{$pkType}}>) => request({
  url: `{{$path}}/${{"{"}}{{$pkName}}[0]}`,
  method: 'DELETE',
  data: {{$pkName}}.length > 1? { list : {{$pkName}}{{"}"}} : undefined
})


// 查询{{.table.Comment}}分页数据
export const queryPaging{{$entityName}} = (page:number, size:number, params: any):Promise<{
  data:{total:number,rows:Array<Paging{{$entityName}}>}}> => request({
  url: '{{$path}}/paging',
  method: 'GET',
  params: { page, size,...params }
})