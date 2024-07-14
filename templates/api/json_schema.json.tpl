#!target:api/json-schema/{{.table.Prefix}}/{{substr_n .table.Name "_" 1}}-json-schema.md

# {{.table.Comment}} -JSON Schema

新增/更新实体

```json
{
    "title": "{{.table.Comment}}",
    "type": "object",
    "properties": {
        {{$columns := .columns}} \
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

新增/更新/删除响应结果

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
