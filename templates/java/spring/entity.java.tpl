#!target:spring/src/main/java/{{.global.pkg}}/entity/{{.table.Title}}{{.global.entity_suffix}}.java
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
{{$entity := join .table.Title .global.entity_suffix}}
/** {{.table.Comment}} */
@Entity
{{/*　@DynamicInsert 排除值为null的字段　*/}} \
@Table(name = "{{.table.Name}}", schema = "{{.table.Schema}}")
public class {{$entity}} {
    {{range $i,$c := .columns}}{{$type := type "java" $c.Type}}
    {{$lowerProp := lower_title $c.Prop}} \
    private {{$type}} {{$lowerProp}};
    public {{$entity}} set{{$c.Prop}}({{$type}} {{$lowerProp}}){
        this.{{$lowerProp}} = {{$lowerProp}};
        return this;
    }

    /** {{$c.Comment}} */{{if $c.IsPk}}
    @Id
    @GeneratedValue(strategy = GenerationType.IDENTITY){{else}}
    @Basic{{end}}
    @Column(name = "{{$c.Name}}"{{if not $c.NotNull}}, nullable = true{{end}} {{if ne $c.Length 0}},length = {{$c.Length}}{{end}})
    public {{$type}} get{{$c.Prop}}() {
        return this.{{$lowerProp}};
    }
    {{end}}


     /** 创建深拷贝  */
    /*
    public {{$entity}} copy(){
        {{$entity}} dst = new {{$entity}}();
        {{range $i,$c := .columns}}
        dst.set{{$c.Prop}}(this.get{{$c.Prop}}());{{end}}
        return dst;
    }
    */

    /** 转换为MAP  */
    /*
    public Map<String,Object> toMap(){
        Map<String,Object> mp = new HashMap<>();\
        {{range $i,$c := .columns}}
        mp.put("{{$c.Name}}",this.{{lower_title $c.Prop}});{{end}}
        return mp;
    }
    */

    /** 从MAP转换 */
    /*
    public static {{$entity}} fromMap(Map<String,Object> data){
        {{$entity}} dst = new {{$entity}}();\
        {{range $i,$c := .columns}}
        {{ $goType := type "java" $c.Type}}\
        {{if eq $goType "int"}}dst.set{{$c.Prop}}(TypeConv.toInt(data.get("{{$c.Prop}}")));\
        {{else if eq $goType "long"}}dst.set{{$c.Prop}}(TypeConv.toLong(data.get("{{$c.Prop}}")));\
        {{else if eq $goType "boolean"}}dst.set{{$c.Prop}}(TypeConv.toBool(data.get("{{$c.Prop}}")));\
        {{else if eq $goType "float"}}dst.set{{$c.Prop}}(TypeConv.toFloat(data.get("{{$c.Prop}}")));\
        {{else if eq $goType "double"}}dst.set{{$c.Prop}}(TypeConv.toDouble(data.get("{{$c.Prop}}")));\
        {{else if eq $goType "BigDecimal"}}dst.set{{$c.Prop}}(TypeConv.toBigDecimal(data.get("{{$c.Prop}}")));\
        {{else if eq $goType "Date"}}dst.set{{$c.Prop}}(TypeConv.toDateTime(data.get("{{$c.Prop}}")));\
        {{else if eq $goType "Byte[]"}}dst.set{{$c.Prop}}(TypeConv.toBytes(data.get("{{$c.Prop}}")));\
        {{else}}dst.set{{$c.Prop}}(TypeConv.toString(data.get("{{$c.Prop}}")));{{end}}{{end}}
        return dst;
    }
    */

    {{/* 通过字段直接给默认值会影响Example.of, 所以通过方法来设置默认值 */}}
    public static {{$entity}} createDefault(){
        {{$entity}} dst = new {{$entity}}();\
        {{range $i,$c := .columns}}
        dst.set{{$c.Prop}}({{default "java" $c.Type}});{{end}}
        return dst;
    }
}
