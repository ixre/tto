#!kind:1#!target:vue2/lang.ts
export const TablesLang = {  {{$tables := .tables}} {{range $i,$table := .tables}}
    "{{$table.Comment}}": "{{$table.Comment}}"{{if not (is_last $i $tables)}},{{end}} \
{{end}}
};