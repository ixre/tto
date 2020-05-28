#!target:kotlin/{{.global.Pkg}}/resources/{{.table.Title}}Resource.kt.gen
package {{pkg "java" .global.Pkg}}.resources

import {{pkg "java" .global.Pkg}}.pojo.{{.table.Title}}Entity
import {{pkg "java" .global.Pkg}}.service.{{.table.Title}}Service
import {{pkg "java" .global.Pkg}}.component.TinyQueryComponent
import net.fze.arch.commons.std.Result
import net.fze.arch.component.report.DataResult
import javax.inject.Inject
import javax.ws.rs.*
import javax.ws.rs.core.MediaType
import javax.enterprise.context.RequestScoped
import javax.annotation.security.PermitAll

{{$tableTitle := .table.Title}}
{{$pkType := type "kotlin" .table.PkType}}

/* {{.table.Comment}}资源 */
@Path("/{{.table.Name}}")
@Produces(MediaType.APPLICATION_JSON)
@Consumes(MediaType.APPLICATION_JSON)
@RequestScoped
class {{.table.Title}}Resource {
    @Inject private lateinit var service:{{.table.Title}}Service
    @Inject private lateinit var queryComponent: TinyQueryComponent

    @GET@Path("/{id}")
    @PermitAll
    fun get(@PathParam("id") id:{{$pkType}}): {{.table.Title}}Entity? {
        return service.findByIdOrNull(id)
    }


    @POST
    @PermitAll
    fun create(entity: {{.table.Title}}Entity):Result {
        val err = this.service.save{{.table.Title}}(entity)
        if(err != null)return Result.create(1,err.message)
        return Result.OK
    }

    @PUT@Path("/{id}")
    @PermitAll
    fun save(@PathParam("id") id:{{$pkType}},entity: {{.table.Title}}Entity):Result {
        entity.{{lower_title .table.PkProp}} = id
        val err = this.service.save{{.table.Title}}(entity)
        if(err != null)return Result.create(1,err.message)
        return Result.OK
    }


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
    fun list(): Set<{{.table.Title}}Entity> {
        return mutableSetOf()
    }

    /** {{.table.Comment}}分页数据 */
    @GET@Path("/paging")
    @PermitAll
    fun paging(@QueryParam("params") params:String,
               @QueryParam("page") page:String,
               @QueryParam("rows") rows:String
    ): DataResult {
        return this.queryComponent.fetchData("default",
                "{{.table.Title}}List", params, page, rows)
    }
}