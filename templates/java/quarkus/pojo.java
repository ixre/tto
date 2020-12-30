#!target:src/main/java/{{.global.pkg}}/pojo/{{.table.Title}}{{.global.entity_suffix}}.java
package {{pkg "java" .global.pkg}}.pojo;

import javax.persistence.Basic;
import javax.persistence.Id;
import javax.persistence.Column;
import javax.persistence.Entity;
import javax.persistence.Table;
import javax.persistence.GenerationType;
import javax.persistence.GeneratedValue;

/** {{.table.Comment}} */
@Entity
@Table(name = "{{.table.Name}}", schema = "{{.table.Schema}}")
public class {{.table.Title}}{{.global.entity_suffix}} {
    {{range $i,$c := .columns}}{{$type := type "java" $c.Type}}

    {{if $c.IsPk}}\
    @Id
    @GeneratedValue(strategy = GenerationType.IDENTITY){{else}}
    @Basic{{end}}
    @Column(name = "{{$c.Name}}"{{if not $c.NotNull}}, nullable = true{{end}} {{if ne $c.Length 0}},length = {{$c.Length}}{{end}})
    private {{$type}} {{$c.Name}};

    /** {{$c.Comment}} */
    public {{$type}} get{{$c.Prop}}() {
        return this.{{$c.Name}};
    }

    public void set{{$c.Prop}}({{$type}} {{$c.Name}}){
        this.{{$c.Name}} = {{$c.Name}};
    }

    {{end}}

    /** 拷贝数据  */
    public {{.table.Title}}{{.global.entity_suffix}} copy({{.table.Title}}{{.global.entity_suffix}} src){
        {{.table.Title}}{{.global.entity_suffix}} dst = new {{.table.Title}}{{.global.entity_suffix}}();
        {{range $i,$c := .columns}}
        dst.set{{$c.Prop}}(src.get{{$c.Prop}}());{{end}}
        return dst;
    }
}
