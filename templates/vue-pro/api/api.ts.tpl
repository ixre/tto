#!kind:1
#!lang:ts＃!name:API和定义文件
#!target:vue-pro/api/api_gen.ts
import {request,Result} from '@/ext/utils'

// 接口清单:
{{range $i,$table := .tables}}\
// {{plus $i 1}}: [{{$table.Comment}}]({{$table.Name}})
{{end}}

// 引用组件,以支持动态路由加载出错时输出错误信息
{{range $i,$table := .tables}}\
//import "@/view/{{$table.Title}}/index.vue"
{{end}}


{{$base_path := .global.base_path}}
{{range $i,$table := .tables}}\

{{$entityName := $table.Title}}
{{$columns := $table.Columns}}
{{$path := join $base_path (path $table.Name) "/"}}\
{{$pkName := lower_title $table.Pk}}
{{$pkType := type "ts" $table.PkType}}

// {{$table.Comment}}对象
export class {{$entityName}} {
    {{range $i,$c := $columns}}// {{$c.Comment}}
    {{lower_title $c.Prop}}: {{type "ts" $c.Type}} = {{if eq $c.Render.Element "radio"}} \
        {{default "ts" $c.Type}} + 1 \
    {{else}} \
      {{default "ts" $c.Type}} \
    {{end}};
    {{end}}
}


/**
 * 获取{{$table.Comment}}
 *
 * @param {{$pkName}} 编号
 * @param params 可选参数
 * @returns {{$table.Comment}}
 */
export const get{{$entityName}} = ({{$pkName}}: {{type "ts" $table.PkType}}, params: any = {}):Promise<{data:{{$entityName}}}> => request({
  url: `{{$path}}/${{"{"}}{{$pkName}}{{"}"}}`,
  method: 'GET',
  params
})

/**
 * 查询{{$table.Comment}}列表
 *
 * @param params 请求参数，默认为空对象
 * @returns {{$table.Comment}}数组
 */
export const query{{$entityName}}List = (params: any = {}):Promise<{data:Array<{{$entityName}}>}> => request({
  url: `{{$path}}`,
  method: 'GET',
  params
})


/**
 * 创建{{$table.Comment}}
 *
 * @param data 数据
 * @returns 请求结果
 */
export const create{{$entityName}} = (data: {{$entityName}}): Promise<{ data: Result }> => request({
  url: `{{$path}}`,
  method: 'POST',
  data
})


/**
 * 更新{{$table.Comment}}
 *
 * @param {{$pkName}} 编号
 * @param data 数据
 * @returns 结果
 */
export const update{{$entityName}} = ({{$pkName}}: {{$pkType}}, data: {{$entityName}}): Promise<{ data: Result }> => request({
  url: `{{$path}}/${{"{"}}{{$pkName}}{{"}"}}`,
  method: 'PUT',
  data
})

/**
 * 删除{{$table.Comment}}
 *
 * @param {{$pkName}} 编号数组
 * @returns 结果
 */
export const delete{{$entityName}} = ({{$pkName}}: Array<{{$pkType}}>): Promise<{ data: Result }> => request({
  url: `{{$path}}/${{"{"}}{{$pkName}}[0]}`,
  method: 'DELETE',
  data: {{$pkName}}.length > 1? { list : {{$pkName}}{{"}"}} : undefined
})


/**
 * 查询分页{{$table.Comment}}
 *
 * @param page 页码
 * @param size 每页显示数量
 * @param params 其他查询参数
 * @returns 分页数据
 */
export const paging{{$entityName}} = (page:number, size:number, params: any):Promise<{
  data:{total:number,rows:Array<any>}}> => request({
  url: `{{$path}}/paging`,
  method: 'GET',
  params: { page, size,...params }
})

{{end}}