#!target:quarkus/src/main/kotlin/{{.global.pkg}}/repo/{{.table.Title}}JpaRepository.kt
package {{pkg "kotlin" .global.pkg}}.repo;

import {{pkg "kotlin" .global.pkg}}.entity.{{.table.Title}}{{.global.entity_suffix}}
import io.quarkus.hibernate.orm.panache.PanacheRepository
import javax.enterprise.context.ApplicationScoped

{{$pkType := type "kotlin" .table.PkType}}
/** {{.table.Comment}}仓储 */
@ApplicationScoped
class {{.table.Title}}JpaRepository : PanacheRepository<{{.table.Title}}{{.global.entity_suffix}}> {
}