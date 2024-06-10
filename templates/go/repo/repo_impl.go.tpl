package impl

#!type:0#!lang:go
#!target:{{.global.pkg}}/repo/impl/{{.table.Name}}_repo_impl.go
{{$title := .table.Title}}
{{$shortTitle := .table.ShortTitle}}
{{$p := substr .table.Name 0 1 }}
{{$structName := join (lower_title $title) "RepoImpl"}}
import(
	"database/sql"
	"fmt"
    "{{pkg "go" .global.pkg}}/repo/model"
    "{{pkg "go" .global.pkg}}/repo"
    "github.com/ixre/gof/db"
    "github.com/ixre/gof/db/orm"
    "log"
)

var _ repo.I{{$title}}Repo = new({{$structName}})
type {{$structName}} struct{
    _orm orm.Orm
}

var {{$structName}}Mapped = false

// New{{$title}}Repo Create new {{$title}}Repo
func New{{$title}}Repo(o orm.Orm) repo.I{{$title}}Repo{
    if !{{$structName}}Mapped{
        _ = o.Mapping(model.{{$shortTitle}}{},"{{.table.Name}}")
        {{$structName}}Mapped = true
    }
    return &{{$structName}}{
        _orm:o,
    }
}
// Get{{$shortTitle}} Get {{.table.Comment}}
func ({{$p}} *{{$structName}}) Get{{$shortTitle}}({{.table.Pk}} {{type "go" .table.PkType}})*model.{{$shortTitle}}{
    e := model.{{$shortTitle}}{}
    err := {{$p}}._orm.Get({{.table.Pk}},&e)
    if err == nil{
        return &e
    }
    if err != sql.ErrNoRows{
      log.Printf("[ Orm][ Error]: %s; Entity:{{$shortTitle}}\n",err.Error())
    }
    return nil
}

// Get{{$shortTitle}}By GetBy {{.table.Comment}}
func ({{$p}} *{{$structName}}) Get{{$shortTitle}}By(where string,v ...interface{})*model.{{$shortTitle}}{
    e := model.{{$shortTitle}}{}
    err := {{$p}}._orm.GetBy(&e,where,v...)
    if err == nil{
        return &e
    }
    if err != sql.ErrNoRows{
      log.Printf("[ Orm][ Error]: %s; Entity:{{$shortTitle}}\n",err.Error())
    }
    return nil
}

// Count{{$shortTitle}} Count {{.table.Comment}} by condition
func ({{$p}} *{{$structName}}) Count{{$shortTitle}}(where string,v ...interface{})(int,error){
   return {{$p}}._orm.Count(model.{{$shortTitle}}{},where,v...)
}

// Select{{$shortTitle}} Select {{.table.Comment}}
func ({{$p}} *{{$structName}}) Select{{$shortTitle}}(where string,v ...interface{})[]*model.{{$shortTitle}} {
    list := make([]*model.{{$shortTitle}},0)
    err := {{$p}}._orm.Select(&list,where,v...)
    if err != nil && err != sql.ErrNoRows{
      log.Printf("[ Orm][ Error]: %s; Entity:{{$shortTitle}}\n",err.Error())
    }
    return list
}

// Save{{$shortTitle}} Save {{.table.Comment}}
func ({{$p}} *{{$structName}}) Save{{$shortTitle}}(v *model.{{$shortTitle}})(int,error){
    id,err := orm.Save({{$p}}._orm,v,int(v.{{title .table.Pk}}))
    if err != nil && err != sql.ErrNoRows{
      log.Printf("[ Orm][ Error]: %s; Entity:{{$shortTitle}}\n",err.Error())
    }
    return id,err
}

// Delete{{$shortTitle}} Delete {{.table.Comment}}
func ({{$p}} *{{$structName}}) Delete{{$shortTitle}}({{.table.Pk}} {{type "go" .table.PkType}}) error {
    err := {{$p}}._orm.DeleteByPk(model.{{$shortTitle}}{}, {{.table.Pk}})
    if err != nil && err != sql.ErrNoRows{
      log.Printf("[ Orm][ Error]: %s; Entity:{{$shortTitle}}\n",err.Error())
    }
    return err
}

// BatchDelete{{$shortTitle}} Batch Delete {{.table.Comment}}
func ({{$p}} *{{$structName}}) BatchDelete{{$shortTitle}}(where string,v ...interface{})(int64,error) {
    r,err := {{$p}}._orm.Delete(model.{{$shortTitle}}{},where,v...)
    if err != nil && err != sql.ErrNoRows{
      log.Printf("[ Orm][ Error]: %s; Entity:{{$shortTitle}}\n",err.Error())
    }
    return r,err
}

// PagingQuery{{$shortTitle}} Query paging data
func ({{$p}} *{{$structName}}) PagingQuery{{$shortTitle}}(begin, end int,where, orderBy string) (total int, rows []map[string]interface{}) {
	if orderBy != "" {
		orderBy = "ORDER BY " + orderBy
	}
	if where == "" {
	    where = "1=1"
	}
	query := fmt.Sprintf(`SELECT COUNT(1) FROM {{.table.Name}} WHERE %s`, where)
	_ = {{$p}}._orm.Connector().ExecScalar(query,&total)
	if total > 0{
	    query = fmt.Sprintf(`SELECT * FROM {{.table.Name}} WHERE %s %s
	        {{if eq .global.db "postgresql"}}LIMIT $2 OFFSET $1{{else}}LIMIT $1,$2{{end}}`,
            where, orderBy)
        err := {{$p}}._orm.Connector().Query(query, func(_rows *sql.Rows) {
            rows = db.RowsToMarshalMap(_rows)
        }, begin, end-begin)
        if err != nil{
            log.Printf("[ Orm][ Error]: %s (table:{{.table.Name}})\n", err.Error())
        }
	}else{
	    rows = make([]map[string]interface{},0)
	}
	return total, rows
}
