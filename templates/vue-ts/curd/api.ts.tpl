#!lang:ts＃!name:API和定义文件
#!target:ts/feature/{{name_path .table.Name}}/api.ts
import request from '@/utils/request'
{{$columns := .columns}}
{{$path := join .global.base_path (name_path .table.Name) "/"}}\

// {{.table.Comment}}对象
export interface I{{.table.Title}} {
    {{range $i,$c := .columns}}// {{$c.Comment}}
    {{$c.Prop}}: {{type "ts" $c.Type}}{{if not (is_last $i $columns)}},{{end}}
    {{end}}
}

export const default{{.table.Title}}:()=>I{{.table.Title}} = () => {
  return {
    {{range $i,$c := .columns}}
    {{if eq $c.Render.Element "radio"}}{{$c.Prop}}: {{default "ts" $c.Type}} + 1,
    {{else}}{{$c.Prop}}: {{default "ts" $c.Type}}{{if not (is_last $i $columns)}},{{end}}\
    {{end}}{{end}}
  }
}

export const get{{.table.Title}} = (id: any, params: any = {}) => request({
  url: `{{$path}}/${id}`,
  method: 'GET',
  params: { ...params }
})

export const query{{.table.Title}}List = (params: any = {}) => request({
  url: '{{$path}}',
  method: 'GET',
  params: { ...params }
})

export const create{{.table.Title}} = (data: any) => request({
  url: '{{$path}}',
  method: 'POST',
  data
})

export const update{{.table.Title}} = (id: any, data: any) => request({
  url: `{{$path}}/${id}`,
  method: 'PUT',
  data
})

export const delete{{.table.Title}} = (id: any) => request({
  url: `{{$path}}/${id}`,
  method: 'DELETE'
})

export const batchDelete{{.table.Title}} = (arr: any[]) => request({
  url: '{{$path}}',
  method: 'DELETE',
  data: arr
})

export const getPaging{{.table.Title}} = (page:number, rows:number, params: any) => request({
  url: '{{$path}}/paging',
  method: 'GET',
  params: { page, rows, params }
})