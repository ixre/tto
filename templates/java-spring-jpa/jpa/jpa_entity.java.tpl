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
import lombok.Data;

{{$entity := join .table.Title .global.entity_suffix}}
/**
 * {{.table.Comment}} 
 */
@Entity
{{/*　@DynamicInsert 排除值为null的字段　*/}} \
@Table(name = "{{.table.Name}}", schema = "{{.table.Schema}}")
@Data
public class {{$entity}} implements Cloneable {
    {{/* 将字段单独生成，以便做裁剪 */}}\
    {{range $i,$c := .columns}}{{$type := orm_type "java" $c.Type}}
    {{$lowerProp := lower_title $c.Prop}} 
    /**
     * {{$c.Comment}}
     */\{{if $c.IsPk}}
    @Id
    @GeneratedValue(strategy = GenerationType.IDENTITY){{else}}
    @Basic{{end}}
    @Column(name = "{{$c.Name}}"{{if not $c.NotNull}}, nullable = true{{end}} {{if ne $c.Length 0}},length = {{$c.Length}}{{end}})
    private {{$type}} {{$lowerProp}}; \
    {{end}}
    
    @Override
    public {{$entity}} clone() {
        try {
            return ({{$entity}}) super.clone();
        } catch (Exception ex) {
            throw new RuntimeException("clone failed:" + ex.getMessage());
        }
    }

    {{/* 通过字段直接给默认值会影响Example.of, 所以通过方法来设置默认值 */}}
    public static {{$entity}} createDefault(){
        {{$entity}} dst = new {{$entity}}();\
        {{range $i,$c := .columns}}
        dst.set{{$c.Prop}}({{default "java" $c.Type}});{{end}}
        return dst;
    }
}
