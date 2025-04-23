#!target:spring/src/main/java/{{.global.pkg}}/service/impl/{{.table.Title}}ServiceImpl.java
package {{pkg "java" .global.pkg}}.service.impl;

import {{pkg "java" .global.pkg}}.entity.{{.table.Title}}{{.global.entity_suffix}};
import {{pkg "java" .global.pkg}}.mapper.{{.table.Title}}Mapper;
import {{pkg "java" .global.pkg}}.service.I{{.table.Title}}Service;
import org.springframework.stereotype.Service;
import net.fze.service.SpringBaseServiceImpl;

{{$tableTitle := .table.Title}}\
{{$shortTitle := .table.ShortTitle}}\
{{$pkName := .table.Pk}}\
{{$pkProp :=  .table.PkProp}}\
{{$pkType := type "java" .table.PkType}}
{{$warpPkType := orm_type "java" .table.PkType}}

/** {{.table.Comment}}服务  */
@Service("{{.table.Name}}_mybatis_service")
public class {{.table.Title}}ServiceImpl extends SpringBaseServiceImpl<{{$tableTitle}}{{.global.entity_suffix}}> implements I{{$tableTitle}}Service {
    /** 保存{{.table.Comment}} */
    @Override
    public void save({{$tableTitle}}{{.global.entity_suffix}} e){
        super.save(e,{{$tableTitle}}{{.global.entity_suffix}}::get{{$pkProp}});
    }
}
