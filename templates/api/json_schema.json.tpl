#!target:api/json-schema/{{.table.Prefix}}/{{substr_n .table.Name "_" 1}}-json-schema.json
{
    "type": "object",
    "properties": {
        "errCode": {
            "type": "integer",
            "title": "错误码",
            "description": "0:正常,1:失败"
        },
        "errMsg": {
            "type": "string",
            "title": "错误信息"
        },
        "data": {
            "type": "object",
            "title": "数据",
            "properties": {
                {{$columns := .columns}} \
                {{range $i,$c := $columns}}
                    "{{$c.Prop}}":{
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
                {{end}}
            }
        }
    },
    "required": [
        "errCode",
        "errMsg",
        "data"
    ]
}