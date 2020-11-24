#!target:kotlin/{{.global.pkg}}/model/{{.table.Title}}Entity.kt
package {{pkg "kotlin" .global.pkg}}.model;

/** {{.table.Comment}} */
class {{.table.Title}}{
    {{range $i,$c:=.columns}}
    /** {{$c.Comment}} */
    var {{lower_title $c.Prop}}:{{type "kotlin" $c.Type}} = {{default "kotlin" $c.Type}} {{end}}

    /** 拷贝数据  */
    fun  copy(src :{{.table.Title}}Entity):{{.table.Title}}Entity{
        val dst = this;
        {{range $i,$c := .columns}}
        dst.{{lower_title $c.Prop}} = src.{{lower_title $c.Prop}}{{end}}
        return dst;
    }
}