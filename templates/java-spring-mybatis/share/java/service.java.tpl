#!target:spring/src/main/java/{{.global.pkg}}/service/I{{.table.Title}}Service.java
package {{pkg "java" .global.pkg}}.service;

{{$tableTitle := .table.Title}}\
{{$shortTitle := .table.ShortTitle}}\
{{$entityType := join .table.Title .global.entity_suffix }}\
{{$pkName := .table.Pk}}\
{{$pkProp :=  .table.PkProp}}\
{{$pkType := type "java" .table.PkType}}\
{{$warpPkType := orm_type "java" .table.PkType}}\

import {{pkg "java" .global.pkg}}.entity.{{.table.Title}}{{.global.entity_suffix}};
import net.fze.service.IBaseService;

/** {{.table.Comment}}服务  */
public interface I{{$tableTitle}}Service extends IBaseService<{{$entityType}}>{
    /** 保存{{.table.Comment}} */
    void save({{$entityType}} e);
}
