#!lang:ts＃!name:API和定义文件
#!target:ts/feature/{{.table.Prefix}}/{{.table.Name}}/api.ts
import request from '@/utils/request'

// {{.table.Comment}}数据库映射类
export interface I{{.table.Title}}Dbo {
    {{range $i,$c := .columns}}// {{$c.Comment}}
    {{$c.Name}}:{{type "ts" $c.Type}}
    {{end}}
}

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

export const get{{.table.Title}}s = (params: any) =>
    request({
        url: '/{{.table.Name}}/list',
        method: 'get',
        params:{...params,portal:"{{.table.Title}}List"}
    })



export const get{{.table.Title}}List = (params: any) =>
    request({
        url: '/{{.table.Name}}',
        method: 'post',
        params:params
    })

export const get{{.table.Title}} = (id: any, params: any) =>
    request({
        url: `/{{.table.Name}}/${id}`,
        method: 'get',
        params
    })

export const create{{.table.Title}} = (data: any) =>
    request({
        url: '/{{.table.Name}}',
        method: 'post',
        data
    })

export const update{{.table.Title}} = (id: any, data: any) =>
    request({
        url: `/{{.table.Name}}/${id}`,
        method: 'put',
        data
    })

export const delete{{.table.Title}} = (id: any) =>
    request({
        url: `/{{.table.Name}}/${id}`,
        method: 'delete'
    });

export const batchDelete{{.table.Title}} = (arr: any[]) =>
    request({
        url: '/{{.table.Name}}',
        method: 'delete',
        data:arr
    });


export const getPaging{{.table.Title}} = (params: any) =>
    request({
        url: '/{{.table.Name}}/paging',
        method: 'get',
        params:{...params,portal:"{{.table.Title}}List"}
    })