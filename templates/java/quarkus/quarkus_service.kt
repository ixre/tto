#!target:src/main/kotlin/{{.global.pkg}}/service/{{.table.Title}}Service.kt
package {{pkg "java" .global.pkg}}.service

import {{pkg "java" .global.pkg}}.pojo.{{.table.Title}}Entity
import {{pkg "java" .global.pkg}}.repo.{{.table.Title}}JpaRepository
import javax.inject.Inject
import javax.enterprise.inject.Default
import javax.enterprise.context.ApplicationScoped
import net.fze.util.catch
import net.fze.commons.Types
import net.fze.commons.TypesConv
import net.fze.util.value
import javax.transaction.Transactional

{{$tableTitle := .table.Title}}
{{$pkName := .table.Pk}}
{{$pkProp := lower_title .table.PkProp}}
{{$pkType := type "kotlin" .table.PkType}}
/** {{.table.Comment}}服务  */
@ApplicationScoped
class {{.table.Title}}Service {
    @Inject@field:Default
    private lateinit var repo: {{$tableTitle}}JpaRepository

    fun parseId(id:Any):Long{return TypesConv.toLong(id)}

    /** 根据ID查找{{.table.Comment}} */
    fun findByIdOrNull(id:{{$pkType}}):{{$tableTitle}}Entity?{
        return this.repo.findByIdOptional(this.parseId(id)).value()
    }

    /** 根据条件查找单个对象 */
    fun findBy(query:String,vararg params:Any):{{$tableTitle}}Entity?{
        return this.repo.find(query,params).singleResultOptional<{{$tableTitle}}Entity>().value()
    }

    /** 根据条件查找并返回列表 */
    fun listBy(query:String,params:List<Any>):List<{{$tableTitle}}Entity>{
        return this.repo.list(query,params)
    }

    /** 保存{{.table.Comment}} */
    @Transactional
    fun save{{$tableTitle}}(e: {{$tableTitle}}Entity):Error? {
        return catch {
            var dst: {{$tableTitle}}Entity
            if (e.{{$pkProp}} > 0) {
                dst = this.repo.findById(this.parseId(e.{{$pkProp}}))!!
            } else {
                dst = {{$tableTitle}}Entity()
                {{$c := try_get .columns "create_time"}}\
                {{if $c }}dst.createTime = Types.time.unix().toLong(){{end}}
            }
            {{range $i,$c := exclude .columns $pkName "create_time" "update_time"}}
            dst.{{lower_title $c.Prop}} = e.{{lower_title $c.Prop}}{{end}}\
            {{$c := try_get .columns "update_time"}}
            {{if $c}}dst.updateTime = Types.time.unix().toLong(){{end}}
            this.repo.persistAndFlush(dst)
            null
        }.error()
    }

    /** 批量保存{{.table.Comment}} */
    @Transactional
    fun saveAll{{$tableTitle}}(entities:Iterable<{{$tableTitle}}Entity>){
        this.repo.persist(entities)
        this.repo.flush()
    }

    /** 删除{{.table.Comment}} */
    @Transactional
    fun deleteById(id:{{$pkType}}):Error? {
        return catch {
            this.repo.deleteById(this.parseId(id))
        }.error()
    }

}
