/** {{.table.Comment}} */
struct S{{.table.Title}}{
    {{range $i,$c:=.T.Columns}}
    /** {{$c.Comment}} */
    {{plus $c.Ordinal 1}}:{{$c.TypeId}} {{$c.Title}}{{end}}
}