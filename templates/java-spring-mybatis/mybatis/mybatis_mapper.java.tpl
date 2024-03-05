#!target:spring/src/main/java/{{.global.pkg}}/mapper/{{.table.Prefix}}/{{.table.Title}}Mapper.java
package {{pkg "java" .global.pkg}}.mapper.{{.table.Prefix}};

import {{pkg "java" .global.pkg}}.entity.{{.table.Title}}{{.global.entity_suffix}};

import net.fze.ext.mybatis.BaseJpaMapper;

{{$pkType := orm_type "java" .table.PkType}}

/** {{.table.Comment}}仓储接口 */
public interface {{.table.Title}}Mapper extends BaseJpaMapper<{{$pkType}},{{.table.Title}}{{.global.entity_suffix}}>{
}
