#!target:spring/src/main/java/{{.global.pkg}}/entity/{{.table.Title}}{{.global.entity_suffix}}.java
package {{pkg "java" .global.pkg}}.entity;


import javax.persistence.Entity;
import javax.persistence.Id;
import javax.persistence.Table;
import com.baomidou.mybatisplus.annotation.IdType;
import com.baomidou.mybatisplus.annotation.TableName;
import com.baomidou.mybatisplus.annotation.TableField;
import com.baomidou.mybatisplus.annotation.TableId;
import java.math.BigDecimal;
import java.util.Date;
import lombok.Data;


{{$entity := join .table.Title .global.entity_suffix}}
/**
 * {{.table.Comment}}(MyBatis)
 */
{{/*　@DynamicInsert 排除值为null的字段,@Entity,@Table,@Id为兼容Hibernate　*/}} \
@Entity
@Table(name = "{{.table.Name}}", schema = "{{.table.Schema}}")
@TableName("{{.table.Name}}")
@Data
public class {{$entity}} implements Cloneable {
    {{/* 将字段单独生成，以便做裁剪,未使用Lombok是因为系统set属性能使用构造者模式 */}}\
    {{range $i,$c := .columns}}{{$type := orm_type "java" $c.Type}}
    {{$lowerProp := lower_title $c.Prop}} 
    /**
     * {{$c.Comment}}
     */\
    {{if $c.IsPk}}
    @Id{{/* 兼容Hibernate*/}}
    @TableId(value="{{$c.Name}}"{{if $c.IsAuto}}, type = IdType.AUTO{{end}}) \
    {{else}}
    @TableField("{{$c.Name}}") \
    {{end}}
    private {{$type}} {{$lowerProp}}; \
    {{end}}


    @Override
    public {{$entity}} clone() {
        try {
            return ({{$entity}}) super.clone();
        } catch (Exception ex) {
            throw new RuntimeException("Entity {{$entity}} clone failed:" + ex.getMessage());
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
