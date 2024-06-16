#!target:{{.global.pkg}}/mongo/restful/{{.table.Name}}_restful.go
package restful
{{$title := .table.Title}}
{{$shortTitle := .table.ShortTitle}}
{{$structName := join (lower_title $title) "Resource"}}
{{$p := substr .table.Name 0 1 }}
{{$namePath := name_path .table.Name}}
{{$pk := .table.Pk}}
{{$pkType := join .table.Title .table.PkProp ""}}
{{$ifaceName := join (title .table.Title) "Service"}}

import (
  "github.com/ixre/gof/typeconv"
  "github.com/labstack/echo/v4"
  "{{pkg "go" .global.pkg}}/model"
  "{{pkg "go" .global.pkg}}/service"
  "net/http"
  "github.com/ixre/gof"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)


// {{.table.Comment}}
type {{$structName}} struct{
  s service.I{{$ifaceName}}
}

func New{{title $structName}}(s service.I{{$ifaceName}})*{{$structName}}{
  return &{{$structName}}{
    s:s,
  }
}

func ({{$p}} {{$structName}}) MapRoute(g *echo.Group) {
  // {{.table.Name}} router
  g.GET("/{{$namePath}}/paging",{{$p}}.paging{{$shortTitle}})
  g.GET("/{{$namePath}}/:id",{{$p}}.get{{$shortTitle}})
  g.POST("/{{$namePath}}",{{$p}}.create{{$shortTitle}})
  g.PUT("/{{$namePath}}/:id",{{$p}}.update{{$shortTitle}})
  g.DELETE("/{{$namePath}}/:id",{{$p}}.delete{{$shortTitle}})
}

func ({{$p}} *{{$structName}}) paging{{$shortTitle}}(ctx echo.Context) error {
  page := typeconv.MustInt(ctx.QueryParam("page"))
  size := typeconv.MustInt(ctx.QueryParam("size"))
  options := []*options.FindOptions{}
  filter := bson.M{}
  total,rows,_ := {{$p}}.s.Paging{{$shortTitle}}(filter,page,size,options...)
  return ctx.JSON(http.StatusOK,map[string]interface{}{
    "total": total,
    "rows": rows,
  })
}

func ({{$p}} *{{$structName}}) error(ctx echo.Context, err error) error {
	if err == nil {
		return ctx.JSON(http.StatusOK, gof.Result{ErrCode: 0, ErrMsg: "success"})
	}
	return ctx.JSON(http.StatusOK, gof.Result{ErrCode: 1, ErrMsg: err.Error()})
}

func ({{$p}} *{{$structName}}) get{{$shortTitle}}(ctx echo.Context) error {
  /** #! 转换主键 */
  {{ $goType := type "protobuf" .table.PkType}}\
  {{if equal_any .table.PkType 3 4 5}}{{$pk}} := typeconv.MustInt(ctx.Param("id"))\
  {{else}}{{$pk}} := ctx.Param("id"){{end}}
  ret := {{$p}}.s.Get{{$shortTitle}}({{$pk}})
  return ctx.JSON(http.StatusOK,ret)
}

func ({{$p}} *{{$structName}}) create{{$shortTitle}}(ctx echo.Context) error {
  dst := model.{{$shortTitle}}{}
  err := ctx.Bind(&dst)
  if err == nil{
    err = {{$p}}.s.Save{{$shortTitle}}(&dst)
  }    
  return {{$p}}.error(ctx, err)
}

func ({{$p}} *{{$structName}}) update{{$shortTitle}}(ctx echo.Context) error {
  dst := model.{{$shortTitle}}{}
  err := ctx.Bind(&dst)
  if err == nil{
    /** #! 转换主键 */
    {{ $goType := type "protobuf" .table.PkType}}\
    {{if equal_any .table.PkType 3 4 5}}{{$pk}} := typeconv.MustInt(ctx.Param("id"))\
    {{else}}{{$pk}} := ctx.Param("id"){{end}}
    dst.{{.table.PkProp}} = {{$pk}}
    err = {{$p}}.s.Save{{$shortTitle}}(&dst)
  }
  return {{$p}}.error(ctx, err)
}

func ({{$p}} *{{$structName}}) delete{{$shortTitle}}(ctx echo.Context) error {
    /** #! 转换主键 */
    {{ $goType := type "protobuf" .table.PkType}}\
    {{if equal_any .table.PkType 3 4 5}}{{$pk}} := typeconv.MustInt(ctx.Param("id"))\
    {{else}}{{$pk}} := ctx.Param("id"){{end}}
    err := {{$p}}.s.Delete{{$shortTitle}}({{$pk}})
    return {{$p}}.error(ctx, err)
}