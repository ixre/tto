#!target:spring/src/main/java/{{.global.pkg}}/restful/{{.table.Prefix}}/{{.table.Title}}Resource.java
package {{pkg "kotlin" .global.pkg}}.restful;

import {{pkg "kotlin" .global.pkg}}.entity.{{.table.Title}}{{.global.entity_suffix}};
import {{pkg "kotlin" .global.pkg}}.service.{{.table.Prefix}}.I{{.table.Title}}Service;
import {{pkg "kotlin" .global.pkg}}.component.ReportDataSource;
import net.fze.annotation.Resource;
import net.fze.common.Result;
import net.fze.common.Standard;
import net.fze.extras.report.DataResult;
import net.fze.extras.report.ReportUtils;
import net.fze.extras.report.Params;
import net.fze.util.Strings;
import net.fze.util.Assert;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.web.bind.annotation.*;
import javax.inject.Inject;
import java.util.List;
import java.util.ArrayList;

{{$tableTitle := .table.Title}}\
{{$shortTitle := .table.ShortTitle}}\
{{$entityType := join .table.Title .global.entity_suffix }}\
{{$pkType := type "java" .table.PkType}}\
{{$pkOrmType := orm_type "java" .table.PkType}}\
{{$resPrefix := replace (name_path .table.Name) "/" ":"}}\
{{$basePath := join .global.base_path (name_path .table.Name) "/"}}\

/* {{.table.Comment}}资源 */
@RestController
@RequestMapping("{{$basePath}}")
public class {{.table.Title}}Resource {
    @Inject private I{{.table.Title}}Service service;
    @Inject private ReportDataSource reportDs;

    /** 获取{{.table.Comment}} */
    @GetMapping("/{id}")
    @Resource(key = "{{$resPrefix}}:get",name="获取{{.table.Comment}}")
    public {{$entityType}} get{{$shortTitle}}(@PathVariable("id") {{$pkType}} id){
        return this.service.find{{$shortTitle}}ById(id);
    }

    /** 创建{{.table.Comment}} */
    @PostMapping
    @Resource(key = "{{$resPrefix}}:create",name="创建{{.table.Comment}}")
    public Result create{{$shortTitle}}(@RequestBody {{$entityType}} entity){
        Error err = Standard.catchError(()->{
            this.validate{{$shortTitle}}(entity);
            return this.service.save{{$shortTitle}}(entity);
        });
        return Result.of(err);
    }

    /** 更新{{.table.Comment}} */
    @PutMapping("/{id}")
    @Resource(key = "{{$resPrefix}}:update",name="更新{{.table.Comment}}")
    public Result update{{$shortTitle}}(@PathVariable("id") {{$pkType}} id,@RequestBody {{$entityType}} entity) {
        Error err = Standard.catchError(()->{
            entity.set{{.table.PkProp}}(id);
            this.validate{{$shortTitle}}(entity);
            return this.service.save{{$shortTitle}}(entity);
        });
        return Result.of(err);
    }

    /** 验证实体 */
    private void validate{{$shortTitle}}({{$entityType}} e){
       {{$validateColumns := exclude .columns .table.Pk "create_time" "update_time" "state"}}
       {{range $i,$c := $validateColumns}}
       {{if $c.NotNull }}\
        Assert.isNullOrEmpty(e.get{{$c.Prop}}(), "{{$c.Comment}}不能为空");\
       {{end}}{{end}}
    }

    /** 删除{{.table.Comment}} */
    @DeleteMapping("/{id}")
    @Resource(key = "{{$resPrefix}}:delete",name="删除{{.table.Comment}}")
    public Result delete{{$shortTitle}}(@PathVariable("id") {{$pkType}} id){
        Error err = this.service.delete{{$shortTitle}}ById(id);
        return Result.of(err);
    }

    /** {{.table.Comment}}分页数据 */
    @PostMapping("/paging")
    @Resource(key = "{{$resPrefix}}:paging",name="查询{{.table.Comment}}分页数据")
    public DataResult paging{{$shortTitle}}(@RequestBody Map<String,Object> params,
               @RequestParam("page") int page,
               @RequestParam("size") int rows){
        Params p = new Params(params);
        //String timeRangeSQL = ReportUtils.timeRangeSQLByJSONTime(p.get("create_time"), "create_time");
        //p.set("create_time", timeRangeSQL);
        return this.reportDs.fetchData("default",
                "{{.table.Prefix}}/{{substr_n .table.Name "_" 1}}_list", p, page, rows);
    }

    /** 查询{{.table.Comment}}列表 */
    @GetMapping
    @Resource(key = "{{$resPrefix}}:list",name="查询{{.table.Comment}}")
    public List<{{.table.Title}}{{.global.entity_suffix}}> query{{$shortTitle}}(@RequestParam(name="params",defaultValue="{}") String params) {
        //val p = ReportUtils.parseParams(params).getValue();
        return this.service.findAll{{$shortTitle}}();
    }

    /** 批量删除{{.table.Comment}} */
    @DeleteMapping("")
    @Resource(key = "{{$resPrefix}}:delete",name="删除{{.table.Comment}}")
    public Result batchDelete{{$shortTitle}}(@RequestBody List<{{$pkOrmType}}> id){
        if(id.isEmpty())return Result.error(2,"没有要删除的行");
        Error err = this.service.batchDelete{{$shortTitle}}(id);
        return Result.of(err);
    }
}
