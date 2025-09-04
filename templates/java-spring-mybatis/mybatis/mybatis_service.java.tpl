#!target:spring/src/main/java/{{.global.pkg}}/service/impl/{{.table.Title}}ServiceImpl.java
package {{pkg "java" .global.pkg}}.service.impl;

import {{pkg "java" .global.pkg}}.entity.{{.table.Title}}{{.global.entity_suffix}};
import {{pkg "java" .global.pkg}}.service.I{{.table.Title}}Service;
import org.springframework.stereotype.Service;
import net.fze.service.SpringBaseServiceImpl;
import net.fze.util.Times;

{{$tableTitle := .table.Title}}\
{{$shortTitle := .table.ShortTitle}}\
{{$pkName := .table.Pk}}\
{{$pkProp :=  .table.PkProp}}\
{{$pkType := type "java" .table.PkType}}
{{$warpPkType := orm_type "java" .table.PkType}}

/** {{.table.Comment}}服务  */
@Service("{{.table.Name}}_mybatis_service")
public class {{.table.Title}}ServiceImpl extends SpringBaseServiceImpl<{{$tableTitle}}{{.global.entity_suffix}}> implements I{{$tableTitle}}Service {
    /** 保存{{.table.Comment}},不采用BeanUtils是因为可以更好的控制哪些字段需更新，哪些需初始化 */
    @Override
    public void save({{$tableTitle}}{{.global.entity_suffix}} e){
        {{$tableTitle}}{{.global.entity_suffix}} dst;
        {{if num_type .table.PkType}}\
        if (e.get{{$pkProp}}() != null && e.get{{$pkProp}}() > 0) {
        {{else}}
        if (e.get{{$pkProp}}() != null && e.get{{$pkProp}}() != "") {
        {{end}}
            dst = this.findById(e.get{{$pkProp}}());
            if(dst == null)throw new IllegalArgumentException("no such data");
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
        dst.setUpdateTime(new java.util.Date());{{end}}{{end}}\
        super.save(dst,{{$tableTitle}}{{.global.entity_suffix}}::get{{$pkProp}});
    }
}
