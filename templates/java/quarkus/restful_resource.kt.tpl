#!target:src/main/kotlin/{{.global.pkg}}/resources/{{.table.Title}}Resource.kt
package {{pkg "kotlin" .global.pkg}}.resources

import {{pkg "kotlin" .global.pkg}}.pojo.{{.table.Title}}Entity
import {{pkg "kotlin" .global.pkg}}.service.{{.table.Title}}Service
import {{pkg "kotlin" .global.pkg}}.component.ReportComponent
import net.fze.common.Result
import net.fze.annotation.Resource
import net.fze.extras.report.DataResult
import net.fze.extras.report.ReportUtils
import javax.inject.Inject
import javax.ws.rs.*
import javax.ws.rs.core.MediaType
import javax.enterprise.context.RequestScoped
import javax.annotation.security.PermitAll

{{$tableTitle := .table.Title}}
{{$shortTitle := .table.ShortTitle}}
{{$pkType := type "kotlin" .table.PkType}}
{{$resPrefix := replace (name_path .table.Name) "/" ":"}}

/* {{.table.Comment}}资源 */
@Path("/{{name_path .table.Name}}")
@Produces(MediaType.APPLICATION_JSON)
@Consumes(MediaType.APPLICATION_JSON)
@RequestScoped
class {{.table.Title}}Resource {
    @Inject private lateinit var service:{{.table.Title}}Service
    @Inject private lateinit var reportComponent: ReportComponent

    /** 获取{{.table.Comment}} */
    @GET@Path("/{id}")
    @PermitAll
    @Resource("{{$resPrefix}}:get",name="获取{{.table.Comment}}")
    fun get(@PathParam("id") id:{{$pkType}}): {{.table.Title}}Entity? {
        return service.find{{$shortTitle}}ByIdOrNull(id)
    }

    /** 创建{{.table.Comment}} */
    @POST
    @PermitAll
    @Resource("{{$resPrefix}}:create",name="创建{{.table.Comment}}")
    fun create(entity: {{.table.Title}}Entity):Result {
        val err = this.service.save{{$shortTitle}}(entity)
        if(err != null)return Result.create(1,err.message)
        return Result.OK
    }

    /** 更新{{.table.Comment}} */
    @PUT@Path("/{id}")
    @PermitAll
    @Resource("{{$resPrefix}}:update",name="更新{{.table.Comment}}")
    fun save(@PathParam("id") id:{{$pkType}},entity: {{.table.Title}}Entity):Result {
        entity.{{lower_title .table.PkProp}} = id
        val err = this.service.save{{$shortTitle}}(entity)
        if(err != null)return Result.create(1,err.message)
        return Result.OK
    }


    /** 删除{{.table.Comment}} */
    @DELETE@Path("/{id}")
    @PermitAll
    @Resource("{{$resPrefix}}:delete",name="删除{{.table.Comment}}")
    fun delete(@PathParam("id") id:{{$pkType}}):Result {
        val err = this.service.delete{{$shortTitle}}ById(id)
        if(err != null)return Result.create(1,err.message)
        return Result.OK
    }

    /** {{.table.Comment}}列表 */
    @GET
    @PermitAll
    @Resource("{{$resPrefix}}:list",name="查询{{.table.Comment}}")
    fun list(@QueryParam("params") params:String="{}"): List<{{.table.Title}}Entity> {
        //val p = ReportUtils.parseParams(params).getValue()
        return mutableListOf()
    }

    /** {{.table.Comment}}分页数据 */
    @GET@Path("/paging")
    @PermitAll
    @Resource("{{$resPrefix}}:paging",name="查询{{.table.Comment}}分页数据")
    fun paging(@QueryParam("params",required = false) params:String = "{}",
               @QueryParam("page") page:String,
               @QueryParam("rows") rows:String
    ): DataResult {
        val p = ReportUtils.parseParams(params)
        return this.reportComponent.fetchData("default",
                "{{.table.Prefix}}/{{substr_n .table.Name "_" 1}}_list", params, page, rows)
    }
}