#!target:domain/src/main/java/{{.global.pkg}}/domain/iface/{{.table.Prefix}}/I{{.table.Title}}AggregateRoot.java
package {{pkg "java" .global.pkg}}.domain.iface.{{.table.Prefix}};

import java.util.List;

{{$comment := .table.Comment}}
{{$pkName := .table.Pk}}
{{$pkProp :=  .table.PkProp}}
{{$pkType := orm_type "java" .table.PkType}}
{{$suffix := .table.ShortTitle}}
{{$tableTitle := .table.Title}}

/** 
 * {{$comment}}聚合根接口
 *
 * @author {{.global.user}}
 */
public interface I{{$tableTitle}}AggregateRoot extends IAggregateRoot<{{$pkType}}>{
    /**
     * 保存
     */
    void save();

    /**
     * 获取值
     */
    {{$tableTitle}}{{.global.entity_suffix}} getValue();

    /**
     * 销毁
     */ 
    void destroy();
}
