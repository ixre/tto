#!target:src/main/kotlin/{{.global.pkg}}/resources/{{.table.Title}}Resource.kt
package {{pkg "kotlin" .global.pkg}}.resources

import {{pkg "kotlin" .global.pkg}}.pojo.{{.table.Title}}Entity
import {{pkg "kotlin" .global.pkg}}.service.{{.table.Title}}Service
import {{pkg "kotlin" .global.pkg}}.component.TinyQueryComponent
import net.fze.common.Result
import net.fze.extras.report.DataResult
import javax.inject.Inject
import javax.ws.rs.*
import javax.ws.rs.core.MediaType
import javax.enterprise.context.RequestScoped
import javax.annotation.security.PermitAll

{{$tableTitle := .table.Title}}
{{$pkType := type "kotlin" .table.PkType}}

/* {{.table.Comment}}资源 */
@Path("/{{name_path .table.Name}}")
@Produces(MediaType.APPLICATION_JSON)
@Consumes(MediaType.APPLICATION_JSON)
@RequestScoped
class {{.table.Title}}Resource {
    @Inject private lateinit var service:{{.table.Title}}Service
    @Inject private lateinit var queryComponent: TinyQueryComponent

    /** 获取{{.table.Comment}} */
    @GET@Path("/{id}")
    @PermitAll
    fun get(@PathParam("id") id:{{$pkType}}): {{.table.Title}}Entity? {
        return service.findByIdOrNull(id)
    }

    /** 创建{{.table.Comment}} */
    @POST
    @PermitAll
    fun create(entity: {{.table.Title}}Entity):Result {
        val err = this.service.save{{.table.Title}}(entity)
        if(err != null)return Result.create(1,err.message)
        return Result.OK
    }

    /** 更新{{.table.Comment}} */
    @PUT@Path("/{id}")
    @PermitAll
    fun save(@PathParam("id") id:{{$pkType}},entity: {{.table.Title}}Entity):Result {
        entity.{{lower_title .table.PkProp}} = id
        val err = this.service.save{{.table.Title}}(entity)
        if(err != null)return Result.create(1,err.message)
        return Result.OK
    }


    /** 删除{{.table.Comment}} */
    @DELETE@Path("/{id}")
    @PermitAll
    fun delete(@PathParam("id") id:{{$pkType}}):Result {
        val err = this.service.deleteById(id)
        if(err != null)return Result.create(1,err.message)
        return Result.OK
    }

    /** {{.table.Comment}}列表 */
    @GET
    @PermitAll
    fun list(): List<{{.table.Title}}Entity> {
        return mutableListOf()
    }

    /** {{.table.Comment}}分页数据 */
    @GET@Path("/paging")
    @PermitAll
    fun paging(@QueryParam("params") params:String,
               @QueryParam("page") page:String,
               @QueryParam("rows") rows:String
    ): DataResult {
        return this.queryComponent.fetchData("default",
                "{{.table.Prefix}}/{{substr_n .table.Name "_" 1}}_list", params, page, rows)
    }
}