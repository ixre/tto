#!kind:1#!target:ts/feature/lang.ts
export const TablesLang = {  {{$tables := .tables}} {{range $i,$table := .tables}}
    "{{$table.Comment}}": "{{$table.Comment}}"{{if not (is_last $i $tables)}},{{end}} \
{{end}}
};