package github.com.jsix.goex.generator.bin.code_templates.java;
// auto generate by gof (http://github.com/ixre/goex)
import java.io.Serializable;
import java.util.List;
import java.util.Map;
import com.gcy.sz.repo.model.*;

/** {{.table.Comment}}仓储 */
public interface I{{.table.Title}}Repo {
    /** 获取{{.table.Comment}} */
    {{.table.Title}}{{.global.entity_suffix}} get(Serializable id);
    /** 根据条件获取单条{{.table.Comment}} */
    {{.table.Title}}{{.global.entity_suffix}} get{{.table.Title}}By(String where,Map<String,Object> params);
    /** 根据条件获取多条{{.table.Comment}} */
    List<{{.table.Title}}{{.global.entity_suffix}}> select{{.table.Title}}(String where,Map<String,Object> params);
    /** 保存{{.table.Comment}} */
    int save{{.table.Title}}({{.table.Title}}{{.global.entity_suffix}} v);
    /** 删除{{.table.Comment}} */
    Error delete{{.table.Title}}(Serializable id);
    /** 批量删除{{.table.Comment}} */
    int BatchDelete{{.table.Title}}(String where,Map<String,Object> params);
}
