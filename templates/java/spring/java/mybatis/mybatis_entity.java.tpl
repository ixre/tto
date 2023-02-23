#!target:spring/src/main/java/{{.global.pkg}}/entity/{{.table.Title}}{{.global.entity_suffix}}.java
package {{pkg "java" .global.pkg}}.entity;


import javax.persistence.Entity;
import javax.persistence.Id;
import javax.persistence.Table;
import com.baomidou.mybatisplus.annotation.TableName;
import com.baomidou.mybatisplus.annotation.TableField;
import com.baomidou.mybatisplus.annotation.TableId;

{{$entity := join .table.Title .global.entity_suffix}}
/** {{.table.Comment}}(MyBatis) */
{{/*　@DynamicInsert 排除值为null的字段,@Entity,@Table,@Id为兼容Hibernate　*/}} \
@Entity
@Table(name = "{{.table.Name}}", schema = "{{.table.Schema}}")
@TableName("{{.table.Name}}")
public class {{$entity}} {
    {{/* 将字段单独生成，以便做裁剪 */}}\
    {{range $i,$c := .columns}}{{$type := orm_type "java" $c.Type}}
    {{$lowerProp := lower_title $c.Prop}} \
    {{if $c.IsPk}}
    @Id{{/* 兼容Hibernate*/}}
    @TableId("{{$c.Name}}") \
    {{else}}
    @TableField("{{$c.Name}}") \
    {{end}}
    private {{$type}} {{$lowerProp}}; // {{$c.Comment}}\
    {{end}}
    
    {{range $i,$c := .columns}}{{$type := type "java" $c.Type}}{{$ormType := orm_type "java" $c.Type}}
    {{$lowerProp := lower_title $c.Prop}} \
    public {{$entity}} set{{$c.Prop}}({{$type}} {{$lowerProp}}){
        this.{{$lowerProp}} = {{$lowerProp}};
        return this;
    }

    /** {{$c.Comment}} */
    public {{$ormType}} get{{$c.Prop}}() {
        return this.{{$lowerProp}};
    }
    {{end}}


    /*
    public {{$entity}} deep(){
        {{$entity}} dst = new {{$entity}}();
        {{range $i,$c := .columns}}
        dst.set{{$c.Prop}}(this.get{{$c.Prop}}());{{end}}
        return dst;
    }
    */

    /*
    public Map<String,Object> toMap(){
        Map<String,Object> mp = new HashMap<>();\
        {{range $i,$c := .columns}}
        mp.put("{{lower_title $c.Prop}}",this.{{lower_title $c.Prop}});{{end}}
        return mp;
    }
    */

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
