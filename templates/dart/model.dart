#!target:dart/lib/entity/{{.table.Name}}_model.dart
{{$title := .table.Title}}
{{$columns := .columns}}

/// {{.table.Comment}}
class {{$title}}{{.global.entity_suffix}} {
    {{range $i,$c := .columns}}
    {{$type := type "dart" $c.Type}}\
    // {{$c.Comment}}
    {{$type}} {{$c.Prop}};{{end}}


    {{$title}}{{.global.entity_suffix}}();

    {{$title}}{{.global.entity_suffix}}.createDefault()
        :{{range $i,$c := .columns}}\
        {{$c.Prop}} = {{default "dart" $c.Type}}{{if is_last $i $columns}};{{else}},{{end}}
        {{end}}

    // 拷贝数据
    {{$title}}{{.global.entity_suffix}}.fromJson(Map<String, dynamic> json)
        :{{range $i,$c := $columns}}\
        {{if eq $c.Type 5}}\
        {{$c.Prop}} = BigInt.from(json["{{$c.Prop}}"]){{if is_last $i $columns}};{{else}},{{end}}
        {{else}}
        {{$c.Prop}} = json["{{$c.Prop}}"]{{if is_last $i $columns}};{{else}},{{end}}
        {{end}}{{end}}

    // 根据字段拷贝数据
    {{$title}}{{.global.entity_suffix}}.fromRawMap(Map<String, dynamic> src)
       :{{range $i,$c := $columns}}\
        {{$c.Prop}} = src["{{$c.Prop}}"]{{if is_last $i $columns}};{{else}},{{end}}
        {{end}}

    // 根据字段拷贝数据
    Map<String,dynamic> toMap(){
        return {
            {{range $i,$c := $columns}}
            "{{$c.Prop}}":this.{{$c.Prop}},{{end}}
        };
    }
}