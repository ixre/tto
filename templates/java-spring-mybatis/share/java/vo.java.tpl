#!target:spring/src/main/java/{{.global.pkg}}/vo/{{.table.Title}}.java
package {{pkg "java" .global.pkg}}.entity;

import net.fze.util.TypeConv;
import java.math.BigDecimal;
import java.util.Date;
{{$entity := join .table.Title}}
/** {{.table.Comment}} */
{{/*　@DynamicInsert 排除值为null的字段　*/}} \
public class {{$entity}} {
    {{/* 将字段单独生成，以便做裁剪 */}}\
    {{range $i,$c := .columns}}{{$type := type "java" $c.Type}}
    {{$lowerProp := lower_title $c.Prop}} \
    private {{$type}} {{$lowerProp}}; // {{$c.Comment}}\
    {{end}}
    
    {{range $i,$c := .columns}}{{$type := type "java" $c.Type}}
    {{$lowerProp := lower_title $c.Prop}} \
    public {{$entity}} set{{$c.Prop}}({{$type}} {{$lowerProp}}){
        this.{{$lowerProp}} = {{$lowerProp}};
        return this;
    }

    /** {{$c.Comment}} */
    public {{$type}} get{{$c.Prop}}() {
        return this.{{$lowerProp}};
    }
    {{end}}

    {{/* 通过字段直接给默认值会影响Example.of, 所以通过方法来设置默认值 */}}
    public static {{$entity}} createDefault(){
        {{$entity}} dst = new {{$entity}}();\
        {{range $i,$c := .columns}}
        dst.set{{$c.Prop}}({{default "java" $c.Type}});{{end}}
        return dst;
    }
}
