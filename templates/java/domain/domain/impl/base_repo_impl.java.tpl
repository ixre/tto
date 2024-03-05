#!target:domain/src/main/java/{{.global.pkg}}/domain/impl/{{.table.Prefix}}/Base{{.table.Title}}Repository.java
package {{pkg "java" .global.pkg}}.domain.impl.{{.table.Prefix}};

import {{pkg "java" .global.pkg}}.domain.iface.{{.table.Prefix}}.{{.table.Title}}{{.global.entity_suffix}};
import {{pkg "java" .global.pkg}}.domain.iface.{{.table.Prefix}}.I{{.table.Title}}Repository;
import {{pkg "java" .global.pkg}}.domain.impl.{{.table.Prefix}}.Base{{.table.Title}}Repository;
import net.fze.util.Systems;
import net.fze.util.Times;

import java.util.List;

{{$comment := .table.Comment}}
{{$suffix := .table.ShortTitle}}
{{$tableTitle := .table.Title}}

/** 
 * {{$comment}}仓储基础类
 *
 * @author {{.global.user}}
 */
public class Base{{.table.Title}}Repository {
    @Inject I{{$tableTitle}}Repository _repo; 
    /**
     * 创建{{.table.Comment}}
     */
    public I{{$tableTitle}}AggregateRoot create{{$suffix}}({{$tableTitle}}{{.global.entity_suffix}} e){
        return new {{$tableTitle}}AggregateRootImpl(this._repo,e);
    }
}
