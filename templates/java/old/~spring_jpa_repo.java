#!target:java/{{.global.pkg}}/repo/{{.table.Title}}Repository.java
package {{pkg "java" .global.pkg}}.repo;

import {{pkg "java" .global.pkg}}.pojo.{{.table.Title}}{{.global.entity_suffix}}
import org.springframework.data.jpa.repository.JpaRepository
{{$pkType := type "java" .table.PkType}}
/** {{.table.Comment}}仓储接口 */
public interface {{.table.Title}}Repository : JpaRepository<{{.table.Title}}{{.global.entity_suffix}}, {{$pkType}}>{

}
