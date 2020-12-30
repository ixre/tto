#!target:java/{{.global.pkg}}/service/{{.table.Title}}Service.java
package {{pkg "java" .global.pkg}}.service

import {{pkg "java" .global.pkg}}.pojo.{{.table.Title}}{{.global.entity_suffix}}
import {{pkg "java" .global.pkg}}.repo.{{.table.Title}}Repository
import org.springframework.stereotype.Service
import java.util.*
import javax.annotation.Resource
{{$tableTitle := .table.Title}}
{{$pkType := type "java" .table.PkType}}
/** {{.table.Comment}}服务  */
@Service
public class {{.table.Title}}Service {
    @Resource
    private {{$tableTitle}}Repository repo;

    /** 保存{{.table.Comment}} */
    public {{$tableTitle}}{{.global.entity_suffix}} save{{$tableTitle}}({{$tableTitle}}{{.global.entity_suffix}} {{.table.Name}}){
        return this.repo.save({{.table.Name}})
    }


    /** 批量保存{{.table.Comment}} */
    public Iterable<{{$tableTitle}}{{.global.entity_suffix}}> saveAll{{$tableTitle}}(Iterable<{{$tableTitle}}{{.global.entity_suffix}}> entities){
        return this.repo.saveAll(entities);
    }

    /** 删除{{.table.Comment}} */
    public void deleteById({{$pkType}} id) {
        this.repo.deleteById(id);
    }

    /** 查找{{.table.Comment}} */
    public Optional<{{$tableTitle}}{{.global.entity_suffix}}> findById({{$pkType}} id){
        return this.repo.findById(id);
    }
}
