#!target:spring/src/main/kotlin/{{.global.pkg}}/repo/{{.table.Title}}JpaRepository.kt
package {{pkg "java" .global.pkg}}.repo;

import {{pkg "kotlin" .global.pkg}}.pojo.{{.table.Title}}Entity
import org.springframework.data.jpa.repository.JpaRepository
{{$pkType := type "kotlin" .table.PkType}}
/** {{.table.Comment}}仓储接口  */
interface {{.table.Title}}JpaRepository : JpaRepository<{{.table.Title}}Entity, {{$pkType}}> {

}