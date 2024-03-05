#!target:domain/src/main/java/{{.global.pkg}}/domain/iface/{{.table.Prefix}}/I{{.table.Title}}Repository.java
package {{pkg "java" .global.pkg}}.domain.iface.{{.table.Prefix}};

import java.util.List;

{{$comment := .table.Comment}}
{{$pkName := .table.Pk}}
{{$pkProp :=  .table.PkProp}}
{{$pkType := orm_type "java" .table.PkType}}
{{$suffix := .table.ShortTitle}}
{{$tableTitle := .table.Title}}

/** 
 * {{$comment}}仓储接口
 *
 * @author {{.global.user}}
 */
public interface I{{$tableTitle}}Repository{

    /**
     * 创建{{$comment}}聚合
     *
     * @param e 数据
     * @return {{$comment}}
     */
    I{{$suffix}}AggregateRoot create{{$suffix}}({{$tableTitle}}{{.global.entity_suffix}} e);

    /**
     * 获取{{$comment}}聚合
     *
     * @param {{$pkName}} 标识
     * @return {{$comment}}
     */
    I{{$suffix}}AggregateRoot get{{$suffix}}({{$pkType}} {{$pkName}});

    /**
     * 保存{{.table.Comment}}
     *
     * @param e 实体
     */
    void save{{$suffix}}({{$tableTitle}}{{.global.entity_suffix}} e);

    /**
     * 根据对象条件查找{{.table.Comment}}
     *
     * @param e 实体
     */
    {{$tableTitle}}{{.global.entity_suffix}} find{{$suffix}}By({{$tableTitle}}{{.global.entity_suffix}} e);

    /**
     * 根据对象条件查找{{.table.Comment}}
     *
     * @param e 实体
     */
    List<{{$tableTitle}}{{.global.entity_suffix}}> find{{$suffix}}ListBy({{$tableTitle}}{{.global.entity_suffix}} e);

    /**
     * 删除{{.table.Comment}}
     *
     * @param id 标识
     */
    int delete{{$suffix}}({{$pkType}} id);
}
