package github.com.jsix.goex.generator.bin.code_templates.java;
// auto generate by gof (http://github.com/ixre/goex)
import java.util.HashMap;
import java.util.List;
import java.util.Map;
import java.io.Serializable;
import com.google.inject.Inject;
import com.line.arch.component.XContext;
import com.line.arch.component.hibernate.TinySession;
import com.gcy.sz.repo.*;
import com.gcy.sz.repo.model.*;

/** {{.table.Comment}}仓储实现 */
public class {{.table.Title}}RepoImpl implements I{{.table.Title}}Repo {
    /** 注入上下文 */
    @Inject XContext ctx;
    /** 获取{{.table.Comment}} */
    public {{.table.Title}}{{.global.entity_suffix}} get(Serializable id){
        TinySession s = this.ctx.hibernate();
        {{.table.Title}}{{.global.entity_suffix}} e = s.get({{.table.Title}}{{.global.entity_suffix}}.class,id);
        return e;
    }
    /** 根据条件获取单条{{.table.Comment}} */
    public {{.table.Title}}{{.global.entity_suffix}} get{{.table.Title}}By(String where,Map<String,Object> params){
        TinySession s = this.ctx.hibernate();
        {{.table.Title}}{{.global.entity_suffix}} e = s.get({{.table.Title}}{{.global.entity_suffix}}.class,where, params);
        s.close();
        return e;
    }

    /** 根据条件获取多条{{.table.Comment}} */
    public List<{{.table.Title}}{{.global.entity_suffix}}> select{{.table.Title}}(String where,Map<String,Object> params){
        TinySession s = this.ctx.hibernate();
        List<{{.table.Title}}{{.global.entity_suffix}}> list = s.select({{.table.Title}}{{.global.entity_suffix}}.class,
        "SELECT * FROM {{.table.Name}} WHERE "+where, params);
        s.close();
        return list;
    }

    /** 保存{{.table.Comment}} */
    public int save{{.table.Title}}({{.table.Title}}{{.global.entity_suffix}} v){
        TinySession s = this.ctx.hibernate();
        s.save(v);
        s.close();
        return v.getId();
    }

    /** 删除{{.table.Comment}} */
    public Error delete{{.table.Title}}(Serializable id){
        TinySession s = this.ctx.hibernate();
        Map<String, Object> data = new HashMap<>();
        data.put("id", id);
        s.execute("DELETE FROM {{.table.Name}} WHERE id=:id", data);
        s.close();
        return null;
    }

    /** 批量删除{{.table.Comment}} */
    public int BatchDelete{{.table.Title}}(String where,Map<String,Object> params){
        TinySession s = this.ctx.hibernate();
        int i = s.execute("DELETE FROM {{.table.Name}} WHERE "+where, params);
        s.close();
        return i;
    }
}
