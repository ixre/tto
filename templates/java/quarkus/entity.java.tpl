#!target:src/main/java/{{.global.pkg}}/entity/{{.table.Title}}{{.global.entity_suffix}}.java
package {{pkg "java" .global.pkg}}.entity;

import net.fze.util.TypeConv;
import javax.persistence.Basic;
import javax.persistence.Id;
import javax.persistence.Column;
import javax.persistence.Entity;
import javax.persistence.Table;
import javax.persistence.GenerationType;
import javax.persistence.GeneratedValue;
import java.math.BigDecimal;
import java.util.Date;
import java.util.HashMap;
import java.util.Map;

/** {{.table.Comment}} */
@Entity
@Table(name = "{{.table.Name}}", schema = "{{.table.Schema}}")
public class {{.table.Title}}{{.global.entity_suffix}} {
    {{range $i,$c := .columns}}{{$type := type "java" $c.Type}}
    private {{$type}} {{$c.Name}};
    public void set{{$c.Prop}}({{$type}} {{$c.Name}}){
        this.{{$c.Name}} = {{$c.Name}};
    }

    /** {{$c.Comment}} */{{if $c.IsPk}}
    @Id
    @GeneratedValue(strategy = GenerationType.IDENTITY){{else}}
    @Basic{{end}}
    @Column(name = "{{$c.Name}}"{{if not $c.NotNull}}, nullable = true{{end}} {{if ne $c.Length 0}},length = {{$c.Length}}{{end}})
    public {{$type}} get{{$c.Prop}}() {
        return this.{{$c.Name}};
    }
    {{end}}

    /** 拷贝数据  */
    public {{.table.Title}}{{.global.entity_suffix}} copy({{.table.Title}}{{.global.entity_suffix}} src){
        {{.table.Title}}{{.global.entity_suffix}} dst = new {{.table.Title}}{{.global.entity_suffix}}();
        {{range $i,$c := .columns}}
        dst.set{{$c.Prop}}(src.get{{$c.Prop}}());{{end}}
        return dst;
    }
    /** 拷贝数据  */
    public Map<String,Object> toMap(){
        Map<String,Object> mp = new HashMap<>();\
        {{range $i,$c := .columns}}
        mp.put("{{$c.Prop}}",this.{{$c.Name}});{{end}}
        return mp;
    }

    public static {{.table.Title}}{{.global.entity_suffix}} fromMap(Map<String,Object> data){
        {{.table.Title}}{{.global.entity_suffix}} dst = new {{.table.Title}}{{.global.entity_suffix}}();\
        {{range $i,$c := .columns}}
        {{ $goType := type "go" $c.Type}}\
        {{if eq $goType "int"}}dst.set{{$c.Prop}}(TypeConv.toInt(data.get("{{$c.Prop}}")));\
        {{else if eq $goType "int16"}}dst.set{{$c.Prop}}(TypeConv.toInt(data.get("{{$c.Prop}}")));\
        {{else if eq $goType "int64"}}dst.set{{$c.Prop}}(TypeConv.toLong(data.get("{{$c.Prop}}")));\
        {{else if eq $goType "bool"}}dst.set{{$c.Prop}}(TypeConv.toBool(data.get("{{$c.Prop}}")));\
        {{else if eq $goType "float32"}}dst.set{{$c.Prop}}(TypeConv.toFloat(data.get("{{$c.Prop}}")));\
        {{else if eq $goType "float64"}}dst.set{{$c.Prop}}(TypeConv.toBigDecimal(data.get("{{$c.Prop}}")));\
        {{else if eq $goType "time.Time"}}dst.set{{$c.Prop}}(TypeConv.toDateTime(data.get("{{$c.Prop}}")));\
        {{else}}dst.set{{$c.Prop}}(TypeConv.toString(data.get("{{$c.Name}}")));{{end}}{{end}}
        return dst;
    }

    public static {{.table.Title}}{{.global.entity_suffix}} createDefault(){
        {{.table.Title}}{{.global.entity_suffix}} dst = new {{.table.Title}}{{.global.entity_suffix}}();\
        {{range $i,$c := .columns}}
        dst.set{{$c.Prop}}({{default "java" $c.Type}});{{end}}
        return dst;
    }
}
