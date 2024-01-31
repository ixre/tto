#!target:quarkus/src/main/kotlin/{{.global.pkg}}/service/{{.table.Title}}Service.kt
package {{pkg "kotlin" .global.pkg}}.service

import {{pkg "kotlin" .global.pkg}}.entity.{{.table.Title}}{{.global.entity_suffix}}
import {{pkg "kotlin" .global.pkg}}.repo.{{.table.Title}}JpaRepository
import javax.inject.Inject
import javax.enterprise.inject.Default
import javax.enterprise.context.ApplicationScoped
import net.fze.common.Types
import net.fze.util.Times
import net.fze.common.TypeConv
import javax.transaction.Transactional

{{$tableTitle := .table.Title}}
{{$shortTitle := .table.ShortTitle}}
{{$pkName := .table.Pk}}
{{$pkProp := lower_title .table.PkProp}}
{{$pkType := type "kotlin" .table.PkType}}
/** {{.table.Comment}}服务  */
@ApplicationScoped
class {{.table.Title}}Service {
    @Inject@field:Default
    private lateinit var repo: {{$tableTitle}}JpaRepository

    fun parse{{$shortTitle}}Id(id:Any):{{$pkType}}{return TypeConv.toLong(id)}

    /** 根据ID查找{{.table.Comment}} */
    fun find{{$shortTitle}}ByIdOrNull(id:{{$pkType}}):{{$tableTitle}}{{.global.entity_suffix}}?{
        return this.repo.findByIdOptional(this.parseId(id))?.get()
    }

    /** 根据条件查找单个对象 */
    fun find{{$shortTitle}}By(query:String,vararg params:Any):{{$tableTitle}}{{.global.entity_suffix}}?{
        val find = this.repo.find(query, *params)
        if (!find.singleResultOptional<{{$tableTitle}}{{.global.entity_suffix}}>().isPresent) {
            return null
        }
        return find.singleResultOptional<{{$tableTitle}}{{.global.entity_suffix}}>()?.get()
    }

    /** 根据条件查找并返回列表 */
    fun list{{$shortTitle}}By(query:String,vararg params:List<Any>):List<{{$tableTitle}}{{.global.entity_suffix}}>{
        return this.repo.list(query,*params)
    }

    /** 保存{{.table.Comment}} */
    @Transactional
    fun save{{$shortTitle}}(e: {{$tableTitle}}{{.global.entity_suffix}}):Error? {
        return catch {
            var dst: {{$tableTitle}}{{.global.entity_suffix}}
            {{if equal_any .table.PkType 3 4 5}}\
            if (e.{{$pkProp}} > 0) {
            {{else}}
            if (e.{{$pkProp}} != "") {
            {{end}}
                dst = this.repo.findById(this.parseId(e.{{$pkProp}}))!!
            } else {
                dst = {{$tableTitle}}{{.global.entity_suffix}}.createDefault()
                {{$c := try_get .columns "create_time"}}\
                {{if $c }}dst.createTime = Times.unix().toLong(){{end}}
            }
            {{range $i,$c := exclude .columns $pkName "create_time" "update_time"}}
            dst.{{lower_title $c.Prop}} = e.{{lower_title $c.Prop}}{{end}}\
            {{$c := try_get .columns "update_time"}}
            {{if $c}}dst.updateTime = Times.unix().toLong(){{end}}
            this.repo.persistAndFlush(dst)
            null
        }.error()
    }

    /** 批量保存{{.table.Comment}} */
    @Transactional
    fun saveAll{{$shortTitle}}(entities:Iterable<{{$tableTitle}}{{.global.entity_suffix}}>){
        this.repo.persist(entities)
        this.repo.flush()
    }

    /** 删除{{.table.Comment}} */
    @Transactional
    fun delete{{$shortTitle}}ById(id:{{$pkType}}):Error? {
        return catch {
            this.repo.deleteById(this.parseId(id))
        }.error()
    }

}
