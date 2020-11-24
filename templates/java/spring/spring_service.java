#!target:spring/src/main/java/{{.global.pkg}}/service/{{.table.Title}}Service.java
package {{pkg "java" .global.pkg}}.service

import {{pkg "java" .global.pkg}}.pojo.{{.table.Title}}Entity
import {{pkg "java" .global.pkg}}.repo.{{.table.Title}}JpaRepository
import org.springframework.stereotype.Service
import org.springframework.data.repository.findByIdOrNull
import java.util.*
import javax.annotation.Resource
{{$tableTitle := .table.Title}}
{{$pkType := type "java" .table.PkType}}
/** {{.table.Comment}}服务  */
@Service
public class {{.table.Title}}Service {
    @Resource
    private {{$tableTitle}}JpaRepository repo;


    /** 查找{{.table.Comment}} */
    public {{$tableTitle}}Entity findByIdOrNull({{$pkType}} id){
        return this.repo.findByIdOrNull(id);
    }

    /** 保存{{.table.Comment}} */
    public {{$tableTitle}}Entity save{{$tableTitle}}({{$tableTitle}}Entity {{.table.Name}}){
        return this.repo.save({{.table.Name}})
    }


    /** 批量保存{{.table.Comment}} */
    public Iterable<{{$tableTitle}}Entity> saveAll{{$tableTitle}}(Iterable<{{$tableTitle}}Entity> entities){
        return this.repo.saveAll(entities);
    }

    /** 删除{{.table.Comment}} */
    public void deleteById({{$pkType}} id) {
        this.repo.deleteById(id);
    }

}
