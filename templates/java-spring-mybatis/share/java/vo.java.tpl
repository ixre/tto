#!target:spring/src/main/java/{{.global.pkg}}/vo/{{.table.Title}}VO.java
package {{pkg "java" .global.pkg}}.entity;

{{$swagger := .global.swagger}} \
import java.math.BigDecimal;
import java.util.Date;
import lombok.Data;
{{if $swagger}} \
import io.swagger.v3.oas.annotations.media.Schema;
{{end}}

{{$entity := join .table.Title}}
/** {{.table.Comment}} */
{{/*　@DynamicInsert 排除值为null的字段　*/}} \
@Data
{{if $swagger}} \
@Schema(name = "{{$entity}}VO",description = "{{.table.Comment}}")
{{end}} \
public class {{$entity}}VO {
    {{/* 将字段单独生成，以便做裁剪 */}}\
    {{range $i,$c := .columns}}{{$type := type "java" $c.Type}}
    {{$lowerProp := lower_title $c.Prop}} \
    /** {{$c.Comment}} */
    {{if $swagger}} \
    @Schema(description = "{{$c.Comment}}")
    {{end}} \
    private {{$type}} {{$lowerProp}} = {{default "java" $c.Type}}; \
    {{end}}
}
