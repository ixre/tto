#!kind:1#!target:vue/api/index.ts
{{range $i,$table := .tables}}\
export * from "./{{$table.Title}}Api";
{{end}}