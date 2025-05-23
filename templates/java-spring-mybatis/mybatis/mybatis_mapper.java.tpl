#!target:spring/src/main/java/{{.global.pkg}}/mapper/{{.table.Title}}Mapper.java
package {{pkg "java" .global.pkg}}.mapper;
{{/*使用Mapper是为了防止MyBatis未添加扫描包目录，造成异常现象*/}}
import {{pkg "java" .global.pkg}}.entity.{{.table.Title}}{{.global.entity_suffix}};
import net.fze.ext.mybatis.BaseJpaMapper;
{{$pkType := orm_type "java" .table.PkType}}

/** {{.table.Comment}}仓储接口 */
public interface {{.table.Title}}Mapper extends BaseJpaMapper<{{.table.Title}}{{.global.entity_suffix}}>{
}
