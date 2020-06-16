#!target:kotlin/{{.global.pkg}}/service/{{.table.Title}}Service.kt.gen
package {{pkg "java" .global.pkg}}.service

import {{pkg "java" .global.pkg}}.pojo.{{.table.Title}}Entity
import {{pkg "java" .global.pkg}}.repo.{{.table.Title}}JpaRepository
import javax.inject.Inject
import javax.enterprise.inject.Default
import javax.enterprise.context.ApplicationScoped
import net.fze.arch.commons.std.catch
import net.fze.arch.commons.std.TypesConv
import javax.transaction.Transactional

{{$tableTitle := .table.Title}}
{{$pkName := .table.Pk}}
{{$pkProp := lower_title .table.PkProp}}
{{$pkType := type "kotlin" .table.PkType}}
/** {{.table.Comment}}服务  */
@ApplicationScoped
class {{.table.Title}}Service {
    @Inject@field:Default
    lateinit var repo: {{$tableTitle}}JpaRepository

    fun parseId(id:Any):Long{return TypesConv.toLong(id)}

    // 根据ID查找{{.table.Comment}}
    fun findByIdOrNull(id:{{$pkType}}):{{$tableTitle}}Entity?{
        return this.repo.findById(this.parseId(id))
    }

    // 保存{{.table.Comment}}
    @Transactional
    fun save{{$tableTitle}}(e: {{$tableTitle}}Entity):Error? {
        return catch {
            var dst: {{$tableTitle}}Entity
            if (e.{{$pkProp}} > 0) {
                dst = this.repo.findById(this.parseId(e.{{$pkProp}}))!!
            } else {
                dst = {{$tableTitle}}Entity()
            }
            {{range $i,$c := .columns}}{{if ne $c.Name $pkName}}
            dst.{{lower_title $c.Prop}} = e.{{lower_title $c.Prop}}{{end}}{{end}}
            this.repo.persist(dst)
            null
        }.error()
    }

    // 批量保存{{.table.Comment}}
    @Transactional
    fun saveAll{{$tableTitle}}(entities:Iterable<{{$tableTitle}}Entity>){
        return this.repo.persist(entities)
    }

    // 删除{{.table.Comment}}
    @Transactional
    fun deleteById(id:{{$pkType}}):Error? {
        return catch {
            this.repo.deleteById(this.parseId(id))
        }.error()
    }

}
