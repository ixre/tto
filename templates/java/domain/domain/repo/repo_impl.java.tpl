#!target:domain/src/main/java/{{.global.pkg}}/domain/repo/{{.table.Prefix}}/{{.table.Title}}RepositoryImpl.java
package {{pkg "java" .global.pkg}}.domain.repo.{{.table.Prefix}};

import {{pkg "java" .global.pkg}}.domain.iface.{{.table.Prefix}}.{{.table.Title}}{{.global.entity_suffix}};
import {{pkg "java" .global.pkg}}.domain.iface.{{.table.Prefix}}.I{{.table.Title}}Repository;
import {{pkg "java" .global.pkg}}.domain.impl.{{.table.Prefix}}.Base{{.table.Title}}Repository;
import net.fze.util.Systems;
import net.fze.util.Times;

import java.util.List;

{{$comment := .table.Comment}}
{{$pkName := .table.Pk}}
{{$pkProp :=  .table.PkProp}}
{{$pkType := orm_type "java" .table.PkType}}
{{$suffix := .table.ShortTitle}}
{{$tableTitle := .table.Title}}

/** 
 * {{$comment}}仓储实现类
 *
 * @author {{.global.user}}
 */
public class {{.table.Title}}RepositoryImpl extends Base{{.table.Title}}Repository implements I{{$tableTitle}}Repository {
    // 如果在Springboot中，则应使用InjectFactory.getInstance直接获取实例
    //{{$tableTitle}}Mapper repo = InjectFactory.getInstance({{$tableTitle}}Mapper.class);

    /**
     * 获取{{$comment}}聚合
     *
     * @param {{$pkName}} 标识
     * @return {{$comment}}
     */
    @Override
    public I{{$suffix}}AggregateRoot get{{$suffix}}({{$pkType}} {{$pkName}}){
        {{$tableTitle}}{{.global.entity_suffix}} e = this.find{{$suffix}}ById({{$pkName}});
        if (e != null) {
            return super.create{{$suffix}}(e);
        }
        return null;
    }

    
    /** 查找{{.table.Comment}} */
    private {{$tableTitle}}{{.global.entity_suffix}} find{{$suffix}}ById({{$pkType}} {{$pkName}}){
        // TODO: not implemented
        throw new RuntimeException("not implemented");
        //return this.repo.findById({{$pkName}}).orElse(null);
    }


    /**
     * 保存{{.table.Comment}}
     *
     * @param e 实体
     */
    @Override
    public void save{{$suffix}}({{$tableTitle}}{{.global.entity_suffix}} e){
        // TODO: not implemented
        throw new RuntimeException("not implemented");
        //this.repo.save(e);
    }

    /**
     * 根据对象条件查找{{.table.Comment}}
     *
     * @param e 实体
     */
    @Override
    public {{$tableTitle}}{{.global.entity_suffix}} find{{$suffix}}By({{$tableTitle}}{{.global.entity_suffix}} e){
        // TODO: not implemented
        throw new RuntimeException("not implemented");
        //return this.repo.findOne(e).orElse(null);
    }

    /**
     * 根据对象条件查找{{.table.Comment}}
     *
     * @param e 实体
     */
    @Override
    public List<{{$tableTitle}}{{.global.entity_suffix}}> find{{$suffix}}ListBy({{$tableTitle}}{{.global.entity_suffix}} e) {
        // TODO: not implemented
        throw new RuntimeException("not implemented");
        //return this.repo.findAll(e);
    }


    /**
     * 删除{{.table.Comment}}
     *
     * @param id 标识
     */
    @Override
    public int delete{{$suffix}}({{$pkType}} id) {
        // TODO: not implemented
        throw new RuntimeException("not implemented");
        //return this.repo.deleteById(id);
    }
}
