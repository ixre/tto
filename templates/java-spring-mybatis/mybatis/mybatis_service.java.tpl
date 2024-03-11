#!target:spring/src/main/java/{{.global.pkg}}/service/{{.table.Prefix}}/impl/{{.table.Title}}ServiceImpl.java
package {{pkg "java" .global.pkg}}.service.impl;

import {{pkg "java" .global.pkg}}.entity.{{.table.Title}}{{.global.entity_suffix}};
import {{pkg "java" .global.pkg}}.mapper.{{.table.Title}}Mapper;
import {{pkg "java" .global.pkg}}.service.{{.table.Prefix}}.I{{.table.Title}}Service;
import com.baomidou.mybatisplus.core.conditions.query.QueryWrapper;
import com.baomidou.mybatisplus.extension.plugins.pagination.Page;
import org.springframework.data.domain.Pageable;
import org.springframework.stereotype.Service;
import net.fze.util.Systems;
import net.fze.common.data.PagingResult;
import net.fze.util.Times;
import javax.inject.Inject;
import java.util.List;
{{$tableTitle := .table.Title}}\
{{$shortTitle := .table.ShortTitle}}\
{{$pkName := .table.Pk}}\
{{$pkProp :=  .table.PkProp}}\
{{$pkType := type "java" .table.PkType}}
{{$warpPkType := orm_type "java" .table.PkType}}

/** {{.table.Comment}}服务  */
@Service("{{.table.Name}}_mybatis_service")
public class {{.table.Title}}ServiceImpl implements I{{.table.Title}}Service{
    @Inject
    private {{$tableTitle}}Mapper repo;

    /** 查找{{.table.Comment}} */
    @Override
    public {{$tableTitle}}{{.global.entity_suffix}} find{{$shortTitle}}ById({{$pkType}} id){
        return this.repo.findById(id).orElse(null);
    }

    /** 查找全部{{.table.Comment}} */
    @Override
    public List<{{$tableTitle}}{{.global.entity_suffix}}> findAll{{$shortTitle}}() {
        return this.repo.findBy(null);
    }

    /** 保存{{.table.Comment}} */
    @Override
    public {{$pkType}} save{{$shortTitle}}({{$tableTitle}}{{.global.entity_suffix}} e){
        {{$tableTitle}}{{.global.entity_suffix}} dst;
        {{if num_type .table.PkType}}\
        if (e.get{{$pkProp}}() > 0) {
        {{else}}
        if (e.get{{$pkProp}}() != "") {
        {{end}}
            dst = this.repo.findById(e.get{{$pkProp}}()).orElse(null);
            if(dst == null){
                throw new IllegalArgumentException("no such data");
            }
        } else {
            dst = {{$tableTitle}}{{.global.entity_suffix}}.createDefault();
            {{$c := try_get .columns "create_time"}}\
            {{if $c}}{{if num_type $c.Type }}\
            dst.setCreateTime(Times.unix());
            {{else}}\
            dst.setCreateTime(new java.util.Date());{{end}}{{end}}
        }\
        {{range $i,$c := exclude .columns $pkName "create_time" "update_time"}}
        dst.set{{$c.Prop}}(e.get{{$c.Prop}}());{{end}}\
        {{$c := try_get .columns "update_time"}}
        {{if $c}}{{if num_type $c.Type }}\
        dst.setUpdateTime(Times.unix());
        {{else}}\
        dst.setUpdateTime(new java.util.Date());{{end}}{{end}}
        dst = this.repo.save(dst,{{$tableTitle}}{{.global.entity_suffix}}::get{{.table.PkProp}});
        return dst.get{{.table.PkProp}}();
    }
    
    /** 根据对象条件查找 */
    @Override
    public {{$tableTitle}}{{.global.entity_suffix}} find{{$shortTitle}}By({{$tableTitle}}{{.global.entity_suffix}} o){
         return this.repo.findOne(o).orElse(null);
    }

    /** 根据对象条件查找 */
    @Override
    public List<{{$tableTitle}}{{.global.entity_suffix}}> find{{$shortTitle}}ListBy({{$tableTitle}}{{.global.entity_suffix}} o) {
         return this.repo.findBy(o);
    }

    /** 根据条件分页查询 */
    @Override
    public PagingResult<{{$tableTitle}}{{.global.entity_suffix}}> findPaging{{$shortTitle}}({{$tableTitle}}{{.global.entity_suffix}} o, Pageable page) {
        Page<{{$tableTitle}}{{.global.entity_suffix}}> p = this.repo.selectPage(new Page<>(page.getPageNumber(), page.getPageSize()),
              new QueryWrapper<>(o));
        return new PagingResult<>(p.getTotal(),p.getRecords());
    }

    /** 批量保存{{.table.Comment}} */
    @Override
    public Iterable<{{$tableTitle}}{{.global.entity_suffix}}> saveAll{{$shortTitle}}(Iterable<{{$tableTitle}}{{.global.entity_suffix}}> entities){
        entities.forEach(a->this.repo.updateById(a));
        return entities;
    }

    /** 删除{{.table.Comment}} */
    @Override
    public void delete{{$shortTitle}}ById({{$pkType}} id) {
        this.repo.deleteById(id);
    }

    /** 批量删除{{.table.Comment}} */
    @Override
    public void batchDelete{{$shortTitle}}(List<{{$warpPkType}}> id){
       this.repo.deleteBatchIds(id);
    }
}
