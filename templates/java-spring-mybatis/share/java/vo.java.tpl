#!target:spring/src/main/java/{{.global.pkg}}/vo/{{.table.Title}}VO.java
package {{pkg "java" .global.pkg}}.entity;

import net.fze.util.TypeConv;
import java.math.BigDecimal;
import java.util.Date;
import lombok.Data;

{{$entity := join .table.Title}}
/** {{.table.Comment}} */
{{/*　@DynamicInsert 排除值为null的字段　*/}} \
@Data
public class {{$entity}}VO {
    {{/* 将字段单独生成，以便做裁剪 */}}\
    {{range $i,$c := .columns}}{{$type := type "java" $c.Type}}
    {{$lowerProp := lower_title $c.Prop}} \
    /** {{$c.Comment}} */
    private {{$type}} {{$lowerProp}} = {{default "java" $c.Type}}; \
    {{end}}
}
