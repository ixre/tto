# SpringBoot+JPA 模板

生成使用`SpringBoot`和`JPA`访问数据的代码模板

## 依赖项

添加`fze-commons`工具库

```kts
// https://mvnrepository.com/artifact/net.fze/fze-commons
implementation("net.fze:fze-commons:0.5.0")
```

添加`javax.persistence`依赖

```xml

<dependency>
    <groupId>javax.persistence</groupId>
    <artifactId>javax.persistence-api</artifactId>
    <version>2.2</version>
</dependency>
```

使用分页时，需配置`MyBatisPlus`分页插件

```java
@Configuration
public class MybatisPlusConfig {
    @Bean
    public MybatisPlusInterceptor mybatisPlusInterceptor() {
        MybatisPlusInterceptor interceptor = new MybatisPlusInterceptor();
        interceptor.addInnerInterceptor(new PaginationInnerInterceptor(DbType.MYSQL));
        return interceptor;
    }
}
```

## 扩展

实体默认使用lombok注解，如需要生成get/set方法，请参考以下代码模板

```go
{{range $i,$c := .columns}}{{$ormType := orm_type "java" $c.Type}}
{{$lowerProp := lower_title $c.Prop}} \
public {{$entity}} set{{$c.Prop}}({{$ormType}} {{$lowerProp}}){
    this.{{$lowerProp}} = {{$lowerProp}};
    return this;
}

/** {{$c.Comment}} */
public {{$ormType}} get{{$c.Prop}}() {
    return this.{{$lowerProp}};
}
{{end}}
```

自定义`save`方法字段

```go
/** 保存{{.table.Comment}} */
@Override
public {{$pkType}} save({{$tableTitle}}{{.global.entity_suffix}} e){
    {{$tableTitle}}{{.global.entity_suffix}} dst;
    {{if num_type .table.PkType}}\
    if (e.get{{$pkProp}}() > 0) {
    {{else}}
    if (e.get{{$pkProp}}() != "") {
    {{end}}
        dst = this.repo.findById(e.get{{$pkProp}}()).orElse(null);
        if(dst == null){
            throw new IllegalArgumentException("no such data");
        }
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
    dst.setUpdateTime(new java.util.Date());{{end}}{{end}}
    dst = this.repo.save(dst,{{$tableTitle}}{{.global.entity_suffix}}::get{{.table.PkProp}});
    return dst.get{{.table.PkProp}}();
}
```
