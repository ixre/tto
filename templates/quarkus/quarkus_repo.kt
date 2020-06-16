#!target:kotlin/{{.global.pkg}}/repo/{{.table.Title}}JpaRepository.kt.gen
package {{pkg "java" .global.pkg}}.repo;

import {{pkg "kotlin" .global.pkg}}.pojo.{{.table.Title}}Entity
import io.quarkus.hibernate.orm.panache.PanacheRepository
import javax.enterprise.context.ApplicationScoped

{{$pkType := type "kotlin" .table.PkType}}
/** {{.table.Comment}}仓储接口  */
@ApplicationScoped
class {{.table.Title}}JpaRepository : PanacheRepository<{{.table.Title}}Entity> {

}