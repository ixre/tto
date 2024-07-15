#!kind:1
#!target:api/json-schema.md

# API接口-JSON Schema

该文档用于前端接口调用及API文档


目录:

{{range $i,$table := .tables}}\
- [{{$table.Comment}}](#{{$table.Comment}})
{{end}}


通用规范

- 响应结果规范

```json
{
    "title": "响应数据",
    "type": "object",
    "properties": {
        "code": {
            "type": "integer",
            "title": "状态码",
            "description": "0:正常,<1:失败和其他错误"
        },
        "msg": {
            "type": "string",
            "title": "错误信息"
        },
        "data": {
            "type": "object",
            "title": "数据",
            "properties": {}
        }
    },
    "required": [
        "code",
        "msg"
    ]
}
```

{{$base_path := .global.base_path}}
{{range $i,$table := .tables}}\
{{$entityName := $table.Title}}\
{{$path := join $base_path (path $table.Name) "/"}}\
{{$pkName := lower_title $table.Pk}}\
{{$pkType := type "ts" $table.PkType}}\
{{$columns := $table.Columns}} \

## {{$table.Comment}}

### 新增接口

- 接口名称: 新增{{$table.Comment}}
- 请求方式: POST
- 接口地址: /{{$path}}

```ts
/**
 * 创建{{$table.Comment}}
 *
 * @param data 数据
 * @returns 请求
 */
export const create{{$entityName}} = (data: {{$entityName}}) => request({
  url: `{{$path}}`,
  method: 'POST',
  data
})
```

### 更新接口

- 接口名称: 修改{{$table.Comment}}
- 请求方式: PUT
- 接口地址: /{{$path}}/{id}

```ts
/**
 * 更新{{$table.Comment}}
 *
 * @param {{$pkName}} 编号
 * @param data 数据
 * @returns 结果
 */
export const update{{$entityName}} = (id: {{$pkType}}, data: {{$entityName}}) => request({
  url: `{{$path}}/${id}`,
  method: 'PUT',
  data
})
```

### 查询接口

- 接口名称: 获取{{$table.Comment}}
- 请求方式: GET
- 接口地址: /{{$path}}/{id}

```ts
/**
 * 获取{{$table.Comment}}
 *
 * @param {{$pkName}} 编号
 * @param params 可选参数
 * @returns {{$table.Comment}}
 */
export const get{{$entityName}} = (id: {{type "ts" $table.PkType}}, params: any = {}):Promise<{data:{{$entityName}}}> => request({
  url: `{{$path}}/${id}`,
  method: 'GET',
  params
})
```

### 查询列表

- 接口名称: 查询{{$table.Comment}}
- 请求方式: GET
- 接口地址: /{{$path}}

```ts
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
```

### 删除接口

- 接口名称: 删除{{$table.Comment}}
- 请求方式: DELETE
- 接口地址: /{{$path}}/{id}

```ts
/**
 * 删除{{$table.Comment}}
 *
 * @param {{$pkName}} 编号
 * @returns 结果
 */
export const delete{{$entityName}} = ({{$pkName}}: {{$pkType}}) => request({
  url: `{{$path}}/${{"{"}}{{$pkName}}}`,
  method: 'DELETE'
})
```

### 实体类型

`ts`语言:

```ts
// {{$table.Comment}}对象
export class {{$table.Title}} {
    {{range $i,$c := $columns}}// {{$c.Comment}}
    {{lower_title $c.Prop}}: {{type "ts" $c.Type}} = {{if eq $c.Render.Element "radio"}} \
        {{default "ts" $c.Type}} + 1 \
    {{else}} \
      {{default "ts" $c.Type}} \
    {{end}};
    {{end}}
}
```

实体文档:

```json
{
    "title": "{{$table.Comment}}",
    "type": "object",
    "properties": {
        {{range $i,$c := $columns}}
        "{{lower_title $c.Prop}}":{
            {{ $javaType := type "java" $c.Type}}\
            {{if eq $javaType "int"}}"type": "integer",\
            {{else if eq $javaType "long"}}"type": "integer",\
            {{else if eq $javaType "boolean"}}"type": "boolean",\
            {{else if eq $javaType "float"}}"type": "number",\
            {{else if eq $javaType "double"}}"type": "number",\
            {{else if eq $javaType "BigDecimal"}}"type": "number",\
            {{else if eq $javaType "Date"}}"type": "object",\
            {{else if eq $javaType "Byte[]"}}"type": "object",\
            {{else if eq $javaType "String"}}"type": "string",\
            {{else}}"type": "object",{{end}}
            "title": "{{$c.Comment}}"
        }\
        {{if not (is_last $i $columns)}},{{end}}
        {{end}}\
    },
    "required": [
        {{range $i,$c := $columns}}\
        "{{lower_title $c.Prop}}"{{if not (is_last $i $columns)}},{{end}}
        {{end}}\
    ]
}
```


{{end}}


