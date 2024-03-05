#!target:domain/src/main/java/{{.global.pkg}}/domain/impl/{{.table.Prefix}}/{{.table.Title}}AggregateRootImpl.java
package {{pkg "java" .global.pkg}}.domain.impl.{{.table.Prefix}};

import java.util.List;

{{$comment := .table.Comment}}
{{$pkName := .table.Pk}}
{{$pkProp :=  .table.PkProp}}
{{$pkType := orm_type "java" .table.PkType}}
{{$suffix := .table.ShortTitle}}
{{$tableTitle := .table.Title}}

/** 
 * {{$comment}}聚合根实现类
 *
 * @author {{.global.user}}
 */
public class {{$tableTitle}}AggregateRootImpl implements I{{$tableTitle}}AggregateRoot{

    private I{{$tableTitle}}Repository _repo; 
    
    private {{$tableTitle}}{{.global.entity_suffix}} _value = null;

    protected {{$tableTitle}}AggregateRootImpl(I{{$tableTitle}}Repository repo, {{$tableTitle}}{{.global.entity_suffix}} value) {
        this._repo = repo;
        this._value = value;
    }

    @Override
    public {{$pkType}} getAggregateRootId() {
        return this._value.get{{$pkProp}}();
    }

    /**
     * 保存
     */
    public void save(){
        this._repo.save{{$suffix}}(this._value);
    }

    /**
     * 获取值
     */
    public {{$tableTitle}}{{.global.entity_suffix}} getValue(){
        return this._value.clone();
    }

    /**
     * 销毁
     */ 
    public void destroy(){
        this._repo.delete{{$suffix}}ById(this.getAggregateRootId());
    }
}
