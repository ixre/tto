#!target:spring/src/main/java/{{.global.pkg}}/mapper/{{.table.Prefix}}/{{.table.Title}}Mapper.java
package {{pkg "java" .global.pkg}}.mapper.{{.table.Prefix}};

import {{pkg "java" .global.pkg}}.entity.{{.table.Title}}{{.global.entity_suffix}};
import com.baomidou.mybatisplus.core.mapper.BaseMapper;

/** {{.table.Comment}}仓储接口 */
public interface {{.table.Title}}Mapper extends BaseMapper<{{.table.Title}}{{.global.entity_suffix}}>{
}
