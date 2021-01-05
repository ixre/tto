#!target:spring/src/main/kotlin/{{.global.pkg}}/service/{{.table.Title}}Service.kt
package {{pkg "java" .global.pkg}}.service

import {{pkg "java" .global.pkg}}.entity.{{.table.Title}}{{.global.entity_suffix}}
import {{pkg "java" .global.pkg}}.repo.{{.table.Title}}JpaRepository
import org.springframework.stereotype.Service
import org.springframework.data.repository.findByIdOrNull
import net.fze.common.catch
import javax.annotation.Resource
{{$tableTitle := .table.Title}}
{{$shortTitle := .table.ShortTitle}}
{{$pkProp :=  .table.PkProp}}
{{$pkType := type "kotlin" .table.PkType}}
/** {{.table.Comment}}服务  */
@Service
class {{.table.Title}}Service {
    @Resource
    lateinit var repo: {{$tableTitle}}JpaRepository

    // 根据ID查找{{.table.Comment}}
    fun find{{$shortTitle}}ById(id:{{$pkType}}):{{$tableTitle}}{{.global.entity_suffix}}?{
        return this.repo.findByIdOrNull(id)
    }

    // 保存{{.table.Comment}}
    fun save{{$shortTitle}}(e: {{$tableTitle}}{{.global.entity_suffix}}):Error? {
        return catch {
            val dst: {{$tableTitle}}{{.global.entity_suffix}}
            {{if equal_any .table.PkType 3 4 5}}\
            if (e.{{$pkProp}} > 0) {
            {{else}}
            if (e.{{$pkProp}} != "") {
            {{end}}
                dst = this.repo.findByIdOrNull(e.{{$pkProp}})!!
            } else {
                dst = {{$tableTitle}}{{.global.entity_suffix}}.createDefault()
            }
            {{range $i,$c := .columns}}
            dst.{{lower_title $c.Prop}} = e.{{lower_title $c.Prop}}{{end}}
            this.repo.save(dst)
            null
        }.error()
    }

    // 批量保存{{.table.Comment}}
    fun saveAll{{$shortTitle}}(entities:Iterable<{{$tableTitle}}{{.global.entity_suffix}}>): Iterable<{{$tableTitle}}{{.global.entity_suffix}}>{
        return this.repo.saveAll(entities)
    }

    // 删除{{.table.Comment}}
    fun delete{{$shortTitle}}ById(id:{{$pkType}}):Error? {
        return catch {
            this.repo.deleteById(id)
            null
        }.error()
    }
}
