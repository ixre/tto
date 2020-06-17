#!lang:ts＃!name:API和定义文件
#!target:ts/feature/{{name_path .table.Name}}/api.ts
import request from '@/utils/request'
{{$path := str_join "/" .global.url_prefix (name_path .table.Name)}}

// {{.table.Comment}}对象
export interface I{{.table.Title}} {
    {{range $i,$c := .columns}}// {{$c.Comment}}
    {{lower_title $c.Prop}}:{{type "ts" $c.Type}}
    {{end}}
}

export const default{{.table.Title}}:()=>I{{.table.Title}}=()=>{
    return {
        {{range $i,$c := .columns}}
        {{lower_title $c.Prop}}:{{default "ts" $c.Type}},{{end}}
    };
}

export const get{{.table.Title}} = (id: any, params: any = {}) =>
    request({
        url: `{{$path}}/${id}`,
        method: 'get',
        params:{...params}
    })

export const get{{.table.Title}}List = (params: any = {}) =>
    request({
        url: '{{$path}}',
        method: 'get',
        params:{...params}
    })

export const create{{.table.Title}} = (data: any) =>
    request({
        url: '{{$path}}',
        method: 'post',
        data
    })

export const update{{.table.Title}} = (id: any, data: any) =>
    request({
        url: `{{$path}}/${id}`,
        method: 'put',
        data
    })

export const delete{{.table.Title}} = (id: any) =>
    request({
        url: `{{$path}}/${id}`,
        method: 'delete'
    });

export const batchDelete{{.table.Title}} = (arr: any[]) =>
    request({
        url: '{{$path}}',
        method: 'delete',
        data:arr
    });


export const getPaging{{.table.Title}} = (page:number,rows:number,params: any) =>
    request({
        url: '{{$path}}/paging',
        method: 'get',
        params:{page,rows,params}
    })