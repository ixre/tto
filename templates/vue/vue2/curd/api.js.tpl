#!lang:ts＃!name:API和定义文件
#!target:vue2/{{name_path .table.Name}}/api.js
import request from '@/utils/request'
{{$columns := .columns}}
{{$path := join .global.base_path (name_path .table.Name) "/"}}\

// {{.table.Comment}}对象
class {{.table.Title}} {
    constructor({
     {{range $i,$c := .columns}}{{$c.Prop}} = {{default "ts" $c.Type}}\
     {{if eq $c.Render.Element "radio"}} + 1{{end}} \
     {{if not (is_last $i $columns)}},{{end}}
     {{end}}\
    } = {}) {
      {{range $i,$c := .columns}} \
      this.{{$c.Prop}} = {{$c.Prop}};// {{$c.Comment}}
      {{end}}
    }
}

export const get{{.table.Title}} = (id, params = {}) => request({
  url: `{{$path}}/${id}`,
  method: 'GET',
  params: { ...params }
})

export const query{{.table.Title}}List = (params = {}) => request({
  url: '{{$path}}',
  method: 'GET',
  params: { ...params }
})

export const create{{.table.Title}} = (data = {}) => request({
  url: '{{$path}}',
  method: 'POST',
  data
})

export const update{{.table.Title}} = (id, data = {}) => request({
  url: `{{$path}}/${id}`,
  method: 'PUT',
  data
})

export const delete{{.table.Title}} = (id) => request({
  url: `{{$path}}/${id}`,
  method: 'DELETE'
})

export const batchDelete{{.table.Title}} = (arr = []) => request({
  url: '{{$path}}',
  method: 'DELETE',
  data: arr
})

export const getPaging{{.table.Title}} = (page = 0, rows = 100, params = {}) => request({
  url: '{{$path}}/paging',
  method: 'GET',
  params: { page, rows, params }
})