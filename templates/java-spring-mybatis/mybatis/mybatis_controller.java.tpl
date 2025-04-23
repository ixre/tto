#!target:spring/src/main/java/{{.global.pkg}}/controller/{{.table.Title}}Controller.java
package {{pkg "kotlin" .global.pkg}}.controller;

import {{pkg "kotlin" .global.pkg}}.entity.{{.table.Title}}{{.global.entity_suffix}};
import {{pkg "kotlin" .global.pkg}}.service.I{{.table.Title}}Service;
import net.fze.common.data.PagingParams;
import net.fze.annotation.Resource;
import net.fze.common.Result;
import net.fze.util.Systems;
import net.fze.common.data.PagingResult;
import net.fze.ext.report.Params;
import net.fze.util.Assert;
import net.fze.ext.mybatis.MyBatisQueryWrapper;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.web.bind.annotation.*;
import java.io.Serializable;
import java.util.List;
import java.util.Map;


{{$tableTitle := .table.Title}}\
{{$shortTitle := .table.ShortTitle}}\
{{$entityType := join .table.Title .global.entity_suffix }}\
{{$pkType := type "java" .table.PkType}}\
{{$pkOrmType := orm_type "java" .table.PkType}}\
{{$resPrefix := replace (name_path .table.Name) "/" ":"}}\
{{$basePath := join .global.base_path (path .table.Name) "/"}}\

/* {{.table.Comment}}资源 */
@RestController
@RequestMapping("{{$basePath}}s")
public class {{.table.Title}}Controller {
    @Autowired private I{{.table.Title}}Service service;

    /** 获取{{.table.Comment}} */
    @GetMapping("/{id}")
    @Resource(key = "{{$resPrefix}}:get",name="获取{{.table.Comment}}")
    public {{$entityType}} get{{$shortTitle}}(@PathVariable("id") {{$pkType}} id){
        return this.service.findById(id);
    }

    /** 创建{{.table.Comment}} */
    @PostMapping
    @Resource(key = "{{$resPrefix}}:create",name="创建{{.table.Comment}}")
    public Result create{{$shortTitle}}(@RequestBody {{$entityType}} entity){
        Error err = Systems.catchError(()->{
            this.validate{{$shortTitle}}(entity);
            this.service.save(entity);
            return null;
        });
        return Result.of(err);
    }

    /** 更新{{.table.Comment}} */
    @PutMapping("/{id}")
    @Resource(key = "{{$resPrefix}}:update",name="更新{{.table.Comment}}")
    public Result update{{$shortTitle}}(@PathVariable("id") {{$pkType}} id,@RequestBody {{$entityType}} entity) {
        Error err = Systems.catchError(()->{
            entity.set{{.table.PkProp}}(id);
            this.validate{{$shortTitle}}(entity);
            this.service.save(entity);
            return null;
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
        this.service.deleteById(id);
        return Result.of(null);
    }

    /** {{.table.Comment}}分页数据 */
    @GetMapping
    @Resource(key = "{{$resPrefix}}:paging",name="查询{{.table.Comment}}分页数据")
    public PagingResult<{{$entityType}}> paging{{$shortTitle}}(@RequestParam Map<String,Object> params,
               @RequestParam("page") int page,
               @RequestParam("size") int size){
        Params p = new Params(params);
        //https://www.zhihu.com/question/586324313/answer/2912413783
        //String timeRangeSql = ReportUtils.timestampSQLByJSONTime(p.get("create_time"), "create_time");
        //p.set("create_time", timeRangeSql);
        MyBatisQueryWrapper<{{$entityType}}> query = new MyBatisQueryWrapper<>();
        //String keyword = TypeConv.toString(params.get("keyword"));
        //query.eq("name",keyword)
        return this.service.selectPaging(query, PagingParams.of(page,size));
    }

    /** 批量删除{{.table.Comment}} */
    @DeleteMapping("")
    @Resource(key = "{{$resPrefix}}:delete",name="删除{{.table.Comment}}")
    public Result batchDelete(@RequestBody List<Serializable> id){
        if(id.isEmpty())return Result.error(2,"没有要删除的行");
        this.service.batchDelete(id);
        return Result.success();
    }
}
