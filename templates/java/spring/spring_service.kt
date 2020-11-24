#!target:spring/src/main/kotlin/{{.global.pkg}}/service/{{.table.Title}}Service.kt
package {{pkg "java" .global.pkg}}.service

import {{pkg "java" .global.pkg}}.pojo.{{.table.Title}}Entity
import {{pkg "java" .global.pkg}}.repo.{{.table.Title}}JpaRepository
import org.springframework.stereotype.Service
import org.springframework.data.repository.findByIdOrNull
import net.fze.util.catch
import javax.annotation.Resource
{{$tableTitle := .table.Title}}
{{$pkProp := lower_title .table.Pk}}
{{$pkType := type "kotlin" .table.PkType}}
/** {{.table.Comment}}服务  */
@Service
class {{.table.Title}}Service {

    @Resource
    lateinit var repo: {{$tableTitle}}JpaRepository

    // 根据ID查找{{.table.Comment}}
    fun findByIdOrNull(id:{{$pkType}}):{{$tableTitle}}Entity?{
        return this.repo.findByIdOrNull(id)
    }

    // 保存{{.table.Comment}}
    fun save{{$tableTitle}}(e: {{$tableTitle}}Entity):Error? {
        return catch {
            var dst: {{$tableTitle}}Entity
            if (e.{{$pkProp}} > 0) {
                dst = this.repo.findByIdOrNull(e.{{$pkProp}})!!
            } else {
                dst = {{$tableTitle}}Entity()
            }
            {{range $i,$c := .columns}}
            dst.{{lower_title $c.Prop}} = e.{{lower_title $c.Prop}}{{end}}
            this.repo.save(dst)
            null
        }.error()
    }

    // 批量保存{{.table.Comment}}
    fun saveAll{{$tableTitle}}(entities:Iterable<{{$tableTitle}}Entity>): Iterable<{{$tableTitle}}Entity>{
        return this.repo.saveAll(entities)
    }

    // 删除{{.table.Comment}}
    fun deleteById(id:{{$pkType}}) {
         this.repo.deleteById(id)
    }

}
