namespace java {{pkg "thrift" .global.pkg}}.rpc
namespace netstd {{pkg "thrift" .global.pkg}}.rpc

{{.global.base_path}}

/** {{.table.Comment}} */
struct S{{.table.Title}}{
    {{range $i,$c:=.T.Columns}}
    /** {{$c.Comment}} */
    {{plus $c.Ordinal 1}}:{{type "thrift" $c.Type}} {{$c.Prop}}{{end}}
}