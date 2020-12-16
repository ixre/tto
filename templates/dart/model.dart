#!target:dart/lib/entity/{{.table.Name}}_model.dart
{{$title := .table.Title}}

class {{$title}}Model {
    {{range $i,$c := .columns}}
    {{$type := type "dart" $c.Type}}\
    // {{$c.Comment}}
    {{$type}} {{$c.Prop}};{{end}}


    {{$title}}Model();

    // 拷贝数据
    {{$title}}Model.fromMap(Map<String, dynamic> src){
        {{range $i,$c := .columns}}
        {{$c.Prop}} = src["{{$c.Prop}}"];{{end}}
    }
    // 根据字段拷贝数据
    {{$title}}Model.fromRawMap(Map<String, dynamic> src){
        {{range $i,$c := .columns}}
        {{$c.Prop}} = src["{{$c.Name}}"];{{end}}
    }
}