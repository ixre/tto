#!target:spring/src/main/java/{{.global.pkg}}/service/{{.table.Prefix}}/I{{.table.Title}}Service.java
package {{pkg "java" .global.pkg}}.service;

import {{pkg "java" .global.pkg}}.entity.{{.table.Title}}{{.global.entity_suffix}};
import org.springframework.data.domain.Pageable;
import net.fze.common.data.PagingResult;

import java.util.List;
{{$tableTitle := .table.Title}}\
{{$shortTitle := .table.ShortTitle}}\
{{$pkName := .table.Pk}}\
{{$pkProp :=  .table.PkProp}}\
{{$pkType := type "java" .table.PkType}}
{{$warpPkType := orm_type "java" .table.PkType}}

/** {{.table.Comment}}服务  */
public interface I{{.table.Title}}Service {

    /** 查找{{.table.Comment}} */
    {{$tableTitle}}{{.global.entity_suffix}} find{{$shortTitle}}ById({{$pkType}} id);

    /** 查找全部{{.table.Comment}} */
    List<{{$tableTitle}}{{.global.entity_suffix}}> findAll{{$shortTitle}}();

    /** 保存{{.table.Comment}} */
    {{$pkType}} save{{$shortTitle}}({{$tableTitle}}{{.global.entity_suffix}} e);

    /** 根据对象条件查找 */
    {{$tableTitle}}{{.global.entity_suffix}} find{{$shortTitle}}By({{$tableTitle}}{{.global.entity_suffix}} o);

    /** 根据对象条件查找 */
    List<{{$tableTitle}}{{.global.entity_suffix}}> find{{$shortTitle}}ListBy({{$tableTitle}}{{.global.entity_suffix}} o);

    /** 根据条件分页查询 */
    PagingResult<{{$tableTitle}}{{.global.entity_suffix}}> findPaging{{$shortTitle}}({{$tableTitle}}{{.global.entity_suffix}} o, Pageable page);

    /** 批量保存{{.table.Comment}} */
    Iterable<{{$tableTitle}}{{.global.entity_suffix}}> saveAll{{$shortTitle}}(Iterable<{{$tableTitle}}{{.global.entity_suffix}}> entities);

    /** 删除{{.table.Comment}} */
    void delete{{$shortTitle}}ById({{$pkType}} id);

    /** 批量删除{{.table.Comment}} */
    void batchDelete{{$shortTitle}}(List<{{$warpPkType}}> id);
}
