#!target:domain/src/main/java/{{.global.pkg}}/service/{{.table.Title}}ServiceImpl.java
package {{pkg "java" .global.pkg}}.service;

import {{pkg "java" .global.pkg}}.entity.{{.table.Title}}{{.global.entity_suffix}};
import {{pkg "java" .global.pkg}}.repo.{{.table.Prefix}}.{{.table.Title}}JpaRepository;
import {{pkg "java" .global.pkg}}.service.{{.table.Prefix}}.I{{.table.Title}}Service;
import net.fze.util.Systems;
import net.fze.util.Times;
import javax.inject.Inject;
import java.util.List;

{{$comment := .table.Comment}}
{{$pkName := .table.Pk}}
{{$pkProp :=  .table.PkProp}}
{{$pkType := orm_type "java" .table.PkType}}
{{$suffix := .table.ShortTitle}}
{{$tableTitle := .table.Title}}

/** {{.table.Comment}}服务  */
@Singleton
public class {{$tableTitle}}ServiceImpl{
    @Inject
    private I{{$tableTitle}}Repository _repo;

    /** 获取{{.table.Comment}} */
    public {{$tableTitle}}{{.global.entity_suffix}} find{{$suffix}}ById({{$pkType}} id){
        I{{$tableTitle}}AggregateRoot ia = this._repo.get{{$suffix}}(id);
        if(ia != null){
            return ia.getValue();
        }
        return null;
    }

    /** 保存{{.table.Comment}} */
    @Override
    public Error save{{$suffix}}({{$tableTitle}}{{.global.entity_suffix}} e){
         I{{$tableTitle}}AggregateRoot ia;
         ITransaction trans = InjectFactory.getInstance(ITransactionManager.class).beginTransaction();
         return Systems.tryCatch(()-> {
            {{$tableTitle}}{{.global.entity_suffix}} dst;
            {{if num_type .table.PkType}}\
            if (e.get{{$pkProp}}() > 0) {
            {{else}}
            if (e.get{{$pkProp}}() != "") {
            {{end}}
                ia = this._repo.get{{$suffix}}(e.get{{$pkProp}}());
                if(ia == null)throw new IllegalArgumentException("no such data");
                dst = ia.getValue();
            } else {
                dst = {{$tableTitle}}{{.global.entity_suffix}}.createDefault();
                {{$c := try_get .columns "create_time"}}\
                {{if $c}}{{if num_type $c.Type }}\
                dst.setCreateTime(Times.unix());
                {{else}}\
                dst.setCreateTime(new java.util.Date());{{end}}{{end}}
                ia = this._repo.create{{$suffix}}(dst);
            }\
            {{range $i,$c := exclude .columns $pkName "create_time" "update_time"}}
            dst.set{{$c.Prop}}(e.get{{$c.Prop}}());{{end}}\
            {{$c := try_get .columns "update_time"}}
            {{if $c}}{{if num_type $c.Type }}\
            dst.setUpdateTime(Times.unix());
            {{else}}\
            dst.setUpdateTime(new java.util.Date());{{end}}{{end}}
            ia.setValue(dst);
            ia.save();
            e.set{{.table.PkProp}}(ia.getAggregateRootId());
            trans.commit();
            return null;
          }).except(it->{
            trans.rollback();
            it.printStackTrace();
            return null;
         }).error();
    }
}
